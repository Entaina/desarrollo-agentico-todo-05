package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"go.yaml.in/yaml/v3"
)

// Etiquetas de control de la skill adw: ver .claude/skills/adw/references/adapters/github.md.
const (
	waitingLabel = "adw:waiting-for-human"
	stagePrefix  = "stage:"
)

// defaultMaxRejections es el tope de rechazos seguidos de la verificación si config/adw.yml no lo fija.
const defaultMaxRejections = 3

// config refleja config/adw.yml. MaxRejections no lo usa el trigger, sino la skill, que lo lee
// del mismo fichero; está aquí para que KnownFields(true) lo admita y para validarlo al arrancar.
type config struct {
	Repo            string `yaml:"repo"`
	Label           string `yaml:"label"`
	IntervalSeconds int    `yaml:"interval_seconds"`
	MaxRejections   int    `yaml:"max_rejections"`
	Prompt          string `yaml:"prompt"`
}

type issue struct {
	URL    string `json:"url"`
	Labels []struct {
		Name string `json:"name"`
	} `json:"labels"`
	Assignees []struct {
		Login string `json:"login"`
	} `json:"assignees"`
}

func readConfig(path string) (config, error) {
	// El decodificador solo cambia los campos presentes: si se omite, queda el valor por defecto.
	cfg := config{MaxRejections: defaultMaxRejections}
	file, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	if err := decoder.Decode(&cfg); err != nil {
		return cfg, fmt.Errorf("configuración YAML: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return cfg, fmt.Errorf("la configuración debe contener un solo documento YAML")
	}
	if strings.TrimSpace(cfg.Label) == "" || cfg.IntervalSeconds < 1 {
		return cfg, fmt.Errorf("label no puede estar vacío e interval_seconds debe ser positivo")
	}
	if cfg.MaxRejections < 1 {
		return cfg, fmt.Errorf("max_rejections debe ser mayor o igual que 1")
	}
	if !strings.Contains(cfg.Prompt, "{{payload}}") {
		return cfg, fmt.Errorf("prompt debe contener {{payload}}")
	}
	return cfg, nil
}

// agentTurn dice si el turno del caso es del agente: nadie espera al humano,
// no hay un stage abierto y la issue está libre o asignada a esta cuenta.
func agentTurn(it issue, me string) bool {
	for _, label := range it.Labels {
		if label.Name == waitingLabel || strings.HasPrefix(label.Name, stagePrefix) {
			return false
		}
	}
	for _, assignee := range it.Assignees {
		if !strings.EqualFold(assignee.Login, me) {
			return false
		}
	}
	return true
}

// say escribe una línea del resumen, con la hora y, si es de un caso, su número.
func say(number string, format string, args ...any) {
	fmt.Print(time.Now().Format("15:04:05"), " ")
	if number != "" {
		fmt.Printf("[#%s] ", number)
	}
	fmt.Printf(format+"\n", args...)
}

func run(args []string) error {
	flags := flag.NewFlagSet("adw", flag.ContinueOnError)
	path := flags.String("config", "config/adw.yml", "fichero YAML de configuración")
	once := flags.Bool("once", false, "procesa una sola consulta y termina")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("configura el prompt y las opciones en %s", *path)
	}
	cfg, err := readConfig(*path)
	if err != nil {
		return err
	}
	for _, tool := range []string{"gh", "claude", "git"} {
		if _, err := exec.LookPath(tool); err != nil {
			return fmt.Errorf("no se encuentra %s en PATH", tool)
		}
	}
	login := exec.Command("gh", "api", "user", "--jq", ".login")
	login.Stderr = os.Stderr
	output, err := login.Output()
	if err != nil {
		return fmt.Errorf("cuenta de gh: %w", err)
	}
	me := strings.TrimSpace(string(output))
	top, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return fmt.Errorf("raíz del checkout: %w", err)
	}
	root := strings.TrimSpace(string(top))

	for {
		queryArgs := []string{"issue", "list"}
		if repo := strings.TrimSpace(cfg.Repo); repo != "" {
			queryArgs = append(queryArgs, "--repo", repo)
		}
		queryArgs = append(queryArgs, "--state", "open", "--label", cfg.Label,
			"--search", "sort:created-desc", "--limit", "100", "--json", "url,number,title,labels,assignees,updatedAt")
		query := exec.Command("gh", queryArgs...)
		query.Stderr = os.Stderr
		output, err := query.Output()
		if err != nil {
			return fmt.Errorf("consulta de issues: %w", err)
		}
		var snapshots []json.RawMessage
		if len(bytes.TrimSpace(output)) > 0 {
			if err := json.Unmarshal(output, &snapshots); err != nil {
				return fmt.Errorf("respuesta de la consulta de issues: %w", err)
			}
		}
		issues := make([]issue, len(snapshots))
		turn := 0
		for i, snapshot := range snapshots {
			if err := json.Unmarshal(snapshot, &issues[i]); err != nil {
				return fmt.Errorf("issue de la consulta: %w", err)
			}
			if agentTurn(issues[i], me) {
				turn++
			}
		}
		// Cada consulta deja una línea en el resumen, aunque no encuentre nada: así se ve que el
		// trigger sigue vivo.
		next := ""
		if !*once {
			next = fmt.Sprintf(" · siguiente en %d s", cfg.IntervalSeconds)
		}
		say("", "consulta: %d issues abiertas con %q · %d por lanzar%s", len(issues), cfg.Label, turn, next)
		for i, it := range issues {
			if !agentTurn(it, me) {
				continue
			}
			snapshot := snapshots[i]
			status, err := exec.Command("git", "status", "--porcelain").Output()
			if err != nil {
				return fmt.Errorf("estado del checkout: %w", err)
			}
			if len(bytes.TrimSpace(status)) > 0 {
				return fmt.Errorf("el checkout tiene cambios sin registrar: déjalo limpio antes de lanzar %s", it.URL)
			}
			var payload bytes.Buffer
			if err := json.Compact(&payload, snapshot); err != nil {
				return err
			}
			prompt := strings.ReplaceAll(cfg.Prompt, "{{payload}}", payload.String())
			agent := exec.Command("claude", "-p", prompt, "--permission-prompts", "none")
			// Sin tope de espera a los subagentes en segundo plano: si el coordinador termina su turno
			// mientras uno trabaja, claude -p corta la sesión a los 10 minutos y deja el caso a medias.
			agent.Env = append(os.Environ(), "CLAUDE_CODE_PRINT_BG_WAIT_CEILING_MS=0")
			// Lo que escribe la sesión va a su log, junto a los de los hooks: la salida del trigger es
			// solo el resumen. Stdin queda cerrado: no hay conversación ni evaluación de shell.
			number := it.URL[strings.LastIndex(it.URL, "/")+1:]
			logs := filepath.Join(".adw", "logs", number)
			if err := os.MkdirAll(filepath.Join(root, logs), 0o755); err != nil {
				return err
			}
			logPath := filepath.Join(logs, "trigger-"+time.Now().Format("20060102-150405")+".txt")
			logFile, err := os.Create(filepath.Join(root, logPath))
			if err != nil {
				return err
			}
			agent.Stdout, agent.Stderr = logFile, logFile
			say(number, "lanzado · su salida, en %s", logPath)
			started := time.Now()
			err = agent.Run()
			logFile.Close()
			elapsed := time.Since(started).Round(time.Second)
			if err != nil {
				say(number, "la sesión falla tras %s: %v (ver %s)", elapsed, err, logPath)
				return fmt.Errorf("Claude para %s: %w", it.URL, err)
			}
			say(number, "terminado en %s", elapsed)
		}
		if *once {
			return nil
		}
		if err := exec.Command("sleep", strconv.Itoa(cfg.IntervalSeconds)).Run(); err != nil {
			return err
		}
	}
}

func exitCode(err error) int {
	var process *exec.ExitError
	if errors.As(err, &process) && process.ExitCode() > 0 {
		return process.ExitCode()
	}
	return 2
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(exitCode(err))
	}
}
