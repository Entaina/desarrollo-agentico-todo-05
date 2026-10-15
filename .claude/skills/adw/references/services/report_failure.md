---
name: report_failure
description: Comunica al humano que un stage ha fallado o ha agotado su tope, y cómo reintentar.
---

# report_failure

Lo invoca el motor cuando el servicio de un stage falla o el stage agota su tope.

## Entradas

- El stage que ha fallado y el error.
- Si el stage agotó su tope: el valor del tope y los intentos de la racha, con sus actas.

## Qué hace

Comunica en la petición de trabajo qué stage ha fallado, por qué y cómo reintentar. No intentes arreglarlo: decidir qué hacer le corresponde al humano.

Si el stage agotó su tope, dilo, indica su valor y resume los rechazos de la racha con el enlace a cada acta. Explica que al reintentar el contador vuelve a cero y el último rechazo llega al plan junto con lo que el humano escriba antes de devolver el turno.

## Resultado

El fallo comunicado al humano.
