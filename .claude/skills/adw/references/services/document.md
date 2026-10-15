---
name: document
description: Escribe o actualiza el documento de cierre del caso y su entrada en el índice condicional, a partir del diff real y del acta de la verificación.
skills: [document]
---

# document

## Entradas

- `work_item`: identidad, tipo y título de la petición.
- `change`: el diff de la rama respecto a su base, incluidos los cambios sin registrar.
- El acta vigente de la verificación, para la sección de pruebas.
- La ruta del documento del caso, que da el adaptador.

## Qué hace

1. Ejecuta la skill `document` (`.claude/skills/document/SKILL.md`) con la petición, el diff, el acta y la ruta.
2. Si el caso ya tiene documento, actualízalo para que describa el diff vigente. No crees otro.
3. No registres ni publiques los cambios: es trabajo de `publish`.
4. Comprueba que el documento y su entrada en el índice existen y que no has tocado ningún otro fichero.

## Resultado

El documento y su entrada en el índice, escritos en la copia de trabajo.
