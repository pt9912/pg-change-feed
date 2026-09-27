#!/usr/bin/env bash
# image-mutation — baut oder entfernt ein Mutations-Image mit eigenem
# Repository-Namen und eigenem Tag, getrennt vom Lauf-Beleg-Pfad
# (`make image`, `:dev`, `harness/image-hash.txt`). Die Kopie liefert der
# Aufrufer (`git archive <Stand>` in ein Scratchpad-Verzeichnis, Mutation
# darauf mit Edit/Write); dieses Skript baut nur, es zieht keine Kopie.
# `build` ruft genau `docker buildx build --load -t
# pg-change-feed-mutation:<TAG> <SRC>` auf; `rm` ruft genau `docker rmi
# pg-change-feed-mutation:<TAG>`. Jede Eingabeprüfung läuft vor dem
# jeweiligen Docker-Aufruf. Vertrag, Grenze und Exit-Codes:
# harness/targets/image-mutation.md (AGENTS.md §3.1).
#
# Aufruf: image-mutation.sh build <SRC> <TAG>
#         image-mutation.sh rm <TAG>
# Exit: 0 gebaut/entfernt · 1 der Docker-Aufruf endet ≠ 0 · 2 Eingabefehler
set -uo pipefail

IMAGE_REPO="pg-change-feed-mutation"
TAG_RE='^[a-z0-9][a-z0-9_.-]{0,62}$'

usage() {
  echo "Aufruf: image-mutation.sh build <SRC> <TAG> | image-mutation.sh rm <TAG>" >&2
}

# check_tag <TAG>: Positiv-Zeichenklasse, dazu die reservierten Namen der
# Lauf-Beleg-Tags. Setzt bei Verstoß eine Meldung auf stderr und gibt 1
# zurück; der Aufrufer bricht dann mit Exit 2 ab, vor jedem Docker-Aufruf.
check_tag() {
  local tag=$1
  if [[ ! $tag =~ $TAG_RE ]]; then
    echo "image-mutation: TAG '$tag' verletzt die Zeichenklasse [a-z0-9][a-z0-9_.-]{0,62}" >&2
    return 1
  fi
  if [ "$tag" = "dev" ] || [ "$tag" = "latest" ]; then
    echo "image-mutation: TAG '$tag' ist reserviert (dev, latest sind die Lauf-Beleg-Tags)" >&2
    return 1
  fi
  return 0
}

if [ $# -lt 1 ]; then
  usage
  exit 2
fi
verb=$1
shift

case "$verb" in
  build)
    if [ $# -ne 2 ]; then
      usage
      exit 2
    fi
    src=$1
    tag=$2
    if [ -z "$src" ] || [ -z "$tag" ]; then
      usage
      exit 2
    fi
    if ! check_tag "$tag"; then
      exit 2
    fi
    if [ ! -d "$src" ]; then
      echo "image-mutation: SRC '$src' ist kein Verzeichnis" >&2
      exit 2
    fi
    src_abs=$(realpath -- "$src") || exit 2
    if [ ! -f "$src_abs/Dockerfile" ]; then
      echo "image-mutation: SRC '$src' trägt kein Dockerfile" >&2
      exit 2
    fi
    if [ ! -f "$src_abs/go.mod" ]; then
      echo "image-mutation: SRC '$src' trägt kein go.mod" >&2
      exit 2
    fi
    root=$(git rev-parse --show-toplevel) || exit 2
    root_abs=$(realpath -- "$root") || exit 2
    if [ "$src_abs" = "$root_abs" ] || [[ $src_abs == "$root_abs"/* ]]; then
      echo "image-mutation: SRC '$src' ist die Repo-Wurzel oder liegt unter ihr (eine Mutation baut nie aus dem Arbeitsbaum, AGENTS.md §3.1)" >&2
      exit 2
    fi
    docker buildx build --load -t "${IMAGE_REPO}:${tag}" "$src_abs"
    rc=$?
    if [ "$rc" -ne 0 ]; then
      echo "image-mutation: docker buildx build endet mit Exit $rc" >&2
      exit 1
    fi
    echo "image-mutation: ${IMAGE_REPO}:${tag} gebaut aus $src_abs"
    ;;
  rm)
    if [ $# -ne 1 ]; then
      usage
      exit 2
    fi
    tag=$1
    if [ -z "$tag" ]; then
      usage
      exit 2
    fi
    if ! check_tag "$tag"; then
      exit 2
    fi
    docker rmi "${IMAGE_REPO}:${tag}"
    rc=$?
    if [ "$rc" -ne 0 ]; then
      echo "image-mutation: docker rmi endet mit Exit $rc" >&2
      exit 1
    fi
    echo "image-mutation: ${IMAGE_REPO}:${tag} entfernt"
    ;;
  *)
    usage
    exit 2
    ;;
esac
