package user_test

import (
	"errors"
	"github.com/google/uuid"
	controlUser "github.com/motojouya/geezer_auth/internal/control/user"
	dbTestUtility "github.com/motojouya/geezer_auth/internal/db/testUtility"
	userTestUtility "github.com/motojouya/geezer_auth/internal/behavior/user/testUtility"
	entryUser "github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	shelterCompany "github.com/motojouya/geezer_auth/internal/shelter/company"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
	shelterText "github.com/motojouya/geezer_auth/internal/shelter/text"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getBehaviorForRegister(userAuthentic *shelterUser.UserAuthentic, refreshToken shelterText.Token, accessToken pkgText.JwtToken) (*userTestUtility.UserCreatorMock, *userTestUtility.EmailSetterMock, *userTestUtility.PasswordSetterMock, *userTestUtility.RefreshTokenIssuerMock, *userTestUtility.AccessTokenIssuerMock) {
	var userCreator = &userTestUtility.UserCreatorMock{
		FakeExecute: func(entry entryUser.UserGetter) (*shelterUser.UserAuthentic, error) {
			return userAuthentic, nil
		},
	}

	var emailSetter = &userTestUtility.EmailSetterMock{
		FakeExecute: func(entry entryUser.EmailGetter, user *shelterUser.UserAuthentic) error {
			return nil
		},
	}

	var passwordSetter = &userTestUtility.PasswordSetterMock{
		FakeExecute: func(entry entryUser.PasswordGetter, user *shelterUser.UserAuthentic) error {
			return nil
		},
	}

	var refreshTokenIssuer = &userTestUtility.RefreshTokenIssuerMock{
		FakeExecute: func(user *shelterUser.UserAuthentic) (shelterText.Token, error) {
			return refreshToken, nil
		},
	}

	var accessTokenIssuer = &userTestUtility.AccessTokenIssuerMock{
		FakeExecute: func(user *shelterUser.UserAuthentic) (pkgText.JwtToken, error) {
			return accessToken, nil
		},
	}

	return userCreator, emailSetter, passwordSetter, refreshTokenIssuer, accessTokenIssuer
}

func getShelterUserAuthenticForRegister(expectId string) *shelterUser.UserAuthentic {
	var userId uint = 1
	var userIdentifier, _ = pkgText.NewIdentifier(expectId)
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

func getRegisterEntry() entryUser.UserRegisterRequest {
	return entryUser.UserRegisterRequest{
		UserRegister: entryUser.UserRegister{
			Email:    "test@example.com",
			Name:     "TestName",
			Bot:      false,
			Password: "password123",
		},
	}
}

func TestRegister(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectUUID, _ = uuid.NewUUID()
	var refreshToken, _ = shelterText.CreateToken(expectUUID)
	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)
	var userAuthentic = getShelterUserAuthenticForRegister(expectIdentifier)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)

	var userCreator, emailSetter, passwordSetter, refreshTokenIssuer, accessTokenIssuer = getBehaviorForRegister(userAuthentic, refreshToken, accessToken)
	var control = controlUser.NewRegisterControl(
		db,
		userCreator,
		emailSetter,
		passwordSetter,
		refreshTokenIssuer,
		accessTokenIssuer,
	)

	var entry = getRegisterEntry()

	var userRegisterResponse, err = controlUser.RegisterExecute(control, entry, nil)

	assert.NoError(t, err)
	assert.Equal(t, expectIdentifier, userRegisterResponse.User.Identifier)
	assert.Equal(t, expectUUID.String(), userRegisterResponse.RefreshToken)
	assert.Equal(t, expectToken, userRegisterResponse.AccessToken)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 1, transactionCalledCount.CommitCalled)
	assert.Equal(t, 0, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)

	t.Logf("User Identifier: %+v", userRegisterResponse)
}

func TestRegisterErrCreator(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectUUID, _ = uuid.NewUUID()
	var refreshToken, _ = shelterText.CreateToken(expectUUID)
	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)
	var userAuthentic = getShelterUserAuthenticForRegister(expectIdentifier)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)

	var userCreator, emailSetter, passwordSetter, refreshTokenIssuer, accessTokenIssuer = getBehaviorForRegister(userAuthentic, refreshToken, accessToken)
	userCreator.FakeExecute = func(entry entryUser.UserGetter) (*shelterUser.UserAuthentic, error) {
		return nil, errors.New("user creation error")
	}

	var control = controlUser.NewRegisterControl(
		db,
		userCreator,
		emailSetter,
		passwordSetter,
		refreshTokenIssuer,
		accessTokenIssuer,
	)

	var entry = getRegisterEntry()

	var _, err = controlUser.RegisterExecute(control, entry, nil)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestRegisterErrEmail(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectUUID, _ = uuid.NewUUID()
	var refreshToken, _ = shelterText.CreateToken(expectUUID)
	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)
	var userAuthentic = getShelterUserAuthenticForRegister(expectIdentifier)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)

	var userCreator, emailSetter, passwordSetter, refreshTokenIssuer, accessTokenIssuer = getBehaviorForRegister(userAuthentic, refreshToken, accessToken)
	emailSetter.FakeExecute = func(entry entryUser.EmailGetter, user *shelterUser.UserAuthentic) error {
		return errors.New("email setting error")
	}

	var control = controlUser.NewRegisterControl(
		db,
		userCreator,
		emailSetter,
		passwordSetter,
		refreshTokenIssuer,
		accessTokenIssuer,
	)

	var entry = getRegisterEntry()

	var _, err = controlUser.RegisterExecute(control, entry, nil)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestRegisterErrPassword(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectUUID, _ = uuid.NewUUID()
	var refreshToken, _ = shelterText.CreateToken(expectUUID)
	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)
	var userAuthentic = getShelterUserAuthenticForRegister(expectIdentifier)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)

	var userCreator, emailSetter, passwordSetter, refreshTokenIssuer, accessTokenIssuer = getBehaviorForRegister(userAuthentic, refreshToken, accessToken)
	passwordSetter.FakeExecute = func(entry entryUser.PasswordGetter, user *shelterUser.UserAuthentic) error {
		return errors.New("password setting error")
	}

	var control = controlUser.NewRegisterControl(
		db,
		userCreator,
		emailSetter,
		passwordSetter,
		refreshTokenIssuer,
		accessTokenIssuer,
	)

	var entry = getRegisterEntry()

	var _, err = controlUser.RegisterExecute(control, entry, nil)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestRegisterErrRefToken(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectUUID, _ = uuid.NewUUID()
	var refreshToken, _ = shelterText.CreateToken(expectUUID)
	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)
	var userAuthentic = getShelterUserAuthenticForRegister(expectIdentifier)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)

	var userCreator, emailSetter, passwordSetter, refreshTokenIssuer, accessTokenIssuer = getBehaviorForRegister(userAuthentic, refreshToken, accessToken)
	refreshTokenIssuer.FakeExecute = func(user *shelterUser.UserAuthentic) (shelterText.Token, error) {
		return "", errors.New("refresh token issuing error")
	}

	var control = controlUser.NewRegisterControl(
		db,
		userCreator,
		emailSetter,
		passwordSetter,
		refreshTokenIssuer,
		accessTokenIssuer,
	)

	var entry = getRegisterEntry()

	var _, err = controlUser.RegisterExecute(control, entry, nil)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestRegisterErrAccToken(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectUUID, _ = uuid.NewUUID()
	var refreshToken, _ = shelterText.CreateToken(expectUUID)
	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)
	var userAuthentic = getShelterUserAuthenticForRegister(expectIdentifier)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)

	var userCreator, emailSetter, passwordSetter, refreshTokenIssuer, accessTokenIssuer = getBehaviorForRegister(userAuthentic, refreshToken, accessToken)
	accessTokenIssuer.FakeExecute = func(user *shelterUser.UserAuthentic) (pkgText.JwtToken, error) {
		return "", errors.New("access token issuing error")
	}

	var control = controlUser.NewRegisterControl(
		db,
		userCreator,
		emailSetter,
		passwordSetter,
		refreshTokenIssuer,
		accessTokenIssuer,
	)

	var entry = getRegisterEntry()

	var _, err = controlUser.RegisterExecute(control, entry, nil)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}
