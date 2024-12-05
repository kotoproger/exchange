package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/kotoproger/exchange/app"
	"github.com/kotoproger/exchange/internal/source"
	"github.com/kotoproger/exchange/internal/source/cbr"
)

func main() {
	godotenv.Load()

	connURL, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		panic("cant get connection url")
	}
	pool, err := pgxpool.New(context.Background(), connURL)
	if err != nil {
		panic("cant create connection pool")
	}
	app := app.NewApp(
		context.Background(),
		[]source.ExchangeSource{
			cbr.NewCbr("https://www.cbr-xml-daily.ru/daily_json.js"),
		},
		pool,
	)

	fmt.Println(app.UpdateRates())
}
