# Деплой на selenoid.autotests.cloud

Публичный Selenoid для курсов и примеров: **Selenium WebDriver** + **Playwright WebSocket**.

**Basic auth** (`user1` / `1234`): настроен в **nginx** для защищённых endpoint'ов.

| Путь | Auth сейчас на сервере | Как подключаться |
|------|------------------------|------------------|
| `/wd/hub` | **да** | `http://user1:1234@selenoid.autotests.cloud/wd/hub` |
| `/playwright/` | **да** (с июня 2026) | `wss://user1:1234@selenoid.autotests.cloud/playwright/chromium/1.61.1` |
| `/status` | нет | без логина |

Чтобы закрыть Playwright тем же паролем — добавьте `auth_basic` в `location /playwright/` (см. [`nginx-selenoid.conf`](nginx-selenoid.conf)).

## Endpoints

| Назначение | URL |
|------------|-----|
| Selenium | `http://user1:1234@selenoid.autotests.cloud/wd/hub` |
| Playwright | `wss://user1:1234@selenoid.autotests.cloud/playwright/chromium/1.61.1` |
| UI | `http://selenoid.autotests.cloud:8080/` |
| Status | `https://selenoid.autotests.cloud/status` |
| Video | `https://selenoid.autotests.cloud/video/` |

### Переменные для тестов

```bash
# Selenium (auth обязателен)
export SELENOID_URL=http://user1:1234@selenoid.autotests.cloud/wd/hub

# Playwright — когда auth включён в nginx для /playwright/:
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
| `SELENOID_DEPLOY_HOST` | `selenoid.autotests.cloud` | SSH-хост |
| `SELENOID_DEPLOY_USER` | `selenoid` | Пользователь в группе `docker` (создаётся `bootstrap.sh`) |
| `SELENOID_DEPLOY_KEY` | `-----BEGIN OPENSSH...` | Приватный SSH-ключ |

Опционально — Variables:

| Variable | Default | Описание |
|----------|---------|----------|
| `SELENOID_CONFIG_DIR` | `/opt/selenoid` | Каталог конфигурации на сервере |

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
curl -sL https://raw.githubusercontent.com/qa-guru/cm/master/deploy/deploy.sh -o deploy.sh
chmod +x deploy.sh
./deploy.sh
```

Или из клона репозитория:

```bash
./deploy/deploy.sh
```

Pin версии (опционально):

```bash
SELENOID_VERSION=v2.0.1 ./deploy/deploy.sh
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
| 443 | `/wd/hub/` | `127.0.0.1:4444` | `auth_basic 'test'`, `/etc/nginx/.htpasswd` |
| 443 | `/playwright/` | `127.0.0.1:4444` | то же (добавлено после `/wd/hub/`) |
| 443 | `/` | `127.0.0.1:8080` (UI) | нет |
| 443 | `/status` | `127.0.0.1:4444` | `auth_basic off` |
| 4445 | `/` | `127.0.0.1:4444` | `auth_basic 'API'`, тот же htpasswd |

До патча `/playwright/` на 443 попадал в `location /` → UI (8080), **без auth**.

Патч (уже применён на проде):

```bash
sudo cp /etc/nginx/sites-available/selenoid /etc/nginx/sites-available/selenoid.bak.$(date +%Y%m%d)
sudo python3 patch-selenoid-nginx-playwright.py   # из deploy/
sudo nginx -t && sudo systemctl reload nginx
```

Скрипт **не трогает** блок `/wd/hub/` — только вставляет `location /playwright/` сразу после него.

Справочные файлы: [`nginx-selenoid.conf`](nginx-selenoid.conf), [`nginx-playwright-snippet.conf`](nginx-playwright-snippet.conf).

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

## Миграция со старого deploy-remote.sh

Старый скрипт с ручным `docker run` и `/root/.aerokube/` **устарел**.

### Одной командой (с `/root/.aerokube` на `/opt/selenoid`)

```bash
# на сервере, от root
sudo ./deploy/migrate-to-opt.sh
```

Скрипт создаёт пользователя `selenoid`, переносит данные, права `selenoid:docker`, `cm` в `/home/selenoid/cm`, затем `deploy.sh`.

### Вручную

```bash
sudo docker stop selenoid selenoid-ui && sudo docker rm selenoid selenoid-ui
sudo DEPLOY_USER=selenoid ./deploy/bootstrap.sh
sudo rsync -a /root/.aerokube/selenoid/ /opt/selenoid/
sudo chown -R selenoid:docker /opt/selenoid
# от пользователя selenoid (после re-login для docker group):
SELENOID_CONFIG_DIR=/opt/selenoid ./deploy/deploy.sh
```

Видео из старого каталога (если не делали migrate-to-opt.sh):

```bash
sudo rsync -a /root/.aerokube/selenoid/video/ /opt/selenoid/video/
sudo chown -R selenoid:docker /opt/selenoid/video
```

После проверки legacy можно удалить: `sudo rm -rf /root/.aerokube /root/cm`.

---

## Релизы стека

| Версия | Документация |
|--------|--------------|
| v2.0.1 | [RELEASE_v2.0.1.md](RELEASE_v2.0.1.md) |
| v2.0.0 | [RELEASE_v2.0.0.md](RELEASE_v2.0.0.md) |
