# Release v2.0.1 — qa-guru/cm

**Дата:** 25 июня 2026  
**Предыдущий:** [v2.0.0](RELEASE_v2.0.0.md)  
**GitHub:** https://github.com/qa-guru/cm/releases/tag/v2.0.1

Патч-релиз: **Go 1.23** и фиксированный **Docker API 1.45** для контейнера Selenoid.

---

## Что нового

| Изменение | Было | Стало |
|-----------|------|-------|
| **Go** | 1.22 | **1.23.x** |
| **DOCKER_API_VERSION** в hub | версия клиента daemon | **1.45** (стабильно на Engine 26.1 и 27.x) |

Рекомендуемый Docker Engine на сервере: **26.1.x**.

---

## Обновление

```bash
curl -sL https://github.com/qa-guru/cm/releases/download/v2.0.1/cm_linux_amd64 -o cm
chmod +x cm

./cm selenoid stop && ./cm selenoid-ui stop
./cm selenoid start -v v2.0.1 -f
./cm selenoid-ui start -v v2.0.1 -f
```

Связанные релизы: [selenoid v2.0.1](https://github.com/qa-guru/selenoid/releases/tag/v2.0.1), [selenoid-ui v2.0.1](https://github.com/qa-guru/selenoid-ui/releases/tag/v2.0.1).
