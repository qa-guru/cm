# Release v3.0.5 — qa-guru/cm

**Дата:** 4 сентября 2026  
**Предыдущий:** [v3.0.4](RELEASE_v3.0.4.md)  
**GitHub:** https://github.com/qa-guru/cm/releases/tag/v3.0.5  
**Stack cut:** hub → **v3.0.15**; cm → **v3.0.5**; UI → **v3.0.54**.  
Prod pin cm на [selenoid.qa.guru](https://selenoid.qa.guru) остаётся **v3.0.4** до отдельного deploy-чата.

## Что нового

| Изменение | Описание |
|-----------|----------|
| **Docker API 1.55** | `cm selenoid start` передаёт hub `DOCKER_API_VERSION=1.55` (литерал в коде был **1.45**). README и v3.0.0 это уже обещали. |

```bash
docker pull qaguru/cm:v3.0.5
```

Prod deploy (отдельный чат): `CM_VERSION=v3.0.5`. Hub prod pin остаётся **v3.0.15**.

Связанные: [selenoid v3.0.15](https://github.com/qa-guru/selenoid/releases/tag/v3.0.15), [selenoid-ui v3.0.54](https://github.com/qa-guru/selenoid-ui/releases/tag/v3.0.54).

---

## Cut checklist

1. `main` green (`ci/test.sh` + govulncheck на Go 1.27).
2. `git tag -a v3.0.5 -m "v3.0.5"` → push tag → GitHub Release (published) → binaries + `qaguru/cm:v3.0.5`.
3. Prod cm pin → отдельный deploy-чат (этот cut его не двигает).
