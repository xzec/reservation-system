-- +goose Up
-- +goose StatementBegin
create table if not exists users
(
    id             uuid      default gen_random_uuid() not null,
    email          text                                not null,
    email_verified boolean   default false             not null,
    created_at     timestamp default current_timestamp not null,
    updated_at     timestamp default current_timestamp not null
);

create unique index email
    on users (email);

create unique index id
    on users (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists users;
-- +goose StatementEnd
