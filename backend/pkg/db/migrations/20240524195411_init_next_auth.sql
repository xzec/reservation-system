-- +goose Up
-- +goose StatementBegin
create table if not exists users
(
    id             uuid primary key default gen_random_uuid(),
    created_at     timestamp        default current_timestamp not null,
    updated_at     timestamp        default current_timestamp not null,
    email          text unique                                not null,
    email_verified timestamptz,
    name           text,
    image          text
);

create table sessions
(
    id            uuid primary key default gen_random_uuid(),
    created_at    timestamp        default current_timestamp not null,
    updated_at    timestamp        default current_timestamp not null,
    user_id       uuid                                       not null references users (id) on delete cascade,
    expires       timestamptz                                not null,
    session_token text                                       not null
);

create table accounts
(
    id                  uuid primary key default gen_random_uuid(),
    created_at          timestamp        default current_timestamp not null,
    updated_at          timestamp        default current_timestamp not null,
    user_id             uuid                                       not null references users (id) on delete cascade,
    type                text                                       not null,
    provider            text                                       not null,
    provider_account_id text                                       not null,
    refresh_token       text,
    access_token        text,
    expires_at          bigint,
    id_token            text,
    scope               text,
    session_state       text,
    token_type          text
);

create table verification_tokens
(
    identifier text        not null,
    expires    timestamptz not null,
    token      text        not null,

    primary key (identifier, token)
);

create or replace function set_updated_at_now() returns trigger as
$$
begin
    if row (new.*) is distinct from row (old.*) then
        new.updated_at = current_timestamp;
    end if;
    return new;
end;
$$ language plpgsql;

create trigger users_set_updated_at_now
    before update
    on users
    for each row
    when (row (new.*) is distinct from row (old.*))
execute procedure set_updated_at_now();

create trigger sessions_set_updated_at_now
    before update
    on sessions
    for each row
    when (row (new.*) is distinct from row (old.*))
execute function set_updated_at_now();

create trigger accounts_set_updated_at_now
    before update
    on accounts
    for each row
    when (row (new.*) is distinct from row (old.*))
execute function set_updated_at_now();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop trigger if exists users_set_updated_at_now on users;
drop trigger if exists sessions_set_updated_at_now on sessions;
drop trigger if exists accounts_set_updated_at_now on accounts;
drop trigger if exists verification_tokens_set_updated_at_now on verification_tokens;

drop function if exists set_updated_at_now();

drop table if exists verification_tokens;
drop table if exists accounts;
drop table if exists sessions;
drop table if exists users;
-- +goose StatementEnd
