# ReviewAssigner

Сервис автоматического назначения ревьюверов для Pull Request’ов в команде.

## Как запустить (одна команда)

```bash
git clone <ваш-репозиторий>
cd ReviewAssigner
docker-compose up --build
```

Сервис будет доступен по адресу http://localhost:8080
PostgreSQL поднимется автоматически, миграции применятся при старте.

## Аутентификация

Сначала получите JWT-токен:

```bash
# Admin — полный доступ
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"user_id": "admin", "password": "admin"}'

# Обычный пользователь — только чтение
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"user_id": "user", "password": "user"}'
```

Ответ:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.xxxxx",
  "role": "admin"
}
```

Все дальнейшие запросы (кроме `/health` и `/auth/login`) требуют заголовок:
```
Authorization: Bearer <ваш_токен>
```

## Полные примеры запросов (curl)

### 1. Health check
```bash
curl http://localhost:8080/health
# → {"status":"ok"}
```

### 2. Создать команду
```bash
curl -X POST http://localhost:8080/team/add \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "backend",
    "members": [
      {"user_id": "u1", "username": "Alice", "team_name": "backend", "is_active": true},
      {"user_id": "u2", "username": "Bob", "team_name": "backend", "is_active": true},
      {"user_id": "u3", "username": "Charlie", "team_name": "backend", "is_active": true},
      {"user_id": "u4", "username": "David", "team_name": "backend", "is_active": true}
    ]
  }' | jq
```

### 3. Получить команду
```bash
curl "http://localhost:8080/team/get?team_name=backend" \
  -H "Authorization: Bearer <token>" | jq
```

### 4. Создать PR (автоматически назначит до 2 активных ревьюверов)
```bash
curl -X POST http://localhost:8080/pullRequest/create \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-1001",
    "pull_request_name": "Add pagination",
    "author_id": "u1"
  }' | jq
```

### 5. Посмотреть назначенных ревьюверов
```bash
curl "http://localhost:8080/users/getReview?user_id=u2" \
  -H "Authorization: Bearer <token>" | jq
```

### 6. Замержить PR (идемпотентно)
```bash
curl -X POST http://localhost:8080/pullRequest/merge \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"pull_request_id": "pr-1001"}' | jq
```

### 7. Переназначить ревьювера (только на OPEN PR)
```bash
curl -X POST http://localhost:8080/pullRequest/reassign \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-1001",
    "old_user_id": "u2"
  }' | jq
```

### 8. Деактивировать пользователя (не будет назначаться на новые PR)
```bash
curl -X POST http://localhost:8080/users/setIsActive \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"user_id": "u3", "is_active": false}' | jq
```

### 9. Статистика назначений (дополнительная фича)
```bash
curl http://localhost:8080/stats \
  -H "Authorization: Bearer <token>" | jq
```

Пример ответа:
```json
{
  "user_assignments": {
    "u2": 5,
    "u3": 3,
    "u4": 7
  },
  "pr_assignments": {
    "pr-1001": 2,
    "pr-1002": 1
  }
}
```

## Особенности реализации

- Полная чистая архитектура (usecase → repository → delivery)
- Два репозитория: Postgres (прод) + in-memory (для тестов)
- JWT + роли admin/user
- Идемпотентный merge
- Автоматические миграции при старте
- Unit-тесты для всех usecase
- Pre-commit hooks + golangci-lint
- Запуск одной командой без дополнительных действий


Команда для запуска `docker-compose up --build`.
