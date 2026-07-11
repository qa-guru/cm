# Configuration Manager (qa-guru fork)


[![Configuration Manager](https://qa-guru.github.io/selenoid-tests/readme/badge-cm.svg)](https://qa-guru.github.io/selenoid-tests/reports/latest/dashboard/)

[![CM stats](https://qa-guru.github.io/selenoid-tests/readme/stats-cm.svg)](https://qa-guru.github.io/selenoid-tests/reports/latest/dashboard/)

[![CM metrics](https://qa-guru.github.io/selenoid-tests/readme/metrics-panel-cm.svg)](https://qa-guru.github.io/selenoid-tests/reports/latest/dashboard/)

| Link | Description |
|------|-------------|
| [Dashboard](https://qa-guru.github.io/selenoid-tests/reports/latest/dashboard/) | Full pyramid — filter epic **cm** in awesome |
| [Awesome](https://qa-guru.github.io/selenoid-tests/reports/latest/awesome/) | CM integration test details |
| [selenoid-tests](https://github.com/qa-guru/selenoid-tests) | Orchestrator + merged Allure |


Установщик стека **qa-guru/selenoid** — форк [aerokube/cm](https://github.com/aerokube/cm). Одна команда на чистом сервере с Docker поднимает hub, UI, `browsers.json` и browser-образы.

[![Build Status](https://github.com/qa-guru/cm/workflows/build/badge.svg)](https://github.com/qa-guru/cm/actions?query=workflow%3Abuild)
[![Go Report Card](https://goreportcard.com/badge/github.com/qa-guru/cm)](https://goreportcard.com/report/github.com/qa-guru/cm)
[![Release](https://img.shields.io/github/release/qa-guru/cm.svg)](https://github.com/qa-guru/cm/releases/latest)
[![Docker Pulls](https://img.shields.io/docker/pulls/qaguru/cm.svg)](https://hub.docker.com/r/qaguru/cm)

| | |
|---|---|
| **GitHub** | [qa-guru/cm](https://github.com/qa-guru/cm) |
| **Docker Hub** | [`qaguru/cm`](https://hub.docker.com/r/qaguru/cm) |
| **Текущий релиз** | **v2.2.1** — [docs/RELEASE_v2.2.1.md](docs/RELEASE_v2.2.1.md) · `qaguru/cm:v2.2.1` |

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

## Экосистема qa-guru Selenoid

| Ресурс | Ссылка | Роль |
|--------|--------|------|
| selenoid | [github.com/qa-guru/selenoid](https://github.com/qa-guru/selenoid) | Hub |
| selenoid-ui | [github.com/qa-guru/selenoid-ui](https://github.com/qa-guru/selenoid-ui) | Web UI |
| **cm** (этот) | [github.com/qa-guru/cm](https://github.com/qa-guru/cm) | Установщик |
| browser-image | [github.com/qa-guru/browser-image](https://github.com/qa-guru/browser-image) | Docker browser nodes |
| selenoid-tests | [github.com/qa-guru/selenoid-tests](https://github.com/qa-guru/selenoid-tests) | E2e/integration ethalon |
| Docker Hub | [hub.docker.com/u/qaguru](https://hub.docker.com/u/qaguru) | Образы `qaguru/*` |

## browsers.json — SSOT (Docker)

| Где | Назначение |
|-----|------------|
| `dev/browsers.json` (selenoid-home) | **SSOT** полный Docker-каталог (`qaguru/*` images) |
| `selenoid/data/browsers.json` (этот репо, go:embed) | Копия SSOT для `cm selenoid configure` / `start` без `-j` |
| [`selenoid/config/browsers.json`](https://github.com/qa-guru/selenoid/blob/main/config/browsers.json) | Копия SSOT на hub |
| [`selenoid.autotests.cloud/deploy/browsers-production.json`](https://github.com/qa-guru/selenoid.autotests.cloud/blob/main/deploy/browsers-production.json) | Prod overlay (только по явному prod-запросу) |

Правки Docker-каталога — в `dev/browsers.json`, затем `dev/scripts/sync-cm-browsers.sh`.

Режим `--use-drivers` (локальные chromedriver/geckodriver без Docker) **не поддерживается из коробки** — нужен свой JSON-каталог и флаг `--drivers-info <url>`. Для стека qa-guru используйте Docker-режим (по умолчанию).

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

Текущий релиз cm: **v2.2.1** — [release notes](docs/RELEASE_v2.2.1.md) (patch: stack semver + ecosystem README).  
Предыдущий: [docs/RELEASE_v2.2.0.md](docs/RELEASE_v2.2.0.md). Docker: `qaguru/cm:v2.2.1`.

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

Флаги `--browsers`, `--browser-env`, `--drivers-info` — только для режима `--use-drivers`; `--drivers-info` обязателен (собственный каталог chromedriver/geckodriver).

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
