#!/usr/bin/env bash
# Post-deploy smoke checks against selenoid.autotests.cloud (or any public base URL).
set -euo pipefail

BASE_URL="${1:-https://selenoid.autotests.cloud}"
BASE_URL="${BASE_URL%/}"
SELENOID_USER="${SELENOID_USER:-user1}"
SELENOID_PASSWORD="${SELENOID_PASSWORD:-1234}"
AUTH=(-u "${SELENOID_USER}:${SELENOID_PASSWORD}")

echo "=== GET $BASE_URL/status (no auth required) ==="
status_json="$(curl -fsSL "$BASE_URL/status")"
echo "$status_json" | (command -v jq >/dev/null && jq . || cat)

if ! command -v jq >/dev/null; then
  echo "jq not found — skipping browser version assertions" >&2
  exit 0
fi

echo "=== browser versions ==="
for pair in "chrome:148.0" "chromium:1.61.1" "firefox:150.0"; do
  browser="${pair%%:*}"
  version="${pair##*:}"
  if jq -e --arg b "$browser" --arg v "$version" '.browsers[$b][$v] != null' <<<"$status_json" >/dev/null; then
    echo "OK  $browser $version"
  else
    echo "FAIL $browser $version not in /status" >&2
    exit 1
  fi
done

echo "=== GET $BASE_URL/ (UI, no auth) ==="
ui_code="$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/")"
if [[ "$ui_code" == "200" ]]; then
  echo "OK  UI is public (HTTP 200)"
else
  echo "FAIL UI should be public without credentials (HTTP $ui_code)" >&2
  exit 1
fi

echo "=== GET $BASE_URL/wd/hub/status without auth (expect 401) ==="
wd_no_auth="$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/wd/hub/status" || true)"
if [[ "$wd_no_auth" == "401" ]]; then
  echo "OK  /wd/hub requires auth (HTTP 401)"
else
  echo "FAIL /wd/hub should require auth (HTTP $wd_no_auth)" >&2
  exit 1
fi

echo "=== GET $BASE_URL/wd/hub/status (with basic auth) ==="
wd_code="$(curl -fsSL "${AUTH[@]}" -o /dev/null -w "%{http_code}" "$BASE_URL/wd/hub/status")"
if [[ "$wd_code" == "200" ]]; then
  echo "OK  /wd/hub with auth (HTTP 200)"
else
  echo "FAIL /wd/hub with auth: HTTP $wd_code" >&2
  exit 1
fi

echo "=== GET $BASE_URL/playwright/... without auth (expect 401) ==="
pw_code="$(curl -s -o /dev/null -w "%{http_code}" --max-time 5 "$BASE_URL/playwright/playwright-chromium/1.61.1" || true)"
if [[ "$pw_code" == "401" ]]; then
  echo "OK  /playwright/ requires auth (HTTP 401)"
else
  echo "FAIL /playwright/ should require auth (HTTP $pw_code)" >&2
  exit 1
fi

echo "Smoke OK: $BASE_URL (auth: $SELENOID_USER:***)"
