# Release v2.1.1 — qa-guru/cm

**Дата:** 28 июня 2026  
**Предыдущий:** [v2.1.0](https://github.com/qa-guru/cm/releases/tag/v2.1.0)  
**GitHub:** https://github.com/qa-guru/cm/releases/tag/v2.1.1

Патч-релиз: дефолтный каталог **`/opt/selenoid`**, зачистка legacy aerokube, актуальные тесты и drivers-mode конфиг под Chrome 148.

---

## Что нового

| Изменение | Описание |
|-----------|----------|
| **Default config dir** | `GetSelenoidConfigDir()` → `/opt/selenoid` (было `~/.aerokube/selenoid`) |
| **Legacy cleanup** | Удалены мёртвый `createConfig()` и флаги `--last-versions`, `--tmpfs`, `--shm-size`, `--vnc` |
| **drivers-mode** | `browsers.json`: Chrome 148 (CfT), geckodriver 0.37, Edge 145 |
| **Тесты** | Под 8 браузеров, `testdata/`, `TestAllUrlsAreValid` снова проверяет URL |
| **Документация** | README: prod `browsers.json`, флаги; удалены upstream asciidoc docs |

---

## Обновление

```bash
curl -sL https://github.com/qa-guru/cm/releases/download/v2.1.1/cm_linux_amd64 -o cm
chmod +x cm

./cm selenoid stop && ./cm selenoid-ui stop
./cm selenoid update -c /opt/selenoid
./cm selenoid-ui update -c /opt/selenoid
```

Prod: [deploy workflow](https://github.com/qa-guru/selenoid.autotests.cloud/actions/workflows/deploy.yml) с `version=v2.1.1`.
