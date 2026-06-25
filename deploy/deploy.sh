#!/usr/bin/env bash
# Deploy qa-guru Selenoid stack via cm (hub + UI + browser images).
# Run on the server as selenoid — not via sudo.
set -euo pipefail

CONFIG_DIR="${SELENOID_CONFIG_DIR:-/opt/selenoid}"
CM_BIN="${CM_BIN:-$HOME/cm}"
CM_URL="${CM_URL:-https://github.com/qa-guru/cm/releases/latest/download/cm_linux_amd64}"
VERSION="${SELENOID_VERSION:-v2.0.1}"
GITHUB_OWNER="${GITHUB_OWNER:-qa-guru}"

version_args=()
if [[ -n "$VERSION" ]]; then
  version_args=(-v "$VERSION")
fi

if ! groups | grep -qw docker; then
  echo "Current user is not in the docker group. Run deploy/bootstrap.sh first." >&2
  exit 1
fi

mkdir -p "$CONFIG_DIR/bin"

if [[ ! -x "$CM_BIN" ]]; then
  echo "Downloading cm to $CM_BIN"
  curl -fsSL "$CM_URL" -o "$CM_BIN"
  chmod +x "$CM_BIN"
fi

download_binary() {
  local repo="$1" dest="$2"
  local tag="${VERSION:-latest}"
  local url="https://github.com/${GITHUB_OWNER}/${repo}/releases/download/${tag}/${repo}_linux_amd64"
  echo "Downloading ${repo} ${tag} → ${dest}"
  curl -fsSL "$url" -o "$dest"
  chmod 755 "$dest"
}

if [[ ! -x "$CONFIG_DIR/bin/selenoid" ]]; then
  download_binary selenoid "$CONFIG_DIR/bin/selenoid"
fi
if [[ ! -x "$CONFIG_DIR/bin/selenoid-ui" ]]; then
  download_binary selenoid-ui "$CONFIG_DIR/bin/selenoid-ui"
fi

echo "=== stop legacy containers ==="
docker stop selenoid selenoid-ui 2>/dev/null || true
docker rm selenoid selenoid-ui 2>/dev/null || true

echo "=== stop cm-managed services ==="
"$CM_BIN" selenoid stop -c "$CONFIG_DIR" 2>/dev/null || true
"$CM_BIN" selenoid-ui stop -c "$CONFIG_DIR" 2>/dev/null || true

echo "=== start hub (config-dir: $CONFIG_DIR) ==="
"$CM_BIN" selenoid start -c "$CONFIG_DIR" "${version_args[@]}"

echo "=== start UI (host network -> 127.0.0.1:4444) ==="
"$CM_BIN" selenoid-ui download -c "$CONFIG_DIR" "${version_args[@]}" 2>/dev/null || true
docker stop selenoid-ui 2>/dev/null || true
docker rm selenoid-ui 2>/dev/null || true
"$CM_BIN" selenoid-ui stop -c "$CONFIG_DIR" 2>/dev/null || true

UI_IMAGE="aerokube/selenoid-ui:latest-release"
docker pull "$UI_IMAGE" >/dev/null 2>&1 || true
docker run -d --name selenoid-ui \
  --restart unless-stopped \
  --network host \
  -v "${CONFIG_DIR}:/etc/selenoid:ro" \
  -v "${CONFIG_DIR}/bin/selenoid-ui:/selenoid-ui:ro" \
  "$UI_IMAGE" \
  /selenoid-ui \
    -selenoid-uri=http://127.0.0.1:4444 \
    -browsers-conf=/etc/selenoid/browsers.json \
    -listen=:8080

echo "=== local hub status ==="
curl -sf "http://127.0.0.1:4444/status" | (command -v jq >/dev/null && jq . || cat)
echo

echo "=== UI backend status ==="
ui_json=""
ui_http="000"
for attempt in 1 2 3 4 5 6; do
  ui_http="$(curl -sS -o /tmp/ui-status.json -w '%{http_code}' "http://127.0.0.1:8080/status" 2>/dev/null || echo "000")"
  if [[ "$ui_http" == "200" ]]; then
    ui_json="$(cat /tmp/ui-status.json)"
    break
  fi
  echo "UI /status HTTP ${ui_http} (attempt ${attempt}/6), waiting..." >&2
  sleep 2
done
rm -f /tmp/ui-status.json
if [[ "$ui_http" != "200" || -z "$ui_json" ]]; then
  echo "FAIL: selenoid-ui /status HTTP ${ui_http}" >&2
  docker logs --tail 40 selenoid-ui 2>&1 || true
  docker inspect selenoid-ui --format '{{json .Config.Cmd}}' 2>&1 || true
  exit 1
fi
if command -v jq >/dev/null; then
  echo "$ui_json" | jq .
  if jq -e '.errors | length > 0' <<<"$ui_json" >/dev/null 2>&1; then
    echo "FAIL: selenoid-ui cannot reach hub (see errors above)" >&2
    docker logs --tail 40 selenoid-ui 2>&1 || true
    docker inspect selenoid-ui --format '{{json .Config.Cmd}}' 2>&1 || true
    exit 1
  fi
  if ! jq -e '.state.total != null' <<<"$ui_json" >/dev/null 2>&1; then
    echo "FAIL: selenoid-ui /status missing .state — check --selenoid-uri" >&2
    docker logs --tail 40 selenoid-ui 2>&1 || true
    docker inspect selenoid-ui --format '{{json .Config.Cmd}}' 2>&1 || true
    exit 1
  fi
else
  echo "$ui_json"
fi
echo
docker ps --filter name=selenoid --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}'
