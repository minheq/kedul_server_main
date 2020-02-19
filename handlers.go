package main

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/minheq/kedul_server_main/auth"
)

type phoneNumberVerifyRequest struct {
	PhoneNumber string `json:"phone_number"`
	CountryCode string `json:"country_code"`
}

func (p *phoneNumberVerifyRequest) Bind(r *http.Request) error {
	return nil
}

type phoneNumberVerifyResponse struct {
	VerificationID string `json:"verification_id"`
}

func (rd *phoneNumberVerifyResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *server) handleLoginVerify(authService auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &phoneNumberVerifyRequest{}

		if err := render.Bind(r, data); err != nil {
			s.respondError(w, r, err)
			return
		}

		verificationID, err := authService.LoginVerify(r.Context(), data.PhoneNumber, data.CountryCode)

		if err != nil {
			s.respondError(w, r, err)
			return
		}

		render.Render(w, r, &phoneNumberVerifyResponse{VerificationID: verificationID})
	}
}

type phoneNumberCheckRequest struct {
	VerificationID string `json:"verification_id"`
	Code           string `json:"code"`
}

func (p *phoneNumberCheckRequest) Bind(r *http.Request) error {
	return nil
}

type loginVerifyCheckResponse struct {
	AccessToken string `json:"access_token"`
}

func (rd *loginVerifyCheckResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *server) handleLoginCheck(authService auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &phoneNumberCheckRequest{}

		if err := render.Bind(r, data); err != nil {
			s.respondError(w, r, err)
			return
		}

		accessToken, err := authService.LoginCheck(r.Context(), data.VerificationID, data.Code)

		if err != nil {
			s.respondError(w, r, err)
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

func (s *server) handleGetCurrentUser(authService auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser, _ := r.Context().Value(userCtxKey).(*auth.User)

		render.Render(w, r, newUserResponse(currentUser))
	}
}

type updateUserProfileRequest struct {
	FullName       string `json:"full_name"`
	ProfileImageID string `json:"image_id"`
}

func (p *updateUserProfileRequest) Bind(r *http.Request) error {
	return nil
}

func (s *server) handleUpdateUserProfile(authService auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser, _ := r.Context().Value(userCtxKey).(*auth.User)
		data := &updateUserProfileRequest{}

		if err := render.Bind(r, data); err != nil {
			s.respondError(w, r, err)
			return
		}

		user, err := authService.UpdateUserProfile(r.Context(), data.FullName, data.ProfileImageID, currentUser)

		if err != nil {
			s.respondError(w, r, err)
			return
		}

		render.Render(w, r, newUserResponse(user))
	}
}

func (s *server) handleUpdatePhoneNumberVerify(authService auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser, _ := r.Context().Value(userCtxKey).(*auth.User)
		data := &phoneNumberVerifyRequest{}

		if err := render.Bind(r, data); err != nil {
			s.respondError(w, r, err)
			return
		}

		verificationID, err := authService.UpdatePhoneNumberVerify(r.Context(), data.PhoneNumber, data.CountryCode, currentUser)

		if err != nil {
			s.respondError(w, r, err)
			return
		}

		render.Render(w, r, &phoneNumberVerifyResponse{VerificationID: verificationID})
	}
}

func (s *server) handleUpdatePhoneNumberCheck(authService auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser, _ := r.Context().Value(userCtxKey).(*auth.User)
		data := &phoneNumberCheckRequest{}

		if err := render.Bind(r, data); err != nil {
			s.respondError(w, r, err)
			return
		}

		user, err := authService.UpdatePhoneNumberCheck(r.Context(), data.VerificationID, data.Code, currentUser)

		if err != nil {
			s.respondError(w, r, err)
			return
		}

		render.Render(w, r, newUserResponse(user))
	}
}
