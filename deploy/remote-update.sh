#!/usr/bin/env bash
set -euo pipefail
LOG=/tmp/selenoid-deploy.log
exec > >(tee -a "$LOG") 2>&1

echo "DEPLOY_START $(date -Is)"

CM=/root/cm
CONF=/root/.aerokube/selenoid
BROWSERS_JSON="${BROWSERS_JSON:-/tmp/browsers-production.json}"

curl -fsSL https://github.com/qa-guru/cm/releases/latest/download/cm_linux_amd64 -o "$CM"
chmod +x "$CM"

"$CM" selenoid stop -c "$CONF" || true
"$CM" selenoid-ui stop -c "$CONF" || true

"$CM" selenoid update -c "$CONF" -j "$BROWSERS_JSON"
"$CM" selenoid-ui update -c "$CONF"

echo "=== status ==="
curl -sf http://127.0.0.1:4444/status || true
echo
docker ps --filter name=selenoid --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}'

echo "DEPLOY_DONE $(date -Is)"
