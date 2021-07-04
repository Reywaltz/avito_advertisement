package main

import (
	"net/http"
	"os"

	"github.com/Reywaltz/avito_advertising/cmd/advert-api/handlers"
	"github.com/Reywaltz/avito_advertising/internal/repositories"
	"github.com/Reywaltz/avito_advertising/pkg/log"
	"github.com/Reywaltz/avito_advertising/pkg/postgres"
	"github.com/gorilla/mux"
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

	repo := repositories.NewRepo(db)
	adHandlers, err := handlers.NewHandlers(log, repo)
	if err != nil {
		log.Fatal("Can't create Ads handlers: %s", err)
	}

	router := mux.NewRouter()
	router.StrictSlash(true)
	adHandlers.Routes(router)

	http.ListenAndServe(":8000", router)
}
