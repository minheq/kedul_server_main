package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/minheq/kedul_server_main/app"
	"github.com/minheq/kedul_server_main/logger"
)

type smsSenderMock struct {
	Text string
}

func (s *smsSenderMock) SendSMS(phoneNumber string, countryCode string, text string) error {
	s.Text = text
	return nil
}

var accessToken = ""

func post(target string, body interface{}, response interface{}, server *server) error {
	var b bytes.Buffer
	json.NewEncoder(&b).Encode(body)

	req := httptest.NewRequest("POST", target, &b)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		return fmt.Errorf("request error. received %v status code. response body: %v", w.Code, w.Body)
	}

	json.NewDecoder(w.Body).Decode(response)

	return nil
}

func setupDB(db *sql.DB) {
	databaseName := os.Getenv("DATABASE_NAME")
	driver, _ := postgres.WithInstance(db, &postgres.Config{})
	migrations, err := migrate.NewWithDatabaseInstance("file://migrations", databaseName, driver)

	if err != nil {
		log.Fatalf("error instantiating migrations: %v", err)
	}

	err = migrations.Down()

	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("error running migrations down, %v", err)
	}

	err = migrations.Up()

	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("error running migrations up, %v", err)
	}
}

func TestIntegration(t *testing.T) {
	router := chi.NewRouter()
	log := logger.NewLogger()
	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		t.Error(err)
	}

	setupDB(db)

	smsSender := &smsSenderMock{}

	server := newServer(db, router, log, smsSender)

	loginVerifyResp := &phoneNumberVerifyResponse{}

	t.Run("login verify", func(t *testing.T) {
		body := phoneNumberVerifyRequest{
			PhoneNumber: "999999999",
			CountryCode: "VN",
		}

		err := post("/auth/login_verify", body, loginVerifyResp, server)

		if err != nil {
			t.Error(err)
			return
		}
	})

	t.Run("login check", func(t *testing.T) {
		body := phoneNumberCheckRequest{
			VerificationID: loginVerifyResp.VerificationID,
			Code:           smsSender.Text,
		}

		loginCheckResp := &loginCheckResponse{}

		err := post("/auth/login_check", body, loginCheckResp, server)

		accessToken = loginCheckResp.AccessToken

		if err != nil {
			t.Error(err)
			return
		}
	})

	business := &app.Business{}

	t.Run("create business", func(t *testing.T) {
		body := createBusinessRequest{
			Name: "my business",
		}

		err := post("/businesses", body, business, server)

		if err != nil {
			t.Error(err)
			return
		}
	})

	fmt.Printf("%v", business)
}
