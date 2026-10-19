---
name: bug
description: Investiga un bug y escribe un plan autocontenido en plans/
allowed-tools: Bash, Read, Glob, Grep, Write
---

# Planifica un bug

Recibe en `$ARGUMENTS` la descripción de un bug. Tu único propósito es investigar y escribir un plan; no implementes la corrección.

1. Comprueba `git status --short --branch` y lee el código, los tests y la configuración relacionados.
2. Consulta el índice condicional `docs/changes/README.md` y lee los documentos cuyas condiciones encajen con el encargo: explican cambios anteriores sobre el mismo código.
3. Reconstruye el comportamiento esperado y el observado a partir de evidencia del repositorio. No presentes una conjetura como un hecho.
4. Define una reproducción mínima y el test de regresión que debe fallar antes del arreglo. Si el fallo está en el navegador, como el JavaScript de `public/assets/`, el test de regresión es un test de navegador en `e2e/`: un test de Go no ejecuta ese código y pasaría también antes del arreglo.
5. Delimita el cambio mínimo que corrige la causa raíz sin alterar comportamientos ajenos.
6. Escribe un único plan Markdown autocontenido en `plans/`. Usa un nombre kebab-case descriptivo, salvo que el argumento indique una ruta concreta dentro de `plans/`; en ese caso usa exactamente esa ruta.

El plan es el único archivo que puedes crear o modificar. No cambies código de aplicación, tests, configuración ni dependencias.

El plan debe incluir:

- `# Bug: …`
- `## Encargo y resultado esperado`
- `## Evidencia y causa raíz` (hechos observados, incertidumbres y reproducción)
- `## Alcance` (incluido y excluido)
- `## Archivos relevantes` (ruta y responsabilidad)
- `## Implementación paso a paso`
- `## Tests y validación`
- `## Riesgos y recuperación`

Usa solo comandos y herramientas que existan en este repositorio. La validación habitual de la aplicación es comprobar formato con `gofmt -l .`, ejecutar `go vet ./...`, `go test -p 1 ./...`, `go build ./...` y, para los tests de navegador, `npx playwright test`; los tests necesitan PostgreSQL y las migraciones de test descritas en `database.yml` y `.github/workflows/ci.yaml`. Incluye únicamente las comprobaciones pertinentes al cambio investigado.

Al terminar, informa solo de la ruta del plan, sus decisiones principales y las incertidumbres pendientes. Sugiere `/implement <ruta>` como siguiente paso.
