package app

import (
	"context"

	"github.com/kotoproger/exchange/internal/repository"
	"github.com/kotoproger/exchange/internal/source"
	"github.com/stretchr/testify/mock"
)

type MockQueries struct {
	mock.Mock
}

func (q MockQueries) ArchiveRate(ctx context.Context, arg repository.ArchiveRateParams) error {
	args := q.Called(ctx, arg)
	return args.Error(0)
}

func (q MockQueries) GetCuurentRate(ctx context.Context, arg repository.GetCuurentRateParams) (*repository.GetCuurentRateRow, error) {
	args := q.Called(ctx, arg)
	returnrow, _ := args.Get(0).(repository.GetCuurentRateRow)
	return &returnrow, args.Error(1)
}

func (q MockQueries) GetRateOnDate(ctx context.Context, arg repository.GetRateOnDateParams) (*repository.GetRateOnDateRow, error) {
	args := q.Called(ctx, arg)
	returnrow, _ := args.Get(0).(repository.GetRateOnDateRow)
	return &returnrow, args.Error(1)
}

func (q MockQueries) UpdateRate(ctx context.Context, arg repository.UpdateRateParams) error {
	args := q.Called(ctx, arg)
	return args.Error(0)
}

type MockWrapper struct {
	mock.Mock
}

func (m MockWrapper) GetRepository(ctx context.Context) (repo repository.Querier, commit func() error, rollback func(), release func(), err error) {
	args := m.Called(ctx)
	argRepo := args.Get(0)
	repo, _ = argRepo.(repository.Querier)
	argCommit := args.Get(1)
	commit, _ = argCommit.(func() error)
	argRollback := args.Get(2)
	rollback, _ = argRollback.(func())
	argRelease := args.Get(3)
	release, _ = argRelease.(func())
	err = args.Error(4)
	return
}

type MockFunc struct {
	mock.Mock
}

func (m MockFunc) call() {
	m.Called()
}

func (m MockFunc) callError() error {
	args := m.Called()
	return args.Error(0)
}

type MockExchangeSource struct {
	mock.Mock
}

func (m MockExchangeSource) Get() <-chan source.ExchangeRate {
	args := m.Called()

	list := []source.ExchangeRate{}
	for i := 1; i <= args.Int(0); i++ {
		sourceItem := args.Get(i).(source.ExchangeRate)
		list = append(list, sourceItem)
	}

	chanel := make(chan source.ExchangeRate, len(list))
	for _, item := range list {
		chanel <- item
	}

	close(chanel)

	return chanel
}
