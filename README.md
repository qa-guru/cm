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

- Docker и доступ пользователя к `docker` (группа `docker`)
- Опубликованные релизы [qa-guru/selenoid](https://github.com/qa-guru/selenoid/releases) и [qa-guru/selenoid-ui](https://github.com/qa-guru/selenoid-ui/releases) (сейчас: **v2.0.0**)

Альтернатива — локальные бинарники:

```bash
./cm selenoid start \
  --selenoid-binary /opt/selenoid/bin/selenoid \
  --selenoid-ui-binary /opt/selenoid/bin/selenoid-ui   # для selenoid-ui start
```

## Установка на selenoid.autotests.cloud

```bash
# пользователь в группе docker
curl -sL https://github.com/qa-guru/cm/releases/latest/download/cm_linux_amd64 -o cm
chmod +x cm

./cm selenoid start -v v2.0.0    # тег релиза qa-guru/selenoid
./cm selenoid-ui start -v v2.0.0
```

Nginx для Playwright WebSocket — см. `deploy/nginx-playwright-snippet.conf`.

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
| `-v, --version` | Тег GitHub Release для бинарников qa-guru (не тег Docker-обёртки) |
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
