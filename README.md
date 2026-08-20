# URL Shortener

Сервис для сокращения URL с PostgreSQL и Docker.

## Стек

- **Go** — API
- **PostgreSQL 18** — хранение данных
- **Redis** — кеширование
- **Docker / Docker Compose** — контейнеризация
- **golang-migrate** — миграции
- **Makefile** — автоматизация

## Запуск

### Требования

- **Go 1.26.4+**
- **Docker / Docker Compose**
- **golang-migrate**
- **make**

### 1. Настройка окружения

Создай `.env` файл в корне проекта:

- Linux/macOS:
```bash
cp .env.example .env
```

- Windows PowerShell:
```bash
Copy-Item .env.example .env
```

### 2. Docker compose

```bash
docker-compose up
```

### После запуска API доступно на: http://localhost:8080

## Регистрация

**POST** `/register`

Тело запроса:
```json
{
    "email": "user@example.com",
    "password": "123456"
}
```

## Авторизация

**POST** `/login`

Тело запроса:
```json
{
    "email": "user@example.com",
    "password": "123456"
}
```

Тело ответа:
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

## Редирект

**GET** `/{code}`

Редирект на оригинальный URL

## Получить short ссылку

**POST** `/short`

Тело запроса:
```json
{
    "url":"long url"
}
```

## Посмотреть все свои ссылки

**GET** `/urls`

Тело ответа:
```json
[
    {
        "OriginalUrl": "https://example.com",
        "ShortCode": "abc123",
        "CreatedAt": "2026-08-12T14:30:00Z"
    },
    {
        "OriginalUrl": "https://google.com",
        "ShortCode": "def456",
        "CreatedAt": "2026-08-12T15:00:00Z"
    }
]
```

## Для разработки

### 1. Поднять PostgreSQL

```bash
make db-up
```

### 2. Выключить PostgreSQL

```bash
make db-down
```

### 3. Накатить последнюю миграцию

```bash
make migrate-up
```

### 4. Откатить миграцию на 1 шаг назад

```bash
make migrate-down
```

### 5. Откатить все миграции

```bash
make migrate-down-all
```

### 6. Создать миграцию (вместо your_name имя миграции)

```bash
make migrate-create name=your_name
```

### 7. Логи

```bash
make logs       # только приложения
make logs-all   # все контейнеры
```

## Тестирование

```bash
go test ./...
```

## CI/CD

При пуше в `main` GitHub Actions:
- Запускает тесты
- Собирает Docker-образ
- Пушит в Docker Hub с тегами `latest` и `<sha>`

При создании Pull Request - запускает тесты.