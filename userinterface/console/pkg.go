package console

import (
	"fmt"
	"math"
	"strconv"

	"github.com/Rhymond/go-money"
	"github.com/kotoproger/exchange/app"
	"github.com/kotoproger/exchange/userinterface"
)

const (
	HELP string = ""
)

type Console struct {
	app app.Exchanger
}

func NewConsole(app app.Exchanger) *Console {
	return &Console{app: app}
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
			c.app.UpdateRates()
			c.print("Update successfully")
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
			}
			amount := money.New(
				int64(math.Round(floatValue*math.Pow(10, float64(currencyFrom.Fraction)))),
				currencyFrom.Code,
			)
			result, exchangeError := c.app.Exchange(amount, currencyTo)
			if exchangeError != nil {
				c.printError(exchangeError)
				continue
			}
			c.print(result.Display())
		default:
			c.print(fmt.Sprintf("Unknown command `%s`", args[0]))
		}
	}
}

func (c *Console) printHelp() {
	c.print(HELP)
}

func (c *Console) print(str string) {
	fmt.Print("< ")
	fmt.Println(str)
}

func (c *Console) readCommand() ([]string, error) {
	fmt.Print("> ")
	comandName, amount, currencyFrom, currencyTo, input := "", "", "", "", ""
	count, inputErr := fmt.Scanln(&comandName, &amount, &currencyFrom, &currencyTo)
	if count == 0 && inputErr != nil {
		return []string{input}, fmt.Errorf("scan input: %w", inputErr)
	}

	return []string{
		comandName, amount, currencyFrom, currencyTo,
	}, nil
}

func (c *Console) printError(err error) {
	fmt.Print("< ")
	fmt.Println(err)
}
