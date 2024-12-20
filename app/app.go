package app

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kotoproger/exchange/internal/repository"
	"github.com/kotoproger/exchange/internal/repositorywrapper"
	"github.com/kotoproger/exchange/internal/source"
)

type App struct {
	ctx         context.Context
	rateSources []source.ExchangeSource
	repoPool    repositorywrapper.RepositoryPool
}

func NewApp(
	ctx context.Context,
	sources []source.ExchangeSource,
	repoPool repositorywrapper.RepositoryPool,
) *App {
	return &App{
		ctx:         ctx,
		rateSources: sources,
		repoPool:    repoPool,
	}
}

func (app *App) convert(amount *money.Money, to *money.Currency, rate float64) *money.Money {
	newAmount := int64(math.Round(float64(amount.Amount()) * rate))
	return money.New(newAmount, to.Code)
}

func (app *App) Exchange(amount *money.Money, to *money.Currency) (*money.Money, error) {
	if to.Code == amount.Currency().Code {
		return amount, nil
	}
	repo, commit, _, release, err := app.repoPool.GetRepository(app.ctx)
	if err != nil {
		return nil, fmt.Errorf("get repository: %w", err)
	}

	rateRow, repoError := repo.GetCuurentRate(
		app.ctx,
		repository.GetCuurentRateParams{CurrencyFrom: amount.Currency().Code, CurrencyTo: to.Code},
	)
	commit()
	release()

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
	if to.Code == amount.Currency().Code {
		return amount, nil
	}
	pgtime := pgtype.Timestamptz{}
	pgtime.Scan(date)
	repo, commit, _, release, err := app.repoPool.GetRepository(app.ctx)
	if err != nil {
		return nil, fmt.Errorf("get repository: %w", err)
	}

	rateRow, repoError := repo.GetRateOnDate(
		app.ctx,
		repository.GetRateOnDateParams{
			CurrencyFrom: amount.Currency().Code,
			CurrencyTo:   to.Code,
			CreatedAt:    pgtime,
		},
	)
	commit()
	release()

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
	for _, source := range app.rateSources {
		for rate := range source.Get() {
			_, ok2 := updatedPairs[rate.From.Code]
			if !ok2 {
				updatedPairs[rate.From.Code] = make(map[string]bool)
			}
			_, ok3 := updatedPairs[rate.From.Code][rate.To.Code]
			if ok3 {
				continue
			}
			var pgRate pgtype.Numeric
			scanErr := pgRate.Scan(strconv.FormatFloat(rate.Rate, 'f', -1, 64))
			if scanErr != nil {
				return fmt.Errorf("convert rate: %w", scanErr)
			}
			updatedPairs[rate.From.Code][rate.To.Code] = true

			repo, commit, rollback, release, err := app.repoPool.GetRepository(app.ctx)
			if err != nil {
				return fmt.Errorf("get repository: %w", err)
			}

			updateerr := repo.UpdateRate(
				app.ctx,
				repository.UpdateRateParams{
					CurrencyFrom: rate.From.Code,
					CurrencyTo:   rate.To.Code,
					Rate:         pgRate,
				},
			)

			if updateerr != nil {
				rollback()
				release()
				return fmt.Errorf("update rate: %w", updateerr)
			}

			archiveerr := repo.ArchiveRate(
				app.ctx,
				repository.ArchiveRateParams{
					CurrencyFrom: rate.From.Code,
					CurrencyTo:   rate.To.Code,
				},
			)
			if archiveerr != nil {
				rollback()
				release()
				return fmt.Errorf("update rate: %w", archiveerr)
			}

			commit()
			release()
		}
	}

	return nil
}
