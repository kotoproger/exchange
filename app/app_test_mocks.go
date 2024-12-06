package app

import (
	"context"

	"github.com/kotoproger/exchange/internal/repository"
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
