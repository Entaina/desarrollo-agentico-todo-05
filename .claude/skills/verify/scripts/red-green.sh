#!/usr/bin/env bash
# Ejecuta tests nuevos sobre la base y sobre el cambio: la evidencia de que podían fallar.
#
#   red-green.sh --out <directorio en .adw/> --name <NN_afirmación-detalle> --base <ref> \
#                --run <patrón de -run> [--extra <fichero>]... <paquete> <fichero de test> ...
#
# Crea un worktree desprendido de <ref> en .adw/base/, copia allí los ficheros de test
# indicados y los --extra que necesiten (fixtures, testdata), tal como están en la copia
# de trabajo, y ejecuta en los dos lados
#   GO_ENV=test go test -v -p 1 -count=1 -run <patrón> <paquete>
# Guarda cada ejecución en <out>/<name>-base.txt y <out>/<name>-cambio.txt y borra el worktree.
#
# Imprime los tests seleccionados, el resultado de cada lado y el veredicto:
#   RED-GREEN: DISCRIMINANTE                rojo en la base porque fallan todos los tests
#                                           seleccionados, y verde en el cambio
#   RED-GREEN: DISCRIMINANTE (compilación)  la base no compila con los tests: vale si lo que falta
#                                           es el comportamiento que pide el criterio
#   RED-GREEN: NO DISCRIMINANTE (<motivo>)
# Sale con 0 si es discriminante, 1 si no, 2 si el uso es incorrecto y 3 si no se puede preparar la base.
set -o pipefail

usage() {
  echo "uso: red-green.sh --out <directorio en .adw/> --name <NN_afirmación-detalle> --base <ref> --run <patrón> [--extra <fichero>]... <paquete> <fichero de test> ..." >&2
  exit 2
}

root=$(git rev-parse --show-toplevel) || exit 2
cd "$root" || exit 2

out="" name="" base="" run=""
extras=()
while (($#)); do
  case $1 in
    --out | --name | --base | --run | --extra)
      [[ $# -ge 2 ]] || usage
      case $1 in
        --out) out=${2%/} ;;
        --name) name=$2 ;;
        --base) base=$2 ;;
        --run) run=$2 ;;
        --extra) extras+=("$2") ;;
      esac
      shift 2
      ;;
    -*) usage ;;
    *) break ;;
  esac
done
pkg=${1:-}
[[ -n $out && -n $name && -n $base && -n $run && -n $pkg ]] || usage
shift
(($#)) || usage

out=${out#"$root"/}
[[ $out == .adw/* && $out != *..* ]] || { echo "red-green.sh: --out debe estar dentro de .adw/" >&2; exit 2; }
[[ $name =~ ^[0-9][0-9]_[a-z0-9-]+$ ]] || { echo "red-green.sh: --name debe tener la forma NN_afirmación-detalle" >&2; exit 2; }
[[ $pkg =~ ^\./[A-Za-z0-9_./-]*$ && $pkg != *..* && -d $pkg ]] || { echo "red-green.sh: el paquete debe ser un directorio del repositorio, como ./actions" >&2; exit 2; }
commit=$(git rev-parse --verify --quiet "$base^{commit}") || { echo "red-green.sh: $base no es un commit" >&2; exit 2; }
for test in "$@"; do
  [[ $test == *_test.go && $test != *..* && $test != /* && -f $test ]] || { echo "red-green.sh: $test no es un fichero de test del repositorio" >&2; exit 2; }
done
for extra in ${extras[@]+"${extras[@]}"}; do
  [[ $extra != *..* && $extra != /* && $extra != .* && -f $extra ]] || { echo "red-green.sh: $extra no es un fichero del repositorio" >&2; exit 2; }
done

mkdir -p "$out" .adw/base || exit 3
# Worktrees que dejó una ejecución interrumpida.
git worktree prune
worktree="$root/.adw/base/$(basename "$out")-$$"
trap 'git worktree remove --force "$worktree" > /dev/null 2>&1' EXIT
git worktree add --detach --quiet "$worktree" "$commit" || { echo "red-green.sh: no se puede crear el worktree de la base" >&2; exit 3; }
for file in "$@" ${extras[@]+"${extras[@]}"}; do
  mkdir -p "$worktree/$(dirname "$file")" && cp "$file" "$worktree/$file" || exit 3
done

shown="GO_ENV=test go test -v -p 1 -count=1 -run '$run' $pkg"

# side <etiqueta> <directorio> <fichero>: ejecuta los tests y guarda la salida. Devuelve su código.
side() {
  local label=$1 dir=$2 file=$3 start output code
  start=$SECONDS
  output=$(cd "$dir" && env GO_ENV=test go test -v -p 1 -count=1 -run "$run" "$pkg" 2>&1)
  code=$?
  {
    echo "\$ $shown"
    echo "en: $label"
    echo "exit: $code"
    echo "duración: $((SECONDS - start))s"
    echo "---"
    printf '%s\n' "$output"
  } > "$file"
  return "$code"
}

base_file="$out/$name-base.txt"
change_file="$out/$name-cambio.txt"
side "base $(git rev-parse --short "$commit") con los tests del cambio" "$worktree" "$base_file"
base_code=$?
side "cambio (copia de trabajo)" "$root" "$change_file"
change_code=$?

# Tests seleccionados: los que se ejecutan a la profundidad del patrón (Suite/Test_X son dos niveles).
depth=$(($(printf '%s' "$run" | tr -cd '/' | wc -c) + 1))
selected=$(grep -E '^ *=== RUN' "$change_file" | awk '{print $3}' | awk -F/ -v d="$depth" 'NF == d' | sort -u)

echo "tests seleccionados: ${selected:-ninguno}" | tr '\n' ' '
echo
echo "base:   exit=$base_code $base_file"
echo "cambio: exit=$change_code $change_file"

if grep -qE '\[(build|setup) failed\]' "$change_file"; then
  echo "RED-GREEN: NO DISCRIMINANTE (el cambio no compila con los tests)"
  exit 1
fi
if [[ -z $selected ]]; then
  echo "RED-GREEN: NO DISCRIMINANTE (el patrón no selecciona ningún test)"
  exit 1
fi
for test in $selected; do
  if ! grep -qF -- "--- PASS: $test " "$change_file"; then
    echo "RED-GREEN: NO DISCRIMINANTE ($test no pasa en el cambio)"
    exit 1
  fi
done
if ((change_code != 0)); then
  echo "RED-GREEN: NO DISCRIMINANTE (el cambio no está en verde)"
  exit 1
fi
if ((base_code == 0)); then
  echo "RED-GREEN: NO DISCRIMINANTE (los tests también pasan en la base)"
  exit 1
fi
if grep -qE '\[(build|setup) failed\]' "$base_file"; then
  echo "RED-GREEN: DISCRIMINANTE (compilación)"
  exit 0
fi
for test in $selected; do
  if ! grep -qF -- "--- FAIL: $test " "$base_file"; then
    echo "RED-GREEN: NO DISCRIMINANTE ($test no falla en la base: revisa $base_file)"
    exit 1
  fi
done
# Un rojo porque a la base le falta un fixture o un fichero no demuestra el comportamiento.
if grep -qiE 'could not find (a )?scenario|no such file or directory' "$base_file"; then
  echo "RED-GREEN: NO DISCRIMINANTE (la base falla por un fichero que no tiene: pásalo con --extra)"
  exit 1
fi
echo "RED-GREEN: DISCRIMINANTE"
exit 0
