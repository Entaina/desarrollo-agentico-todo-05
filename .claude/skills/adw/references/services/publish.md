---
name: publish
description: Registra los cambios del caso y publica la propuesta para revisión, o la actualiza si ya existe.
---

# publish

## Entradas

- Los cambios de la copia de trabajo, con el documento de cierre, y el informe de `implement`.
- El plan vigente.
- El acta vigente de la verificación, con el cambio que validó, y el plan de verificación vigente.
- El tipo de trabajo.

## Qué hace

1. Registra todos los cambios con un mensaje convencional según el tipo, documento incluido. El checkout estaba limpio al empezar, así que todos los cambios son del caso.
2. Comprueba que lo registrado es lo que se validó: respecto al cambio validado del acta solo pueden diferir el documento de cierre y su índice. Si difiere algo más, no publiques: es un fallo, porque la verificación ya no vale.
3. Publícalos.
4. Si la propuesta no existe, créala. Su descripción la vincula a la petición e incluye un resumen, el enlace al plan vigente, el informe de `implement`, con los comandos que fallaron si los hubo, una sección de evidencia y la lista de puntos para el humano. Después avisa en la petición.
5. Si la propuesta ya existe, actualízala con los cambios y renueva en su descripción la sección de evidencia y la lista de puntos para el humano. Si guardaba trabajo sin validar, déjala lista para revisión.

La sección de evidencia resume el acta, enlaza el acta y el plan de verificación y dice en qué commit se validó el cambio, con las capturas y vídeos que el acta ya publicó. La lista de puntos para el humano recoge los `should-fix`, los juicios de confianza baja y las verificaciones humanas.

## Resultado

La propuesta publicada con todos los cambios del caso.
