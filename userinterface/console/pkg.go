package console

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/kotoproger/exchange/app"
	"github.com/kotoproger/exchange/userinterface"
)

const (
	HELP string = ""
)

type Console struct {
	app app.Exchanger
	in  io.Reader
	out io.Writer
}

func NewConsole(app app.Exchanger, in io.Reader, out io.Writer) *Console {
	return &Console{app: app, in: in, out: out}
}

func (c *Console) Run() {
	for {
		args, readerr := c.readCommand()
		if readerr != nil {
			c.printError(readerr)
			return
		}
		switch args[0] {
		case string(userinterface.EXIT):
			return
		case string(userinterface.HELP):
			c.printHelp()
		case string(userinterface.UPDATE):
			err := c.app.UpdateRates()
			if err == nil {
				c.print("Update successfully")
			} else {
				c.printError(err)
			}
		case string(userinterface.EXCHANGE):
			currencyFrom := money.GetCurrency(args[2])
			if currencyFrom == nil {
				c.print(fmt.Sprintf("wrong currency code `%s`", args[2]))
				continue
			}
			currencyTo := money.GetCurrency(args[3])
			if currencyTo == nil {
				c.print(fmt.Sprintf("wrong currency code `%s`", args[3]))
				continue
			}
			floatValue, tofloatError := strconv.ParseFloat(args[1], 64)
			if tofloatError != nil {
				c.printError(tofloatError)
				continue
			}
			amount := money.New(
				int64(math.Round(floatValue*math.Pow(10, float64(currencyFrom.Fraction)))),
				currencyFrom.Code,
			)
			var result *money.Money
			var exchangeError error
			if args[4] == "" {
				result, exchangeError = c.app.Exchange(amount, currencyTo)
			} else {
				dt, tperr := time.Parse(time.RFC3339, args[4])
				if tperr != nil {
					c.printError(tperr)
					continue
				}
				result, exchangeError = c.app.ExchangeToDate(amount, currencyTo, dt)
			}
			if exchangeError != nil {
				c.printError(fmt.Errorf("date time parse error: %w", exchangeError))
				continue
			} else if result == nil {
				c.print("cant find rate")
				continue
			}
			fmt.Println(result)
			c.print(fmt.Sprintf("%s -> %s", amount.Display(), result.Display()))
		default:
			c.print(fmt.Sprintf("Unknown command `%s`", args[0]))
		}
	}
}

func (c *Console) printHelp() {
	c.print(HELP)
}

func (c *Console) print(str string) {
	fmt.Fprint(c.out, "< ")
	fmt.Fprintln(c.out, str)
}

func (c *Console) readCommand() ([]string, error) {
	fmt.Fprint(c.out, "> ")

	comandName, amount, currencyFrom, currencyTo, date := "", "", "", "", ""
	count, inputErr := fmt.Fscanln(c.in, &comandName, &amount, &currencyFrom, &currencyTo, &date)
	if count == 0 && inputErr != nil {
		return []string{comandName}, fmt.Errorf("scan input: %w", inputErr)
	}

	return []string{
		comandName, amount, currencyFrom, currencyTo, date,
	}, nil
}

func (c *Console) printError(err error) {
	fmt.Fprint(c.out, "< ")
	fmt.Fprintln(c.out, err)
}
