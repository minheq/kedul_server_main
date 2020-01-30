package handlers

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/minheq/kedul_server_main/errors"
	"github.com/minheq/kedul_server_main/models"
)

type loginVerifyCheckRequest struct {
	ClientState string `json:"clientState"`
	Code        string `json:"code"`
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
func HandleLoginVerifyCheck(store models.Store, tokenAuth *jwtauth.JWTAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &loginVerifyCheckRequest{}

		if err := render.Bind(r, data); err != nil {
			_ = render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		accessToken, err := LoginVerifyCheck(data.ClientState, data.Code, store, tokenAuth)

		if err != nil {
			_ = render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		render.Render(w, r, &loginVerifyCheckResponse{AccessToken: accessToken})
	}
}

// LoginVerifyCheck returns accessToken given the clientState and code match the persisted verification code
func LoginVerifyCheck(clientState string, code string, store models.Store, tokenAuth *jwtauth.JWTAuth) (string, error) {
	const op errors.Op = "handlers/login_verify_check.LoginVerifyCheck"

	verificationCode, err := store.GetVerificationCodeByClientStateAndCode(clientState, code)

	if err != nil {
		return "", errors.E(op, err, "verification code not found")
	}

	if verificationCode.ExpiredAt.Before(time.Now()) {
		return "", errors.E(op, "verification code expired", errors.Forbidden)
	}

	account, err := store.GetAccountByID(verificationCode.AccountID)

	if err != nil {
		return "", errors.E(op, err, "account not found")
	}

	account.IsPhoneNumberVerified = true
	account.UpdatedAt = time.Now()

	err = store.UpdateAccount(account)

	if err != nil {
		return "", errors.E(op, err, "could not update account")
	}

	err = store.RemoveVerificationCodeByID(verificationCode.ID)

	if err != nil {
		return "", errors.E(op, err, "could not remove verification code")
	}

	_, accessToken, err := tokenAuth.Encode(jwt.MapClaims{"account_id": account.ID})

	if err != nil {
		return "", errors.E(op, err, "failed to generate access token")
	}

	return accessToken, nil
}
