![GoLint](https://github.com/kotoproger/exchange/actions/workflows/linter.yml/badge.svg?branch=master) ![Tests](https://github.com/kotoproger/exchange/actions/workflows/tests.yml/badge.svg?branch=master) ![CodeQl](https://github.com/kotoproger/exchange/actions/workflows/github-code-scanning/codeql/badge.svg?branch=master) ![Coverage](https://github.com/kotoproger/exchange/actions/workflows/coverage.yml/badge.svg?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/kotoproger/exchange)](https://goreportcard.com/report/github.com/kotoproger/exchange) ![Snyk](https://snyk.io/test/github/kotoproger/exchange/badge.svg)

Запустить приложение
```
make docker-app 
```

## требования к работе 
### Обязательные требования для каждого проекта
* Наличие юнит-тестов на ключевые алгоритмы (core-логику) сервиса. **покрытие тестами > 80%**
* Наличие валидных Dockerfile для сервиса. **в наличии**
* Ветка master успешно проходит пайплайн в CI-CD системе **проходит**
(на ваш вкус, GitHub Actions, Circle CI, Travis CI, Jenkins, GitLab CI и пр.).
**Пайплайн должен в себе содержать**:
    - запуск последней версии `golangci-lint` на весь проект с **некоторые линтеры отключены для тестов, так же некоторые линтеры из требуемых ругаются что они отключены, пришлось их убрать**
    [конфигом, представленным в данном репозитории](./.golangci.yml);
    - запуск юнит тестов командой вида `go test -race -count 100`; **go test -v -count=100 -race -timeout=1m ./...**
    - сборку бинаря сервиса для версии Go не ниже 1.14. **1.22**
