---
name: clarify
description: Comprueba si una petición de trabajo está clara y clasificada. Si no lo está, publica las preguntas que bloquean el trabajo.
---

# clarify

## Entradas

- La petición de trabajo: título, descripción, comentarios y tipo, si ya lo tiene.
- `questions`: dudas funcionales planteadas desde otros artefactos, si las hay.

## Qué hace

1. Lee la petición completa y las dudas recibidas.
2. Decide si está clara: sabes qué comportamiento se quiere y cómo reconocer que está hecho. Las decisiones técnicas son tuyas: no preguntes por ellas.
3. Si está clara y no tiene tipo, clasifícala como `feature`, `bug` o `chore`. Si el humano ya la clasificó, respeta su tipo.
4. Si no está clara, publica solo las preguntas que bloquean el trabajo, numeradas y concretas. No repitas las que ya tienen respuesta e incluye las dudas recibidas que sigan abiertas.
5. Si decides tú una duda recibida, porque es técnica o la petición ya la responde, publica la decisión en la petición: es una aclaración más, y quien planteó la duda la leerá al volver a planificar.

## Resultado

La petición clara y clasificada, con las dudas recibidas resueltas por escrito, o las preguntas publicadas.
