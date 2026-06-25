#!/usr/bin/env bash

set -e

if [ -z "${DOCKER_USERNAME:-}" ] || [ -z "${DOCKER_PASSWORD:-}" ]; then
	echo "ERROR: Docker push for qaguru/cm requires DOCKER_USERNAME and DOCKER_PASSWORD repository secrets" >&2
	exit 1
fi

docker login -u="$DOCKER_USERNAME" -p="$DOCKER_PASSWORD"
docker build --pull -t "qaguru/cm:${1}" .
docker push "qaguru/cm:${1}"
