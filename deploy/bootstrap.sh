#!/usr/bin/env bash
# One-time server bootstrap for selenoid.autotests.cloud.
# Run with sudo on a fresh Ubuntu host with Docker installed.
set -euo pipefail

DEPLOY_USER="${DEPLOY_USER:-${SUDO_USER:-}}"
CONFIG_DIR="${SELENOID_CONFIG_DIR:-/opt/selenoid}"
CM_BIN="/home/${DEPLOY_USER}/cm"

if [[ -z "$DEPLOY_USER" || "$DEPLOY_USER" == "root" ]]; then
  echo "Set DEPLOY_USER to a non-root account (e.g. DEPLOY_USER=selenoid sudo -E ./bootstrap.sh)" >&2
  exit 1
fi

if [[ "$(id -u)" -ne 0 ]]; then
  echo "Run as root: sudo DEPLOY_USER=$DEPLOY_USER ./bootstrap.sh" >&2
  exit 1
fi

echo "=== docker group for $DEPLOY_USER ==="
usermod -aG docker "$DEPLOY_USER"

echo "=== config dir $CONFIG_DIR ==="
mkdir -p "$CONFIG_DIR"/{video,logs,bin}
chown -R "$DEPLOY_USER:docker" "$CONFIG_DIR"
chmod 775 "$CONFIG_DIR"

echo "=== cm binary for $DEPLOY_USER ==="
sudo -u "$DEPLOY_USER" bash -c "
  curl -fsSL https://github.com/qa-guru/cm/releases/latest/download/cm_linux_amd64 -o '$CM_BIN'
  chmod +x '$CM_BIN'
"

echo "=== docker network selenoid (if missing) ==="
docker network inspect selenoid >/dev/null 2>&1 || docker network create selenoid

echo "Bootstrap complete."
echo "Next (as $DEPLOY_USER, new login shell for docker group):"
echo "  ./deploy/deploy.sh"
