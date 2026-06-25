# Configuration Manager (qa-guru fork)

Форк [aerokube/cm](https://github.com/aerokube/cm) для установки [qa-guru/selenoid](https://github.com/qa-guru/selenoid) с Playwright, twilio/selenoid и qaguru/playwright.

## Что делает

Одна команда на чистом сервере с Docker:

```bash
./cm selenoid start
./cm selenoid-ui start
```

1. Скачивает обёрточные образы `aerokube/selenoid:latest-release` и `aerokube/selenoid-ui:latest-release`
2. Скачивает бинарники из GitHub Releases [qa-guru/selenoid](https://github.com/qa-guru/selenoid) и [qa-guru/selenoid-ui](https://github.com/qa-guru/selenoid-ui)
3. Записывает встроенный `browsers.json` (Chrome/Firefox + Playwright) в `~/.aerokube/selenoid/`
4. Тянет все browser-образы из конфига
5. Запускает контейнеры с примонтированными бинарниками qa-guru

## Предварительные условия

- **Docker Engine 26.1.x** (API 1.45) — `cm selenoid start` передаёт hub `DOCKER_API_VERSION=1.45`
- **Go 1.23.x** — для сборки cm
- Docker и доступ пользователя к `docker` (группа `docker`)
- Опубликованные релизы [qa-guru/selenoid](https://github.com/qa-guru/selenoid/releases) и [qa-guru/selenoid-ui](https://github.com/qa-guru/selenoid-ui/releases) (сейчас: **v2.0.2** hub)

Альтернатива — локальные бинарники:

```bash
./cm selenoid start \
  --selenoid-binary /opt/selenoid/bin/selenoid \
  --selenoid-ui-binary /opt/selenoid/bin/selenoid-ui   # для selenoid-ui start
```

## Установка на selenoid.autotests.cloud

Basic auth для WebDriver и Playwright: **`user1` / `1234`** (nginx, `/wd/hub` и `/playwright/`), см. [`deploy/nginx-selenoid.conf`](deploy/nginx-selenoid.conf).

```bash
# пользователь в группе docker
curl -sL https://github.com/qa-guru/cm/releases/latest/download/cm_linux_amd64 -o cm
chmod +x cm

./cm selenoid start
./cm selenoid-ui start
```

Подключение тестов к публичному hub:

```bash
export SELENOID_URL=http://user1:1234@selenoid.autotests.cloud/wd/hub
export PW_TEST_CONNECT_WS_ENDPOINT=wss://user1:1234@selenoid.autotests.cloud/playwright/chromium/1.61.1
# Альтернатива, если клиент не принимает user:pass в URL:
# export PW_TEST_CONNECT_WS_ENDPOINT=wss://selenoid.autotests.cloud/playwright/chromium/1.61.1
# export PW_TEST_CONNECT_HEADERS='{"Authorization":"Basic dXNlcjE6MTIzNA=="}'
```

По умолчанию `cm` скачивает **последний** GitHub Release qa-guru/selenoid и qa-guru/selenoid-ui (`-v` / `--version` не нужен). Как в [документации Aerokube](https://aerokube.com/cm/latest/): `./cm selenoid start`, `./cm selenoid-ui start`.

Явная версия — только если нужен конкретный тег:

```bash
./cm selenoid start -v v2.0.2
```

Обновление уже установленного стека:

```bash
./cm selenoid update
./cm selenoid-ui update
# или перекачать бинарники принудительно:
./cm selenoid start -f
```

Nginx для Playwright WebSocket — см. [`deploy/nginx-playwright-snippet.conf`](deploy/nginx-playwright-snippet.conf).

Автодеплой на **selenoid.autotests.cloud** — см. [`deploy/README.md`](deploy/README.md) (GitHub Actions + ручной `./deploy/deploy.sh`).

## Обновление browsers.json в cm

После изменения `selenoid-src/config/browsers.json`:

```bash
./scripts/sync-cm-browsers.sh
```

## Сборка

```bash
cd cm-src
go build -o ../bin/cm .
```

## Флаги

| Флаг | Описание |
|------|----------|
| `-v, --version` | Опционально: тег GitHub Release (по умолчанию **latest** — последний релиз) |
| `-j, --browsers-json` | Свой `browsers.json` вместо встроенного |
| `--selenoid-binary` | Путь к бинарнику hub |
| `--selenoid-ui-binary` | Путь к бинарнику UI |
| `-f, --force` | Перекачать образы и бинарники |

## Структура на сервере

```
~/.aerokube/selenoid/
  browsers.json
  bin/selenoid
  bin/selenoid-ui
  video/
  logs/
```
