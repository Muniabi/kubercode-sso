# KuberCode SSO Microservice

![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-green.svg)

Микросервис аутентификации и авторизации, построенный на Go с использованием Gin framework. Предоставляет REST API для управления пользователями, аутентификацией и сессиями.

## 🚀 Возможности

-   Регистрация и аутентификация пользователей
-   JWT-based аутентификация с refresh токенами
-   Двухфакторная аутентификация (2FA) с OTP
-   Управление сессиями и устройствами
-   Восстановление пароля
-   Изменение email и пароля
-   MongoDB для хранения данных
-   Redis для кэширования и управления сессиями
-   NATS для асинхронной коммуникации
-   SMTP для отправки email

## 📋 Предварительные требования

-   Go 1.24 или выше
-   MongoDB 4.4+
-   Redis 6+
-   NATS Server
-   SMTP сервер (для отправки email)

## 🛠 Установка

1. Клонируйте репозиторий:

```bash
git clone https://github.com/your-username/kubercode-sso.git
cd kubercode-sso
```

2. Установите зависимости:

```bash
go mod download
```

3. Создайте необходимые директории и файлы:

```bash
mkdir -p certs
```

4. Сгенерируйте RSA ключи для JWT:

```bash
openssl genrsa -out certs/jwtRSA256-private.pem 2048
openssl rsa -in certs/jwtRSA256-private.pem -pubout -out certs/jwtRSA256-public.pem
```

## ⚙️ Конфигурация

Основные настройки находятся в файле `config/config.yaml`. Основные параметры:

```yaml
env: "local"
JWTAlgorithm: "RSA256"
refreshTokenDurationDays: 30
accessTokenDurationMinutes: 10
mongoDbConnectionString: "mongodb://localhost:27017"
redisAddress: "localhost:6379"
```

## 🚀 Запуск

1. Запуск в режиме разработки:

```bash
go run cmd/sso/main.go
```

2. Сборка и запуск:

```bash
go build -o sso cmd/sso/main.go
./sso
```

## 📚 API Endpoints

### Аутентификация

#### Регистрация

```http
POST /api/v1/auth/signup
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "securepassword",
    "firstName": "John",
    "lastName": "Doe"
}
```

#### Вход

```http
POST /api/v1/auth/login
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "securepassword"
}
```

#### Выход

```http
POST /api/v1/auth/logout
Authorization: Bearer <access_token>
```

#### Обновление токена

```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
    "refreshToken": "<refresh_token>"
}
```

### Управление пользователем

#### Изменение пароля

```http
POST /api/v1/auth/change-password
Authorization: Bearer <access_token>
Content-Type: application/json

{
    "oldPassword": "oldpassword",
    "newPassword": "newpassword"
}
```

#### Изменение email

```http
POST /api/v1/auth/change-email
Authorization: Bearer <access_token>
Content-Type: application/json

{
    "newEmail": "newemail@example.com",
    "password": "currentpassword"
}
```

### Двухфакторная аутентификация

#### Отправка OTP

```http
POST /api/v1/auth/otp/send
Authorization: Bearer <access_token>
```

#### Проверка OTP

```http
POST /api/v1/auth/otp/verify
Content-Type: application/json

{
    "otp": "123456"
}
```

## 🔒 Безопасность

-   Все пароли хешируются с использованием bcrypt
-   JWT токены подписываются с использованием RSA-256
-   Поддержка CORS
-   Rate limiting для защиты от брутфорс атак
-   Сессии хранятся в Redis с TTL
-   Поддержка blacklist для revoked токенов

## 📦 Зависимости

-   [Gin](https://github.com/gin-gonic/gin) - Web framework
-   [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver) - MongoDB клиент
-   [Redis](https://github.com/go-redis/redis) - Redis клиент
-   [NATS](https://github.com/nats-io/nats.go) - NATS клиент
-   [JWT-Go](https://github.com/golang-jwt/jwt) - JWT реализация

## 🏗 Архитектура

Микросервис построен с использованием чистой архитектуры:

```
├── cmd/
│   └── sso/
│       └── main.go
├── internal/
│   ├── domain/
│   │   └── auth/
│   └── infrastructure/
│       ├── http/
│       │   ├── handlers/
│       │   └── middleware/
│       └── repository/
│           └── mongodb/
└── config/
    ├── config.go
    └── config.yaml
```

## 📝 Логирование

Сервис использует структурированное логирование с различными уровнями:

-   INFO - стандартные операции
-   WARN - предупреждения
-   ERROR - ошибки
-   DEBUG - отладочная информация

## 🔄 CI/CD

Рекомендуется настроить CI/CD пайплайн для:

-   Линтинга кода
-   Запуска тестов
-   Сборки Docker образа
-   Деплоя в Kubernetes

## 📄 Лицензия

MIT License - см. файл [LICENSE](LICENSE) для деталей.

## 👥 Команда

-   Разработка и поддержка: KuberCode Team

## 🤝 Вклад в проект

Мы приветствуем вклад в проект! Пожалуйста, следуйте этим шагам:

1. Форкните репозиторий
2. Создайте ветку для ваших изменений
3. Внесите изменения
4. Создайте Pull Request

## 📞 Поддержка

По всем вопросам обращайтесь:

-   Email: support@kubercode.com
-   Issues: GitHub Issues
