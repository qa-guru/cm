# Release v2.1.0 — qa-guru/cm

**Дата:** 28 июня 2026  
**Предыдущий:** [v2.0.9](https://github.com/qa-guru/cm/releases/tag/v2.0.9)  
**GitHub:** https://github.com/qa-guru/cm/releases/tag/v2.1.0

Инфраструктурный релиз: **единая версия стека v2.1.0**, prod-деплой вынесен в [qa-guru/selenoid.autotests.cloud](https://github.com/qa-guru/selenoid.autotests.cloud). Код cm без изменений относительно v2.0.9.

---

## Что нового

| Изменение | Описание |
|-----------|----------|
| **Prod deploy** | Скрипты, nginx, CI — в [selenoid.autotests.cloud](https://github.com/qa-guru/selenoid.autotests.cloud) |
| **cm** | Только сборка, релиз и Docker-образ установщика |
| **Стек v2.1.0** | cm, hub, UI, deploy — одна версия **v2.1.0** |

---

## Обновление

```bash
curl -sL https://github.com/qa-guru/cm/releases/download/v2.1.0/cm_linux_amd64 -o cm
chmod +x cm
```

Prod: [deploy workflow](https://github.com/qa-guru/selenoid.autotests.cloud/actions/workflows/deploy.yml) с `version=v2.1.0`, `ref=v2.1.0`.
