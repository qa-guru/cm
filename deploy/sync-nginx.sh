#!/usr/bin/env bash
# Apply nginx-selenoid.conf on the server (requires sudo).
# Preserves ssl_certificate* lines from the existing site config when present.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONF_SRC="${SCRIPT_DIR}/nginx-selenoid.conf"
SITE_NAME="${NGINX_SITE_NAME:-selenoid}"
SITE_PATH="/etc/nginx/sites-available/${SITE_NAME}"
TMP="/tmp/nginx-selenoid.generated"

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

cp "$CONF_SRC" "$TMP"

if [[ -f "$SITE_PATH" ]]; then
  mapfile -t ssl_lines < <(grep -E '^\s*ssl_' "$SITE_PATH" || true)
  if ((${#ssl_lines[@]} > 0)); then
    python3 - "$TMP" "${ssl_lines[@]}" <<'PY'
import sys
path = sys.argv[1]
ssl_lines = [line.rstrip() for line in sys.argv[2:]]
text = open(path, encoding="utf-8").read()
marker = "    # ssl_certificate ...;"
replacement = "\n".join(ssl_lines)
if marker in text:
    text = text.replace(marker, replacement, 1)
    text = text.replace(marker, replacement, 1)  # two server blocks
open(path, "w", encoding="utf-8").write(text)
PY
  fi
fi

cp "$TMP" "$SITE_PATH"
ln -sf "$SITE_PATH" "/etc/nginx/sites-enabled/${SITE_NAME}"
nginx -t
systemctl reload nginx
echo "OK: nginx reloaded ($SITE_PATH)"
