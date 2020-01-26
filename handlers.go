package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/graphql-go/handler"
)

func (s server) handleLoginStart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			PhoneNumber string `json:"phoneNumber"`
			CountryCode string `json:"countryCode"`
		}

		err := s.Decode(r, &input)

		if err != nil {
			s.logger.Error("failed to decode: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		state, err := s.accountService.LogInStart(input.PhoneNumber, input.CountryCode)

		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		type payload struct {
			State string `json:"state"`
		}

		s.Respond(w, r, &payload{State: state}, http.StatusAccepted)
	}
}

func (s server) handleLoginVerify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			ClientState string `json:"clientState"`
			Code        string `json:"code"`
		}

		err := json.NewDecoder(r.Body).Decode(&input)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		state, err := s.accountService.LogInVerify(input.ClientState, input.Code)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Fprint(w, state)
	}
}

func (s server) handleGraphQL() http.HandlerFunc {

	gql := &gql{
		accountService:  s.accountService,
		businessService: s.businessService,
	}

	schema := gql.schema()

	graphQLHandler := handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		GraphiQL:   false,
		Playground: true,
	})

	return func(w http.ResponseWriter, r *http.Request) {
		graphQLHandler.ServeHTTP(w, r)
	}
}
