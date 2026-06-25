# qa-guru Selenoid Stack — Release v2.0.0

Общие release notes для стека **selenoid.autotests.cloud** и локальной разработки.

| Компонент | Release | Документация |
|-----------|---------|--------------|
| **cm** | [v2.0.0](https://github.com/qa-guru/cm/releases/tag/v2.0.0) | установщик стека |
| Selenoid hub | [v2.0.0](https://github.com/qa-guru/selenoid/releases/tag/v2.0.0) | [selenoid v2.0.0](https://github.com/qa-guru/selenoid/releases/tag/v2.0.0) |
| Selenoid UI | [v2.0.0](https://github.com/qa-guru/selenoid-ui/releases/tag/v2.0.0) | [selenoid-ui v2.0.0](https://github.com/qa-guru/selenoid-ui/releases/tag/v2.0.0) |
| Playwright image | [qaguru/playwright](https://hub.docker.com/r/qaguru/playwright) | [playwright-image](https://github.com/qa-guru/playwright-image) |
| Примеры тестов | [selenoid_selenium_playwright_tests](https://github.com/qa-guru/selenoid_selenium_playwright_tests) | коллекция примеров |

---

## Endpoints

**Basic auth:** `user1` / `1234` (для WebDriver и Playwright)

| Назначение | URL |
|------------|-----|
| Selenium | `http://user1:1234@selenoid.autotests.cloud/wd/hub` |
| Playwright | `wss://user1:1234@selenoid.autotests.cloud/playwright/chromium/1.61.1` |
| UI | `http://selenoid.autotests.cloud:8080/` |
| Status | `https://selenoid.autotests.cloud/status` |
| Video | `https://selenoid.autotests.cloud/video/` |

---

## Переменные для тестов

```bash
export PW_TEST_CONNECT_WS_ENDPOINT=wss://user1:1234@selenoid.autotests.cloud/playwright/chromium/1.61.1
export SELENOID_URL=http://user1:1234@selenoid.autotests.cloud/wd/hub
```

Деплой: [README.md](README.md)
