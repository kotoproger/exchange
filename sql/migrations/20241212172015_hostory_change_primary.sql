-- +goose Up
-- +goose StatementBegin
ALTER TABLE "general".history_rates DROP CONSTRAINT history_rates_pkey;
CREATE UNIQUE INDEX history_rates_currency_from_idx ON "general".history_rates (currency_from,currency_to,created_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX "general".history_rates_currency_from_idx;
ALTER TABLE "general".history_rates ADD CONSTRAINT history_rates_pkey PRIMARY KEY (currency_from,currency_to,created_at);
-- +goose StatementEnd
