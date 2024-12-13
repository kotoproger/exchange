![GoLint](https://github.com/kotoproger/exchange/actions/workflows/linter.yml/badge.svg?branch=master) ![Tests](https://github.com/kotoproger/exchange/actions/workflows/tests.yml/badge.svg?branch=master) ![CodeQl](https://github.com/kotoproger/exchange/actions/workflows/github-code-scanning/codeql/badge.svg?branch=master) ![Coverage](https://github.com/kotoproger/exchange/actions/workflows/coverage.yml/badge.svg?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/kotoproger/exchange)](https://goreportcard.com/report/github.com/kotoproger/exchange) ![Snyk](https://snyk.io/test/github/kotoproger/exchange/badge.svg)

Запустить приложение
```
make docker-app 
```
# Проектная работа

## Конвертер валют

### Описание: 
В этом проекте студентам предлагается написать программу командной строки, которая конвертирует сумму из одной валюты в другую. Для получения курсов валют программа может использовать стороннее API, например, API открытых данных Центрального банка или сторонних сервисов.

### Требования:
- [x] Программа должна принимать ввод суммы и валюты, из которой нужно конвертировать, а также валюты, в которую нужно сконвертировать. 
- [x] Программа должна использовать сторонний API для получения курсов валют. **так же логика работы расчитана на использование нескольких источников**
- [x] Реализована обработка ошибок, связанных с соединением с API или некорректными входными данными. 
- [x] Результаты конвертации должны быть точными и учитывать валютные курсы. 

### Развертывание
- [x] Развертывание сервиса должно осуществляться с использованием docker compose в директории с проектом.

### Тестирование
- [x] Написаны юнит-тесты на core логику приложения. Плюсом будут тесты на транспортном уровне и на уровне хранения. **>80% включая тесты на клиента поставщика курсов**

## Обязательные требования для каждого проекта
- [x] Наличие юнит-тестов на ключевые алгоритмы (core-логику) сервиса. **> 70%**
- [x] Наличие валидных Dockerfile для сервиса. 
- [x] Ветка master успешно проходит пайплайн в CI-CD системе (на ваш вкус, GitHub Actions, Circle CI, Travis CI, Jenkins, GitLab CI и пр.). 
### Пайплайн должен в себе содержать
- [x] запуск последней версии `golangci-lint` на весь проект **некоторые линтеры отключены для тестов, так же некоторые линтеры из требуемых ругаются что они отключены, пришлось их убрать**
- [x] запуск юнит тестов командой вида `go test -race -count 100`; **go test -v -count=100 -race -timeout=1m ./...**
- [x] сборку бинаря сервиса для версии Go не ниже 1.14. **1.22**
