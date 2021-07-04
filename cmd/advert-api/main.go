package main

import (
	"os"

	"github.com/Reywaltz/avito_advertising/pkg/log"
	"github.com/Reywaltz/avito_advertising/pkg/postgres"
)

func main() {
	log, err := log.NewLogger(os.Getenv("DEV"))
	if err != nil {
		panic("Can't init logger")
	}

	db, err := postgres.NewDB("postgres://advert_user:pass@localhost:5433/advert")
	if err != nil {
		log.Fatalf("Can't connect to db: %v", err)
	}
}
