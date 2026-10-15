# Documentation

El documento de cierre del caso: qué se construyó, cómo se usa, qué configuración necesita, cómo se probó y qué conviene saber para trabajar sobre ello. Se escribe en la rama y viaja dentro de la propuesta: se revisa y se integra con el código, no después.

## Información

- El documento del caso. Hay uno por caso: si el cambio vuelve a pasar por la verificación, se actualiza el mismo.
- Su entrada en el índice condicional, con las condiciones en que conviene cargarlo al planificar otro trabajo.

## Contrato

- **Entradas:**
  - `work_item`: la petición de trabajo (identidad, tipo, título y criterios).
  - `change`: el cambio que se documenta y el acta de su verificación, y si está listo.
  - `environment`: dónde se escribe el documento, y si la copia de trabajo está en la rama del caso.
- **Expone:** `documented`.

## Ciclo de vida

| Stage | Guard | Servicio | Milestones |
| --- | --- | --- | --- |
| `Documenting` | El `change` está listo, la copia de trabajo está en la rama del caso y la Documentation no tiene `documented` | [`document`](../services/document.md) | `documented` |

| Milestone | Se alcanza cuando | Se invalida cuando |
| --- | --- | --- |
| `documented` | El documento, que cita el acta vigente, y su entrada en el índice existen y describen el diff vigente | El `change` deja de estar listo |
