# Release v3.0.1 — qa-guru/cm

**Дата:** 26 июля 2026  
**Предыдущий:** [v3.0.0](RELEASE_v3.0.0.md)  
**GitHub:** https://github.com/qa-guru/cm/releases/tag/v3.0.1  
**Stack cut:** hub + cm → **v3.0.1**; UI — **v3.0.8** на `selenoid-ui` `main`.

---

## Что нового

| Изменение | Описание |
|-----------|----------|
| **govulncheck** | `golang.org/x/text` → **v0.39.0** (GO-2026-5970) |
| **Android embed tests** | Configure tests ожидают android в embedded `browsers.json` (9 top-level browsers) |

Нет изменений CLI UX относительно v3.0.0.

---

## Обновление

```bash
curl -sL https://github.com/qa-guru/cm/releases/download/v3.0.1/cm_linux_amd64 -o cm
chmod +x cm

./cm selenoid stop && ./cm selenoid-ui stop
./cm selenoid update -v v3.0.1
./cm selenoid-ui update -v v3.0.8
```

Prod deploy (отдельный чат): `CM_VERSION=v3.0.1`, hub `SELENOID_VERSION=v3.0.1`, UI `SELENOID_UI_VERSION=v3.0.8`.

Связанные: [selenoid v3.0.1](https://github.com/qa-guru/selenoid/releases/tag/v3.0.1), [selenoid-ui v3.0.8](https://github.com/qa-guru/selenoid-ui/releases/tag/v3.0.8).

---

## Cut checklist

1. `main` green (CI test/build).
2. `git tag -a v3.0.1 -m "v3.0.1"` → push tag → GitHub Release (published).
3. Prod deploy pins → v3.0.1 / UI v3.0.8.
