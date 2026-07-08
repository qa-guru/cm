# Release note — cm embedded browsers v2.3.0

**Дата:** 8 июля 2026  
**Область:** `selenoid/data/browsers-qaguru.json` (go:embed)  
**cm binary tag:** без обязательного релиза бинарника — достаточно `git pull` / rebuild embed при следующей сборке cm

Синхронизация с hub [`config/browsers.json`](https://github.com/qa-guru/selenoid/blob/main/config/browsers.json):

| Ключ | Default | Образы |
|------|---------|--------|
| chrome | 149.0 | `qaguru/webdriver-chrome:149` (+ min, 148.*) |
| firefox | 151.0 | `qaguru/webdriver-firefox:151` (+ min, 150.*) |
| msedge | 145.0 | `qaguru/webdriver-msedge:145` (+ min, 144.*) |
| playwright-* | 1.61.1 | без изменений |

После обновления репо: `cm selenoid configure` / prod `browsers-production.json` уже содержат тот же WebDriver-каталог.
