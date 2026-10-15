# GSM: cómo se lee un ciclo de vida

Cada artefacto describe su ciclo de vida con el modelo Guard-Stage-Milestone (GSM). No hay estados ni transiciones dibujadas. El estado de un artefacto es qué milestones tiene alcanzados y qué stages tiene abiertos, y el proceso avanza porque se cumplen condiciones sobre ese estado.

## Elementos

- **Milestone:** un resultado con significado de negocio, como `clarified` o `published`. Está alcanzado o no. Se alcanza cuando se cumple su condición de logro y se invalida cuando se cumple su condición de invalidación.
- **Stage:** una zona de trabajo que persigue uno o varios milestones ejecutando un servicio. Está abierto o cerrado.
- **Guard:** la condición que abre un stage.
- **Sentry:** la forma de cualquier condición: «al ocurrir un evento, si se cumple una condición». El evento es opcional. Guards, condiciones de logro y condiciones de invalidación son sentries.

## Contrato

Un artefacto solo conoce su propio estado. Lo que necesita de fuera lo declara en su contrato:

- **Entradas:** lo que recibe de otros artefactos, con un nombre abstracto (`specification`, `change_request`).
- **Expone:** los milestones y la información que otros artefactos pueden usar.

Sus sentries solo miran su propio estado y sus entradas. A qué artefacto concreto corresponde cada entrada lo deciden las [asociaciones](associations.md).

## Reglas

1. Un stage cerrado se abre cuando se cumple su guard.
2. Al abrirse, sus milestones dejan de estar alcanzados: el stage existe para volver a conseguirlos.
3. Un stage se cierra al alcanzar uno de sus milestones.
4. Un milestone alcanzado se invalida cuando se cumple su condición de invalidación. La invalidación se propaga: lo que dependía de él, en el mismo artefacto o a través de las asociaciones, se reevalúa.
5. Si el servicio de un stage falla, el stage se cierra alcanzando `failed`, un milestone que tienen todos los artefactos. También se cierra con `failed` cuando el artefacto agota un tope que fija su definición, como los rechazos seguidos de la Verification.
6. Algunos milestones son terminales: `failed` y los que cada artefacto marque como tales. Mientras un artefacto del caso tenga uno alcanzado, no se abre ningún stage del caso. `failed` se invalida cuando el humano devuelve el turno (ver [control](control.md)): es un reintento.

## Business step

Ante cada evento, el motor aplica las reglas hasta que no cambia nada. Primero reevalúa invalidaciones y su propagación, después logros y por último guards. Abre los stages que tocan, ejecuta sus servicios y vuelve a empezar con el resultado. Eso es una invocación.

## Quiescencia

Cuando no hay stages abiertos ni guards que se cumplan, el caso queda quieto hasta el siguiente evento externo. Si el caso no ha terminado, le toca al humano: el motor le cede el turno (ver [control](control.md)).
