#!/usr/bin/env bash
# Post-deploy smoke checks against selenoid.autotests.cloud (or any public base URL).
set -euo pipefail

BASE_URL="${1:-https://selenoid.autotests.cloud}"
BASE_URL="${BASE_URL%/}"

echo "=== GET $BASE_URL/status ==="
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

echo "=== GET $BASE_URL/wd/hub/status (Selenium) ==="
curl -fsSL -o /dev/null -w "HTTP %{http_code}\n" "$BASE_URL/wd/hub/status" || \
  curl -fsSL -o /dev/null -w "HTTP %{http_code}\n" "$BASE_URL/status"

echo "Smoke OK: $BASE_URL"
