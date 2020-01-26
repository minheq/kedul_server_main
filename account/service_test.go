package account

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/minheq/kedulv2/service_salon/errors"
)

type SMSSenderMock struct {
	text string
}

func (s *SMSSenderMock) SendSMS(phoneNumber string, countryCode string, text string) error {
	s.text = text
	return nil
}

var (
	db             *sql.DB
	accountStore   Store
	accountService Service
	smsSender      *SMSSenderMock
)

func TestMain(m *testing.M) {
	db, _ := sql.Open("postgres", "postgres://postgres@127.0.0.1:5432/kedul?sslmode=disable")
	driver, _ := postgres.WithInstance(db, &postgres.Config{})
	migrations, _ := migrate.NewWithDatabaseInstance("file://../migrations", "kedul", driver)

	migrations.Down()
	migrations.Up()

	smsSender = &SMSSenderMock{}
	accountStore = NewAccountStore(db)
	accountService = NewAccountService(accountStore, smsSender, jwtauth.New("HS256", []byte("secret"), nil))

	code := m.Run()

	db.Close()

	os.Exit(code)
}

func TestLogInHappyPath(t *testing.T) {
	var code string
	var clientState string
	var err error

	t.Run("should send code and return state when login start", func(t *testing.T) {
		clientState, err = accountService.LogInStart("999111333", "VN")
		code = smsSender.text

		if err != nil {
			t.Error(err)
		}

		if code == "" {
			t.Error("missing code")
		}
	})

	t.Run("should return access token when login verified", func(t *testing.T) {
		accessToken, err := accountService.LogInVerify(clientState, code)

		if err != nil {
			t.Error(err)
		}

		if accessToken == "" {
			t.Error("missing access token")
		}
	})
}

func TestLoginStartTwice(t *testing.T) {
	var codeOne string
	var clientStateOne string
	var codeTwo string
	var clientStateTwo string
	var err error

	t.Run("should send code and return state when login start first time", func(t *testing.T) {
		clientStateOne, err = accountService.LogInStart("999111333", "VN")
		codeOne = smsSender.text

		if err != nil {
			t.Error(err)
		}

		if codeOne == "" {
			t.Error("missing code")
		}
	})

	// This behaves like "resending"
	t.Run("should send different code and state when login start second time", func(t *testing.T) {
		clientStateTwo, err = accountService.LogInStart("999111333", "VN")
		codeTwo = smsSender.text

		if err != nil {
			t.Error(err)
		}

		if codeOne == codeTwo {
			t.Error("same code")
		}

		if clientStateOne == clientStateTwo {
			t.Error("same client state")
		}
	})
}

func TestLoginWithExpiredVerificationCode(t *testing.T) {
	now := time.Now()

	account := NewAccount("999999999", "VN")

	err := accountStore.StoreAccount(account)

	if err != nil {
		t.Error(err)
		return
	}

	expiredVerificationCode := &VerificationCode{
		ID:          uuid.Must(uuid.New(), nil).String(),
		AccountID:   account.ID,
		Code:        "111111",
		CodeType:    "LOGIN",
		ClientState: "ABC",
		PhoneNumber: account.PhoneNumber,
		CountryCode: account.CountryCode,
		ExpiredAt:   now.Add(time.Duration(-1) * time.Minute),
		CreatedAt:   now,
	}

	err = accountStore.StoreVerificationCode(expiredVerificationCode)

	if err != nil {
		t.Error(err)
		return
	}

	t.Run("should return error when log in verify with expired verification code", func(t *testing.T) {
		_, err := accountService.LogInVerify(expiredVerificationCode.ClientState, expiredVerificationCode.Code)

		if !errors.Is(errors.Forbidden, err) {
			t.Error("should forbid")
		}
	})
}
