#!/usr/bin/env bash
# Deploy qa-guru Selenoid stack via cm (hub + UI + browser images).
# Run on the server as a user in the docker group — not via sudo.
set -euo pipefail

CONFIG_DIR="${SELENOID_CONFIG_DIR:-/opt/selenoid}"
CM_BIN="${CM_BIN:-$HOME/cm}"
CM_URL="${CM_URL:-https://github.com/qa-guru/cm/releases/latest/download/cm_linux_amd64}"
VERSION="${SELENOID_VERSION:-}"
FORCE="${SELENOID_FORCE:-1}"

version_args=()
if [[ -n "$VERSION" ]]; then
  version_args=(-v "$VERSION")
fi

force_args=()
if [[ "$FORCE" == "1" ]]; then
  force_args=(-f)
fi

if ! groups | grep -qw docker; then
  echo "Current user is not in the docker group. Run deploy/bootstrap.sh first." >&2
  exit 1
fi

mkdir -p "$CONFIG_DIR"

if [[ ! -x "$CM_BIN" ]]; then
  echo "Downloading cm to $CM_BIN"
  curl -fsSL "$CM_URL" -o "$CM_BIN"
  chmod +x "$CM_BIN"
fi

echo "=== stop legacy containers (root / manual installs) ==="
docker stop selenoid selenoid-ui 2>/dev/null || true
docker rm selenoid selenoid-ui 2>/dev/null || true

echo "=== stop cm-managed services ==="
"$CM_BIN" selenoid stop -c "$CONFIG_DIR" 2>/dev/null || true
"$CM_BIN" selenoid-ui stop -c "$CONFIG_DIR" 2>/dev/null || true

echo "=== update hub (config-dir: $CONFIG_DIR) ==="
"$CM_BIN" selenoid update -c "$CONFIG_DIR" "${force_args[@]}" "${version_args[@]}"

echo "=== update UI ==="
"$CM_BIN" selenoid-ui update -c "$CONFIG_DIR" "${force_args[@]}" "${version_args[@]}"

echo "=== local status ==="
curl -sf "http://127.0.0.1:4444/status" | (command -v jq >/dev/null && jq . || cat)
echo
docker ps --filter name=selenoid --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}'
