# ReviewAssigner

Сервис для автоматического назначения ревьюверов Pull Request'ам в командах.

## Запуск
1. `docker-compose up --build`
2. Сервис доступен на http://localhost:8080
3. API-документация: http://localhost:8080/swagger/index.html (без авторизации)

## API
- Health: GET /health
- Команды: POST /team/add, GET /team/get
- Пользователи: POST /users/setIsActive, GET /users/getReview
- PR: POST /pullRequest/create, POST /pullRequest/merge, POST /pullRequest/reassign
- Статистика: GET /stats (дополнительно)
- Массовая деактивация: POST /team/deactivateUsers (дополнительно)

Авторизация: AdminToken для POST/PUT/DELETE, UserToken для GET (кроме /health и /swagger).

## Тестирование
- Unit-тесты: `go test ./internal/usecase/...`
- Интеграционные: `go test ./tests/...`
- Нагрузочное: `k6 run tests/load_test.js` (5 RPS, SLI <300 мс, 99.9% успешности)
- Линтер: `golangci-lint run`

## Дополнительные фичи
- Статистика назначений.
- Массовая деактивация с переназначением PR.
- Pre-commit hooks: Установите pre-commit, затем `pre-commit install`.
- Безопасность: Параметризованные SQL-запросы (sqlx), проверка gosec.

## Деплой
- Для продакшена: Используйте Kubernetes/Docker Swarm. Переменные окружения для БД.
- Миграции: Автоматически применяются при `docker-compose up`.
- Мониторинг: Добавьте Prometheus для метрик (RPS, latency).