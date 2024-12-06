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
	returnrow, _ := args.Get(0).(money.Money)
	return &returnrow, args.Error(1)
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
	testCases := []struct {
		name        string
		input       string
		output      string
		updateRates []any
		exchanges   []struct {
			input  []any
			output []any
		}
	}{
		{
			name:        "only exit",
			input:       "exit\n",
			output:      "> ",
			updateRates: []any{},
		},
		{
			name:        "unknown command + exit",
			input:       "sfdasf\nexit\n",
			output:      "> < Unknown command `sfdasf`\n> ",
			updateRates: []any{},
		},
		{
			name:        "update rates + exit",
			input:       "update\nexit\n",
			output:      "> < Update successfully\n> ",
			updateRates: []any{nil},
		},
		{
			name:        "double update rates + exit",
			input:       "update\nupdate\nexit\n",
			output:      "> < Update successfully\n> < Update successfully\n> ",
			updateRates: []any{nil, nil},
		},
		{
			name:        "errorneus update + exit",
			input:       "update\nexit\n",
			output:      "> < some update error\n> ",
			updateRates: []any{fmt.Errorf("some update error")},
		},
		{
			name:        "exchange wrong from surrency + exit",
			input:       "exchange 100 aas usd\nexit\n",
			output:      "> < wrong currency code `aas`\n> ",
			updateRates: []any{},
		},
		{
			name:        "exchange wrong to surrency + exit",
			input:       "exchange 100 rub aas\nexit\n",
			output:      "> < wrong currency code `aas`\n> ",
			updateRates: []any{},
		},
		{
			name:        "exchange float parsing error + exit",
			input:       "exchange asd rub usd\nexit\n",
			output:      "> < strconv.ParseFloat: parsing \"asd\": invalid syntax\n> ",
			updateRates: []any{},
		},
		{
			name:  "successfully exchange + exit",
			input: "exchange 100 rub usd\nexit\n",
			output: fmt.Sprintf(
				"> < %s -> %s\n> ",
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

			console := NewConsole(&mock, &in, &out)
			console.Run()
			assert.Equal(t, testCase.output, out.String())
			mock.AssertExpectations(t)
		})
	}
}
