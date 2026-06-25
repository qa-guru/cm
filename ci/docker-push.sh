#!/usr/bin/env bash

set -e

if [ -z "${DOCKER_USERNAME:-}" ] || [ -z "${DOCKER_PASSWORD:-}" ]; then
	echo "Skipping Docker push for cm: DOCKER_USERNAME/DOCKER_PASSWORD not set"
	exit 0
fi

docker login -u="$DOCKER_USERNAME" -p="$DOCKER_PASSWORD"
docker build --pull -t "qaguru/cm:${1}" .
docker push "qaguru/cm:${1}"
