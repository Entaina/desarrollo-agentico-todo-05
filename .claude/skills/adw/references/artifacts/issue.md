# Issue

La petición de trabajo. Es el artefacto raíz del caso.

## Información

- Título y descripción.
- Preguntas pendientes y respuestas del humano.
- Tipo de trabajo: `feature`, `bug` o `chore`. Si el humano no lo ha indicado, lo decide el agente.

## Contrato

- **Entradas:** `questions`, dudas funcionales que plantean otros artefactos al planificar el cambio o su verificación.
- **Expone:** `clarified`, y su identidad, tipo y título.

## Ciclo de vida

| Stage | Guard | Servicio | Milestones |
| --- | --- | --- | --- |
| `Clarification` | Sin `clarified` ni `info_requested`; o con `info_requested` y una respuesta humana posterior a la última pregunta | [`clarify`](../services/clarify.md) | `clarified`, `info_requested` |

| Milestone | Se alcanza cuando | Se invalida cuando |
| --- | --- | --- |
| `clarified` | La issue está clasificada y dice qué comportamiento se quiere y cómo reconocer que está hecho, sin dudas pendientes | Un humano cambia el alcance, o llega una duda por `questions` |
| `info_requested` | El agente ha publicado las preguntas que le bloquean | Se reabre `Clarification` |
| `closed` | La issue se cierra. Es terminal | — |
