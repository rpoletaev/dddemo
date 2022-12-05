-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS subscriptions (
    id SERIAL PRIMARY KEY,
    user_id bigint not null unique,
    status varchar(10) not null,
    status_changed_at datetime not null,
    is_need_to_prolong boolean,
    active_since datetime,
    active_until datetime
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS subscriptions;
-- +goose StatementEnd
