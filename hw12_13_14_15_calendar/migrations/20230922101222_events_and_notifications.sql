-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
    id          VARCHAR(255) PRIMARY KEY,
    title       VARCHAR(255) NOT NULL DEFAULT 'new event',
    date_time   TIMESTAMP WITH TIME ZONE NOT NULL,
                              duration    BIGINT NOT NULL DEFAULT 0,
                              description TEXT NOT NULL DEFAULT '',
                              user_id     VARCHAR NOT NULL CHECK (title <> ''),
    reminder    BIGINT
    );
create index id_index on events (id);
create index user_id_index on events (user_id);

CREATE TABLE notifications (
    event_id    VARCHAR(255) PRIMARY KEY,
    title       VARCHAR(255),
    date_time   TIMESTAMP NOT NULL,
    user_id     VARCHAR(255) NOT NULL,
    FOREIGN KEY (event_id) REFERENCES events(id)
    );
create index event_id_index on events (id);
create index notifications_event_id_index on notifications (event_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS events;
-- +goose StatementEnd
