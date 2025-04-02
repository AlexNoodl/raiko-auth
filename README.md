# Raiko Auth

Raiko Auth — это RESTful сервис авторизации на Go для регистрации и входа пользователей. Использует MongoDB, JWT, bcrypt и Swagger.

## Основные возможности

- Регистрация и вход пользователей (email или username).
- Валидация паролей (8-20 символов, 1 заглавная, 1 строчная, 1 цифра, 1 спецсимвол).
- Логирование через `logrus`.
- Конфигурация через `.env`.
- Документация API в Swagger UI.

## Технологии

- Go 1.23
- Gin
- MongoDB
- JWT
- bcrypt
- logrus
- godotenv
- Swagger (swaggo)

## Установка и запуск

### Локальный запуск

1. **Клонируйте репозиторий**:
   ```bash
   git clone <repository-url>
   cd raiko-auth
   ```

2. **Установите зависимости**:
   ```bash
   go mod tidy
   ```

3. **Настройте `.env`**:
   ```
   PORT=8080
   MONGODB_URI=mongodb://user:password@localhost:27017/auth_db?authSource=admin
   MONGODB_DATABASE=auth_db
   JWT_SECRET=your-secret-key
   ```

4. **Сгенерируйте Swagger**:
   ```bash
   swag init
   ```

5. **Запустите**:
   ```bash
   go run cmd/main.go
   ```

### Запуск в Docker


1.**Запустите приложение**:
   ```bash
   docker-compose up --build
   ```
2.**Проверьте логи**:
   ```bash
   docker logs raiko-auth
   ```

## Swagger Документация

- Локально: `http://localhost:8080/swagger/index.html`
- В Docker: `http://localhost:5001/swagger/index.html`