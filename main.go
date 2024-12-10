package main

import (
	"context"
	"log"
	"os"
	"sync"

	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/kotoproger/exchange/app"
	"github.com/kotoproger/exchange/internal/repository"
	"github.com/kotoproger/exchange/internal/repositorywrapper"
	"github.com/kotoproger/exchange/internal/source"
	"github.com/kotoproger/exchange/internal/source/cbr"
	"github.com/kotoproger/exchange/userinterface/console"
	"github.com/pressly/goose/v3"
)

func main() {

	godotenv.Load()

	connURL, ok := os.LookupEnv("APP_DATABASE_URL")
	if !ok {
		panic("cant get connection url")
	}

	pool, err := pgxpool.New(context.Background(), connURL)
	if err != nil {
		panic("cant create connection pool")
	}

	goose.SetDialect("postgres")

	db := stdlib.OpenDBFromPool(pool)
	if db == nil {
		panic("cannot open db")
	}

	println("Migrating")

	err = goose.Up(db, "./sql/migrations")
	if err != nil {
		log.Fatalf(err.Error())
	}
	println("Done")

	sourceURL, ok := os.LookupEnv("APP_CB_API")
	if !ok {
		panic("cant get exchange source url")
	}

	app := app.NewApp(
		context.Background(),
		[]source.ExchangeSource{
			cbr.NewCbr(sourceURL),
		},
		&repositorywrapper.Wrapper{
			Pool: pool,
			Repo: &repository.Queries{},
		},
	)

	var wg sync.WaitGroup
	controller := console.NewConsole(app, os.Stdin, os.Stdout)
	wg.Add(1)
	go func(wg *sync.WaitGroup, controller *console.Console) {
		defer wg.Done()
		controller.Run()
	}(&wg, controller)

	wg.Wait()
}
