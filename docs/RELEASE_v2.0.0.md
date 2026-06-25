# Release v2.0.0 — qa-guru/cm

**Дата:** 25 июня 2026  
**База:** форк [aerokube/cm](https://github.com/aerokube/cm) v1.8.8  
**GitHub:** https://github.com/qa-guru/cm/releases/tag/v2.0.0

Configuration Manager для установки стека **qa-guru/selenoid** v2.0.0 с Playwright одной командой.

---

## Что делает

```bash
cm selenoid start -v v2.0.0
cm selenoid-ui start -v v2.0.0
```

1. Пишет встроенный `browsers.json` (twilio/selenoid + qaguru/playwright) в `~/.aerokube/selenoid/`
2. Тянет browser-образы из конфига
3. Скачивает бинарники hub и UI из GitHub Releases `qa-guru/selenoid` и `qa-guru/selenoid-ui`
4. Запускает Docker-контейнеры с примонтированными бинарниками qa-guru

Обёрточный образ UI: `qaguru/selenoid-ui:latest-release`. Hub в production — нативный бинарник qa-guru/selenoid; при `cm selenoid start` в Docker — `qaguru/selenoid:latest-release`.

---

## Установка

```bash
curl -sL https://github.com/qa-guru/cm/releases/download/v2.0.0/cm_linux_amd64 -o cm
chmod +x cm
./cm selenoid start -v v2.0.0
./cm selenoid-ui start -v v2.0.0
```

Ассеты: `cm_linux_amd64`, `cm_darwin_arm64`, `cm_windows_amd64.exe`, …

---

## Отличия от aerokube/cm

| aerokube/cm | qa-guru/cm v2.0.0 |
|-------------|-------------------|
| `selenoid/chrome`, `selenoid/firefox` | twilio/selenoid + qaguru/playwright |
| stock aerokube/selenoid hub | бинарник из qa-guru/selenoid Release |
| без Playwright | встроенный browsers.json с Playwright |
| GitHub owner aerokube | **qa-guru** |

---

## Флаги

| Флаг | Описание |
|------|----------|
| `-v v2.0.0` | Тег Release для скачивания selenoid / selenoid-ui |
| `-j path` | Свой `browsers.json` |
| `--selenoid-binary` | Локальный бинарник hub |
| `--selenoid-ui-binary` | Локальный бинарник UI |
| `-f` | Перекачать образы и бинарники |

---

## Связанные релизы

- [qa-guru/selenoid v2.0.0](https://github.com/qa-guru/selenoid/releases/tag/v2.0.0)
- [qa-guru/selenoid-ui v2.0.0](https://github.com/qa-guru/selenoid-ui/releases/tag/v2.0.0)
- [deploy/RELEASE_v2.0.0.md](../deploy/RELEASE_v2.0.0.md) — общий чеклист стека
