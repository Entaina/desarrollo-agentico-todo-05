# Plan

La especificación de un cambio de código. Todo cambio de código parte de un plan.

## Información

- Alcance, pasos, archivos y tests, como los producen las skills `feature`, `bug` y `chore`.
- Revisiones: el plan inicial y una por cada petición de cambio. Las anteriores se conservan, y cada revisión cita las peticiones que recoge y su origen: «revisión del humano» o «rechazo de validación, intento k», con los defectos de cada dimensión.

## Contrato

- **Entradas:**
  - `specification`: lo que hay que construir, y si está clara.
  - `change_request`: peticiones de cambio sobre el plan vigente: las del humano al revisar la propuesta y los rechazos de la validación.
  - `environment`: dónde se trabaja, y si está listo.
- **Expone:** `ready`, `question_raised`.

## Ciclo de vida

| Stage | Guard | Servicio | Milestones |
| --- | --- | --- | --- |
| `Planning` | La `specification` está clara, el `environment` está listo y el Plan no tiene `ready` ni `question_raised` | [`write_plan`](../services/write_plan.md) | `ready`, `question_raised` |

| Milestone | Se alcanza cuando | Se invalida cuando |
| --- | --- | --- |
| `ready` | Hay un plan, o una revisión, que cubre la `specification` y todas las `change_request` recibidas | Llega una `change_request`, o la `specification` deja de estar clara |
| `question_raised` | El planificador encuentra una duda funcional que la `specification` no resuelve | La `specification` vuelve a estar clara |

Ninguna revisión del plan necesita aprobación humana.
