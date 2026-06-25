#!/usr/bin/env bash
# Осторожно включает basic auth на /playwright/ в СУЩЕСТВУЮЩЕМ nginx-конфиге.
# Блок /wd/hub не переписывается — только читается путь к htpasswd оттуда.
#
# Запуск на сервере:
#   sudo ./nginx-enable-playwright-auth.sh --dry-run   # посмотреть план
#   sudo ./nginx-enable-playwright-auth.sh             # применить
set -euo pipefail

DRY_RUN=0
NGINX_CONF=""
HTPASSWD_FILE=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --dry-run) DRY_RUN=1; shift ;;
    --conf) NGINX_CONF="${2:?}"; shift 2 ;;
    --htpasswd) HTPASSWD_FILE="${2:?}"; shift 2 ;;
    -h|--help)
      sed -n '2,12p' "$0"
      exit 0
      ;;
    *) echo "Unknown option: $1" >&2; exit 2 ;;
  esac
done

if [[ "$(id -u)" -ne 0 ]]; then
  echo "Run as root: sudo $0 [--dry-run]" >&2
  exit 1
fi

find_nginx_conf() {
  if [[ -n "$NGINX_CONF" ]]; then
    echo "$NGINX_CONF"
    return
  fi
  local f
  while IFS= read -r f; do
    if grep -q 'location /wd/hub' "$f" 2>/dev/null; then
      echo "$f"
      return
    fi
  done < <(grep -rl 'selenoid\.autotests\.cloud\|location /wd/hub' /etc/nginx 2>/dev/null || true)
  echo ""
}

detect_htpasswd_from_wd_hub() {
  local conf="$1"
  awk '
    /location[[:space:]]+\/wd\/hub/ { in_wd=1 }
    in_wd && /auth_basic_user_file[[:space:]]+/ {
      gsub(/;/, "", $2)
      print $2
      exit
    }
    in_wd && /^[[:space:]]*}/ { in_wd=0 }
  ' "$conf"
}

playwright_block_has_auth() {
  local conf="$1"
  awk '
    /location[[:space:]]+\/playwright\// { in_pw=1 }
    in_pw && /auth_basic[[:space:]]/ { found=1 }
    in_pw && /^[[:space:]]*}/ { if (found) exit 0; in_pw=0 }
    END { exit found ? 0 : 1 }
  ' "$conf"
}

playwright_block_exists() {
  grep -q 'location /playwright/' "$1"
}

CONF="$(find_nginx_conf)"
if [[ -z "$CONF" || ! -f "$CONF" ]]; then
  echo "Не найден nginx-конфиг с location /wd/hub. Укажите: --conf /path/to/site.conf" >&2
  exit 1
fi

if [[ -z "$HTPASSWD_FILE" ]]; then
  HTPASSWD_FILE="$(detect_htpasswd_from_wd_hub "$CONF")"
fi
if [[ -z "$HTPASSWD_FILE" || ! -f "$HTPASSWD_FILE" ]]; then
  echo "Не найден htpasswd (из /wd/hub). Укажите: --htpasswd /path/to/file" >&2
  exit 1
fi

echo "=== nginx-enable-playwright-auth ==="
echo "Config:   $CONF"
echo "Htpasswd: $HTPASSWD_FILE (тот же, что у /wd/hub)"
echo "Dry-run:  $DRY_RUN"
echo

if playwright_block_has_auth "$CONF"; then
  echo "OK: auth_basic уже есть в location /playwright/ — ничего не делаем."
  exit 0
fi

BACKUP="${CONF}.bak.$(date +%Y%m%d%H%M%S)"
PATCHED="$(mktemp)"

if playwright_block_exists "$CONF"; then
  echo "План: добавить auth_basic в существующий location /playwright/ (блок /wd/hub не трогаем)."
  awk -v htpasswd="$HTPASSWD_FILE" '
    /location[[:space:]]+\/playwright\// {
      print
      print "        auth_basic \"Selenoid\";"
      print "        auth_basic_user_file " htpasswd ";"
      in_pw=1
      next
    }
    in_pw && /auth_basic[[:space:]]/ { next }
    in_pw && /auth_basic_user_file[[:space:]]/ { next }
    { print }
  ' "$CONF" > "$PATCHED"
else
  echo "План: добавить новый location /playwright/ с auth (блок /wd/hub не трогаем)."
  SNIPPET="$(mktemp)"
  cat > "$SNIPPET" <<EOF

    # Playwright WebSocket — basic auth как у /wd/hub (добавлено nginx-enable-playwright-auth.sh)
    location /playwright/ {
        auth_basic "Selenoid";
        auth_basic_user_file ${HTPASSWD_FILE};

        proxy_pass http://127.0.0.1:4444;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_read_timeout 3600s;
        proxy_send_timeout 3600s;
    }
EOF
  awk -v snippetfile="$SNIPPET" '
    /location[[:space:]]+\/wd\/hub/ { saw_wd=1 }
    saw_wd && /^[[:space:]]*}/ {
      print
      while ((getline line < snippetfile) > 0) print line
      close(snippetfile)
      saw_wd=0
      next
    }
    { print }
  ' "$CONF" > "$PATCHED"
  rm -f "$SNIPPET"
fi

echo "--- diff (фрагмент) ---"
diff -u "$CONF" "$PATCHED" | head -80 || true
echo "--- end diff ---"
echo

if [[ "$DRY_RUN" -eq 1 ]]; then
  echo "Dry-run: файл не изменён. Для применения запустите без --dry-run."
  rm -f "$PATCHED"
  exit 0
fi

cp -a "$CONF" "$BACKUP"
mv "$PATCHED" "$CONF"
echo "Backup: $BACKUP"

nginx -t
systemctl reload nginx
echo "OK: nginx reloaded. Проверка:"
echo "  curl -u user1:1234 -o /dev/null -w '%{http_code}' https://selenoid.autotests.cloud/wd/hub/status"
echo "  # Playwright — из теста с wss://user1:1234@selenoid.autotests.cloud/playwright/..."
