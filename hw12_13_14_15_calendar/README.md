#### Результатом выполнения следующих домашних заданий является сервис «Календарь»:
- [Домашнее задание №12 «Заготовка сервиса Календарь»](./docs/12_README.md)
- [Домашнее задание №13 «Внешние API от Календаря»](./docs/13_README.md)
- [Домашнее задание №14 «Кроликизация Календаря»](./docs/14_README.md)
- [Домашнее задание №15 «Докеризация и интеграционное тестирование Календаря»](./docs/15_README.md)

#### Ветки при выполнении
- `hw12_calendar` (от `master`) -> Merge Request в `master`
- `hw13_calendar` (от `hw12_calendar`) -> Merge Request в `hw12_calendar` (если уже вмержена, то в `master`)
- `hw14_calendar` (от `hw13_calendar`) -> Merge Request в `hw13_calendar` (если уже вмержена, то в `master`)
- `hw15_calendar` (от `hw14_calendar`) -> Merge Request в `hw14_calendar` (если уже вмержена, то в `master`)

**Домашнее задание не принимается, если не принято ДЗ, предшедствующее ему.**


### Заметки:
```text
./calendar --config=/configs/config.toml
./calendar version

для запуска указать:
./calendar           --config=configs/calendar_config.toml
./calendar_scheduler --config=configs/scheduler_config.yaml
./calendar_sender    --config=configs/sender_config.yaml


Создание Docker контейнера с PostgreSQL:
docker run -d --name pg -e POSTGRES_PASSWORD=password -p 5432:5432 postgres
docker run -d --name pg -e POSTGRES_PASSWORD=postgres -e PGDATA=/var/lib/postgresql/data/pgdata -v pg_data:/var/lib/postgresql/data -p 5432:5432 postgres


Проверка работы PostgreSQL:
psql -h localhost -p 5432 -U postgres -d postgres
Далее нужно ввести пароль, который установили при запуске контейнера.

Остановка и удаление контейнера:
docker stop pg
docker rm pg


Подключиться к БД
docker exec -it pg psql -Upostgres -dpostgres


Создание БД, создание пользователя и выдача прав:
create database exampledb; 
create user otus_user with encrypted password 'otus_password'; 
grant all privileges on database exampledb to otus_user;


DROP TABLE notifications;
DROP TABLE events;
select * from events;
select * from notifications;
\d events

Выполнение sql скрипта:
psql 'host=localhost user=postgres password=postgres dbname=postgres' < 20230921223047_events_and_notifications.sql


установить goose
go install github.com/pressly/goose/v3/cmd/goose@latest
goose --version
goose create events_and_notifications sql


golangci-lint run
golangci-lint run --fix


Запустить rabbitmq:
$ docker run -d --name rb -p 15672:15672 -p 5672:5672 rabbitmq:3-management 
docker rm <идентификатор>

Админка rabbitmq:
http://localhost:15672/ guest:guest


```

В папку docs положила коллекцию для postman и для grpc добавила в docs/for_grpc.md json объектов.





