# Release v3.0.0 — qa-guru/cm

**Дата:** 23 июля 2026  
**Предыдущий:** [v2.3.0](RELEASE_v2.3.0.md)  
**GitHub:** https://github.com/qa-guru/cm/releases/tag/v3.0.0  
**Stack cut:** hub + cm → **v3.0.0**; UI — **v3.0.x** на `selenoid-ui` `main`.

---

## Что нового

| Изменение | Описание |
|-----------|----------|
| **Selenoid 3 line** | cm semver **v3.0.0** с `main`; 2.x pin-ветки заморожены |
| **Browser SSOT** | Embedded `browsers.json`: android **16.0** + полный WD/PW каталог |
| **Docker API 1.55** | cm передаёт hub `DOCKER_API_VERSION=1.55` под Engine 29.x |
| **Go 1.26** | Toolchain align с hub |

---

## Обновление

```bash
curl -sL https://github.com/qa-guru/cm/releases/download/v3.0.0/cm_linux_amd64 -o cm
chmod +x cm

./cm selenoid stop && ./cm selenoid-ui stop
./cm selenoid update -v v3.0.0
./cm selenoid-ui update -v v3.0.7
```

Prod deploy: `CM_VERSION=v3.0.0`, hub `SELENOID_VERSION=v3.0.0`, UI `SELENOID_UI_VERSION=v3.0.7+`.

Связанные: [selenoid v3.0.0](https://github.com/qa-guru/selenoid/releases/tag/v3.0.0), [selenoid-ui](https://github.com/qa-guru/selenoid-ui/releases).

---

## Cut checklist

1. `main` green (CI test/build).
2. `git tag -a v3.0.0 -m "v3.0.0"` → push tag → GitHub Release (published).
3. Prod deploy pins → v3.0.0.
