# Release v2.1.7 — qa-guru/cm

**Дата:** 8 июля 2026  
**Предыдущий:** [v2.1.6](RELEASE_v2.1.6.md)  
**GitHub:** https://github.com/qa-guru/cm/releases/tag/v2.1.7

Релиз стека **Go 1.26.5 + moby/moby client**: миграция с `docker/docker`, blocking `govulncheck`, embedded catalog с WebDriver Firefox/Edge.

---

## Что нового

| Изменение | Описание |
|-----------|----------|
| **Go** | **1.26.5** (`toolchain`, CI `GOTOOLCHAIN`) |
| **Docker SDK** | `moby/moby/client` вместо legacy `docker/docker` |
| **govulncheck** | blocking в CI (`ci/test.sh`) |
| **browsers-qaguru** | WebDriver **firefox** и **msedge** в embedded catalog |
| **Тесты** | `TestConfigureDocker` / `TestLimitNoPull` — 8 browser entries |

---

## Обновление

```bash
curl -sL https://github.com/qa-guru/cm/releases/download/v2.1.7/cm_linux_amd64 -o cm
chmod +x cm

./cm selenoid stop && ./cm selenoid-ui stop
./cm selenoid update -v v2.1.8
./cm selenoid-ui update -v v2.1.7
```

Prod deploy: `CM_VERSION=v2.1.7`, hub `v2.1.8`, UI `v2.1.7`.

Связанные релизы: [selenoid v2.1.8](https://github.com/qa-guru/selenoid/releases/tag/v2.1.8), [selenoid-ui v2.1.7](https://github.com/qa-guru/selenoid-ui/releases/tag/v2.1.7).
