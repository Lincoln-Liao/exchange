#!/bin/bash
export IMAGE_VERSION
IMAGE_VERSION=latest
REPO_ROOT=$($(which git) rev-parse --show-toplevel)


PWD="${REPO_ROOT}" docker compose \
    -f "${REPO_ROOT}"/docker-compose.yaml \
    up --remove-orphans --build "$@"