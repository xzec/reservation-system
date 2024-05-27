-- +goose Up
-- +goose StatementBegin
alter table accounts
add constraint unique_provider_account unique (provider, provider_account_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table accounts
drop constraint unique_provider_account;
-- +goose StatementEnd
