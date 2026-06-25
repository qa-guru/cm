#!/usr/bin/env bash
# Deploy qa-guru Selenoid stack via cm (hub + UI + browser images).
# Run on the server as selenoid — not via sudo.
set -euo pipefail

CONFIG_DIR="${SELENOID_CONFIG_DIR:-/opt/selenoid}"
CM_BIN="${CM_BIN:-$HOME/cm}"
CM_URL="${CM_URL:-https://github.com/qa-guru/cm/releases/latest/download/cm_linux_amd64}"
VERSION="${SELENOID_VERSION:-v2.0.3}"
UI_VERSION="${SELENOID_UI_VERSION:-v2.0.1}"
GITHUB_OWNER="${GITHUB_OWNER:-qa-guru}"
UI_USER="${SELENOID_UI_USER:-user1}"
UI_PASSWORD="${SELENOID_UI_PASSWORD:-1234}"
UI_AUTH=(-u "${UI_USER}:${UI_PASSWORD}")

ensure_htpasswd() {
  local dest="$1"
  if [[ -f "$dest" ]]; then
    return 0
  fi
  if [[ -r /etc/nginx/selenoid.htpasswd ]]; then
    cp /etc/nginx/selenoid.htpasswd "$dest"
    chmod 600 "$dest"
    return 0
  fi
  if command -v htpasswd >/dev/null 2>&1; then
    htpasswd -cb "$dest" "$UI_USER" "$UI_PASSWORD"
    chmod 600 "$dest"
    return 0
  fi
  if command -v openssl >/dev/null 2>&1; then
    printf '%s:%s\n' "$UI_USER" "$(openssl passwd -apr1 "$UI_PASSWORD")" >"$dest"
    chmod 600 "$dest"
    return 0
  fi
  echo "FAIL: cannot create $dest (install htpasswd or openssl)" >&2
  exit 1
}

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
  local repo="$1" dest="$2" tag="${3:-${VERSION:-latest}}"
  local url="https://github.com/${GITHUB_OWNER}/${repo}/releases/download/${tag}/${repo}_linux_amd64"
  local tmp="${dest}.new.$$"
  echo "Downloading ${repo} ${tag} → ${dest}"
  curl -fsSL "$url" -o "$tmp"
  chmod 755 "$tmp"
  mv "$tmp" "$dest"
}

echo "=== stop legacy containers ==="
docker stop selenoid selenoid-ui 2>/dev/null || true
docker rm selenoid selenoid-ui 2>/dev/null || true

echo "=== stop cm-managed services ==="
"$CM_BIN" selenoid stop -c "$CONFIG_DIR" 2>/dev/null || true
"$CM_BIN" selenoid-ui stop -c "$CONFIG_DIR" 2>/dev/null || true

if pgrep -f "${CONFIG_DIR}/bin/selenoid" >/dev/null 2>&1; then
  pkill -f "${CONFIG_DIR}/bin/selenoid" || true
  sleep 1
fi

echo "=== download hub binaries (selenoid ${VERSION:-latest}, selenoid-ui ${UI_VERSION:-latest}) ==="
download_binary selenoid "$CONFIG_DIR/bin/selenoid" "$VERSION"
download_binary selenoid-ui "$CONFIG_DIR/bin/selenoid-ui" "$UI_VERSION"

echo "=== configure hub (pull browser images, write browsers.json) ==="
"$CM_BIN" selenoid configure -c "$CONFIG_DIR" "${version_args[@]}"

echo "=== force refresh Playwright browser image ==="
docker pull qaguru/playwright:v1.61.1-noble

mkdir -p "$CONFIG_DIR/video" "$CONFIG_DIR/logs"

echo "=== start hub (native binary on host — hub-in-docker breaks browser port bindings) ==="
export DOCKER_API_VERSION="${DOCKER_API_VERSION:-1.45}"
nohup "${CONFIG_DIR}/bin/selenoid" \
  -conf "${CONFIG_DIR}/browsers.json" \
  -limit 5 \
  -container-network selenoid \
  -video-output-dir "${CONFIG_DIR}/video/" \
  -video-recorder-image selenoid/video-recorder:latest-release \
  -log-output-dir "${CONFIG_DIR}/logs/" \
  -listen :4444 \
  >> "${CONFIG_DIR}/logs/selenoid.log" 2>&1 &

for attempt in 1 2 3 4 5 6 7 8 9 10; do
  if curl -sf "http://127.0.0.1:4444/status" >/dev/null 2>&1; then
    break
  fi
  echo "hub /status not ready (attempt ${attempt}/10)..." >&2
  sleep 2
done

echo "=== start UI (host network -> 127.0.0.1:4444) ==="
"$CM_BIN" selenoid-ui download -c "$CONFIG_DIR" "${version_args[@]}" 2>/dev/null || true
docker stop selenoid-ui 2>/dev/null || true
docker rm selenoid-ui 2>/dev/null || true
"$CM_BIN" selenoid-ui stop -c "$CONFIG_DIR" 2>/dev/null || true

HTPASSWD="${CONFIG_DIR}/htpasswd"
ensure_htpasswd "$HTPASSWD"

UI_IMAGE="aerokube/selenoid-ui:latest-release"
docker pull "$UI_IMAGE" >/dev/null 2>&1 || true
docker run -d --name selenoid-ui \
  --restart unless-stopped \
  --network host \
  -v "${CONFIG_DIR}:/etc/selenoid:ro" \
  -v "${CONFIG_DIR}/bin/selenoid-ui:/selenoid-ui:ro" \
  --entrypoint /selenoid-ui \
  "$UI_IMAGE" \
    -selenoid-uri=http://127.0.0.1:4444 \
    -browsers-conf=/etc/selenoid/browsers.json \
    -users=/etc/selenoid/htpasswd \
    -listen=:8080

echo "=== local hub status ==="
curl -sf "http://127.0.0.1:4444/status" | (command -v jq >/dev/null && jq . || cat)
echo

echo "=== UI backend status ==="
ui_json=""
ui_http="000"
for attempt in 1 2 3 4 5 6; do
  ui_http="$(curl -sS -o /tmp/ui-status.json -w '%{http_code}' "${UI_AUTH[@]}" "http://127.0.0.1:8080/status" 2>/dev/null || echo "000")"
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

ui_no_auth="$(curl -s -o /dev/null -w '%{http_code}' "http://127.0.0.1:8080/" 2>/dev/null || echo "000")"
if [[ "$ui_no_auth" == "401" ]]; then
  echo "OK  UI basic auth enabled (HTTP 401 without credentials)"
else
  echo "FAIL: selenoid-ui without auth returned HTTP ${ui_no_auth} (expected 401)" >&2
  docker logs --tail 40 selenoid-ui 2>&1 || true
  docker inspect selenoid-ui --format '{{json .Config.Cmd}}' 2>&1 || true
  exit 1
fi
echo

echo "=== smoke: create chrome session ==="
session_json="$(curl -sS -m 120 -X POST "http://127.0.0.1:4444/wd/hub/session" \
  -H 'Content-Type: application/json' \
  -d '{"capabilities":{"alwaysMatch":{"browserName":"chrome","browserVersion":"148.0","selenoid:options":{"sessionTimeout":"30s","name":"deploy-smoke"}}}}' || true)"
if command -v jq >/dev/null; then
  session_id="$(jq -r '.value.sessionId // .sessionId // empty' <<<"$session_json")"
  if [[ -z "$session_id" ]]; then
    echo "FAIL: could not create chrome session: $session_json" >&2
    tail -30 "${CONFIG_DIR}/logs/selenoid.log" 2>&1 || true
    exit 1
  fi
  echo "OK  session ${session_id}"
  curl -sS -X DELETE "http://127.0.0.1:4444/wd/hub/session/${session_id}" >/dev/null || true
else
  echo "$session_json"
fi
echo
docker ps --filter name=selenoid --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}'
pgrep -af "${CONFIG_DIR}/bin/selenoid" || true

echo "=== nginx (basic auth on UI) ==="
NGINX_CONF="${NGINX_CONF_SRC:-/tmp/nginx-selenoid.conf}"
NGINX_SYNC="${NGINX_SYNC_SCRIPT:-/tmp/sync-nginx.sh}"
if [[ ! -f "$NGINX_CONF" || ! -f "$NGINX_SYNC" ]]; then
  echo "WARN: nginx config not found ($NGINX_CONF / $NGINX_SYNC) — skip"
elif sudo -n true 2>/dev/null; then
  if NGINX_CONF_SRC="$NGINX_CONF" sudo "$NGINX_SYNC"; then
    echo "OK  nginx config applied"
  else
    echo "WARN: nginx sync failed — run: sudo NGINX_CONF_SRC=$NGINX_CONF $NGINX_SYNC" >&2
  fi
else
  echo "WARN: passwordless sudo unavailable — run on server as root:" >&2
  echo "  sudo ./deploy/bootstrap.sh   # once, installs NOPASSWD for sync-nginx.sh" >&2
  echo "  sudo NGINX_CONF_SRC=$NGINX_CONF $NGINX_SYNC" >&2
fi

exit 0
