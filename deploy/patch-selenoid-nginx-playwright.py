#!/usr/bin/env python3
"""Patch /etc/nginx/sites-available/selenoid on selenoid.autotests.cloud.

Inserts location /playwright/ after /wd/hub/ in the listen-443 server block.
Matches existing style: auth_basic 'test', /etc/nginx/.htpasswd, $connection_upgrade map.

Usage on server:
  sudo cp /etc/nginx/sites-available/selenoid /etc/nginx/sites-available/selenoid.bak.$(date +%Y%m%d)
  sudo python3 patch-selenoid-nginx-playwright.py
  sudo nginx -t && sudo systemctl reload nginx
"""
from pathlib import Path

CONF = Path("/etc/nginx/sites-available/selenoid")
text = CONF.read_text()

if "location /playwright/" in text:
    print("SKIP: location /playwright/ already exists")
    raise SystemExit(0)

PLAYWRIGHT_BLOCK = """
        location /playwright/ {
                auth_basic           'test';
                auth_basic_user_file /etc/nginx/.htpasswd;
                proxy_pass http://127.0.0.1:4444/playwright/;
                proxy_http_version 1.1;
                proxy_set_header Upgrade $http_upgrade;
                proxy_set_header Connection $connection_upgrade;
                proxy_set_header Host $host;
                proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
                proxy_set_header X-Real-IP $remote_addr;
                proxy_read_timeout 3600s;
                proxy_send_timeout 3600s;
        }

"""

needle = """        location /wd/hub/ {
                auth_basic           'test';
                auth_basic_user_file /etc/nginx/.htpasswd;
                proxy_pass http://127.0.0.1:4444/wd/hub/;
                proxy_http_version 1.1;
                proxy_set_header Upgrade $http_upgrade;
                proxy_set_header Connection $connection_upgrade;
                proxy_set_header Host $host;
                proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
                proxy_set_header X-Real-IP $remote_addr;
        }


        location / {"""

if needle not in text:
    raise SystemExit("ERROR: anchor after /wd/hub/ not found — inspect /etc/nginx/sites-available/selenoid manually")

replacement = needle.replace(
    "\n\n        location / {",
    "\n" + PLAYWRIGHT_BLOCK + "        location / {",
    1,
)
CONF.write_text(text.replace(needle, replacement, 1))
print("OK: inserted location /playwright/ after /wd/hub/")
