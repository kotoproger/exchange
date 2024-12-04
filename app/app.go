package app

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kotoproger/exchange/internal/repository"
	"github.com/kotoproger/exchange/internal/source"
)

type App struct {
	repository  repository.Queries
	ctx         context.Context
	rateSources []source.ExchangeSource
	conn        pgx.Conn
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

	return app.convert(amount, to, float64(rateRow.Rate.Exp)), nil
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

	return app.convert(amount, to, float64(rateRow.Rate.Exp)), nil
}

func (app *App) UpdateRates() error {
	updatedPairs := make(map[string]map[string]bool)
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
			scanErr := pgRate.Scan(rate.Rate)
			if scanErr != nil {
				return fmt.Errorf("convert rate: %w", scanErr)
			}
			updatedPairs[rate.From.Code][rate.To.Code] = true
			transaction, err := app.conn.BeginTx(
				app.ctx,
				pgx.TxOptions{IsoLevel: pgx.ReadCommitted},
			)
			if err != nil {
				return fmt.Errorf("start transaction: %w", err)
			}

			transactionRepository := app.repository.WithTx(transaction)
			transactionRepository.UpdateRate(
				app.ctx,
				repository.UpdateRateParams{
					CurrencyFrom: rate.From.Code,
					CurrencyTo:   rate.To.Code,
					Rate:         pgRate,
				},
			)
			transactionRepository.ArchiveRate(
				app.ctx,
				repository.ArchiveRateParams{
					CurrencyFrom: rate.From.Code,
					CurrencyTo:   rate.To.Code,
				},
			)
			transaction.Commit(app.ctx)
		}
	}

	return nil
}