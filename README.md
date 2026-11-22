# Restaurant Reservation API

Учебный REST API для управления ресторанами, столами, гостями и бронированиями. Данные хранятся в PostgreSQL; пересечения активных бронирований одного стола запрещаются ограничением БД.

## Требования

- Go версии, указанной в `go.mod`
- PostgreSQL
- Docker для интеграционных тестов

## Запуск

Создайте базу данных `reservation`, скопируйте пример настроек и запустите сервер:

```sh
cp .env.example .env
```

При необходимости измените `DATABASE_URL` и `HTTP_ADDR` в `.env`, затем:

```sh
make run
```

При запуске приложение применяет миграции автоматически. По умолчанию API слушает `:8080`.

## Команды

```sh
make build             # собрать все пакеты
make test              # обычные тесты
make test-integration  # интеграционные тесты с PostgreSQL в Testcontainers
make migrate-up        # применить миграции вручную (нужна CLI golang-migrate)
make migrate-down      # откатить последнюю миграцию
```

Интеграционные тесты запускают PostgreSQL через Docker. Для ручных команд миграции укажите полный DSN в `DATABASE_URL`, включая нужные параметры подключения.

## API

Все маршруты имеют префикс `/api/v1`.

| Метод | Маршрут | Назначение |
| --- | --- | --- |
| `POST` | `/restaurants` | Создать ресторан |
| `GET` | `/restaurants/{id}` | Получить ресторан |
| `POST` | `/restaurants/{id}/tables` | Добавить стол |
| `GET` | `/restaurants/{id}/tables` | Список столов ресторана |
| `POST` | `/guests` | Создать гостя |
| `GET` | `/guests/{id}` | Получить гостя |
| `GET` | `/guests/by-phone?phone=...` | Найти гостя по телефону |
| `POST` | `/bookings` | Создать бронирование |
| `GET` | `/bookings/{id}` | Получить бронирование |
| `POST` | `/bookings/{id}/cancel` | Отменить бронирование |

Пример создания бронирования:

```sh
curl -X POST http://localhost:8080/api/v1/bookings \
  -H 'Content-Type: application/json' \
  -d '{"table_id":1,"guest_id":1,"starts_at":"2030-01-15T18:00:00Z","guests_count":2}'
```

Бронирование длится два часа. Время передаётся в формате RFC 3339.
