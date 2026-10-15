# Adaptador: GitHub

La definición dice qué artefactos hay, cómo evolucionan y qué servicios existen. Este adaptador dice cómo se representan, se observan, se proyectan y se materializan en GitHub, en el repositorio Git que aloja y en su copia de trabajo local.

## Representación

| Artefacto | En GitHub | Identidad |
| --- | --- | --- |
| Issue | Una issue. Su tipo es la etiqueta `feature`, `bug` o `chore` | Su URL |
| Environment | La rama del caso, registrada en el comentario de control, y la copia de trabajo local en esa rama | El nombre de la rama |
| Plan | Comentarios en la issue que empiezan por `<!-- adw:plan -->`, uno por revisión. El vigente es el más reciente | La URL del comentario |
| Change | La rama del Environment y su PR, cuya descripción incluye `Closes <URL de la issue>`. Un PR en borrador guarda trabajo sin validar: no está publicado para revisión | La rama y la URL del PR |
| Verification · plan | Comentarios en la issue que empiezan por `<!-- adw:verification-plan -->`, uno por revisión. El vigente es el más reciente | La URL del comentario |
| Verification · acta | Comentarios en la issue que empiezan por `<!-- adw:evidence -->`, uno por intento, con sus capturas y vídeos adjuntos. La vigente es la más reciente | La URL del comentario |
| Verification · ficheros | `.adw/evidence/issue-N/<n>/` en la copia de trabajo local, donde `<n>` es el número del acta en la issue (1 la primera, y así sucesivamente aunque el contador de rechazos vuelva a cero). Los ficheros se numeran `01_<afirmación>-<detalle>.<ext>`, como `02_a1-contador.png`; los registros de los gates se llaman `gate-<gate>.txt` | La ruta |
| Documentation | `docs/changes/<tipo>-<N>-<slug>.md`, con el slug de la rama, y su entrada en `docs/changes/README.md`, en la rama del caso | La ruta |

Para ir del PR a su issue: `gh pr view <URL> --json closingIssuesReferences`. Para ir de la issue a su PR: lee la rama en el comentario de control y usa `gh pr list --repo OWNER/REPO --head <rama> --state all`.

## Eventos

El payload puede venir de un webhook, de GitHub Actions, de la API de eventos o directamente de un humano, en lenguaje natural: «quiero hacer la issue #12», «he añadido un comentario en la issue X». Si viene de un humano, el mensaje es el evento: identifica la issue o el PR que menciona. El trigger del proyecto (`adws/`) envía la foto JSON de una issue abierta y admitida cuyo turno es tuyo: sin `adw:waiting-for-human`, sin stage abierto y sin otro dueño. No dice qué ha cambiado: lo averiguas al reconstruir el estado.

Lo que ocurre en una issue afecta a la Issue y lo que ocurre en un PR, a su Change. Un comentario en la conversación de un PR llega como comentario de una issue que tiene `pull_request`: afecta al Change. Un PR en borrador guarda trabajo sin validar y todavía no está en revisión: lo que el humano escriba en él cuenta como escrito en la issue.

## Comentarios

- Todo comentario del agente empieza por una marca: `<!-- adw:control -->` en el comentario de control, `<!-- adw:plan -->` en los planes, `<!-- adw:verification-plan -->` en los planes de verificación, `<!-- adw:evidence -->` en las actas y `<!-- adw -->` en el resto. El agente y el humano suelen usar la misma cuenta, así que el autor no sirve para distinguirlos. Un comentario con marca lo ha producido el agente; uno sin ella, el humano.
- En un PR hay tres tipos de comentario y hay que leerlos todos: conversación y revisiones (`gh pr view <URL> --json comments,reviews`) y comentarios en línea (`gh api repos/OWNER/REPO/pulls/N/comments`).
- Si el humano es el autor del PR, GitHub no le deja pedir cambios ni aprobar: su revisión llega siempre como comentarios. Decide si pide cambios por lo que dice, no por el tipo de revisión.

## Control

- **Admisión:** la issue lleva la etiqueta `adw`. Si una petición directa admite la issue, pon la etiqueta.
- **Propiedad:** la asignación de la issue. El dueño es la cuenta de GitHub con la que trabaja el sistema (`@me`), y sin asignar significa libre. Para reclamarla, `gh issue edit <URL> --add-assignee @me` y relee: si hay otro asignado además de ti, quítate con `--remove-assignee @me` y no actúes. Dos máquinas con la misma cuenta cuentan como el mismo dueño. El humano libera el caso desasignando la issue.
- **Turno:** la issue lleva `adw:waiting-for-human` mientras el turno es del humano. Al cederlo, publica primero un comentario que diga qué esperas del humano y que quite esa etiqueta cuando termine, y después pon la etiqueta. El humano devuelve el turno quitándola. Si lo devuelve con una petición directa, quítala tú.
- **Trabajo sin validar:** si al ceder el turno quedan cambios sin registrar, haz primero commit con `wip: <título de la issue> (sin validar)`: así la copia de trabajo queda limpia aunque falle lo demás. Después, `git push -u origin <rama>`. Si el caso no tiene PR, ábrelo en borrador con `gh pr create --draft --repo OWNER/REPO --base <base> --head <rama> --body-file <fichero>`, con `Closes <URL de la issue>` y una descripción que diga que el trabajo no está validado y enlace la última acta, si la hay. Si ya tiene PR y no está en borrador (`gh pr view <URL> --json isDraft`), pásalo a borrador con `gh pr ready <URL> --undo` y comenta en él lo mismo. Si el plan del repositorio no admite borradores, no abras PR: enlaza la rama. Menciona el PR o la rama en el comentario con el que cedes el turno. Si el push falla, dilo en ese comentario: el commit queda en local.
- **Comentario de control:** un único comentario en la issue que empieza por `<!-- adw:control -->`. Recoge la rama del caso y su base. Lo crea el dueño la primera vez que lo necesita y después lo edita en el sitio, sin publicar otro. Solo lo escribe el dueño, así que no hay escrituras concurrentes.

## Estado

**Proyección.** Se proyecta en la issue con etiquetas. Todas van en kebab-case: minúsculas y palabras separadas por guiones.

- Cada stage abierto es `stage:` más el nombre del stage (`ReviewResponse` → `stage:review-response`). Es también el candado del caso. Si hay varios stages abiertos a la vez, como `VerificationPlanning` y `Planning`, lleva una etiqueta por cada uno.
- Cada milestone alcanzado es el artefacto, dos puntos y el milestone (`changes_requested` del Change → `change:changes-requested`).
- El tipo de la issue es `feature`, `bug` o `chore`. Si la issue ya trae una, respétala.
- Los milestones terminales se observan sin etiqueta: `issue:closed` es la issue cerrada, `change:merged` es el PR integrado y `change:discarded`, el PR cerrado sin integrar.
- `environment:ready` y `environment:released` no se proyectan: describen la máquina, no el caso. Se observan en local: la rama actual, que esté al día con `origin` y si `http://localhost:3000/healthz` responde.

**Catálogo.** Son las únicas etiquetas que usa el adaptador. No inventes otras ni cambies su forma.

| Familia | Etiquetas |
| --- | --- |
| Control | `adw`, `adw:waiting-for-human` |
| Tipo | `feature`, `bug`, `chore` |
| Stages | `stage:clarification`, `stage:setup`, `stage:teardown`, `stage:planning`, `stage:implementation`, `stage:verification-planning`, `stage:validation`, `stage:documenting`, `stage:publication`, `stage:review-response` |
| Issue | `issue:clarified`, `issue:info-requested` |
| Plan | `plan:ready`, `plan:question-raised` |
| Change | `change:implemented`, `change:published`, `change:review-answered`, `change:changes-requested` |
| Verification | `verification:planned`, `verification:question-raised`, `verification:validated`, `verification:rejected` |
| Documentation | `documentation:documented` |
| Fallos | `issue:failed`, `environment:failed`, `plan:failed`, `change:failed`, `verification:failed`, `documentation:failed` |

**Observación.** Las etiquetas son el último estado proyectado. Contrástalas con lo observable: comentarios humanos posteriores al último comentario del agente, comentarios de plan, commits sin subir, si el PR contiene la cabeza de la rama. Para los artefactos de este adaptador:

- **`change:published`:** el PR existe, no está en borrador (`gh pr view <URL> --json isDraft,headRefOid,body`), contiene la cabeza de la rama y su descripción dice «validada en» ese commit, con el enlace al acta vigente.

- **`verification:planned`:** hay un plan de verificación posterior a la última vez que la Issue perdió `clarified` y que recoge todas las peticiones de cambio del humano. Una petición posterior a él lo invalida.
- **Acta vigente:** la última `adw:evidence`. Su veredicto da `validated` o `rejected`, que valen mientras el Change conserve `implemented`. `publish` ata el acta al commit que crea, y la descripción del PR dice «validada en `<sha>`». La evidencia sigue valiendo mientras ese commit sea la cabeza del PR. Los ficheros de `docs/changes/` no la invalidan.
- **Rechazos seguidos:** las actas con veredicto «rechazado» posteriores a la última intervención humana, que es lo más reciente entre el último comentario humano en la issue o el PR y la última vez que se quitó `adw:waiting-for-human` (`gh api repos/OWNER/REPO/issues/N/timeline --paginate`, eventos `unlabeled`). El intento en curso es ese número más uno.
- **Rechazo pendiente:** si la última acta es un rechazo y ninguna revisión del plan posterior la cita, el rechazo todavía no ha llegado al plan. Pasa cuando el humano devuelve el turno tras agotarse el tope.
- **`documentation:documented`:** el documento y su entrada en el índice están en la rama respecto a la base, registrados o en la copia de trabajo (`git diff --name-only origin/<base>` y `git status --porcelain`), y el documento enlaza el acta vigente. Si enlaza una anterior, describe otro diff y no cuenta.

**Escritura.** Cambia etiquetas con `gh issue edit <URL> --add-label … --remove-label …`, sin reemplazar el resto. Si una etiqueta no existe, créala con `gh label create`. Al abrir un stage, pon su etiqueta antes de nada. En lo demás, publica primero y cambia las etiquetas al final: quien vea una etiqueta nueva ya encontrará en GitHub lo que anuncia.

**Reintentos.** Si un stage falla o agota su tope, proyecta `<artifact>:failed`, comenta el fallo y cede el turno. El humano reintenta devolviéndolo: al quitar `adw:waiting-for-human`, se invalida `failed`, quitas su etiqueta y el contador de rechazos vuelve a cero. Si el fallo fue por agotar el tope, proyecta además `verification:rejected`: el último rechazo llega al plan. Si una ejecución se interrumpió y dejó una etiqueta `stage:*`, el humano la quita para reintentar.

## Acta

El acta es la evidencia que ve el humano. Los ficheros de `.adw/evidence/` se quedan en la máquina que validó; lo que el acta no publique no lo puede comprobar nadie más.

```markdown
<!-- adw:evidence -->
### ADW · acta de verificación (intento k/<max_rejections>)
Plan de verificación: <enlace a la revisión vigente> · Cambio validado: `<rama>` @ `<commit corto>` con los cambios sin registrar · árbol `<huella corta>`
Veredicto: validado | rechazado — <regla de decisión aplicada>

#### Funcional
| # | Afirmación | Método | Esperado | Observado | Veredicto | Evidencia | Reproducir |
| --- | --- | --- | --- | --- | --- | --- | --- |

#### Fidelidad
Matriz criterio → afirmaciones y test asociado · cambios sin justificar · criterios sin cubrir

#### Estructural
Automatizado: <gate · código de salida> · Juicios: regla · cita · razonamiento · refutación · confianza · severidad

#### Para el humano
`should-fix`, juicios de confianza baja y verificaciones humanas
```

- Las capturas y los vídeos se citan en el cuerpo con la misma ruta que se pasa a `--attach`, como `![A1 · contador con una tarea completada](.adw/evidence/issue-12/1/02_a1-contador.png)`: `gh` sube el fichero y reescribe la referencia con la URL publicada. Un vídeo citado así queda como su URL en una línea propia, que GitHub muestra como reproductor. Lo que no se cite se añade al final, en el orden de los `--attach`.
- La evidencia textual (registros, transcripciones, consultas) va recortada en bloques de código, dentro de `<details>` si es larga, citando el fichero del que sale.
- Si el rechazo agota el tope, el acta termina con «Tope alcanzado (`max_rejections`: <valor>): el caso pasa al humano.».
- Un comentario de GitHub admite 65 536 caracteres. Si el acta se acerca, recorta los extractos de texto: los registros completos están en `.adw/evidence/`.

## Servicios

El plan, el plan de verificación y el acta viven en la issue, no en la rama. `plans/issue-N.md` es solo una copia de trabajo para las skills: nunca se commitea. `.adw/` está en `.gitignore`: lo que se escribe ahí no ensucia el checkout.

| Servicio | En GitHub |
| --- | --- |
| `clarify` | Publica las preguntas como comentario en la issue. Si la clasificas, pon la etiqueta del tipo. |
| `setup_environment` | Si el comentario de control no tiene rama, asígnala como `<tipo>/issue-N-<slug>` (el slug sale del título, en minúsculas y con guiones), con la rama por defecto del repositorio como base (`gh repo view OWNER/REPO --json defaultBranchRef`), y regístrala. Si la rama no existe en `origin`, `git fetch origin <base>` y `git switch -c <rama> origin/<base>`. Si existe, `git fetch origin <rama>`, `git switch <rama>` y `git merge --ff-only origin/<rama>`. Después ejecuta `env:start`. |
| `teardown_environment` | Ejecuta `env:stop`. |
| `write_plan` | El plan se escribe en `plans/issue-N.md`. Publícalo en la issue como comentario que empiece por `<!-- adw:plan -->`, indicando la revisión, las peticiones que recoge y su origen, con el enlace al acta si es un rechazo de validación, y borra el fichero. |
| `implement` | Antes de lanzarlo, escribe el plan vigente en `plans/issue-N.md`. Bórralo cuando termine: no forma parte del Change. |
| `plan_verification` | Lee el código de la base, no de la copia de trabajo, que puede estar en la rama de otro caso: `git show origin/<base>:<ruta>`. El plan de verificación se escribe en `.adw/evidence/issue-N/verification-plan.md`. Publícalo en la issue como comentario que empiece por `<!-- adw:verification-plan -->` y por `### ADW · plan de verificación (revisión r)`, indicando lo que recoge: la especificación de la issue y las peticiones del humano, con su enlace. |
| `validate` | Usa `.adw/evidence/issue-N/<n>/`. Si ese directorio existe de un intento que no llegó a publicar su acta, vacíalo antes. Antes de empezar, escribe ahí el plan de verificación vigente (`verification-plan.md`) y el plan de implementación vigente (`plan.md`), que solo lee la fidelidad. El cambio validado es la rama (`git branch --show-current`), el commit de partida (`git rev-parse --short HEAD`) y la huella completa del árbol de la copia de trabajo, cambios sin registrar incluidos (`bash .claude/skills/adw/scripts/tree.sh`). Antes de verificar, comprueba que la aplicación sirve la copia de trabajo: en modo host, `buffalo dev` recompila solo al cambiar el código: espera a que su salida muestre una compilación posterior al último cambio, o reinícialo; en modo Docker completo (`docker compose ps` muestra `app`), reconstrúyela con `docker compose up --build -d app` y espera a que `/healthz` responda. El navegador usa la sesión `-s=issue-N`, que se cierra al terminar. El acta se escribe en `acta.md` y se publica con `gh issue comment <URL> --body-file .adw/evidence/issue-N/<n>/acta.md --attach <fichero> …`, con un `--attach` por cada captura o vídeo. |
| `document` | Escribe `docs/changes/<tipo>-<N>-<slug>.md` y su entrada en `docs/changes/README.md` en la copia de trabajo. No hace commit. |
| `publish` | Hace commit con un mensaje convencional según el tipo (`feat:`, `fix:` o `chore:`), documento incluido. Comprueba con `git diff --name-only <huella del acta> HEAD` que solo difieren `docs/changes/<tipo>-<N>-<slug>.md` y `docs/changes/README.md`. Después, `git push -u origin <rama>`. Si no hay PR, lo abre con `gh pr create --repo OWNER/REPO --base <base> --head <rama> --body-file <fichero>`, con `Closes <URL de la issue>` y el enlace al plan vigente en la descripción, y comenta el enlace en la issue. Si ya existe, el push lo actualiza, `gh pr edit <URL> --body-file <fichero>` renueva la descripción y, si estaba en borrador, `gh pr ready <URL>` lo deja listo para revisión. La descripción lleva `## Evidencia` (resumen del acta, enlaces al acta y al plan de verificación y «validada en `<sha>`», con el commit que acaba de crear) y `## Para el humano`. Las capturas y vídeos reutilizan las URLs que el acta ya publicó (`gh api repos/OWNER/REPO/issues/comments/<id> --jq .body`): no se vuelven a subir. |
| `answer_review` | Responde con comentarios en el PR. |
| `report_failure` | Comenta en la issue qué stage ha fallado y por qué, y que el humano quite `adw:waiting-for-human` para reintentar. Si se agotó el tope, lo dice, indica el valor de `max_rejections` y resume los rechazos de la racha con el enlace a cada acta. |

El workflow compuesto se publica como comentario en la issue o en el PR, según dónde haya ocurrido el evento.

## Límites de `--attach`

`gh` 2.99 o posterior sube imágenes y vídeos con `--attach` en `gh issue comment`, `gh pr create` y `gh pr edit`, hasta 50 ficheros por comando. Acepta PNG, JPEG, GIF, WebP, SVG, MP4, MOV y WebM. Las imágenes admiten hasta 10 MB; los vídeos, hasta 10 MB en el plan gratuito y 100 MB en los de pago, así que los vídeos deben ser cortos: comprueba su tamaño antes de adjuntarlos. Solo funciona en github.com y GitHub Enterprise Cloud, y requiere permiso de push en el repositorio. Si una subida falla, `gh` termina con error: no es un rechazo del cambio, pero el acta queda incompleta y `validate` falla.
