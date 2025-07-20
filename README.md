# ActionLogger

Система мониторинга пользовательских действий

Задача: Реализовать gRPC API для системы логов пользовательских действий с возможностью фильтрации и потоковой выдачи данных.

Требования:

1. Сервис на Connect:  
Создать proto-файл с описанием сервиса ActionLogger Реализовать сервер с использованием connect-go

2. PostgreSQL + pgx:  
Создать таблицу user_actions : id, user_id, action_type, timestamp, details (JSONB) Использовать pgx v5 для работы с БД

3. Фильтрация:  
Реализовать суммарные запросы с комбинацией фильтров:
- по user_id  
- по action_type  
- временной диапазон
- поиск по details (JSONB поле)

4. Потоковая выдача:  
- Для GetActions использовать server-side streaming
- Лимитировать выдачу пачками по 100 записей

5. Мониторинг:  
- WatchActions должен отслеживать новые события (использовать LISTEN/NOTIFY PostgreSQL)

# Установка

```sh
cp .env.example .env
make install-deps
docker compose up
make local-migration-up
```

# Проверка

```sh
go run cmd/server/main.go
go run cmd/client/main.go
```

 Скрипт [client/main.go](cmd/client/main.go) запускает по очереди:
 - LogAction
 - GetActions без фильтра (в результате 10 записей - 9 из миграции [seed_user_actions](migrations/20250720001325_seed_user_actions.sql) и 1 которую только что записал LogAction)
 - GetActions с фильтром (из имеющихся под фильтр подпадают 3 записи из миграции, отмеченные в [файле миграции](migrations/20250720001325_seed_user_actions.sql) комментариями)
 - WatchActions

Для проверки WatchActions можно запустить [LogAction.http](test/LogAction.http) , создаваемые записи будут выводиться в консоли