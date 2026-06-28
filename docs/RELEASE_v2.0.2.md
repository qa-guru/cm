# Release v2.0.2 — qa-guru/cm

**Дата:** 25 июня 2026  
**Предыдущий:** [v2.0.1](RELEASE_v2.0.1.md)  
**GitHub:** https://github.com/qa-guru/cm/releases/tag/v2.0.2

Патч-релиз: деплой стека с hub **v2.0.2**, исправления `cm selenoid start` / `cm selenoid-ui start` в Docker-образах Aerokube.

---

## Что нового

| Изменение | Описание |
|-----------|----------|
| **Hub v2.0.2 по умолчанию** | `cm selenoid start` тянет hub v2.0.2 (UI остаётся на latest / v2.0.1) |
| **Docker ENTRYPOINT** | Не дублировать бинарник в `cmd` — флаги hub/UI снова применяются в образах aerokube |
| **Production deploy** | Hub как нативный процесс, UI на host network, `/opt/selenoid` |

Связанный hub: [selenoid v2.0.2](https://github.com/qa-guru/selenoid/releases/tag/v2.0.2) — DELETE Playwright session, Docker inspect на хосте.

---

## Обновление

```bash
curl -sL https://github.com/qa-guru/cm/releases/latest/download/cm_linux_amd64 -o cm
chmod +x cm

./cm selenoid stop && ./cm selenoid-ui stop
./cm selenoid update -c /opt/selenoid
./cm selenoid-ui update -c /opt/selenoid
```
