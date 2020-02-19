package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
	"github.com/minheq/kedul_server_main/logger"
	"github.com/minheq/kedul_server_main/phone"
	"github.com/sirupsen/logrus"
)

func main() {
	router := chi.NewRouter()
	log := logger.NewLogger()
	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)
	smsSender := phone.NewSMSSender()

	if err != nil {
		log.WithFields(logrus.Fields{
			"DATABASE_URL": dbURL,
			"error":        err.Error(),
		}).Fatal("error opening database")
	}

	server := newServer(db, router, log, smsSender)

	fmt.Println("Server listening at localhost:4000")

	http.ListenAndServe(":4000", server.router)
}
