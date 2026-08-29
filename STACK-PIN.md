# Stack pin: main / v3.0.x (Selenoid 3)

**Репозиторий:** Configuration manager (qa-guru/cm)

Этот файл на **`main`** описывает живой cm toolchain. Pin-ветки 2.x — отдельные frozen `STACK-PIN.md`.

| Поле | Значение |
|------|----------|
| Линия | Selenoid 3 |
| Stack semver | cm cut **v3.0.3** (latest tag; prod pin — deploy-чат) |
| Docker API | TBD (paired с hub) |
| Docker Engine | TBD (prod: Debian 12 · Docker 29.6) |
| Go | 1.27.0+ |
| Go (примечание) | Факт `go.mod` + `toolchain go1.27.0` |
| Prod | [selenoid.qa.guru](https://selenoid.qa.guru) |
| Git anchor | `main` |
| Docker image | `qaguru/cm:v3.0.x` |
| browsers.json | embed `selenoid/data/browsers.json` = копия SSOT из `dev/browsers.json` |

## Selenoid 2 maintenance pin (не путать)

Maintenance **v2.3.0** — только pin-ветка
[`selenoid2-1.55-engine29.6-go1.26-react18`](https://github.com/qa-guru/cm/tree/selenoid2-1.55-engine29.6-go1.26-react18).
Rollback **v2.2.1** —
[`selenoid2-1.45-engine26.1-go1.26-react16`](https://github.com/qa-guru/cm/tree/selenoid2-1.45-engine26.1-go1.26-react16).

См. также: [`projects/selenoid-home/README.md`](https://github.com/qa-guru/zero-design-system/blob/master/projects/selenoid-home/README.md) (monorepo SSOT).
