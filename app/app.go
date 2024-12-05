package app

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kotoproger/exchange/internal/repository"
	"github.com/kotoproger/exchange/internal/source"
)

type App struct {
	repository  repository.Querier
	ctx         context.Context
	rateSources []source.ExchangeSource
	pool        *pgxpool.Pool
}

func NewApp(
	ctx context.Context,
	sources []source.ExchangeSource,
	pool *pgxpool.Pool,
) *App {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		panic(fmt.Errorf("cant acquire connection: %w", err))
	}
	return &App{
		repository:  repository.New(conn),
		ctx:         ctx,
		rateSources: sources,
		pool:        pool,
	}
}

func (app *App) convert(amount *money.Money, to *money.Currency, rate float64) *money.Money {
	newAmount := int64(math.Round(float64(amount.Amount()) * rate))
	return money.New(newAmount, to.Code)
}

func (app *App) Exchange(amount *money.Money, to *money.Currency) (*money.Money, error) {
	rateRow, repoError := app.repository.GetCuurentRate(
		app.ctx,
		repository.GetCuurentRateParams{CurrencyFrom: amount.Currency().Code, CurrencyTo: to.Code},
	)

	if repoError != nil {
		return nil, fmt.Errorf("find Exchange rate: %w", repoError)
	}

	pgfloat, err := rateRow.Rate.Float64Value()
	if err != nil {
		return nil, fmt.Errorf("get float 64 value from db: %w", err)
	}
	return app.convert(amount, to, pgfloat.Float64), nil
}

func (app *App) ExchangeToDate(amount *money.Money, to *money.Currency, date time.Time) (*money.Money, error) {
	rateRow, repoError := app.repository.GetRateOnDate(
		app.ctx,
		repository.GetRateOnDateParams{
			CurrencyFrom: amount.Currency().Code,
			CurrencyTo:   to.Code,
			CreatedAt:    pgtype.Timestamptz{Time: date},
		},
	)

	if repoError != nil {
		return nil, fmt.Errorf("find Exchange rate on date %s: %w", date.String(), repoError)
	}
	pgfloat, err := rateRow.Rate.Float64Value()
	if err != nil {
		return nil, fmt.Errorf("get float 64 value from db: %w", err)
	}
	return app.convert(amount, to, pgfloat.Float64), nil
}

func (app *App) UpdateRates() error {
	updatedPairs := make(map[string]map[string]bool)
	conn, acquireerr := app.pool.Acquire(app.ctx)
	if acquireerr != nil {
		return fmt.Errorf("acquire connection from pool: %w", acquireerr)
	}
	defer conn.Release()
	repositoryobject, ok := app.repository.(*repository.Queries)
	if !ok {
		return fmt.Errorf("get repository from interface")
	}
	for _, source := range app.rateSources {
		for rate := range source.Get() {
			_, ok := updatedPairs[rate.From.Code]
			if !ok {
				updatedPairs[rate.From.Code] = make(map[string]bool)
			}
			_, ok = updatedPairs[rate.From.Code][rate.To.Code]
			if ok {
				continue
			}
			var pgRate pgtype.Numeric
			scanErr := pgRate.Scan(strconv.FormatFloat(rate.Rate, 'f', -1, 64))
			if scanErr != nil {
				return fmt.Errorf("convert rate: %w", scanErr)
			}
			updatedPairs[rate.From.Code][rate.To.Code] = true
			transaction, err := conn.BeginTx(
				app.ctx,
				pgx.TxOptions{IsoLevel: pgx.ReadCommitted},
			)
			if err != nil {
				return fmt.Errorf("start transaction: %w", err)
			}

			transactionRepository := repositoryobject.WithTx(transaction)
			updateerr := transactionRepository.UpdateRate(
				app.ctx,
				repository.UpdateRateParams{
					CurrencyFrom: rate.From.Code,
					CurrencyTo:   rate.To.Code,
					Rate:         pgRate,
				},
			)
			if updateerr != nil {
				return fmt.Errorf("update rate: %w", updateerr)
			}

			archiveerr := transactionRepository.ArchiveRate(
				app.ctx,
				repository.ArchiveRateParams{
					CurrencyFrom: rate.From.Code,
					CurrencyTo:   rate.To.Code,
				},
			)
			if archiveerr != nil {
				return fmt.Errorf("update rate: %w", archiveerr)
			}

			commiterr := transaction.Commit(app.ctx)
			if commiterr != nil {
				return fmt.Errorf("commit: %w", commiterr)
			}
		}
	}

	return nil
}
