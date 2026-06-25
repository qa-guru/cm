# qa-guru Selenoid Stack — Release v2.0.1

| Компонент | Release | Документация |
|-----------|---------|--------------|
| **cm** | [v2.0.1](https://github.com/qa-guru/cm/releases/tag/v2.0.1) | [cm-src/docs/RELEASE_v2.0.1.md](../docs/RELEASE_v2.0.1.md) |
| Selenoid hub | [v2.0.1](https://github.com/qa-guru/selenoid/releases/tag/v2.0.1) | [selenoid v2.0.1](https://github.com/qa-guru/selenoid/releases/tag/v2.0.1) |
| Selenoid UI | [v2.0.1](https://github.com/qa-guru/selenoid-ui/releases/tag/v2.0.1) | [selenoid-ui v2.0.1](https://github.com/qa-guru/selenoid-ui/releases/tag/v2.0.1) |

Предыдущий стек: [v2.0.0](RELEASE_v2.0.0.md)

---

## Что нового в v2.0.1

Патч-релиз toolchain — без изменений протоколов Playwright / WebDriver.

| | v2.0.0 | v2.0.1 |
|---|--------|--------|
| **Go** | 1.22 | **1.23** |
| **Docker API (hub)** | не зафиксирован | **1.45** |
| **Docker Engine (сервер)** | — | **26.1.x** рекомендуется |

---

## Деплой на selenoid.autotests.cloud

Автоматически: GitHub Actions **deploy** в [qa-guru/cm](https://github.com/qa-guru/cm/actions/workflows/deploy.yml) при публикации релиза cm.

Вручную на сервере:

```bash
curl -sL https://github.com/qa-guru/cm/releases/latest/download/cm_linux_amd64 -o ~/cm
chmod +x ~/cm
./deploy/deploy.sh   # из клона cm или curl raw deploy.sh
```

Подробнее: [deploy/README.md](README.md)

---

## Обновление на сервере

```bash
./cm selenoid stop && ./cm selenoid-ui stop
./cm selenoid update -c /opt/selenoid
./cm selenoid-ui update -c /opt/selenoid
```

Pin Docker Engine 26.1.x (Ubuntu): [docker-engine-pin-ubuntu.sh](https://github.com/qa-guru/selenoid/blob/master/scripts/docker-engine-pin-ubuntu.sh)
