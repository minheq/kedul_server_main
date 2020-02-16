package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/minheq/kedul_server_main/auth"
	"github.com/minheq/kedul_server_main/errors"
)

type loginVerifyRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	CountryCode string `json:"countryCode"`
}

func (p *loginVerifyRequest) Bind(r *http.Request) error {
	return nil
}

type loginVerifyResponse struct {
	VerificationID string `json:"verificationID"`
}

func (rd *loginVerifyResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// HandleLoginVerify handles login verification initialization
func HandleLoginVerify(authService auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &loginVerifyRequest{}

		if err := render.Bind(r, data); err != nil {
			_ = render.Render(w, r, errors.NewErrResponse(err))
			return
		}

		state, err := authService.LoginVerify(data.PhoneNumber, data.CountryCode)

		if err != nil {
			fmt.Println(err)
			_ = render.Render(w, r, errors.NewErrResponse(err))
			return
		}

		render.Render(w, r, &loginVerifyResponse{VerificationID: state})
	}
}

type loginVerifyCheckRequest struct {
	VerificationID string `json:"verificationID"`
	Code           string `json:"code"`
}

func (p *loginVerifyCheckRequest) Bind(r *http.Request) error {
	return nil
}

type loginVerifyCheckResponse struct {
	AccessToken string `json:"accessToken"`
}

func (rd *loginVerifyCheckResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// HandleLoginVerifyCheck handles login verification checking
func HandleLoginVerifyCheck(authService auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &loginVerifyCheckRequest{}

		if err := render.Bind(r, data); err != nil {
			_ = render.Render(w, r, errors.NewErrResponse(err))
			return
		}

		accessToken, err := authService.LoginVerifyCheck(data.VerificationID, data.Code)

		if err != nil {
			_ = render.Render(w, r, errors.NewErrResponse(err))
			return
		}

		render.Render(w, r, &loginVerifyCheckResponse{AccessToken: accessToken})
	}
}
