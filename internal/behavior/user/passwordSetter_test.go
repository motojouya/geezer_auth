package user_test

import (
	"github.com/motojouya/geezer_auth/internal/behavior/testUtility"
	"github.com/motojouya/geezer_auth/internal/behavior/user"
	dbUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
	shelterCompany "github.com/motojouya/geezer_auth/internal/shelter/company"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	shelterText "github.com/motojouya/geezer_auth/internal/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"errors"
)

type passwordSetterDBMock struct {
	addPassword func(userPassword *dbUser.UserPassword, now time.Time) (*dbUser.UserPassword, error)
}

func (mock passwordSetterDBMock) AddPassword(userPassword *dbUser.UserPassword, now time.Time) (*dbUser.UserPassword, error) {
	return mock.addPassword(userPassword, now)
}

type passwordGetEntryMock struct {
	getPassword func() (shelterText.Password, error)
}

func (mock passwordGetEntryMock) GetPassword() (shelterText.Password, error) {
	return mock.getPassword()
}

func getShelterUserAuthenticForPassword() *shelterUser.UserAuthentic {
	var userId uint = 1
	var userIdentifier, _ = pkgText.NewIdentifier("TestIdentifier")
	var emailId, _ = pkgText.NewEmail("test@example.com")
	var userName, _ = pkgText.NewName("TestName")
	var botFlag = false
	var userRegisteredDate = time.Now()
	var updateDate = time.Now()
	var userValue = shelterUser.NewUser(userId, userIdentifier, userName, emailId, botFlag, userRegisteredDate, updateDate)

	var companyIdentifier, _ = pkgText.NewIdentifier("CP-TESTES")
	var companyId uint = 1
	var companyName, _ = pkgText.NewName("TestCompany")
	var companyRegisteredDate = time.Now()
	var company = shelterCompany.NewCompany(companyId, companyIdentifier, companyName, companyRegisteredDate)

	var label, _ = pkgText.NewLabel("TEST_ROLE")
	var roleName, _ = pkgText.NewName("TestRole")
	var description, _ = shelterText.NewText("Role for testing")
	var roleRegisteredDate = time.Now()

	var roles = []shelterRole.Role{shelterRole.NewRole(roleName, label, description, roleRegisteredDate)}
	var companyRole = shelterUser.NewCompanyRole(company, roles)

	var email, _ = pkgText.NewEmail("test_2@gmail.com")
	return shelterUser.NewUserAuthentic(userValue, companyRole, &email)
}

func getLocalerMockForPassword(t *testing.T, now time.Time) *testUtility.LocalerMock {
	var getNow = func() time.Time {
		return now
	}
	return &testUtility.LocalerMock{
		FakeGetNow: getNow,
	}
}

func getPasswordSetDbMock(t *testing.T, expectPassword string, firstNow time.Time) passwordSetterDBMock {
	var addPassword = func(userPassword *dbUser.UserPassword, now time.Time) (*dbUser.UserPassword, error) {
		var verifyPassword = shelterText.VerifyPassword(shelterText.HashedPassword(userPassword.Password), shelterText.Password(expectPassword))
		assert.NoError(t, verifyPassword)
		assert.WithinDuration(t, now, firstNow, time.Second, "Expected 'now' to be within 1 second of current time")
		return userPassword, nil
	}
	return passwordSetterDBMock{
		addPassword: addPassword,
	}
}

func getGetPasswordEntryMock(t *testing.T, expectPassword string, firstNow time.Time) passwordGetEntryMock {
	var password, _ = shelterText.NewPassword(expectPassword)
	var getPassword = func() (shelterText.Password, error) {
		return password, nil
	}
	return passwordGetEntryMock{
		getPassword: getPassword,
	}
}

func TestPasswordSetter(t *testing.T) {
	var expectPassword = "password01"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForPassword()

	var localerMock = getLocalerMockForPassword(t, firstNow)
	var dbMock = getPasswordSetDbMock(t, expectPassword, firstNow)
	var entryMock = getGetPasswordEntryMock(t, expectPassword, firstNow)

	setter := user.NewPasswordSet(localerMock, dbMock)
	err := setter.Execute(entryMock, userAuthentic)

	assert.NoError(t, err)
}

func TestPasswordSetterErrGetPass(t *testing.T) {
	var expectPassword = "password01"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForPassword()

	var localerMock = getLocalerMockForPassword(t, firstNow)
	var dbMock = getPasswordSetDbMock(t, expectPassword, firstNow)
	var entryMock = getGetPasswordEntryMock(t, expectPassword, firstNow)

	entryMock.getPassword = func() (shelterText.Password, error) {
		return shelterText.Password(""), errors.New("error getting password")
	}

	setter := user.NewPasswordSet(localerMock, dbMock)
	err := setter.Execute(entryMock, userAuthentic)

	assert.Error(t, err)
}

func TestPasswordSetterErrAddPass(t *testing.T) {
	var expectPassword = "password01"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForPassword()

	var localerMock = getLocalerMockForPassword(t, firstNow)
	var dbMock = getPasswordSetDbMock(t, expectPassword, firstNow)
	var entryMock = getGetPasswordEntryMock(t, expectPassword, firstNow)

	dbMock.addPassword = func(userPassword *dbUser.UserPassword, now time.Time) (*dbUser.UserPassword, error) {
		return nil, errors.New("error adding password")
	}

	setter := user.NewPasswordSet(localerMock, dbMock)
	err := setter.Execute(entryMock, userAuthentic)

	assert.Error(t, err)
}
