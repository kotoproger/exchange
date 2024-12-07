package console

import (
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockExchanger struct {
	mock.Mock
}

func (m *MockExchanger) Exchange(amount *money.Money, to *money.Currency) (*money.Money, error) {
	args := m.Called(amount, to)
	returnrow, _ := args.Get(0).(*money.Money)
	return returnrow, args.Error(1)
}

func (m *MockExchanger) ExchangeToDate(amount *money.Money, to *money.Currency, date time.Time) (*money.Money, error) {
	args := m.Called(amount, to, date)
	returnrow, _ := args.Get(0).(*money.Money)
	return returnrow, args.Error(1)
}

func (m *MockExchanger) UpdateRates() error {
	args := m.Called()
	return args.Error(0)
}

type TestStringReader struct {
	buffer []byte
}

func (b *TestStringReader) Read(p []byte) (n int, err error) {
	for index := range p {
		if len(b.buffer) > 0 {
			n++
			p[index] = b.buffer[0]
			b.buffer = b.buffer[1:]
		} else {
			err = io.EOF
			p[index] = 0
		}
	}

	return
}

func TestRun(t *testing.T) { 
	dt, _ := time.Parse(time.RFC3339, "2026-01-02T15:04:05Z")
	testCases := []struct {
		name        string
		input       string
		output      string
		updateRates []any
		exchanges   []struct {
			input  []any
			output []any
		}
		exchangesToDate []struct {
			input  []any
			output []any
		}
	}{
		{
			name:        "only exit",
			input:       "exit\n",
			output:      fmt.Sprintf("< %s\n> ", HELP),
			updateRates: []any{},
		},
		{
			name:        "only exit",
			input:       "exit\n",
			output:      fmt.Sprintf("< %s\n> ", HELP),
			updateRates: []any{},
		},
		{
			name:        "unknown command + exit",
			input:       "sfdasf\nexit\n",
			output:      fmt.Sprintf("< %s\n> < Unknown command `sfdasf`\n> ", HELP),
			updateRates: []any{},
		},
		{
			name:        "update rates + exit",
			input:       "update\nexit\n",
			output:      fmt.Sprintf("< %s\n> < Update successfully\n> ", HELP),
			updateRates: []any{nil},
		},
		{
			name:        "double update rates + exit",
			input:       "update\nupdate\nexit\n",
			output:      fmt.Sprintf("< %s\n> < Update successfully\n> < Update successfully\n> ", HELP),
			updateRates: []any{nil, nil},
		},
		{
			name:        "errorneus update + exit",
			input:       "update\nexit\n",
			output:      fmt.Sprintf("< %s\n> < some update error\n> ", HELP),
			updateRates: []any{fmt.Errorf("some update error")},
		},
		{
			name:        "exchange wrong from surrency + exit",
			input:       "exchange 100 aas usd\nexit\n",
			output:      fmt.Sprintf("< %s\n> < wrong currency code `aas`\n> ", HELP),
			updateRates: []any{},
		},
		{
			name:        "exchange wrong to surrency + exit",
			input:       "exchange 100 rub aas\nexit\n",
			output:      fmt.Sprintf("< %s\n> < wrong currency code `aas`\n> ", HELP),
			updateRates: []any{},
		},
		{
			name:        "exchange float parsing error + exit",
			input:       "exchange asd rub usd\nexit\n",
			output:      fmt.Sprintf("< %s\n> < strconv.ParseFloat: parsing \"asd\": invalid syntax\n> ", HELP),
			updateRates: []any{},
		},
		{
			name:  "successfully exchange + exit",
			input: "exchange 100 rub usd\nexit\n",
			output: fmt.Sprintf(
				"< %s\n> < %s -> %s\n> ",
				HELP,
				money.New(
					int64(10000),
					"rub",
				).Display(),
				money.New(
					int64(97),
					"usd",
				).Display(),
			),
			updateRates: []any{},
			exchanges: []struct {
				input  []any
				output []any
			}{
				{
					input: []any{
						money.New(
							int64(10000),
							"rub",
						),
						money.GetCurrency("usd"),
					},
					output: []any{
						money.New(
							int64(97),
							"usd",
						),
						nil,
					},
				},
			},
		},
		{
			name:        "not found rate + exit",
			input:       "exchange 100 rub usd\nexit\n",
			output:      fmt.Sprintf("< %s\n> < cant find rate\n> ", HELP),
			updateRates: []any{},
			exchanges: []struct {
				input  []any
				output []any
			}{
				{
					input: []any{
						money.New(
							int64(10000),
							"rub",
						),
						money.GetCurrency("usd"),
					},
					output: []any{
						nil,
						nil,
					},
				},
			},
		},
		{
			name:  "successfully exchange to date + exit",
			input: "exchange 100 rub usd 2026-01-02T15:04:05Z\nexit\n",
			output: fmt.Sprintf(
				"< %s\n> < %s -> %s\n> ",
				HELP,
				money.New(
					int64(10000),
					"rub",
				).Display(),
				money.New(
					int64(97),
					"usd",
				).Display(),
			),
			updateRates: []any{},
			exchangesToDate: []struct {
				input  []any
				output []any
			}{
				{
					input: []any{
						money.New(
							int64(10000),
							"rub",
						),
						money.GetCurrency("usd"),
						dt,
					},
					output: []any{
						money.New(
							int64(97),
							"usd",
						),
						nil,
					},
				},
			},
		},
		{
			name:        "date pare error + exit",
			input:       "exchange 100 rub usd 2026-01-02T1\nexit\n",
			output:      fmt.Sprintf("< %s\n> < parsing time \"2026-01-02T1\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"\" as \":\"\n> ", HELP),
			updateRates: []any{},
		},
		{
			name:        "not found rate to date + exit",
			input:       "exchange 100 rub usd 2026-01-02T15:04:05Z\nexit\n",
			output:      fmt.Sprintf("< %s\n> < cant find rate\n> ", HELP),
			updateRates: []any{},
			exchangesToDate: []struct {
				input  []any
				output []any
			}{
				{
					input: []any{
						money.New(
							int64(10000),
							"rub",
						),
						money.GetCurrency("usd"),
						dt,
					},
					output: []any{
						nil,
						nil,
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mock := MockExchanger{}
			out := strings.Builder{}
			in := TestStringReader{buffer: []byte(testCase.input)}

			for _, someResult := range testCase.updateRates {
				mock.On("UpdateRates").Return(someResult).Once()
			}
			for _, conf := range testCase.exchanges {
				mock.On("Exchange", conf.input...).Return(conf.output...)
			}
			for _, conf := range testCase.exchangesToDate {
				mock.On("ExchangeToDate", conf.input...).Return(conf.output...)
			}

			console := NewConsole(&mock, &in, &out)
			console.Run()
			assert.Equal(t, testCase.output, out.String())
			mock.AssertExpectations(t)
		})
	}
}
