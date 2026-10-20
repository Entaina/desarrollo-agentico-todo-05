---
description: Start the TODO development environment and verify that the application responds
allowed-tools: Bash, Read
---

# Arranca el entorno de desarrollo

Levanta la aplicación del repositorio con los scripts del proyecto, `bin/setup` y `bin/dev`, y comprueba que responde. No mates procesos existentes ni elimines volúmenes.

La aplicación usa el puerto `3000` y PostgreSQL, el `5432`.

## Pasos

1. Si `http://localhost:3000/healthz` ya responde y `docker compose ps` muestra en marcha los servicios de este proyecto, el entorno ya está arrancado: informa y termina.
2. Ejecuta `bin/setup`. Comprueba las herramientas, instala lo que falta y deja PostgreSQL en marcha con la base de test migrada. Si falla, informa de lo que dice y detente: es un fallo del entorno.
3. Lanza `bin/dev` en segundo plano, con la herramienta Bash en segundo plano, y conserva su salida para poder diagnosticar el arranque. Usa `bin/dev --docker` si `$ARGUMENTS` contiene `--docker`. Sin buffalo, `bin/dev` pasa solo a Docker completo y lo dice.
4. Si `bin/dev` termina enseguida con un error, informa de lo que dice y detente. Si el error es que otro proceso ocupa un puerto, no lo finalices.
5. Espera hasta que `http://localhost:3000/healthz` responda correctamente. En modo Docker completo, la primera vez se construye la imagen y tarda más.

## Resultado

Informa de:

- modo utilizado, el que dice `bin/dev`;
- servicios iniciados;
- URL de la aplicación;
- respuesta de `/healthz`;
- cualquier fallback, advertencia o proceso que deba mantenerse activo: `bin/dev` sigue en segundo plano.
