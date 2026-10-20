# Desarrollo

Cómo se arranca la aplicación, cómo se prueba y qué comprueba la CI.

## Requisitos

- Docker, para PostgreSQL y, si se quiere, para la aplicación.
- Go 1.26 y Node.js LTS.

Lo demás lo instala `bin/setup`: la CLI de Buffalo (`buffalo@v0.18.14`) y `soda` (`soda@v6.4.1`), con `go install` si faltan, las dependencias de Node (`npm ci`) y el Chromium de los tests de navegador. El devcontainer (`.devcontainer/devcontainer.json`) lo instala todo con su propia lista de comandos y arranca la base de datos.

## Preparar, arrancar y parar

```bash
bin/setup        # deja el repositorio listo: dependencias, PostgreSQL y la base de test migrada
bin/dev          # arranca la aplicación en primer plano; Ctrl+C la para
```

`bin/setup` se puede repetir: solo hace lo que falta. `bin/dev` levanta PostgreSQL si no está y arranca la aplicación en uno de los dos modos:

| Modo | Arrancar | Parar |
| --- | --- | --- |
| Host (por defecto) | `bin/dev`: PostgreSQL en su contenedor y `buffalo dev`, que recompila al cambiar el código Go | Ctrl+C para la aplicación; `docker compose stop db`, PostgreSQL |
| Docker completo | `bin/dev --docker`, o `bin/dev` si no está `buffalo`: la aplicación también en su contenedor | Ctrl+C para la aplicación; `docker compose stop`, todo |

`/env:start` hace lo mismo para el agente (`bin/setup` y `bin/dev` en segundo plano) y comprueba que responde; `/env:stop` lo para sin borrar datos.

La aplicación queda en http://localhost:3000, y `/healthz` dice si responde. Usa los puertos 3000 y 5432: si los ocupa otro proyecto, `bin/dev` lo dice y para, y hay que pararlo antes.

## Bases de datos

| Base | Para qué | Cómo se crea |
| --- | --- | --- |
| `todo_development` | La aplicación en desarrollo | La crea `compose.yaml` |
| `todo_test` | Los tests de Go y de navegador | `bin/setup`, que la crea y la migra (`soda create -e test` y `soda migrate -e test`) |

Al arrancar, la aplicación aplica las migraciones pendientes de su entorno. Los tests de Go necesitan la base de test migrada antes: `bin/setup` o `soda migrate -e test`.

## Tests de Go

```bash
GO_ENV=test go test -p 1 ./...
```

`-p 1` porque todos los paquetes comparten la base de test. Las suites están en `actions/` (peticiones y pantallas) y `models/` (reglas del modelo).

## Tests de navegador

```bash
npx playwright test
npx playwright test e2e/tareas/anadir-una-tarea.spec.ts
```

Playwright arranca la aplicación con `GO_ENV=test` en `127.0.0.1:3001` (otro puerto con `E2E_PORT`), contra `todo_test`, y la para al terminar. Cada test empieza con la lista vacía y crea las tareas que necesita con la fixture `tasks` de `e2e/fixtures.ts`. Si un test falla, el informe queda en `playwright-report/` y la traza en `test-results/`.

## Gates

```bash
bash .claude/skills/check/scripts/check.sh
```

Ejecuta en local lo mismo que la CI, en su orden: `gofmt`, `go vet`, migraciones (borra y vuelve a crear `todo_test`), tests de Go, `go build`, tests de navegador y `docker build`. El último solo informa: en local no cambia el veredicto. Se puede pedir solo algunos gates, por ejemplo `check.sh migrate test`.

## CI

`.github/workflows/ci.yaml` se ejecuta en cada PR y en cada push a `main`, con tres jobs:

| Job | Qué hace |
| --- | --- |
| `test` | Formato, `go vet`, migraciones, tests de Go y build, con PostgreSQL como servicio |
| `e2e` | Los tests de navegador. Si fallan, sube el informe como artefacto |
| `docker` | Construye la imagen |

## El workflow `adw`

El trabajo sobre issues de GitHub lo puede hacer el workflow `adw`, que planifica, implementa, verifica con evidencias, documenta y abre el PR. Cómo se configura y se lanza está en [adws/README.md](../../adws/README.md).
