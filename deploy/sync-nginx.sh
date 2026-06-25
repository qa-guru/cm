#!/usr/bin/env bash
# Apply nginx-selenoid.conf on the server (requires sudo).
# Preserves ssl_certificate* lines from the existing site config when present.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONF_SRC="${NGINX_CONF_SRC:-${SCRIPT_DIR}/nginx-selenoid.conf}"
SITE_NAME="${NGINX_SITE_NAME:-selenoid}"
SITE_PATH="/etc/nginx/sites-available/${SITE_NAME}"
TMP="/tmp/nginx-selenoid.generated"
SSL_SNIPPET="/tmp/nginx-selenoid.ssl-snippet"

if [[ ! -f "$CONF_SRC" ]]; then
  echo "Missing $CONF_SRC (set NGINX_CONF_SRC or place nginx-selenoid.conf next to this script)" >&2
  exit 1
fi

if [[ "$(id -u)" -ne 0 ]]; then
  if sudo -n true 2>/dev/null; then
    exec sudo env NGINX_CONF_SRC="$CONF_SRC" NGINX_SITE_NAME="$SITE_NAME" "$0" "$@"
  fi
  echo "Run as root or with passwordless sudo: sudo $0" >&2
  exit 1
fi

cp "$CONF_SRC" "$TMP"

if [[ -f "$SITE_PATH" ]]; then
  grep -E '^\s*ssl_' "$SITE_PATH" >"$SSL_SNIPPET" || true
  if [[ -s "$SSL_SNIPPET" ]]; then
    awk -v sslfile="$SSL_SNIPPET" '
      /# ssl_certificate \.\.\.;/ {
        while ((getline line < sslfile) > 0) print line
        close(sslfile)
        next
      }
      { print }
    ' "$TMP" >"${TMP}.patched"
    mv "${TMP}.patched" "$TMP"
  fi
fi

cp "$TMP" "$SITE_PATH"
ln -sf "$SITE_PATH" "/etc/nginx/sites-enabled/${SITE_NAME}"
nginx -t
systemctl reload nginx
echo "OK: nginx reloaded ($SITE_PATH)"
