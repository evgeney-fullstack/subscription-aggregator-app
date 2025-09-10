Сервис управления подписками

Описание

RESTful API для управления подписками пользователей. Сервис предоставляет возможности CRUDL-операций над подписками, а также расчет суммарной стоимости подписок за выбранный период с фильтрацией.

Требования
Go 1.19+

PostgreSQL 13+

Docker и Docker Compose (для запуска через контейнеры)

Установка и запуск

1. Клонирование репозитория

git clone git@github.com:evgeney-fullstack/subscription-aggregator-app.git

cd subscription-aggregator-app

3. Запуск с помощью Docker Compose
   
docker-compose up -d

Сервис будет доступен по адресу: https://localhost:8080

5. Локальный запуск (без Docker)
   
Настройка базы данных

Создайте базу данных PostgreSQL и настройте подключение в файле config.env:

DB_HOST=localhost

DB_PORT=5432

DB_USER=postgres

DB_PASSWORD=your_password

DB_NAME=subscriptions

DB_SSL_MODE=disable

Запуск сервера

go mod download

go run cmd/main.go

API Endpoints

Подписки (CRUDL)

GET /subscriptions - Получить список всех подписок

GET /subscriptions/{id} - Получить подписку по ID

POST /subscriptions - Создать новую подписку

PUT /subscriptions/{id} - Обновить подписку

DELETE /subscriptions/{id} - Удалить подписку

Суммарная стоимость

GET /subscriptions/total-cost - Получить суммарную стоимость подписок за период

Примеры запросов

Создание подписки

curl -X POST "http://localhost:8080/subscriptions" \

  -H "Content-Type: application/json" \
  
  -d '{
  
    "service_name": "Yandex Plus",    
    "price": 400,    
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",    
    "start_date": "07-2025"    
  }'
  
Получение суммарной стоимости

curl -X GET "http://localhost:8080/subscriptions/total-cost"

  -H "Content-Type: application/json" \
  
  -d '{  
  
     "period": {  
          "start_date": "01-2025",
          "finish_date": "10-2025"    
      },  
     "filters": {  
          "user_id": "60631fee-2bf1-4721-ae6c-7636679a0cba",
          "service_name": "Yandex Plus"    
      }  
  }'
  
Структура базы данных

Таблица subscriptions:

id - SERIAL PRIMARY KEY

service_name - VARCHAR(255) NOT NULL

price - INTEGER NOT NULL

user_id - UUID NOT NULL

start_date - DATE NOT NULL

finish_date - DATE GENERATED ALWAYS AS (start_date + INTERVAL '1 month') STORED

Конфигурация

Настройки сервиса вынесены в файл config.env:

DB_HOST - Хост PostgreSQL

DB_PORT - Порт PostgreSQL

DB_USER - Пользователь PostgreSQL

DB_PASSWORD - Пароль PostgreSQL

DB_NAME - Имя базы данных

DB_SSL_MODE - Режим SSL (disable/require/verify-ca/verify-full)

Логирование

Сервис использует структурированное логирование с уровнями:

INFO - информационные сообщения

ERROR - ошибки

Swagger документация

После запуска сервиса документация доступна по адресу: http://localhost:8080/swagger/index.html

Миграции

Миграции базы данных находятся в директории migrations/. Для применения миграций используйте github.com/golang-migrate/migrate/v4.

Тестирование

Запуск интеграционных тестов:

go test -v ./test -run TestIntegration

Автор
Разработчик: Evgeney Kovalev
