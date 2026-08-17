# Release v3.0.3 — qa-guru/cm

**Дата:** 17 августа 2026  
**Предыдущий:** [v3.0.2](RELEASE_v3.0.2.md)  
**GitHub:** https://github.com/qa-guru/cm/releases/tag/v3.0.3  
**Stack cut:** hub → **v3.0.12** (`-warm-pool-url`); cm → **v3.0.3**; UI → **v3.0.36**.  
Prod pin cm на [selenoid.qa.guru](https://selenoid.qa.guru) остаётся **v3.0.2** до отдельного deploy-чата.

---

## Что нового

| Изменение | Описание |
|-----------|----------|
| **`--pool` / `--warm-pool`** | `cm selenoid start --pool` поднимает sidecar [selenoid-pool](https://github.com/qa-guru/selenoid-pool): warm 4/4 + orchestrator `:9090`. Хабу передаётся `-warm-pool-url http://127.0.0.1:9090` (host network). Без флага — cold hub, как раньше. |
| **`--hot-pool`** | Тот же оркестратор + compose profile `hot` (2/2 `-min`). Не третий бинарь. |
| **Образ** | Compose тянет опубликованный `qaguru/selenoid-pool:min` (без `build:`). DNS alias `selenoid-warm-pool`. |
| **Android catalog** | Embedded `browsers.json`: android **5.1 + 10–16**. |
| **Go** | Toolchain **1.26.6** (govulncheck stdlib vs 1.26.5). |

`cm selenoid stop` гасит sidecar, если его поднял этот cm (маркер в `-c …/warm-pool`).

Hub flag alias `-pool-url` живёт на [qa-guru/selenoid](https://github.com/qa-guru/selenoid) `main`; этот cut хаба не делает (prod auto-deploy). cm по-прежнему передаёт канон `-warm-pool-url`.

Не в этом релизе: вливать lease в бинарь hub · GGR · Jenkins · prod deploy.

---

## Обновление

```bash
curl -sL https://github.com/qa-guru/cm/releases/download/v3.0.3/cm_linux_amd64 -o cm
chmod +x cm

./cm selenoid stop && ./cm selenoid-ui stop
./cm selenoid update -v v3.0.3
./cm selenoid-ui update -v v3.0.36
```

Warm/hot sidecar на чистом Docker-хосте:

```bash
./cm selenoid start --pool          # alias: --warm-pool
./cm selenoid start --hot-pool      # + hot 2/2
```

Prod deploy (отдельный чат): `CM_VERSION=v3.0.3`. Hub prod pin остаётся **v3.0.12**.

Связанные: [selenoid v3.0.12](https://github.com/qa-guru/selenoid/releases/tag/v3.0.12), [selenoid-ui v3.0.36](https://github.com/qa-guru/selenoid-ui/releases/tag/v3.0.36), [selenoid-pool](https://github.com/qa-guru/selenoid-pool).

---

## Cut checklist

1. `main` green (`ci/test.sh` + govulncheck на Go 1.26.6).
2. `git tag -a v3.0.3 -m "v3.0.3"` → push tag → GitHub Release (published) → binaries + `qaguru/cm:v3.0.3`.
3. Prod cm pin → отдельный deploy-чат (этот cut его не двигает).
