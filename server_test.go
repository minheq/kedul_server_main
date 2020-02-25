package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi"
	"github.com/minheq/kedul_server_main/logger"
)

type smsSenderMock struct {
	Text string
}

func (s *smsSenderMock) SendSMS(phoneNumber string, countryCode string, text string) error {
	s.Text = text
	return nil
}

func post(target string, body interface{}, response interface{}, server *server) error {
	var b bytes.Buffer
	json.NewEncoder(&b).Encode(body)

	req := httptest.NewRequest("POST", target, &b)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		return fmt.Errorf("request error. received %v status code. response body: %v", w.Code, w.Body)
	}

	json.NewDecoder(w.Body).Decode(response)

	return nil
}

func TestIntegration(t *testing.T) {
	router := chi.NewRouter()
	log := logger.NewLogger()
	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		t.Error(err)
	}

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

	loginCheckResp := &loginCheckResponse{}

	t.Run("login check", func(t *testing.T) {
		body := phoneNumberCheckRequest{
			VerificationID: loginVerifyResp.VerificationID,
			Code:           smsSender.Text,
		}

		err := post("/auth/login_check", body, loginCheckResp, server)

		if err != nil {
			t.Error(err)
			return
		}
	})

	fmt.Printf("%v", loginCheckResp)
}
