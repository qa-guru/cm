# Release v2.1.5 — qa-guru/cm

**Дата:** 2 июля 2026  
**Предыдущий:** [v2.1.1](RELEASE_v2.1.1.md)  
**GitHub:** https://github.com/qa-guru/cm/releases/tag/v2.1.5

Патч-релиз стека **v2.1.5**: выравнивание с hub и UI, синхронизация `browsers-qaguru.json` с активным каталогом qaguru.

---

## Что нового

| Изменение | Описание |
|-----------|----------|
| **browsers-qaguru.json** | На момент релиза: `qaguru/webdriver-chrome:148` (+ `148-min`), Playwright `1.60.0` (+ `-min`); legacy twilio/firefox/edge убраны. Активный PW сейчас — **1.61.1** (см. [v2.1.6](RELEASE_v2.1.6.md)) |
| **Тесты** | Go unit: readable display names через `t.Run`; `docker_test` под новый каталог |
| **Документация** | README и release notes синхронизированы с hub |

---

## Обновление

```bash
curl -sL https://github.com/qa-guru/cm/releases/download/v2.1.5/cm_linux_amd64 -o cm
chmod +x cm

./cm selenoid stop && ./cm selenoid-ui stop
./cm selenoid update -c /opt/selenoid
./cm selenoid-ui update -c /opt/selenoid
```

Prod: deploy workflow с `version=v2.1.5`.
