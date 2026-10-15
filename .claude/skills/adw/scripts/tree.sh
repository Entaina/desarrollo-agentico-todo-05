#!/usr/bin/env bash
# Huella del árbol de la copia de trabajo, cambios sin registrar incluidos.
#
#   tree.sh    Imprime el hash completo del árbol que tendría un commit de todo lo que hay ahora.
#
# Usa un índice aparte (.adw/index), que parte de HEAD en cada llamada, así que no toca el
# índice real ni la copia de trabajo, y cuenta igual que un commit los ficheros registrados
# aunque estén ignorados. validate la anota en el acta y publish la compara con lo que registra:
#   git diff --name-only <huella> HEAD
set -euo pipefail

root=$(git rev-parse --show-toplevel)
mkdir -p "$root/.adw"
export GIT_INDEX_FILE="$root/.adw/index"
rm -f "$GIT_INDEX_FILE"
git -C "$root" read-tree HEAD
git -C "$root" add -A
git -C "$root" write-tree
