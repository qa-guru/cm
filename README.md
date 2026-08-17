# Configuration Manager (qa-guru fork)

[![Configuration Manager](https://qa-guru.github.io/selenoid-tests/readme/badge-cm.svg)](https://qa-guru.github.io/selenoid-tests/reports/latest/dashboard/)

[![CM stats](https://qa-guru.github.io/selenoid-tests/readme/stats-cm.svg)](https://qa-guru.github.io/selenoid-tests/reports/latest/dashboard/)

[![CM metrics](https://qa-guru.github.io/selenoid-tests/readme/metrics-panel-cm.svg)](https://qa-guru.github.io/selenoid-tests/reports/latest/dashboard/)

<a href="https://qa-guru.github.io/selenoid-tests/reports/latest/dashboard/">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://qa-guru.github.io/selenoid-tests/readme/dashboard-preview-dark.png">
    <img
      src="https://qa-guru.github.io/selenoid-tests/readme/dashboard-preview.png"
      alt="Allure 3 dashboard — pyramid, stability, success distribution"
      width="800"
    />
  </picture>
</a>
Dashboard PNG updates after each orchestrator run on `main` (Playwright screenshot of Allure 3 dashboard).

| Link | Description |
|------|-------------|
| [Dashboard](https://qa-guru.github.io/selenoid-tests/reports/latest/dashboard/) | Full pyramid — filter epic **cm** in awesome |
| [Awesome](https://qa-guru.github.io/selenoid-tests/reports/latest/awesome/) | CM integration test details |
| [selenoid-tests](https://github.com/qa-guru/selenoid-tests) | Orchestrator + merged Allure |

<!-- stack-branches-note:start -->
> ## Стабильные билды
>
> **Prod cm/hub:** релизы с **`main`** → **v3.0.0+** ([selenoid.qa.guru](https://selenoid.qa.guru)). UI — отдельная v3.x линия на `selenoid-ui` `main`.
>
> Pin-ветки **2.x** (`selenoid2-…-react16` / `react18`) — **заморожены** (rollback reference only).
>
> | Ветка | Semver | Назначение |
> |-------|--------|------------|
> | **`main`** | **v3.0.0+** | Активная prod-линия cm |
> | `selenoid2-1.55-…-react18` | v2.3.0 | frozen |
> | `selenoid2-1.45-…-react16` | v2.2.1 | frozen |
>
> Monorepo SSOT: [`projects/selenoid-home/README.md`](https://github.com/qa-guru/zero-design-system/blob/master/projects/selenoid-home/README.md).
<!-- stack-branches-note:end -->


Установщик стека **qa-guru/selenoid** — форк [aerokube/cm](https://github.com/qa-guru/cm). Одна команда на чистом сервере с Docker поднимает hub, UI, `browsers.json` и browser-образы.

[![Build Status](https://github.com/qa-guru/cm/workflows/build/badge.svg)](https://github.com/qa-guru/cm/actions?query=workflow%3Abuild)
[![Go Report Card](https://goreportcard.com/badge/github.com/qa-guru/cm)](https://goreportcard.com/report/github.com/qa-guru/cm)
[![Release](https://img.shields.io/github/release/qa-guru/cm.svg)](https://github.com/qa-guru/cm/releases/latest)
[![Docker Pulls](https://img.shields.io/docker/pulls/qaguru/cm.svg)](https://hub.docker.com/r/qaguru/cm)

| | |
|---|---|
| **GitHub** | [qa-guru/cm](https://github.com/qa-guru/cm) |
| **Docker Hub** | [`qaguru/cm`](https://hub.docker.com/r/qaguru/cm) |
| **Текущий релиз** | **v3.0.1** — [docs/RELEASE_v3.0.1.md](docs/RELEASE_v3.0.1.md) · `qaguru/cm:v3.0.1` |

## Роль в экосистеме

`cm` не содержит логику hub — он **оркестрирует установку** компонентов из других репозиториев:

```
cm selenoid start
    ├── скачивает qaguru/selenoid (hub) из GitHub Releases
    ├── записывает browsers.json (Chrome/Firefox + Playwright)
    ├── docker pull qaguru/webdriver-chrome + qaguru/playwright-*
    └── запускает контейнер hub

cm selenoid start --pool          # alias: --warm-pool
    ├── то же, что start
    ├── docker compose: warm 4/4 + orchestrator :9090
    └── hub: -warm-pool-url http://127.0.0.1:9090  (host network)

cm selenoid start --hot-pool
    └── то же + compose profile `hot` (2/2 min slots, тот же оркестратор)

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
| selenoid-pool | [github.com/qa-guru/selenoid-pool](https://github.com/qa-guru/selenoid-pool) | Warm/hot sidecar |
| Docker Hub | [hub.docker.com/u/qaguru](https://hub.docker.com/u/qaguru) | Образы `qaguru/*` |

## browsers.json — SSOT (Docker)

| Где | Назначение |
|-----|------------|
| `dev/browsers.json` (selenoid-home) | **SSOT** полный Docker-каталог (`qaguru/*` images) |
| `selenoid/data/browsers.json` (этот репо, go:embed) | Копия SSOT для `cm selenoid configure` / `start` без `-j` |
| [`selenoid/config/browsers.json`](https://github.com/qa-guru/selenoid/blob/main/config/browsers.json) | Копия SSOT на hub |
| [`selenoid.qa.guru/deploy/browsers-production.json`](https://github.com/qa-guru/selenoid.qa.guru/blob/main/deploy/browsers-production.json) | Prod overlay (только по явному prod-запросу) |

Правки Docker-каталога — в `dev/browsers.json`, затем `dev/scripts/sync-cm-browsers.sh`.

Режим `--use-drivers` (локальные chromedriver/geckodriver без Docker) **не поддерживается из коробки** — нужен свой JSON-каталог и флаг `--drivers-info <url>`. Для стека qa-guru используйте Docker-режим (по умолчанию).

## Установка

**Предварительные условия:** Docker Engine **29.x** (API **1.55**, moby/moby client), доступ пользователя к `docker`, опубликованные релизы [qa-guru/selenoid](https://github.com/qa-guru/selenoid/releases) и [qa-guru/selenoid-ui](https://github.com/qa-guru/selenoid-ui/releases). Для сборки cm — Go **1.26.5**.

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

Warm 4/4 sidecar (orchestrator `:9090`, hub `-warm-pool-url http://127.0.0.1:9090`):

```bash
./cm selenoid start --pool
# alias: --warm-pool
# hot 2/2, same orchestrator (not a third binary):
./cm selenoid start --hot-pool
```

Без флага — как раньше: только cold hub. `cm selenoid stop` гасит sidecar, если его поднял этот cm (маркер в `-c …/warm-pool`).

Локально без прав на `/opt` — свой каталог:

```bash
./cm selenoid start -c "$HOME/selenoid"
./cm selenoid-ui start -c "$HOME/selenoid"
```

`cm selenoid start` передаёт hub `DOCKER_API_VERSION=1.55`.

Обновление:

```bash
./cm selenoid update
./cm selenoid-ui update
./cm selenoid start -f   # принудительно перекачать образы и бинарники
```

Текущий релиз cm: **v3.0.1** — [release notes](docs/RELEASE_v3.0.1.md).  
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
| `--pool` | Sidecar [selenoid-pool](https://github.com/qa-guru/selenoid-pool): warm 4/4 + orchestrator; hub `-warm-pool-url http://127.0.0.1:9090` |
| `--warm-pool` | Alias for `--pool` |
| `--hot-pool` | Compose profile `hot` (2/2 `-min`); тот же оркестратор, implies `--pool` |

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
  warm-pool/          # только при --pool / --warm-pool / --hot-pool
    docker-compose.yml
    config.yaml
```
