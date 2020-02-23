package auth

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/google/uuid"
	"github.com/minheq/kedul_server_main/errors"
	"github.com/minheq/kedul_server_main/phone"
	"github.com/minheq/kedul_server_main/testutils"
)

type mockAuthStore struct {
	users             []*User
	verificationCodes []*VerificationCode
}

func (s *mockAuthStore) GetVerificationCodeByIDAndCode(ctx context.Context, verificationID string, code string) (*VerificationCode, error) {
	for _, v := range s.verificationCodes {
		if v.VerificationID == verificationID && v.Code == code {
			return v, nil
		}
	}

	return nil, nil
}

func (s *mockAuthStore) StoreVerificationCode(ctx context.Context, vc *VerificationCode) error {
	s.verificationCodes = append(s.verificationCodes, vc)

	return nil
}

func (s *mockAuthStore) RemoveVerificationCodeByPhoneNumber(ctx context.Context, phoneNumber string, countryCode string) error {
	for i, v := range s.verificationCodes {
		if v.PhoneNumber == phoneNumber && v.CountryCode == countryCode {
			s.verificationCodes = append(s.verificationCodes[:i], s.verificationCodes[i+1:]...)
			break
		}
	}

	return nil
}

func (s *mockAuthStore) RemoveVerificationCodeByID(ctx context.Context, id string) error {
	for i, v := range s.verificationCodes {
		if v.ID == id {
			s.verificationCodes = append(s.verificationCodes[:i], s.verificationCodes[i+1:]...)
			break
		}
	}

	return nil
}

func (s *mockAuthStore) GetUserByID(ctx context.Context, id string) (*User, error) {
	for _, u := range s.users {
		if u.ID == id {
			return u, nil
		}
	}

	return nil, nil
}

func (s *mockAuthStore) GetUserByPhoneNumber(ctx context.Context, phoneNumber string, countryCode string) (*User, error) {
	for _, u := range s.users {
		if u.PhoneNumber == phoneNumber && u.CountryCode == countryCode {
			return u, nil
		}
	}

	return nil, nil
}

func (s *mockAuthStore) StoreUser(ctx context.Context, user *User) error {
	s.users = append(s.users, user)

	return nil
}

func (s *mockAuthStore) UpdateUser(ctx context.Context, user *User) error {
	for i, u := range s.users {
		if u.ID == user.ID {
			s.users[i] = user
			break
		}
	}

	return nil
}

var (
	mockStore   Store
	authService Service
	smsSender   *testutils.SmsSenderMock
)

func TestMain(m *testing.M) {
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	mockStore = &mockAuthStore{}
	smsSender = &testutils.SmsSenderMock{}

	authService = NewService(mockStore, tokenAuth, smsSender)

	code := m.Run()

	os.Exit(code)
}

func TestLoginHappyPath(t *testing.T) {
	var code string
	var verificationID string
	var err error

	t.Run("should send code and return verificationID when login start", func(t *testing.T) {
		verificationID, err = authService.LoginVerify(context.Background(), "999111333", "VN")
		code = smsSender.Text

		if err != nil {
			t.Error(err)
			return
		}

		if code == "" {
			t.Error("missing code")
			return
		}
	})

	t.Run("should return access token when login verified", func(t *testing.T) {
		accessToken, err := authService.LoginCheck(context.Background(), verificationID, code)

		if err != nil {
			t.Error(err)
			return
		}

		if accessToken == "" {
			t.Error("missing access token")
		}
	})
}

func TestLoginWithExpiredVerificationCode(t *testing.T) {
	now := time.Now()

	phoneNumber, _ := phone.FormatPhoneNumber("999999999", "VN")
	user := NewUser(phoneNumber, "VN")

	err := mockStore.StoreUser(context.Background(), user)

	if err != nil {
		t.Error(err)
		return
	}

	expiredVerificationCode := &VerificationCode{
		ID:             uuid.Must(uuid.New(), nil).String(),
		UserID:         user.ID,
		Code:           "111111",
		CodeType:       "LOGIN",
		VerificationID: "ABC",
		PhoneNumber:    user.PhoneNumber,
		CountryCode:    user.CountryCode,
		ExpiredAt:      now.Add(time.Duration(-1) * time.Minute),
		CreatedAt:      now,
	}

	err = mockStore.StoreVerificationCode(context.Background(), expiredVerificationCode)

	if err != nil {
		t.Error(err)
		return
	}

	t.Run("should return error when log in verify with expired verification code", func(t *testing.T) {
		_, err := authService.LoginCheck(context.Background(), expiredVerificationCode.VerificationID, expiredVerificationCode.Code)

		if errors.Is(errors.KindInvalid, err) == false {
			t.Error("error should be invalid kind")
		}
	})
}

func TestLoginVerifyTwice(t *testing.T) {
	var codeOne string
	var verificationIDOne string
	var codeTwo string
	var verificationIDTwo string
	var err error

	t.Run("should send code and return verificationID when login start first time", func(t *testing.T) {
		verificationIDOne, err = authService.LoginVerify(context.Background(), "999111334", "VN")
		codeOne = smsSender.Text

		if err != nil {
			t.Error(err)
			return
		}

		if codeOne == "" {
			t.Error("missing code")
		}
	})

	// This behaves like "resending"
	t.Run("should send different code and verificationID when login start second time", func(t *testing.T) {
		verificationIDTwo, err = authService.LoginVerify(context.Background(), "999111334", "VN")
		codeTwo = smsSender.Text

		if err != nil {
			t.Error(err)
			return
		}

		if codeOne == codeTwo {
			t.Error("same code")
			return
		}

		if verificationIDOne == verificationIDTwo {
			t.Error("same verification id")
		}
	})
}

func TestUpdatePhoneNumberHappyPath(t *testing.T) {
	var code string
	var verificationID string
	var err error

	prevPhoneNumber, err := phone.FormatPhoneNumber("999111335", "VN")
	newPhoneNumber, err := phone.FormatPhoneNumber("999111336", "VN")

	currentUser := NewUser(prevPhoneNumber, "VN")
	err = mockStore.StoreUser(context.Background(), currentUser)

	if err != nil {
		t.Error(err)
		return
	}

	t.Run("should send code and return verificationID when login start", func(t *testing.T) {
		verificationID, err = authService.UpdatePhoneNumberVerify(context.Background(), newPhoneNumber, "VN", currentUser)
		code = smsSender.Text

		if err != nil {
			t.Error(err)
			return
		}

		if code == "" {
			t.Error("missing code")
			return
		}
	})

	t.Run("should return access token when login verified", func(t *testing.T) {
		user, err := authService.UpdatePhoneNumberCheck(context.Background(), verificationID, code, currentUser)

		if err != nil {
			t.Error(err)
			return
		}

		if user.PhoneNumber != newPhoneNumber {
			t.Error("user failed to update")
		}
	})
}

func TestUpdateUserProfileHappyPath(t *testing.T) {
	phoneNumber, err := phone.FormatPhoneNumber("999111337", "VN")

	currentUser := NewUser(phoneNumber, "VN")
	err = mockStore.StoreUser(context.Background(), currentUser)
	newFullName := "new_name"
	newProfileImageID := "new_profile_image_id"

	if err != nil {
		t.Error(err)
		return
	}

	t.Run("should update user", func(t *testing.T) {
		user, err := authService.UpdateUserProfile(context.Background(), newFullName, newProfileImageID, currentUser)

		if err != nil {
			t.Error(err)
			return
		}

		if user.FullName != newFullName {
			t.Error("user failed to update full name")
		}

		if user.ProfileImageID != newProfileImageID {
			t.Error("user failed to update profile image id")
		}
	})
}
