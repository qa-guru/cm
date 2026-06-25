#!/usr/bin/env bash
# Apply nginx-selenoid.conf on the server (requires sudo).
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONF_SRC="${SCRIPT_DIR}/nginx-selenoid.conf"
SITE_NAME="${NGINX_SITE_NAME:-selenoid}"
SITE_PATH="/etc/nginx/sites-available/${SITE_NAME}"

if [[ ! -f "$CONF_SRC" ]]; then
  echo "Missing $CONF_SRC" >&2
  exit 1
fi

if [[ "$(id -u)" -ne 0 ]]; then
  if sudo -n true 2>/dev/null; then
    exec sudo "$0" "$@"
  fi
  echo "Run as root or with passwordless sudo: sudo $0" >&2
  exit 1
fi

cp "$CONF_SRC" "$SITE_PATH"
ln -sf "$SITE_PATH" "/etc/nginx/sites-enabled/${SITE_NAME}"
nginx -t
systemctl reload nginx
echo "OK: nginx reloaded ($SITE_PATH)"
