-- +goose Up
-- +goose StatementBegin
CREATE TABLE if not exists general.history_rates (
    currency_from varchar(3) not null,
    currency_to varchar(3) not null,
    created_at timestamptz not null,
    rate decimal,
    PRIMARY KEY (currency_from, currency_to, created_at)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists general.history_rates;
-- +goose StatementEnd
