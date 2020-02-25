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

type testHTTPClient struct {
	accessToken string
	server      *server
}

func (t *testHTTPClient) post(target string, body interface{}, response interface{}) error {
	var b bytes.Buffer
	json.NewEncoder(&b).Encode(body)

	req := httptest.NewRequest("POST", target, &b)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.accessToken))

	w := httptest.NewRecorder()
	t.server.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		return fmt.Errorf("request error. received %v status code. response body: %v", w.Code, w.Body)
	}

	json.NewDecoder(w.Body).Decode(response)

	return nil
}

func (t *testHTTPClient) get(target string, response interface{}) error {
	req := httptest.NewRequest("GET", target, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.accessToken))

	w := httptest.NewRecorder()
	t.server.router.ServeHTTP(w, req)

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

func TestEndToEnd(t *testing.T) {
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

	// Auth
	client := &testHTTPClient{server: server, accessToken: ""}

	t.Run("login verify", func(t *testing.T) {
		body := phoneNumberVerifyRequest{
			PhoneNumber: "999999999",
			CountryCode: "VN",
		}

		err := client.post("/auth/login_verify", body, loginVerifyResp)

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

		err := client.post("/auth/login_check", body, loginCheckResp)

		client.accessToken = loginCheckResp.AccessToken

		if err != nil {
			t.Error(err)
			return
		}
	})

	// App
	business := &app.Business{}
	location := &app.Location{}

	t.Run("create business", func(t *testing.T) {
		body := createBusinessRequest{
			Name: "my business",
		}

		err := client.post("/businesses", body, business)

		if err != nil {
			t.Error(err)
			return
		}
	})

	t.Run("update business", func(t *testing.T) {
		body := updateBusinessRequest{
			Name: "better business",
		}

		err := client.post(fmt.Sprintf("/businesses/%s", business.ID), body, business)

		if err != nil {
			t.Error(err)
			return
		}
	})

	t.Run("get business", func(t *testing.T) {
		resp := &app.Business{}
		err := client.get(fmt.Sprintf("/businesses/%s", business.ID), resp)

		if err != nil {
			t.Error(err)
			return
		}

		if resp.Name != business.Name {
			t.Error(fmt.Errorf("business name does not match. expected=%s, received=%s", business.Name, resp.Name))
		}
	})

	t.Run("create location", func(t *testing.T) {
		body := createLocationRequest{
			BusinessID: business.ID,
			Name:       "my location",
		}

		err := client.post("/locations", body, location)

		if err != nil {
			t.Error(err)
			return
		}
	})

	t.Run("update location", func(t *testing.T) {
		body := updateLocationRequest{
			Name: "better location",
		}

		err := client.post(fmt.Sprintf("/locations/%s", location.ID), body, location)

		if err != nil {
			t.Error(err)
			return
		}
	})

	t.Run("get location", func(t *testing.T) {
		resp := &app.Location{}
		err := client.get(fmt.Sprintf("/locations/%s", location.ID), resp)

		if err != nil {
			t.Error(err)
			return
		}

		if resp.Name != location.Name {
			t.Error(fmt.Errorf("location name does not match. expected=%s, received=%s", location.Name, resp.Name))
		}
	})
}
