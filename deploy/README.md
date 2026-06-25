# Деплой на selenoid.autotests.cloud

Публичный Selenoid для курсов и примеров: **Selenium WebDriver** + **Playwright WebSocket**.

**Basic auth** (`user1` / `1234`): nginx на **`/wd/hub`** и **`/playwright/`**. UI (`/`) открыт без пароля.

| Путь | Auth | Как подключаться |
|------|------|------------------|
| `/` (UI) | нет | `https://selenoid.autotests.cloud` |
| `/wd/hub` | **да** | `http://user1:1234@selenoid.autotests.cloud/wd/hub` или Create Session в UI |
| `/playwright/` | **да** | Create Session в UI (basic auth) или `wss://user1:1234@.../playwright/chromium/1.61.1` |
| `/status` | нет | `https://selenoid.autotests.cloud/status` |
| `:4445` | **да** | прямой hub API для CI (`Authorization: Basic …`) |

Справочный полный конфиг: [`nginx-selenoid.conf`](nginx-selenoid.conf).

## Endpoints

| Назначение | URL |
|------------|-----|
| Selenium | `http://user1:1234@selenoid.autotests.cloud/wd/hub` |
| Playwright | `wss://user1:1234@selenoid.autotests.cloud/playwright/chromium/1.61.1` |
| UI | `https://selenoid.autotests.cloud/` (без auth) |
| Status | `https://selenoid.autotests.cloud/status` |
| Video | `https://selenoid.autotests.cloud/video/` |

### Переменные для тестов

```bash
# Selenium (auth обязателен)
export SELENOID_URL=http://user1:1234@selenoid.autotests.cloud/wd/hub

# Playwright — auth на /playwright/ (как у /wd/hub):
export PW_TEST_CONNECT_WS_ENDPOINT=wss://user1:1234@selenoid.autotests.cloud/playwright/chromium/1.61.1

# Альтернатива для Playwright (если клиент не принимает user:pass в URL):
export PW_TEST_CONNECT_WS_ENDPOINT=wss://selenoid.autotests.cloud/playwright/chromium/1.61.1
export PW_TEST_CONNECT_HEADERS='{"Authorization":"Basic dXNlcjE6MTIzNA=="}'

# Или отдельно для WebDriver:
export SELENOID_HOST=selenoid.autotests.cloud
export SELENOID_USER=user1
export SELENOID_PASSWORD=1234
```

Примеры всех стеков: [selenoid_selenium_playwright_tests](https://github.com/qa-guru/selenoid_selenium_playwright_tests).

---

## Автодеплой (GitHub Actions)

Workflow [`.github/workflows/deploy.yml`](../.github/workflows/deploy.yml) в репозитории **qa-guru/cm**:

| Триггер | Когда |
|---------|-------|
| `release: published` | Новый релиз cm (стек v2.x) |
| `workflow_dispatch` | Ручной запуск из Actions |
| `repository_dispatch: deploy-selenoid` | Вызов из других репозиториев qa-guru |

### Secrets (Settings → Secrets → Actions в qa-guru/cm)

| Secret | Пример | Описание |
|--------|--------|----------|
| `SELENOID_DEPLOY_HOST` | `136.243.89.21` | SSH-хост (**IP сервера**, не CNAME — DNS может указывать на прокси) |
| `SELENOID_DEPLOY_USER` | `selenoid` | Пользователь в группе `docker` |
| `SELENOID_DEPLOY_KEY` | ed25519 private key | Ключ только для Actions → `/home/selenoid/.ssh/authorized_keys` |

Опционально — Variables:

| Variable | Default | Описание |
|----------|---------|----------|
| `SELENOID_CONFIG_DIR` | `/opt/selenoid` | Каталог конфигурации на сервере |
| `SELENOID_PUBLIC_URL` | `https://selenoid.autotests.cloud` | URL для smoke test (не IP — иначе nginx отдаёт чужой cert) |

После настройки secrets: **Actions → deploy → Run workflow** — обновит сервер до latest release.

---

## Ручной деплой на сервере

### Первый раз (bootstrap)

```bash
# на сервере, от root
sudo DEPLOY_USER=selenoid ./deploy/bootstrap.sh
# перелогиниться, чтобы применилась группа docker
```

### Обновление стека

```bash
curl -sL https://raw.githubusercontent.com/qa-guru/cm/main/deploy/deploy.sh -o deploy.sh
chmod +x deploy.sh
./deploy.sh
```

Или из клона репозитория:

```bash
./deploy/deploy.sh
```

Быстрое обновление без полного `deploy.sh`:

```bash
./deploy/remote-update.sh
```

Pin версии (опционально, по умолчанию **v2.0.2**):

```bash
SELENOID_VERSION=v2.0.2 ./deploy/deploy.sh
```

### Проверка

```bash
./deploy/smoke-remote.sh https://selenoid.autotests.cloud
```

---

## Nginx (selenoid.autotests.cloud)

Реальный конфиг на сервере: **`/etc/nginx/sites-available/selenoid`**

| Порт | `location` | Куда | Auth |
|------|------------|------|------|
| 443 | `/` | `127.0.0.1:8080` (UI) | нет |
| 443 | `/wd/hub` | `127.0.0.1:8080` (UI → hub) | **да** |
| 443 | `/playwright/` | `127.0.0.1:8080` (UI → hub) | **да** |
| 443 | `/status` | `127.0.0.1:4444` | нет |
| 4445 | `/` | `127.0.0.1:4444` (hub) | **да** — CI / Playwright с заголовком `Authorization` |

Не проксируйте `/wd/hub` и `/playwright/` напрямую на hub:443 — иначе WebSocket Playwright не получит basic auth в браузере. Проксируйте через selenoid-ui.

Справочные файлы: [`nginx-selenoid.conf`](nginx-selenoid.conf), [`sync-nginx.sh`](sync-nginx.sh).

Применить вручную (если CI не смог из‑за sudo):

```bash
curl -fsSL https://raw.githubusercontent.com/qa-guru/cm/main/deploy/nginx-selenoid.conf -o /tmp/nginx-selenoid.conf
curl -fsSL https://raw.githubusercontent.com/qa-guru/cm/main/deploy/sync-nginx.sh -o /opt/selenoid/bin/sync-nginx.sh
chmod +x /opt/selenoid/bin/sync-nginx.sh
sudo NGINX_CONF_SRC=/tmp/nginx-selenoid.conf /opt/selenoid/bin/sync-nginx.sh
```

После `bootstrap.sh` пользователь `selenoid` может вызывать `sync-nginx.sh` без пароля.

### Очистка видео на сервере

Скрипт [`cleanup-videos.sh`](cleanup-videos.sh) удаляет `.mp4` старше 6 месяцев из `/opt/selenoid/video`. На проде — в root crontab (ежемесячно).

---

## Структура на сервере

```
/opt/selenoid/          # SELENOID_CONFIG_DIR
  browsers.json
  bin/selenoid
  bin/selenoid-ui
  video/
  logs/
/home/selenoid/cm       # бинарник cm (только у пользователя selenoid)
```

Деплой и `cm` — **только от пользователя `selenoid`**, не от root и не из home других пользователей.

---

## Релизы стека

| Версия | Документация |
|--------|--------------|
| v2.0.2 | [RELEASE_v2.0.2.md](RELEASE_v2.0.2.md) |
| v2.0.1 | [RELEASE_v2.0.1.md](RELEASE_v2.0.1.md) |
| v2.0.0 | [RELEASE_v2.0.0.md](RELEASE_v2.0.0.md) |
