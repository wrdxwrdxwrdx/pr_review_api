# PR Review API

REST API сервис для управления ревью pull request'ов, командами и пользователями. Построен на Go, фреймворке Gin и PostgreSQL.

## Требования

- Docker и Docker Compose
- Go 1.23+ (если запускаете локально без Docker)

## Быстрый старт с Docker

1. **Клонируйте репозиторий** (если еще не сделано):

   ```bash
   git clone <repository-url>
   cd pr_review_api
   ```

2. **Создайте файл окружения**:

   ```bash
   cp .env.example .env
   ```

   Отредактируйте `.env` и при необходимости измените значения.

3. **Запустите сервисы**:

   ```bash
   docker-compose up --build
   ```

   Это выполнит:

   - Запуск контейнера PostgreSQL 16
   - Сборку и запуск Go API сервера
   - Автоматическую инициализацию схемы базы данных
   - Открытие API на `http://localhost:8080` (или на вашем настроенном `SERVER_PORT`)

4. **Проверьте, что сервисы запущены**:

   ```bash
   docker-compose ps
   ```

5. **Доступ к базе данных** (опционально):

   ```bash
   docker exec -it pr_review_db psql -U postgres -d pr_review
   ```

## API Endpoints

### Endpoints аутентификации

#### Регистрация пользователя

Создать новый аккаунт пользователя и получить JWT токен.

**Endpoint**: `POST /auth/register`

**Аутентификация**: Не требуется (публичный)

**Тело запроса**:

```json
{
  "user_id": "u1",
  "username": "Alice"
}
```

**Ответ** (200 OK):

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_id": "u1",
  "expires_in": 86400
}
```

#### Вход

Аутентифицировать существующего пользователя и получить JWT токен.

**Endpoint**: `POST /auth/login`

**Аутентификация**: Не требуется (публичный)

**Тело запроса**:

```json
{
  "user_id": "u1",
  "username": "Alice"
}
```

**Ответ** (200 OK):

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_id": "u1",
  "expires_in": 86400
}
```

### Endpoints команд

#### Создать команду

Создать новую команду с участниками. Этот endpoint создает/обновляет пользователей и назначает их в команду.

**Endpoint**: `POST /team/add`

**Аутентификация**: JWT Bearer токен (обязательно)

**Заголовки**:

```
Authorization: Bearer <jwt_token>
```

**Тело запроса**:

```json
{
  "team_name": "backend",
  "members": [
    {
      "user_id": "u1",
      "username": "Alice",
      "is_active": true
    },
    {
      "user_id": "u2",
      "username": "Bob",
      "is_active": true
    }
  ]
}
```

**Ответ** (201 Created):

```json
{
  "team": {
    "team_name": "backend",
    "members": [
      {
        "user_id": "u1",
        "username": "Alice",
        "is_active": true
      },
      {
        "user_id": "u2",
        "username": "Bob",
        "is_active": true
      }
    ]
  }
}
```

#### Получить команду

Получить информацию о команде со всеми участниками.

**Endpoint**: `GET /team/get?TeamNameQuery=<team_name>`

**Аутентификация**: Админ токен ИЛИ JWT Bearer токен

**Заголовки** (выберите один):

```
X-Admin-Token: <admin_token>
```

ИЛИ

```
Authorization: Bearer <jwt_token>
```

**Параметры запроса**:

- `TeamNameQuery` (обязательно): Имя команды

**Пример запроса**:

```bash
curl -H "X-Admin-Token: admin-secret-token" \
  "http://localhost:8080/team/get?TeamNameQuery=backend"
```

**Ответ** (200 OK):

```json
{
  "team_name": "backend",
  "members": [
    {
      "user_id": "u1",
      "username": "Alice",
      "is_active": true
    },
    {
      "user_id": "u2",
      "username": "Bob",
      "is_active": true
    }
  ]
}
```

### Endpoints пользователей

#### Установить статус активности пользователя

Обновить статус активности пользователя (только для админа).

**Endpoint**: `POST /users/setIsActive`

**Аутентификация**: Админ токен (обязательно)

**Заголовки**:

```
X-Admin-Token: <admin_token>
```

**Тело запроса**:

```json
{
  "user_id": "u2",
  "is_active": false
}
```

**Ответ** (200 OK):

```json
{
  "user": {
    "user_id": "u2",
    "username": "Bob",
    "team_name": "backend",
    "is_active": false
  }
}
```

#### Получить ревью пользователя

Получить все pull request'ы, где пользователь назначен ревьювером.

**Endpoint**: `GET /users/getReview?UserIdQuery=<user_id>`

**Аутентификация**: Админ токен ИЛИ JWT Bearer токен

**Заголовки** (выберите один):

```
X-Admin-Token: <admin_token>
```

ИЛИ

```
Authorization: Bearer <jwt_token>
```

**Параметры запроса**:

- `UserIdQuery` (обязательно): ID пользователя

**Примечание**: Обычные пользователи могут просматривать только свои ревью. Админы могут просматривать ревью любого пользователя.

**Пример запроса**:

```bash
curl -H "Authorization: Bearer <jwt_token>" \
  "http://localhost:8080/users/getReview?UserIdQuery=u2"
```

**Ответ** (200 OK):

```json
{
  "user_id": "u2",
  "pull_requests": [
    {
      "pull_request_id": "pr-1001",
      "pull_request_name": "Add search",
      "author_id": "u1",
      "status": "OPEN"
    }
  ]
}
```

### Endpoints Pull Request'ов

#### Создать Pull Request

Создать новый pull request и автоматически назначить до 2 ревьюверов из команды автора.

**Endpoint**: `POST /pullRequest/create`

**Аутентификация**: Админ токен (обязательно)

**Заголовки**:

```
X-Admin-Token: <admin_token>
```

**Тело запроса**:

```json
{
  "pull_request_id": "pr-1001",
  "pull_request_name": "Add search",
  "author_id": "u1"
}
```

**Ответ** (201 Created):

```json
{
  "pr": {
    "pull_request_id": "pr-1001",
    "pull_request_name": "Add search",
    "author_id": "u1",
    "status": "OPEN",
    "assigned_reviewers": ["u2", "u3"],
    "createdAt": "2025-01-15T10:30:00Z",
    "mergedAt": null
  }
}
```

#### Смержить Pull Request

Пометить pull request как смерженный (идемпотентная операция).

**Endpoint**: `POST /pullRequest/merge`

**Аутентификация**: Админ токен (обязательно)

**Заголовки**:

```
X-Admin-Token: <admin_token>
```

**Тело запроса**:

```json
{
  "pull_request_id": "pr-1001"
}
```

**Ответ** (200 OK):

```json
{
  "pr": {
    "pull_request_id": "pr-1001",
    "pull_request_name": "Add search",
    "author_id": "u1",
    "status": "MERGED",
    "assigned_reviewers": ["u2", "u3"],
    "createdAt": "2025-01-15T10:30:00Z",
    "mergedAt": "2025-01-15T14:20:00Z"
  }
}
```

#### Переназначить ревьювера

Заменить конкретного ревьювера на другого активного пользователя из той же команды.

**Endpoint**: `POST /pullRequest/reassign`

**Аутентификация**: Админ токен (обязательно)

**Заголовки**:

```
X-Admin-Token: <admin_token>
```

**Тело запроса**:

```json
{
  "pull_request_id": "pr-1001",
  "old_reviewer_id": "u2"
}
```

**Ответ** (201 Created):

```json
{
  "pr": {
    "pull_request_id": "pr-1001",
    "pull_request_name": "Add search",
    "author_id": "u1",
    "status": "OPEN",
    "assigned_reviewers": ["u3", "u5"],
    "createdAt": "2025-01-15T10:30:00Z",
    "mergedAt": null
  },
  "replaced_by": "u5"
}
```

## Дополнительное задание

### Получить статистику Pull Request'ов

Получить статистику о количестве pull request'ов, созданных каждым пользователем.

**Endpoint**: `GET /pullRequest/statistics`

**Аутентификация**: Админ токен (обязательно)

**Заголовки**:

```
X-Admin-Token: <admin_token>
```

**Пример запроса**:

```bash
curl -H "X-Admin-Token: admin-secret-token" \
  "http://localhost:8080/pullRequest/statistics"
```

**Ответ** (200 OK):

```json
{
  "stats": [
    {
      "user_id": "u1",
      "pull_request_number": 5
    },
    {
      "user_id": "u2",
      "pull_request_number": 3
    },
    {
      "user_id": "u3",
      "pull_request_number": 2
    }
  ]
}
```

Ответ содержит массив статистики, где каждая запись показывает:

- `user_id`: ID пользователя, создавшего pull request'ы
- `pull_request_number`: Общее количество pull request'ов, созданных этим пользователем

## Пример рабочего процесса

Вот полный пример рабочего процесса с использованием curl:

1. **Зарегистрировать пользователя**:

   ```bash
   curl -X POST http://localhost:8080/auth/register \
     -H "Content-Type: application/json" \
     -d '{
       "user_id": "u1",
       "username": "Alice"
     }'
   ```

   Сохраните `token` из ответа.

2. **Создать команду** (используя JWT токен из шага 1):

   ```bash
   curl -X POST http://localhost:8080/team/add \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <token_from_step_1>" \
     -d '{
       "team_name": "backend",
       "members": [
         {"user_id": "u1", "username": "Alice", "is_active": true},
         {"user_id": "u2", "username": "Bob", "is_active": true},
         {"user_id": "u3", "username": "Charlie", "is_active": true}
       ]
     }'
   ```

3. **Создать pull request** (используя админ токен):

   ```bash
   curl -X POST http://localhost:8080/pullRequest/create \
     -H "Content-Type: application/json" \
     -H "X-Admin-Token: admin-secret-token" \
     -d '{
       "pull_request_id": "pr-1001",
       "pull_request_name": "Add search feature",
       "author_id": "u1"
     }'
   ```

4. **Получить ревью пользователя** (используя JWT токен):

   ```bash
   curl -X GET "http://localhost:8080/users/getReview?UserIdQuery=u2" \
     -H "Authorization: Bearer <token_from_step_1>"
   ```

5. **Смержить pull request** (используя админ токен):

   ```bash
   curl -X POST http://localhost:8080/pullRequest/merge \
     -H "Content-Type: application/json" \
     -H "X-Admin-Token: admin-secret-token" \
     -d '{
       "pull_request_id": "pr-1001"
     }'
   ```

## Структура проекта

```
pr_review_api/
├── api/
│   └── openapi.yaml          # OpenAPI спецификация
├── cmd/
│   └── api/
│       └── main.go           # Точка входа приложения
├── internal/
│   ├── config/               # Управление конфигурацией
│   ├── domain/               # Доменные сущности и ошибки
│   ├── handlers/             # HTTP обработчики запросов
│   ├── middleware/           # HTTP middleware (auth, logging, etc.)
│   ├── repository/           # Слой доступа к данным
│   │   ├── postgres/         # Реализация PostgreSQL
│   │   └── interfaces/       # Интерфейсы репозиториев
│   └── services/             # Слой бизнес-логики
├── pkg/
│   ├── auth/                 # JWT аутентификация
│   └── validator/            # Утилиты валидации
├── Dockerfile                # Определение Docker образа
├── docker-compose.yml        # Конфигурация Docker Compose
├── go.mod                    # Зависимости Go модуля
└── README.md                 # Этот файл
```

## Схема базы данных

Приложение автоматически инициализирует следующую схему базы данных:

- **teams**: Хранит информацию о командах
- **users**: Хранит информацию о пользователях со связями с командами
- **pull_requests**: Хранит информацию о pull request'ах с назначенными ревьюверами

Все таблицы включают автоматическое управление временными метками для полей `created_at` и `updated_at`.

