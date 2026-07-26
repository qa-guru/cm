# Release v2.2.0 — qa-guru/cm

**Дата:** 10 июля 2026  
**Предыдущий:** [v2.1.7](RELEASE_v2.1.7.md)  
**GitHub:** https://github.com/qa-guru/cm/releases/tag/v2.2.0 *(tag при cut)*  
**Stack cut:** hub + UI + cm → единый **v2.2.0**.

---

## Что нового

| Изменение | Описание |
|-----------|----------|
| **embed browsers** | `selenoid/data/browsers.json` (rename с `browsers-qaguru.json`) = catalog **v2.2.0** |
| **video sidecar** | CM-managed Selenoid → `qaguru/video-recorder:latest` |
| **Docker UI port** | fix listen port selenoid-ui внутри CM Docker container |
| **env merge** | container override env merge с host environ (не replace) |
| **tests / CI** | `browsers_test.go`; cmd unit; smoke `workflow_call` |
| **Go** | **1.26.5**; Docker Engine **26.1.x** / API **1.45** |

| Ключ | Default | Образы |
|------|---------|--------|
| chrome | 149.0 | `qaguru/webdriver-chrome:149` (+ min, 148.*) |
| firefox | 151.0 | `qaguru/webdriver-firefox:151` (+ min, 150.*) |
| msedge | 145.0 | `qaguru/webdriver-msedge:145` (+ min, 144.*) |
| playwright-* | 1.61.1 | без изменений |

---

## Обновление

```bash
curl -sL https://github.com/qa-guru/cm/releases/download/v2.2.0/cm_linux_amd64 -o cm
chmod +x cm

./cm selenoid stop && ./cm selenoid-ui stop
./cm selenoid update -v v2.2.0
./cm selenoid-ui update -v v2.2.0
```

Prod deploy: `CM_VERSION=v2.2.0`, hub `v2.2.0`, UI `v2.2.0`.

Связанные: [selenoid v2.2.0](https://github.com/qa-guru/selenoid/releases/tag/v2.2.0), [selenoid-ui v2.2.0](https://github.com/qa-guru/selenoid-ui/releases/tag/v2.2.0).

---

## Cut checklist (ручной)

1. Commit на `main` (embed + docs).
2. `git tag -a v2.2.0 -m "v2.2.0"` → push tags *(по команде)*.
3. Release assets `dist/cm_*`; Docker `qaguru/cm:v2.2.0` (если publish workflow).
4. OUT: `selenoid-warm-pool/`.
