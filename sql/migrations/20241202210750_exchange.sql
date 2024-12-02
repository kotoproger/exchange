-- +goose Up
-- +goose StatementBegin
CREATE TABLE if not exists general.current_rates (
    currency_from varchar(3) not null,
    currency_to varchar(3) not null,
    updated_at timestamptz not null,
    rate numeric,
    PRIMARY KEY (currency_from, currency_to)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists general.current_rates;
-- +goose StatementEnd
