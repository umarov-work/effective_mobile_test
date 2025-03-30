# Тестовое задание Effective Mobile

## Описание проекта

Тестовое задание для компании Effective Mobile. Этот проект представляет собой REST API, разработанное на языке Go с использованием фреймворка Gin и ORM GORM для управления данными о людях. API позволяет создавать, получать, обновлять и удалять записи о людях, обогащая их данными о возрасте, поле и национальности с помощью внешних сервисов (agify.io, genderize.io, nationalize.io). Реализованы следующие возможности:
- Документация API через Swagger.
- Поддержка пагинации и фильтрации данных.
- Логирование операций с использованием Logrus.

### API
- POST /person - Create a new person
- GET /persons - Get list of persons
- PUT /person/{id} - Update a person
- DELETE /person/{id} - Delete a person
- GET /swagger/index.html - Swagger

## Требования

- **Go**: версия 1.21 или выше
- **PostgreSQL**: для хранения данных
- **Интернет**: для обогащения данных через внешние API

## Установка и запуск

### 1. Клонирование репозитория
Склонируйте репозиторий на свой компьютер:
```bash
git clone https://github.com/umarov-work/effective_mobile_test.git
cd effective_mobile_test
```
### 2. Настройка окружения
В файле ```effective_mobile_test/config/.env``` укажите корректные для Вас данные.
Например,
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=
DB_NAME=effective_mobile_test_db
PORT=3333
```
### 3. Установка зависимостей
```
cd effective_mobile_test
go mod tidy
```

### 4. Генерация Swagger
```
cd cmd
swag init --dir .,../internal/handlers,../internal/models --output ../docs
```

### 5. Запуск сервера
```
go run main.go
```

