# Домашнее задание 5: `shopx`

В этом задании нужно реализовать многосервисное приложение на тему **`Интернет-магазин X`**.
Тематика магазина может быть любой: книги, электроника, мерч, настольные игры, кофе, одежда и так далее.
Главное, чтобы технический контракт задания был соблюдён.

## Описание задачи

Нужно реализовать учебную микросервисную систему из трёх процессов:

1. `shop-backend` — основной HTTP-бэкенд магазина;
2. `auth-service` — сервис авторизации, доступный для backend только по `gRPC`;
3. `order-worker` — многопоточный сервис обработки очереди заказов.

Архитектура должна выглядеть так:

- пользователь работает с `shop-backend` по HTTP/JSON;
- `shop-backend` вызывает `auth-service` по `gRPC`;
- `auth-service` хранит пользователей в `Postgres`;
- `auth-service` хранит пользовательские сессии в `Redis`;
- после оформления заказа `shop-backend` публикует событие в `Kafka`;
- `order-worker` читает сообщения из Kafka и обрабатывает заказы конкурентно, используя `goroutines` и `channels`.

Фронтенд, HTML и браузерный интерфейс не требуются. Достаточно HTTP API.

## Требования к сдаче

1. Следуй стандартной структуре Go-проекта с несколькими исполняемыми файлами.
2. В проекте должны быть три собираемых бинарных файла:
   - `shop-backend`
   - `auth-service`
   - `order-worker`
3. Основной backend обязан работать по HTTP/JSON.
4. Сервис авторизации обязан работать по `gRPC`.
5. Хранение пользователей обязательно должно быть в `Postgres`.
6. Хранение пользовательских сессий обязательно должно быть в `Redis`.
7. После оформления заказа сообщение обязательно должно публиковаться в `Kafka`.
8. Обработка заказов должна выполняться отдельным многопоточным сервисом `order-worker`.
9. В `order-worker` обязательно должно быть осмысленное конкурентное решение: например, `worker pool`, `pipeline`, `fan-out/fan-in`, `producer-consumer`.
10. Юнит-тесты прикладывать в проект не нужно.
11. Можно использовать сторонние библиотеки, но проверь, что нужные версии доступны через `https://proxy.golang.org`.
12. Все сервисы должны корректно завершаться по `SIGINT` и `SIGTERM`.
13. Схема БД и нужные таблицы должны создаваться автоматически при старте приложения или через встроенную миграцию, запускаемую самим приложением.

## Исполняемые файлы

### 1. `auth-service`

**Сборка из:** `cmd/auth-service`

**Назначение:** регистрация пользователя, логин, валидация сессии.

**Синтаксис запуска:**

```bash
auth-service \
  --listen 127.0.0.1:9090 \
  --postgres "postgres://shopx:shopx@localhost:5432/shopx?sslmode=disable" \
  --redis "127.0.0.1:6379"
```

### 2. `shop-backend`

**Сборка из:** `cmd/shop-backend`

**Назначение:** HTTP API магазина.

**Синтаксис запуска:**

```bash
shop-backend \
  --listen 127.0.0.1:8080 \
  --auth-grpc 127.0.0.1:9090 \
  --postgres "postgres://shopx:shopx@localhost:5432/shopx?sslmode=disable" \
  --kafka "localhost:9092" \
  --orders-topic "orders.created"
```

### 3. `order-worker`

**Сборка из:** `cmd/order-worker`

**Назначение:** конкурентная обработка заказов из Kafka.

**Синтаксис запуска:**

```bash
order-worker \
  --postgres "postgres://shopx:shopx@localhost:5432/shopx?sslmode=disable" \
  --kafka "localhost:9092" \
  --orders-topic "orders.created" \
  --group-id "shopx-workers" \
  --workers 3
```

## Архитектурные требования

### Общая модель

- `shop-backend` не должен хранить пользователей и пароли локально;
- все операции регистрации, логина и проверки токена должны проходить через `auth-service` по `gRPC`;
- пароли должны храниться в `Postgres` только в виде хеша;
- сессии должны храниться в `Redis`;
- заказ после HTTP-оформления не должен синхронно обрабатываться в backend;
- backend должен записать заказ в БД со статусом `new`, опубликовать событие в Kafka и быстро вернуть ответ клиенту;
- `order-worker` должен получать сообщение из Kafka и менять статус заказа минимум по цепочке `new -> processing -> processed` либо `new -> processing -> failed`.

### Конкурентность

`order-worker` — это ключевая многопоточная часть задания.

Минимально требуется:

1. чтение сообщений из Kafka;
2. передача задач на обработку в несколько горутин;
3. безопасное обновление статусов заказов без гонок данных;
4. корректное завершение воркеров по сигналу;
5. использование как минимум одного узнаваемого concurrency-паттерна.

Примеры допустимых решений:

- `worker pool` для параллельной обработки заказов;
- `pipeline` вида `consume -> decode -> process -> persist`;
- `fan-out/fan-in` для распределения сообщений и сбора результатов;
- `producer-consumer` на каналах;
- `context cancellation` для graceful shutdown.

## Обязательный gRPC-контракт

В проекте должен быть proto-файл:

```text
api/auth/v1/auth.proto
```

Пакет:

```text
auth.v1
```

Сервис:

```text
AuthService
```

Минимально обязательные RPC:

1. `Register`
2. `Login`
3. `Validate`

Рекомендуемый смысл методов:

- `Register` — создать пользователя;
- `Login` — проверить пароль, создать сессию в Redis и вернуть токен;
- `Validate` — проверить токен сессии и вернуть данные пользователя.

Для удобства тестирования `auth-service` должен включать **gRPC reflection**.

## Обязательный HTTP API backend

### 1. `POST /api/register`

**Тело запроса:**

```json
{
  "email": "alice@example.com",
  "password": "secret123",
  "name": "Alice"
}
```

**Успешный ответ:** `201 Created`

```json
{
  "id": "user-1",
  "email": "alice@example.com",
  "name": "Alice"
}
```

### 2. `POST /api/login`

**Тело запроса:**

```json
{
  "email": "alice@example.com",
  "password": "secret123"
}
```

**Успешный ответ:** `200 OK`

```json
{
  "sessionToken": "token-abc"
}
```

После логина токен должен быть сохранён в `Redis`.

### 3. `GET /api/me`

Требует HTTP-заголовок:

```text
X-Session-Token: token-abc
```

**Успешный ответ:** `200 OK`

```json
{
  "id": "user-1",
  "email": "alice@example.com",
  "name": "Alice"
}
```

Backend должен валидировать токен не самостоятельно, а через `gRPC`-вызов `auth-service`.

### 4. `POST /api/orders`

Требует HTTP-заголовок:

```text
X-Session-Token: token-abc
```

**Тело запроса:**

```json
{
  "items": [
    { "sku": "book-1", "qty": 2 },
    { "sku": "pen-7", "qty": 1 }
  ]
}
```

**Успешный ответ:** `202 Accepted`

```json
{
  "id": "order-1",
  "status": "new"
}
```

Поведение:

- backend валидирует пользователя через `auth-service`;
- создаёт запись о заказе в таблице `orders`;
- сохраняет заказ со статусом `new`;
- публикует сообщение в Kafka topic `orders.created`;
- не обрабатывает заказ синхронно в HTTP-обработчике.

### 5. `GET /api/orders/{id}`

Требует HTTP-заголовок:

```text
X-Session-Token: token-abc
```

**Успешный ответ:** `200 OK`

```json
{
  "id": "order-1",
  "userId": "user-1",
  "status": "processed",
  "items": [
    { "sku": "book-1", "qty": 2 },
    { "sku": "pen-7", "qty": 1 }
  ]
}
```

Допустимые статусы:

- `new`
- `processing`
- `processed`
- `failed`

## Требования к данным

### Таблица `users`

В `Postgres` обязательно должна быть таблица `users` минимум с такими полями:

- `id`
- `email`
- `name`
- `password_hash`
- `created_at`

Требования:

- `email` должен быть уникальным;
- пароль в открытом виде хранить запрещено;
- для хранения хеша рекомендуется `bcrypt`.

### Таблица `orders`

В `Postgres` обязательно должна быть таблица `orders` минимум с такими полями:

- `id`
- `user_id`
- `status`
- `items`
- `created_at`
- `updated_at`

Поле `items` может быть реализовано как `JSON/JSONB`, строка JSON или иным понятным способом, если данные
корректно читаются и возвращаются через API.

### Redis

В `Redis` должны храниться сессии пользователя.

Минимальное ожидаемое поведение:

- после логина появляется запись с токеном;
- по токену можно определить пользователя;
- при проверке `/api/me` и защищённых заказов backend получает пользователя через `auth-service`.

Рекомендуемый ключ:

```text
session:<token>
```

## Требования к обработке заказов

`order-worker` обязан:

1. подключаться к Kafka;
2. читать сообщения из topic `orders.created`;
3. обрабатывать сообщения конкурентно с количеством воркеров `--workers`;
4. менять статус заказа на `processing`;
5. после успешной обработки менять статус на `processed`;
6. при ошибке помечать заказ как `failed`;
7. корректно закрывать consumer и завершать воркеры по сигналу.

Сама бизнес-логика обработки может быть упрощённой. Например:

- искусственная задержка;
- проверка состава заказа;
- имитация резервирования;
- расчёт стоимости;
- запись служебного события в лог.

Главное, чтобы было видно асинхронную и конкурентную обработку очереди.

## Требования к завершению

Все три сервиса должны корректно реагировать на `SIGINT` и `SIGTERM`.

Минимально требуется:

1. перестать принимать новые запросы или новые сообщения;
2. завершить текущие операции корректно;
3. закрыть соединения с Postgres, Redis, Kafka и gRPC;
4. завершить процесс без паники и зависания.

## Что не требуется

В этом задании не нужно:

- делать frontend;
- реализовывать оплату;
- интегрироваться с внешними платёжными системами;
- делать полноценный каталог товаров из отдельной БД;
- реализовывать роли администратора;
- писать unit-тесты внутри студенческого проекта.

## Пример локального сценария

```bash
# 1. auth-service
./auth-service \
  --listen 127.0.0.1:9090 \
  --postgres "postgres://shopx:shopx@localhost:5432/shopx?sslmode=disable" \
  --redis "127.0.0.1:6379"

# 2. shop-backend
./shop-backend \
  --listen 127.0.0.1:8080 \
  --auth-grpc 127.0.0.1:9090 \
  --postgres "postgres://shopx:shopx@localhost:5432/shopx?sslmode=disable" \
  --kafka "localhost:9092" \
  --orders-topic "orders.created"

# 3. order-worker
./order-worker \
  --postgres "postgres://shopx:shopx@localhost:5432/shopx?sslmode=disable" \
  --kafka "localhost:9092" \
  --orders-topic "orders.created" \
  --group-id "shopx-workers" \
  --workers 3
```

Примеры запросов:

```bash
curl -X POST http://127.0.0.1:8080/api/register \
  -H 'Content-Type: application/json' \
  -d '{"email":"alice@example.com","password":"secret123","name":"Alice"}'

curl -X POST http://127.0.0.1:8080/api/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"alice@example.com","password":"secret123"}'

curl http://127.0.0.1:8080/api/me \
  -H 'X-Session-Token: token-abc'

curl -X POST http://127.0.0.1:8080/api/orders \
  -H 'Content-Type: application/json' \
  -H 'X-Session-Token: token-abc' \
  -d '{"items":[{"sku":"book-1","qty":2},{"sku":"pen-7","qty":1}]}'
```

## Запуск и проверка

Сборка:

```bash
go build -o auth-service ./cmd/auth-service
go build -o shop-backend ./cmd/shop-backend
go build -o order-worker ./cmd/order-worker
```

Минимум, что должно проходить в автопроверке:

- сборка всех трёх бинарников;
- регистрация пользователя через backend;
- логин через backend;
- появление пользователя в `Postgres`;
- появление сессии в `Redis`;
- успешная валидация токена через `gRPC`;
- создание заказа через backend;
- асинхронный перевод заказа в `processed` сервисом `order-worker`;
- graceful shutdown сервисов.
