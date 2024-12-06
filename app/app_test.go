package app

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kotoproger/exchange/internal/repository"
	"github.com/kotoproger/exchange/internal/source"
	"github.com/stretchr/testify/assert"
)

func TestExchangeSuccessfuly(t *testing.T) {
	var pgRate pgtype.Numeric
	pgRate.Scan("10.0")
	repositoryMock := MockQueries{}
	ctx := context.Background()
	repositoryMock.
		On("GetCuurentRate", ctx, repository.GetCuurentRateParams{CurrencyFrom: "RUB", CurrencyTo: "USD"}).
		Return(repository.GetCuurentRateRow{CurrencyFrom: "RUB", CurrencyTo: "USD", Rate: pgRate}, nil).
		Once()

	app := App{
		repository:  repositoryMock,
		ctx:         ctx,
		rateSources: []source.ExchangeSource{},
		pool:        &pgxpool.Pool{},
	}

	actualMoney, err := app.Exchange(money.New(100, "rub"), money.GetCurrency("usd"))
	assert.Nil(t, err)
	assert.Equal(t, money.New(1000, "usd"), actualMoney)
	repositoryMock.AssertExpectations(t)
}

func TestExchangeErrorneus(t *testing.T) {
	var pgRate pgtype.Numeric
	pgRate.Scan("10.0")
	repositoryMock := MockQueries{}
	ctx := context.Background()
	repositoryMock.
		On("GetCuurentRate", ctx, repository.GetCuurentRateParams{CurrencyFrom: "RUB", CurrencyTo: "USD"}).
		Return(nil, fmt.Errorf("some error")).
		Once()

	app := App{
		repository:  repositoryMock,
		ctx:         ctx,
		rateSources: []source.ExchangeSource{},
		pool:        &pgxpool.Pool{},
	}

	actualMoney, err := app.Exchange(money.New(100, "rub"), money.GetCurrency("usd"))
	assert.Nil(t, actualMoney)
	assert.NotNil(t, err)
	repositoryMock.AssertExpectations(t)
}

func TestExchangeToDateErrorneus(t *testing.T) {
	var pgRate pgtype.Numeric
	pgRate.Scan("10.0")
	repositoryMock := MockQueries{}
	ctx := context.Background()
	time := time.Now()
	pgtime := pgtype.Timestamptz{}
	pgtime.Scan(time)
	repositoryMock.
		On("GetRateOnDate", ctx, repository.GetRateOnDateParams{CurrencyFrom: "RUB", CurrencyTo: "USD", CreatedAt: pgtime}).
		Return(nil, fmt.Errorf("some error")).
		Once()

	app := App{
		repository:  repositoryMock,
		ctx:         ctx,
		rateSources: []source.ExchangeSource{},
		pool:        &pgxpool.Pool{},
	}

	actualMoney, err := app.ExchangeToDate(money.New(100, "rub"), money.GetCurrency("usd"), time)
	assert.Nil(t, actualMoney)
	assert.NotNil(t, err)
	repositoryMock.AssertExpectations(t)
}

func TestExchangeToDateSuccessfuly(t *testing.T) {
	var pgRate pgtype.Numeric
	pgRate.Scan("10.0")
	repositoryMock := MockQueries{}
	ctx := context.Background()
	time := time.Now()
	pgtime := pgtype.Timestamptz{}
	pgtime.Scan(time)
	repositoryMock.
		On("GetRateOnDate", ctx, repository.GetRateOnDateParams{CurrencyFrom: "RUB", CurrencyTo: "USD", CreatedAt: pgtime}).
		Return(repository.GetRateOnDateRow{CurrencyFrom: "RUB", CurrencyTo: "USD", Rate: pgRate}, nil).
		Once()

	app := App{
		repository:  repositoryMock,
		ctx:         ctx,
		rateSources: []source.ExchangeSource{},
		pool:        &pgxpool.Pool{},
	}

	actualMoney, err := app.ExchangeToDate(money.New(100, "rub"), money.GetCurrency("usd"), time)
	assert.Nil(t, err)
	assert.Equal(t, money.New(1000, "usd"), actualMoney)
	repositoryMock.AssertExpectations(t)
}
