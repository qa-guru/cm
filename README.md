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
| **Текущий релиз** | **v2.2.0** — [docs/RELEASE_v2.2.0.md](docs/RELEASE_v2.2.0.md) · `qaguru/cm:v2.2.0` |

## Роль в экосистеме

`cm` не содержит логику hub — он **оркестрирует установку** компонентов из других репозиториев:

```
cm selenoid start
    ├── скачивает qaguru/selenoid (hub) из GitHub Releases
    ├── записывает browsers.json (Chrome/Firefox + Playwright)
    ├── docker pull qaguru/webdriver-chrome + qaguru/playwright-*
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
| [browser-image](https://github.com/qa-guru/browser-image) | `docker pull qaguru/webdriver-*` / `qaguru/playwright-*` по `browsers.json` |

## browsers.json — один SSOT (Docker) + drivers catalog

| Где | Назначение |
|-----|------------|
| `dev/browsers.json` (selenoid-home) | **SSOT** полный Docker-каталог (`qaguru/*` images) |
| `selenoid/data/browsers.json` (этот репо, go:embed) | Копия SSOT для `cm selenoid configure` / `start` без `-j` |
| [`selenoid/config/browsers.json`](https://github.com/qa-guru/selenoid/blob/main/config/browsers.json) | Копия SSOT на hub |
| [`browsers.json`](browsers.json) (корень этого репо) | **Не** SSOT: drivers-mode (`cm selenoid start --drivers`), CfT/geckodriver URLs; chrome CfT **149.0.7827.55** |
| [`selenoid.autotests.cloud/deploy/browsers-production.json`](https://github.com/qa-guru/selenoid.autotests.cloud/blob/main/deploy/browsers-production.json) | Prod overlay (только по явному prod-запросу) |

Правки Docker-каталога — в `dev/browsers.json`, затем `dev/scripts/sync-cm-browsers.sh` (не трогает корневой drivers `browsers.json`).

## Установка

**Предварительные условия:** Docker Engine **26.1.x** (API **1.45**, moby/moby client), доступ пользователя к `docker`, опубликованные релизы [qa-guru/selenoid](https://github.com/qa-guru/selenoid/releases) и [qa-guru/selenoid-ui](https://github.com/qa-guru/selenoid-ui/releases). Для сборки cm — Go **1.26.5**.

Скачать бинарник из GitHub Releases или собрать локально (см. [Сборка](#сборка)):

```bash
curl -sL https://github.com/qa-guru/cm/releases/latest/download/cm_linux_amd64 -o cm
chmod +x cm
```

Запуск (по умолчанию каталог **`/opt/selenoid`**, как на prod):

```bash
sudo mkdir -p /opt/selenoid
./cm selenoid start
./cm selenoid-ui start
```

Локально без прав на `/opt` — свой каталог:

```bash
./cm selenoid start -c "$HOME/selenoid"
./cm selenoid-ui start -c "$HOME/selenoid"
```

`cm selenoid start` передаёт hub `DOCKER_API_VERSION=1.45`.

Обновление:

```bash
./cm selenoid update
./cm selenoid-ui update
./cm selenoid start -f   # принудительно перекачать образы и бинарники
```

Текущий релиз cm: **v2.2.0** — [release notes](docs/RELEASE_v2.2.0.md) (binary + embed browsers).  
Предыдущий: **v2.1.7** — [GitHub Releases](https://github.com/qa-guru/cm/releases/tag/v2.1.7) · [notes](docs/RELEASE_v2.1.7.md). Docker: `qaguru/cm:v2.2.0`.

## Сборка

```bash
go build -o cm .
```

## Флаги (Docker-режим, по умолчанию)

| Флаг | Описание |
|------|----------|
| `-v, --version` | Тег GitHub Release (по умолчанию **latest**) |
| `-j, --browsers-json` | Свой `browsers.json` вместо встроенного |
| `--selenoid-binary` | Путь к бинарнику hub вместо скачивания из Release |
| `--selenoid-ui-binary` | Путь к бинарнику UI вместо скачивания из Release |
| `-c, --config-dir` | Каталог данных (default: **`/opt/selenoid`**) |
| `-f, --force` | Перекачать образы и бинарники |
| `-n, --no-download` | Только записать `browsers.json`, без `docker pull` |

Флаги `--browsers`, `--browser-env`, `--drivers-info` — только для режима `--use-drivers` (локальные WebDriver-бинарники без Docker).

Пример с локальными бинарниками:

```bash
./cm selenoid start --selenoid-binary /opt/selenoid/bin/selenoid
./cm selenoid-ui start --selenoid-ui-binary /opt/selenoid/bin/selenoid-ui
```

## Структура на сервере

```
/opt/selenoid/
  browsers.json
  bin/selenoid
  bin/selenoid-ui
  video/
  logs/
```
