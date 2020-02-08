package handlers

import (
	"testing"

	"github.com/minheq/kedul_server_main/models"
	"github.com/minheq/kedul_server_main/testutils"
)

func TestLoginVerifyTwice(t *testing.T) {
	var codeOne string
	var verificationIDOne string
	var codeTwo string
	var verificationIDTwo string
	var err error
	db, cleanup := testutils.SetupDB()
	store := models.NewStore(db)
	smsSender := &testutils.SmsSenderMock{}

	defer cleanup()

	t.Run("should send code and return state when login start first time", func(t *testing.T) {
		verificationIDOne, err = LoginVerify("999111333", "VN", store, smsSender)
		codeOne = smsSender.Text

		if err != nil {
			t.Error(err)
		}

		if codeOne == "" {
			t.Error("missing code")
		}
	})

	// This behaves like "resending"
	t.Run("should send different code and state when login start second time", func(t *testing.T) {
		verificationIDTwo, err = LoginVerify("999111333", "VN", store, smsSender)
		codeTwo = smsSender.Text

		if err != nil {
			t.Error(err)
		}

		if codeOne == codeTwo {
			t.Error("same code")
		}

		if verificationIDOne == verificationIDTwo {
			t.Error("same verification id")
		}
	})
}
