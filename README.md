# Configuration Manager (qa-guru fork)

Установщик стека **qa-guru/selenoid** — форк [aerokube/cm](https://github.com/aerokube/cm). Одна команда на чистом сервере с Docker поднимает hub, UI, `browsers.json` и browser-образы.

[![Build Status](https://github.com/qa-guru/cm/workflows/build/badge.svg)](https://github.com/qa-guru/cm/actions?query=workflow%3Abuild)
[![Go Report Card](https://goreportcard.com/badge/github.com/qa-guru/cm)](https://goreportcard.com/report/github.com/qa-guru/cm)
[![Release](https://img.shields.io/github/release/qa-guru/cm.svg)](https://github.com/qa-guru/cm/releases/latest)
[![Docker Pulls](https://img.shields.io/docker/pulls/qaguru/cm.svg)](https://hub.docker.com/r/qaguru/cm)

| | |
|---|---|
| **GitHub** | [qa-guru/cm](https://github.com/qa-guru/cm) |
| **Docker Hub** | [`qaguru/cm`](https://hub.docker.com/r/qaguru/cm) |

## Роль в экосистеме

`cm` не содержит логику hub — он **оркестрирует установку** компонентов из других репозиториев:

```
cm selenoid start
    ├── скачивает qaguru/selenoid (hub) из GitHub Releases
    ├── записывает browsers.json (Chrome/Firefox + Playwright)
    ├── docker pull twilio/selenoid + qaguru/playwright-*
    └── запускает контейнер hub

cm selenoid-ui start
    ├── скачивает qaguru/selenoid-ui из GitHub Releases
    └── запускает контейнер UI, связанный с hub
```

## Связанные репозитории

| GitHub | Что делает cm с ним |
|--------|---------------------|
| [selenoid](https://github.com/qa-guru/selenoid) | Скачивает бинарник hub, синхронизирует `browsers.json` |
| [selenoid-ui](https://github.com/qa-guru/selenoid-ui) | Скачивает бинарник UI |
| **cm** (этот) | Установщик |
| [playwright-image](https://github.com/qa-guru/playwright-image) | `docker pull qaguru/playwright-*` по `browsers.json` |

После изменения [`config/browsers.json` в qa-guru/selenoid](https://github.com/qa-guru/selenoid/blob/main/config/browsers.json) обновите копии в этом репозитории:

- `selenoid/data/browsers-qaguru.json`
- `deploy/browsers-production.json`

## Установка

**Предварительные условия:** Docker Engine 26.1.x (API 1.45), доступ пользователя к `docker`, опубликованные релизы [qa-guru/selenoid](https://github.com/qa-guru/selenoid/releases) и [qa-guru/selenoid-ui](https://github.com/qa-guru/selenoid-ui/releases). Для сборки cm — Go 1.23.x.

Скачать бинарник из GitHub Releases или собрать локально (см. [Сборка](#сборка)):

```bash
curl -sL https://github.com/qa-guru/cm/releases/latest/download/cm_linux_amd64 -o cm
chmod +x cm
```

Запуск:

```bash
./cm selenoid start
./cm selenoid-ui start
```

`cm selenoid start` передаёт hub `DOCKER_API_VERSION=1.45`.

Обновление:

```bash
./cm selenoid update
./cm selenoid-ui update
./cm selenoid start -f   # принудительно перекачать образы и бинарники
```

Production-стек **v2.0.9** (verified): [deploy/RELEASE_v2.0.9.md](deploy/RELEASE_v2.0.9.md). Endpoints и nginx — [`deploy/README.md`](deploy/README.md).

## Сборка

```bash
go build -o cm .
```

## Флаги

| Флаг | Описание |
|------|----------|
| `-v, --version` | Тег GitHub Release (по умолчанию **latest**) |
| `-j, --browsers-json` | Свой `browsers.json` вместо встроенного |
| `--selenoid-binary` | Путь к бинарнику hub вместо скачивания из Release |
| `--selenoid-ui-binary` | Путь к бинарнику UI вместо скачивания из Release |
| `-f, --force` | Перекачать образы и бинарники |

Пример с локальными бинарниками:

```bash
./cm selenoid start --selenoid-binary /opt/selenoid/bin/selenoid
./cm selenoid-ui start --selenoid-ui-binary /opt/selenoid/bin/selenoid-ui
```

## Структура на сервере

```
~/.aerokube/selenoid/
  browsers.json
  bin/selenoid
  bin/selenoid-ui
  video/
  logs/
```
