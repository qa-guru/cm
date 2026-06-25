# Деплой на selenoid.autotests.cloud

Публичный Selenoid для курсов и примеров: **Selenium WebDriver** + **Playwright WebSocket**.

## Endpoints

| Назначение | URL |
|------------|-----|
| Selenium | `http://selenoid.autotests.cloud/wd/hub` |
| Playwright | `wss://selenoid.autotests.cloud/playwright/chromium/1.61.1` |
| UI | `http://selenoid.autotests.cloud:8080/` |
| Status | `https://selenoid.autotests.cloud/status` |
| Video | `https://selenoid.autotests.cloud/video/` |

### Переменные для тестов

```bash
# Playwright
export PW_TEST_CONNECT_WS_ENDPOINT=wss://selenoid.autotests.cloud/playwright/chromium/1.61.1

# Selenium
export SELENOID_URL=http://selenoid.autotests.cloud/wd/hub
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
| `SELENOID_DEPLOY_USER` | `selenoid` | Пользователь в группе `docker` |
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

## Nginx

Playwright требует WebSocket-прокси. Фрагмент для nginx: [`nginx-playwright-snippet.conf`](nginx-playwright-snippet.conf).

```nginx
server {
    listen 443 ssl;
    server_name selenoid.autotests.cloud;

    location / {
        proxy_pass http://127.0.0.1:4444;
        # ... стандартные proxy_set_header ...
    }

    # вставить содержимое nginx-playwright-snippet.conf
}
```

---

## Структура на сервере

```
/opt/selenoid/          # SELENOID_CONFIG_DIR (рекомендуется)
  browsers.json
  bin/selenoid
  bin/selenoid-ui
  video/
  logs/
~/cm                    # бинарник cm
```

Не используйте `sudo ./cm` — иначе данные окажутся в `/root/.aerokube/`.

---

## Миграция со старого deploy-remote.sh

Старый скрипт с ручным `docker run` и `/root/.aerokube/` **устарел**. Новый путь:

```bash
sudo docker stop selenoid selenoid-ui && sudo docker rm selenoid selenoid-ui
sudo DEPLOY_USER=selenoid ./deploy/bootstrap.sh
# от пользователя selenoid:
SELENOID_CONFIG_DIR=/opt/selenoid ./deploy/deploy.sh
```

Видео из старого каталога (если нужны):

```bash
sudo cp -a /root/.aerokube/selenoid/video/* /opt/selenoid/video/ 2>/dev/null || true
```

---

## Релизы стека

| Версия | Документация |
|--------|--------------|
| v2.0.1 | [RELEASE_v2.0.1.md](RELEASE_v2.0.1.md) |
| v2.0.0 | [RELEASE_v2.0.0.md](RELEASE_v2.0.0.md) |
