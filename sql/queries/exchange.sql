-- name: GetCuurentRate :one
select currency_from, currency_to, rate from general.current_rates where currency_from = $1 AND currency_to = $2;


-- name: GetRateOnDate :one
select hr1.currency_from, hr1.currency_to, hr1.rate from general.history_rates as hr1
where
	currency_from = $1
	and currency_to = $2
	and created_at < $3
order by
	currency_from asc,
	currency_to asc,
	created_at desc
limit 1;

-- name: UpdateRate :exec
insert into general.current_rates (currency_from, currency_to, rate, updated_at) 
values ($1, $2, $3, now())
on conflict (currency_from, currency_to) do update 
    set rate = EXCLUDED.rate, updated_at=EXCLUDED.updated_at
    where current_rates.rate != EXCLUDED.rate
;

-- name: ArchiveRate :exec
insert into general.history_rates(currency_from, currency_to, created_at, rate) 
(select cur.currency_from, cur.currency_to, cur.updated_at, cur.rate from general.current_rates as cur where cur.currency_from = $1 AND cur.currency_to = $2)
ON CONFLICT (currency_from, currency_to, created_at) DO NOTHING
;
