package main

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/minheq/kedul_server_main/auth"
	"github.com/minheq/kedul_server_main/errors"
)

type loginVerifyRequest struct {
	PhoneNumber string `json:"phone_number"`
	CountryCode string `json:"country_code"`
}

func (p *loginVerifyRequest) Bind(r *http.Request) error {
	return nil
}

type loginVerifyResponse struct {
	VerificationID string `json:"verification_id"`
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

		verificationID, err := authService.LoginVerify(r.Context(), data.PhoneNumber, data.CountryCode)

		if err != nil {
			_ = render.Render(w, r, errors.NewErrResponse(err))
			return
		}

		render.Render(w, r, &loginVerifyResponse{VerificationID: verificationID})
	}
}

type loginVerifyCheckRequest struct {
	VerificationID string `json:"verification_id"`
	Code           string `json:"code"`
}

func (p *loginVerifyCheckRequest) Bind(r *http.Request) error {
	return nil
}

type loginVerifyCheckResponse struct {
	AccessToken string `json:"access_token"`
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

		accessToken, err := authService.LoginVerifyCheck(r.Context(), data.VerificationID, data.Code)

		if err != nil {
			_ = render.Render(w, r, errors.NewErrResponse(err))
			return
		}

		render.Render(w, r, &loginVerifyCheckResponse{AccessToken: accessToken})
	}
}

type userResponse struct {
	ID                    string    `json:"id"`
	FullName              string    `json:"full_name"`
	PhoneNumber           string    `json:"phone_number"`
	CountryCode           string    `json:"country_code"`
	IsPhoneNumberVerified bool      `json:"is_phone_number_verified"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

func newUserResponse(user *auth.User) *userResponse {
	return &userResponse{
		ID:                    user.ID,
		FullName:              user.FullName,
		PhoneNumber:           user.PhoneNumber,
		CountryCode:           user.CountryCode,
		IsPhoneNumberVerified: user.IsPhoneNumberVerified,
		CreatedAt:             user.CreatedAt,
		UpdatedAt:             user.UpdatedAt,
	}
}

func (rd *userResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// HandleGetCurrentUser get currently authenticated user
func HandleGetCurrentUser(authService auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := authService.GetCurrentUser(r.Context())

		if err != nil {
			_ = render.Render(w, r, errors.NewErrResponse(err))
			return
		}

		render.Render(w, r, newUserResponse(user))
	}
}
