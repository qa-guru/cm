#!/usr/bin/env bash
# One-time migration: /root/.aerokube/selenoid → /opt/selenoid (owner: selenoid)
# Run on the server as root: sudo ./deploy/migrate-to-opt.sh
set -euo pipefail

DEPLOY_USER="${DEPLOY_USER:-selenoid}"
LEGACY_DIR="/root/.aerokube/selenoid"
CONFIG_DIR="${SELENOID_CONFIG_DIR:-/opt/selenoid}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [[ "$(id -u)" -ne 0 ]]; then
  echo "Run as root: sudo ./deploy/migrate-to-opt.sh" >&2
  exit 1
fi

echo "=== stop legacy root-managed stack ==="
if [[ -x /root/cm ]]; then
  /root/cm selenoid stop -c "$LEGACY_DIR" 2>/dev/null || true
  /root/cm selenoid-ui stop -c "$LEGACY_DIR" 2>/dev/null || true
fi
docker stop selenoid selenoid-ui 2>/dev/null || true
docker rm selenoid selenoid-ui 2>/dev/null || true

echo "=== bootstrap $CONFIG_DIR for $DEPLOY_USER ==="
DEPLOY_USER="$DEPLOY_USER" SELENOID_CONFIG_DIR="$CONFIG_DIR" "$SCRIPT_DIR/bootstrap.sh"

if [[ -d "$LEGACY_DIR" ]]; then
  echo "=== rsync $LEGACY_DIR → $CONFIG_DIR (config + bin only; video if missing) ==="
  for item in browsers.json browsers.json.bak; do
    [[ -e "$LEGACY_DIR/$item" ]] && rsync -a "$LEGACY_DIR/$item" "$CONFIG_DIR/$item"
  done
  for item in bin/selenoid bin/selenoid-ui; do
    [[ -e "$LEGACY_DIR/$item" ]] && rsync -a "$LEGACY_DIR/$item" "$CONFIG_DIR/$item"
  done
  if [[ -d "$LEGACY_DIR/video" && -z "$(ls -A "$CONFIG_DIR/video" 2>/dev/null)" ]]; then
    rsync -a "$LEGACY_DIR/video/" "$CONFIG_DIR/video/"
  fi
  if [[ -d "$LEGACY_DIR/logs" ]]; then
    rsync -a "$LEGACY_DIR/logs/" "$CONFIG_DIR/logs/"
  fi
fi

echo "=== permissions on $CONFIG_DIR ==="
chown -R "$DEPLOY_USER:docker" "$CONFIG_DIR"
chmod 775 "$CONFIG_DIR" "$CONFIG_DIR"/video "$CONFIG_DIR"/logs "$CONFIG_DIR"/bin
[[ -f "$CONFIG_DIR/browsers.json" ]] && chmod 664 "$CONFIG_DIR/browsers.json"
find "$CONFIG_DIR/bin" -type f -exec chmod 755 {} \;

echo "=== deploy as $DEPLOY_USER (docker group via sg) ==="
sudo -u "$DEPLOY_USER" sg docker -c "
  export SELENOID_CONFIG_DIR='$CONFIG_DIR'
  export CM_BIN='/home/$DEPLOY_USER/cm'
  bash '$SCRIPT_DIR/deploy.sh'
"

echo "=== verify ==="
curl -sf http://127.0.0.1:4444/status | head -c 200 || true
echo
curl -sf -o /dev/null -w 'UI HTTP %{http_code}\n' http://127.0.0.1:8080/ || true
docker ps --filter name=selenoid --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}'
ls -la "$CONFIG_DIR/bin" 2>/dev/null || true

echo
echo "Migration complete."
echo "  config: $CONFIG_DIR"
echo "  owner:  $DEPLOY_USER:docker"
echo "  cm:     /home/$DEPLOY_USER/cm"
echo "Legacy data left at $LEGACY_DIR (remove manually after verification)."
