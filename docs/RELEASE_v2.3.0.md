# Release v2.3.0 — qa-guru/cm

**Дата:** 13 июля 2026  
**Предыдущий:** [v2.2.1](RELEASE_v2.2.1.md)  
**GitHub:** https://github.com/qa-guru/cm/releases/tag/v2.3.0 *(tag при cut)*  
**Stack cut:** hub + UI + cm → единый **v2.3.0**.  
**Pin ветка:** `selenoid2-1.55-engine29.6-go1.26-react18`

---

## Что нового

| Изменение | Описание |
|-----------|----------|
| **Stack semver** | Единый cut **v2.3.0** с hub/ui (для cm — docs/align; runtime без изменений) |
| **Docker API 1.55** | cm передаёт hub `DOCKER_API_VERSION=1.55` под Engine 29.x (было 1.45 / Engine 26.1.x) |
| **Go 1.26** | Toolchain align с hub/ui |
| **Browser SSOT** | Каталог без изменений: chrome **149.0**, firefox **151.0**, msedge **145.0**, Playwright **1.61.1** |

---

## Обновление

```bash
curl -sL https://github.com/qa-guru/cm/releases/download/v2.3.0/cm_linux_amd64 -o cm
chmod +x cm

./cm selenoid stop && ./cm selenoid-ui stop
./cm selenoid update -v v2.3.0
./cm selenoid-ui update -v v2.3.0
```

Prod deploy: `CM_VERSION=v2.3.0`, hub `v2.3.0`, UI `v2.3.0`.

Связанные: [selenoid v2.3.0](https://github.com/qa-guru/selenoid/releases/tag/v2.3.0), [selenoid-ui v2.3.0](https://github.com/qa-guru/selenoid-ui/releases/tag/v2.3.0), [selenoid-tests](https://github.com/qa-guru/selenoid-tests).

---

## Cut checklist (ручной)

1. Commit docs на pin-ветке → `git tag -a v2.3.0 -m "v2.3.0"` → push tags *(по команде)*.
2. GitHub Release (published) → `release.yml`: assets `dist/cm_*`; Docker `qaguru/cm:v2.3.0` + `latest-release`.
3. OUT: `warm-pool-orchestrator/`.
