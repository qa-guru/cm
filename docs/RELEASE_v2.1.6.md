# Release v2.1.6 — qa-guru/cm

**Дата:** 7 июля 2026  
**Предыдущий:** [v2.1.5](RELEASE_v2.1.5.md)  
**GitHub:** https://github.com/qa-guru/cm/releases/tag/v2.1.6

Патч-релиз стека **v2.1.6**: синхронизация `browsers-qaguru.json` с hub v2.1.6 (Playwright 1.61.1, shmSize, hosts).

---

## Что нового

| Изменение | Описание |
|-----------|----------|
| **browsers-qaguru.json** | Playwright **1.61.1** (+ `-min`); `shmSize`, `hosts` как в hub `config/browsers.json` |

---

## Обновление

```bash
curl -sL https://github.com/qa-guru/cm/releases/download/v2.1.6/cm_linux_amd64 -o cm
chmod +x cm

./cm selenoid stop && ./cm selenoid-ui stop
./cm selenoid update -c /opt/selenoid
./cm selenoid-ui update -c /opt/selenoid
```

Docker: `docker pull qaguru/cm:v2.1.6`  
Prod: deploy workflow с `version=v2.1.6`.
