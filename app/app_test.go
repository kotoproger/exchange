package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kotoproger/exchange/internal/repository"
	"github.com/kotoproger/exchange/internal/source"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestApp(t *testing.T) { //nolint:funlen
	ctx := context.Background()
	SomeTime := time.Now()
	testCases := []struct {
		name           string
		run            func(app App) []any
		expectedResult []any
		querier        map[string]struct {
			args []any
			res  []any
		}
		resources []struct {
			res []any
		}
		repoPool    error
		release     bool
		commit      bool
		rollback    bool
		commitError error
	}{
		////////////////////////// exchange
		{
			name: "exchange rub -> rub",
			run: func(app App) []any {
				money, err := app.Exchange(
					money.New(10, "RUB"),
					money.GetCurrency("RUB"),
				)
				return []any{money, err}
			},
			expectedResult: []any{
				money.New(10, "RUB"),
				nil,
			},
			querier: make(map[string]struct {
				args []any
				res  []any
			}),
			resources: []struct {
				res []any
			}{},
		},
		{
			name: "exchange rub -> usd",
			run: func(app App) []any {
				money, err := app.Exchange(
					money.New(10, "RUB"),
					money.GetCurrency("USD"),
				)
				return []any{money, err}
			},
			expectedResult: []any{
				money.New(5, "USD"),
				nil,
			},
			querier: map[string]struct {
				args []any
				res  []any
			}{
				"GetCuurentRate": {
					args: []any{
						ctx,
						repository.GetCuurentRateParams{CurrencyFrom: "RUB", CurrencyTo: "USD"},
					},
					res: []any{
						repository.GetCuurentRateRow{CurrencyFrom: "RUB", CurrencyTo: "USD", Rate: getpgtype("0.5")},
						nil,
					},
				},
			},
			resources: []struct {
				res []any
			}{},
			release:  true,
			commit:   true,
			rollback: false,
		},
		{
			name: "exchange rub -> usd error on repository acquire",
			run: func(app App) []any {
				money, err := app.Exchange(
					money.New(10, "RUB"),
					money.GetCurrency("USD"),
				)
				return []any{money, err}
			},
			expectedResult: []any{
				nullMoney(),
				mock.Anything,
			},
			repoPool: errors.New("some error"),
			release:  false,
			commit:   false,
			rollback: false,
		},
		{
			name: "exchange rub -> usd erro on rate search",
			run: func(app App) []any {
				money, err := app.Exchange(
					money.New(10, "RUB"),
					money.GetCurrency("USD"),
				)
				return []any{money, err}
			},
			expectedResult: []any{
				nullMoney(),
				errors.New("some error"),
			},
			querier: map[string]struct {
				args []any
				res  []any
			}{
				"GetCuurentRate": {
					args: []any{
						ctx,
						repository.GetCuurentRateParams{CurrencyFrom: "RUB", CurrencyTo: "USD"},
					},
					res: []any{
						nil,
						errors.New("some error"),
					},
				},
			},
			resources: []struct {
				res []any
			}{},
			release:  true,
			commit:   true,
			rollback: false,
		},
		////////////////////// exchange on date
		{
			name: "exchange on date rub -> rub",
			run: func(app App) []any {
				money, err := app.ExchangeToDate(
					money.New(10, "RUB"),
					money.GetCurrency("RUB"),
					SomeTime,
				)
				return []any{money, err}
			},
			expectedResult: []any{
				money.New(10, "RUB"),
				nil,
			},
			querier: make(map[string]struct {
				args []any
				res  []any
			}),
			resources: []struct {
				res []any
			}{},
		},
		{
			name: "exchange on date rub -> usd",
			run: func(app App) []any {
				money, err := app.ExchangeToDate(
					money.New(10, "RUB"),
					money.GetCurrency("USD"),
					SomeTime,
				)
				return []any{money, err}
			},
			expectedResult: []any{
				money.New(5, "USD"),
				nil,
			},
			querier: map[string]struct {
				args []any
				res  []any
			}{
				"GetRateOnDate": {
					args: []any{
						ctx,
						repository.GetRateOnDateParams{CurrencyFrom: "RUB", CurrencyTo: "USD", CreatedAt: getpgdate(SomeTime)},
					},
					res: []any{
						repository.GetRateOnDateRow{CurrencyFrom: "RUB", CurrencyTo: "USD", Rate: getpgtype("0.5")},
						nil,
					},
				},
			},
			resources: []struct {
				res []any
			}{},
			release:  true,
			commit:   true,
			rollback: false,
		},
		{
			name: "exchange rub -> usd error on repository acquire",
			run: func(app App) []any {
				money, err := app.ExchangeToDate(
					money.New(10, "RUB"),
					money.GetCurrency("USD"),
					SomeTime,
				)
				return []any{money, err}
			},
			expectedResult: []any{
				nullMoney(),
				mock.Anything,
			},
			repoPool: errors.New("some error"),
			release:  false,
			commit:   false,
			rollback: false,
		},
		{
			name: "exchange rub -> usd erro on rate search",
			run: func(app App) []any {
				money, err := app.ExchangeToDate(
					money.New(10, "RUB"),
					money.GetCurrency("USD"),
					SomeTime,
				)
				return []any{money, err}
			},
			expectedResult: []any{
				nullMoney(),
				errors.New("some error"),
			},
			querier: map[string]struct {
				args []any
				res  []any
			}{
				"GetRateOnDate": {
					args: []any{
						ctx,
						repository.GetRateOnDateParams{CurrencyFrom: "RUB", CurrencyTo: "USD", CreatedAt: getpgdate(SomeTime)},
					},
					res: []any{
						nil,
						errors.New("some error"),
					},
				},
			},
			resources: []struct {
				res []any
			}{},
			release:  true,
			commit:   true,
			rollback: false,
		},
		//////////////////////////
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			repoPool := MockWrapper{}
			mocksource := MockExchangeSource{}
			mockQuerier := MockQueries{}
			mockCommit := MockFunc{}
			mockRollback := MockFunc{}
			mockRelease := MockFunc{}
			if testCase.release {
				mockRelease.On("call")
			}
			if testCase.rollback {
				mockRollback.On("call")
			}
			if testCase.commit {
				mockCommit.On("callError").Return(testCase.commitError)
			}

			for method, params := range testCase.querier {
				mockQuerier.On(method, params.args...).Return(params.res...)
				repoPool.On("GetRepository", ctx).Return(
					mockQuerier,
					func(m MockFunc) func() error {
						return func() error { //nolint:gocritic
							return m.callError()
						}
					}(mockCommit),
					func(m MockFunc) func() {
						return func() {
							m.call()
						}
					}(mockRollback),
					func(m MockFunc) func() {
						return func() {
							m.call()
						}
					}(mockRelease),
					nil,
				)
			}
			if testCase.repoPool != nil {
				repoPool.On("GetRepository", ctx).Return(
					nil,
					func(m MockFunc) func() error {
						return func() error { //nolint:gocritic
							return m.callError()
						}
					}(mockCommit),
					func(m MockFunc) func() {
						return func() {
							m.call()
						}
					}(mockRollback),
					func(m MockFunc) func() {
						return func() {
							m.call()
						}
					}(mockRelease),
					testCase.repoPool,
				)
			}

			app := App{
				ctx: ctx,
				rateSources: []source.ExchangeSource{
					&mocksource,
				},
				repoPool: repoPool,
			}

			res := testCase.run(app)

			assert.Equal(t, testCase.expectedResult[0], res[0])
			if testCase.expectedResult[1] != nil {
				assert.NotNil(t, res[1])
			}

			repoPool.AssertExpectations(t)
			mocksource.AssertExpectations(t)
			mockQuerier.AssertExpectations(t)
			mockRelease.AssertExpectations(t)
			mockCommit.AssertExpectations(t)
			mockRollback.AssertExpectations(t)
		})
	}
}

func getpgtype(value string) (ret pgtype.Numeric) {
	ret.Scan(value)
	return
}

func nullMoney() *money.Money {
	return nil
}

func getpgdate(someTime time.Time) (ret pgtype.Timestamptz) {
	ret.Scan(someTime)
	return
}
