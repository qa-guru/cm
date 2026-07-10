# Release v2.2.1 — qa-guru/cm

**Дата:** 10 июля 2026  
**Предыдущий:** [v2.2.0](RELEASE_v2.2.0.md)  
**GitHub:** https://github.com/qa-guru/cm/releases/tag/v2.2.1 *(tag при cut)*  
**Stack cut:** hub + UI + cm → единый **v2.2.1** (patch после v2.2.0).

---

## Что нового

| Изменение | Описание |
|-----------|----------|
| **Stack semver** | Единый binary cut **v2.2.1** с hub/ui (docs-only для cm: код = v2.2.0 tag) |
| **README** | Единый блок «Экосистема qa-guru Selenoid» + [selenoid-tests](https://github.com/qa-guru/selenoid-tests) + [Docker Hub qaguru](https://hub.docker.com/u/qaguru) |
| **Runtime** | Без изменений embed browsers / install logic относительно v2.2.0 |

---

## Обновление

```bash
curl -sL https://github.com/qa-guru/cm/releases/download/v2.2.1/cm_linux_amd64 -o cm
chmod +x cm

./cm selenoid stop && ./cm selenoid-ui stop
./cm selenoid update -v v2.2.1
./cm selenoid-ui update -v v2.2.1
```

Prod deploy: `CM_VERSION=v2.2.1`, hub `v2.2.1`, UI `v2.2.1`.

Связанные: [selenoid v2.2.1](https://github.com/qa-guru/selenoid/releases/tag/v2.2.1), [selenoid-ui v2.2.1](https://github.com/qa-guru/selenoid-ui/releases/tag/v2.2.1), [selenoid-tests](https://github.com/qa-guru/selenoid-tests).

---

## Cut checklist (ручной)

1. Commit docs на `main` → `git tag -a v2.2.1 -m "v2.2.1"` → push tags *(по команде)*.
2. Release assets `dist/cm_*`; Docker `qaguru/cm:v2.2.1`.
3. OUT: `warm-pool-orchestrator/`.
