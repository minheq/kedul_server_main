package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/minheq/kedulv2/service_salon/account"
	"github.com/minheq/kedulv2/service_salon/business"
	"github.com/minheq/kedulv2/service_salon/errors"
	"github.com/minheq/kedulv2/service_salon/sms"
	"github.com/sirupsen/logrus"
)

type server struct {
	db     *sql.DB
	router chi.Router
	logger *logrus.Logger

	accountStore account.Store

	accountService  account.Service
	businessService business.Service

	smsSender sms.Sender
}

func newServer(db *sql.DB, router chi.Router, logger *logrus.Logger, tokenAuth *jwtauth.JWTAuth) server {
	accountStore := account.NewAccountStore(db)

	accountService := account.NewAccountService(accountStore, sms.NewSender(), tokenAuth)
	businessService := business.NewBusinessService()

	s := server{
		db:     db,
		router: router,
		logger: logger,

		accountStore: accountStore,

		accountService:  accountService,
		businessService: businessService,
	}

	s.routes()

	return s
}

func (s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s server) Decode(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)

	return errors.E(err)
}

func (s server) Respond(w http.ResponseWriter, r *http.Request, v interface{}, status int) {
	json.NewEncoder(w).Encode(v)
}

func (s server) routes() {
	s.router.Post("/loginStart", s.handleLoginStart())
	s.router.Post("/loginVerify", s.handleLoginVerify())
	s.router.HandleFunc("/graphql", s.handleGraphQL())
}
