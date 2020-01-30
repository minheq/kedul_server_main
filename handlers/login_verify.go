package handlers

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/minheq/kedul_server_main/errors"
	"github.com/minheq/kedul_server_main/models"
	"github.com/minheq/kedul_server_main/random"
	"github.com/minheq/kedul_server_main/sms"
	"github.com/nyaruka/phonenumbers"
)

type loginVerifyRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	CountryCode string `json:"countryCode"`
}

func (p *loginVerifyRequest) Bind(r *http.Request) error {
	return nil
}

type loginVerifyResponse struct {
	ClientState string `json:"clientState"`
}

func (rd *loginVerifyResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// HandleLoginVerify handles login verification initialization
func HandleLoginVerify(store models.Store, smsSender sms.Sender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &loginVerifyRequest{}

		if err := render.Bind(r, data); err != nil {
			_ = render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		state, err := LoginVerify(data.PhoneNumber, data.CountryCode, store, smsSender)

		if err != nil {
			_ = render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		render.Render(w, r, &loginVerifyResponse{ClientState: state})
	}
}

// LoginVerify login verification initialization core logic
func LoginVerify(phoneNumber string, countryCode string, store models.Store, smsSender sms.Sender) (string, error) {
	const op errors.Op = "handlers/login_verify.LoginVerify"

	parsedPhoneNumber, err := phonenumbers.Parse(phoneNumber, countryCode)

	if err != nil {
		return "", errors.E(op, err, "invalid phone number")
	}

	formattedPhoneNumber := phonenumbers.Format(parsedPhoneNumber, phonenumbers.NATIONAL)

	clientState := random.String(50)
	code := random.Number(6)

	account, err := store.GetAccountByPhoneNumber(formattedPhoneNumber, countryCode)

	if err != nil && errors.Is(errors.NotFound, err) == false {
		return "", errors.E(op, err, "account not found")
	}

	// Create and persist new account if it didn't exist yet
	if account == nil {
		account = models.NewAccount(formattedPhoneNumber, countryCode)

		err = store.StoreAccount(account)

		if err != nil {
			return "", errors.E(op, err, "could not store account")
		}
	}

	err = store.RemoveVerificationCodeByPhoneNumber(formattedPhoneNumber, countryCode)

	if err != nil {
		return "", errors.E(op, err, "could not remove verification code")
	}

	newVerificationCode := models.NewVerificationCode(clientState, code, account.ID, formattedPhoneNumber, countryCode, "LOGIN")

	err = store.StoreVerificationCode(newVerificationCode)

	if err != nil {
		return "", errors.E(op, err, "could not store verification code")
	}

	err = smsSender.SendSMS(formattedPhoneNumber, countryCode, code)

	if err != nil {
		return "", errors.E(op, err, "send sms failed")
	}

	return clientState, nil
}
