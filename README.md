# PR Reviewer Service

## Запуск

1. `docker-compose up --build`
2. Сервис доступен на http://localhost:8080

## API

- Health: GET /health
- См. docs/openapi.yaml или /swagger/index.html (если добавите swag UI).

## Тестирование
- Линтер: `golangci-lint run`
