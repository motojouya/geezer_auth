package user_test

import (
	// "errors"
	"github.com/motojouya/geezer_auth/internal/behavior/user"
	dbUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	shelterText "github.com/motojouya/geezer_auth/internal/shelter/text"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type passwordCheckerDBMock struct {
	getUserPassword func(identifier string) (*dbUser.UserPasswordFull, error)
	getUserPasswordOfEmail func(email string) (*dbUser.UserPasswordFull, error)
}

func (mock passwordCheckerDBMock) GetUserPassword(identifier string) (*dbUser.UserPasswordFull, error) {
	return mock.getUserPassword(identifier)
}

func (mock passwordCheckerDBMock) GetUserPasswordOfEmail(email string) (*dbUser.UserPasswordFull, error) {
	return mock.getUserPasswordOfEmail(email)
}

type loginnerEntryMock struct {
	getPassword func() (shelterText.Password, error)
	getIdentifier func() (*pkgText.Identifier, error)
	getEmailIdentifier func() (*pkgText.Email, error)
}

func (mock loginnerEntryMock) GetPassword() (shelterText.Password, error) {
	return mock.getPassword()
}

func (mock loginnerEntryMock) GetIdentifier() (*pkgText.Identifier, error) {
	return mock.getIdentifier()
}

func (mock loginnerEntryMock) GetEmailIdentifier() (*pkgText.Email, error) {
	return mock.getEmailIdentifier()
}

func getDbUserPassword(expectId string, expectEmail string, expectPassword string) dbUser.UserPasswordFull {
	var now = time.Now()
	var expireDate = now.Add(1 * time.Hour)
	return dbUser.UserPasswordFull{
		UserPassword: dbUser.UserPassword{
			PersistKey:     1,
			UserPersistKey: 2,
			Password:       expectPassword,
			RegisteredDate: now.Add(3 * time.Hour),
			ExpireDate:     &expireDate,
		},
		UserIdentifier:     expectId,
		UserExposeEmailId:  expectEmail,
		UserName:           "TestUserName",
		UserBotFlag:        false,
		UserRegisteredDate: now,
		UserUpdateDate:     now.Add(2 * time.Hour),
	}
}

func getPasswordCheckDbMock(t *testing.T, expectId string, expectEmail string, expectPassword string) passwordCheckerDBMock {
	var userPassword = getDbUserPassword(expectId, expectEmail, expectPassword)
	var getUserPassword = func(identifier string) (*dbUser.UserPasswordFull, error) {
		assert.Equal(t, expectId, identifier)
		return &userPassword, nil
	}
	var getUserPasswordOfEmail = func(email string) (*dbUser.UserPasswordFull, error) {
		assert.Equal(t, expectEmail, email)
		return &userPassword, nil
	}
	return passwordCheckerDBMock{
		getUserPassword: getUserPassword,
		getUserPasswordOfEmail: getUserPasswordOfEmail,
	}
}

func getLoginEntryMock(t *testing.T, expectId string, expectEmail string, expectPassword string) loginnerEntryMock {
	var password, _ = shelterText.NewPassword(expectPassword)
	var email, _ = pkgText.NewEmail(expectEmail)
	var identifier, _ = pkgText.NewIdentifier(expectId)
	var getPassword = func() (shelterText.Password, error) {
		return password, nil
	}
	var getIdentifier = func() (*pkgText.Identifier, error) {
		return &identifier, nil
	}
	var getEmailIdentifier = func() (*pkgText.Email, error) {
		return &email, nil
	}
	return loginnerEntryMock{
		getPassword:   getPassword,
		getIdentifier: getIdentifier,
		getEmailIdentifier:      getEmailIdentifier,
	}
}

func TestPasswordChecker(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectEmail = "test@example.com"

	var expectPassword = "password01"
	var password, _ = shelterText.NewPassword(expectPassword)
	var hashedPassword, _ = shelterText.HashPassword(password)

	var dbMock = getPasswordCheckDbMock(t, expectIdentifier, expectEmail, string(hashedPassword))
	var entryMock = getLoginEntryMock(t, expectIdentifier, expectEmail, expectPassword)

	checker := user.NewPasswordCheck(dbMock)
	err := checker.Execute(entryMock)

	assert.NoError(t, err)
}

// TODO working error cases
