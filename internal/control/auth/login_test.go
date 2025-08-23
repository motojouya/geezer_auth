package auth_test

import (
	"errors"
	"github.com/google/uuid"
	userTestUtility "github.com/motojouya/geezer_auth/internal/behavior/user/testUtility"
	controlAuth "github.com/motojouya/geezer_auth/internal/control/auth"
	dbTestUtility "github.com/motojouya/geezer_auth/internal/db/testUtility"
	entryAuth "github.com/motojouya/geezer_auth/internal/entry/transfer/auth"
	shelterCompany "github.com/motojouya/geezer_auth/internal/shelter/company"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
	shelterText "github.com/motojouya/geezer_auth/internal/shelter/text"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getBehaviorForLogin(t *testing.T, userAuthentic *shelterUser.UserAuthentic, refreshToken shelterText.Token, accessToken pkgText.JwtToken) (*userTestUtility.UserGetterMock, *userTestUtility.PasswordCheckerMock, *userTestUtility.RefreshTokenIssuerMock, *userTestUtility.AccessTokenIssuerMock) {
	var userGetter = &userTestUtility.UserGetterMock{
		FakeExecute: func(identifier pkgText.Identifier) (*shelterUser.UserAuthentic, error) {
			assert.Equal(t, userAuthentic.Identifier, identifier)
			return userAuthentic, nil
		},
	}

	var passwordChecker = &userTestUtility.PasswordCheckerMock{
		FakeExecute: func(entry entryAuth.AuthLoginner) (pkgText.Identifier, error) {
			return userAuthentic.Identifier, nil
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

	return userGetter, passwordChecker, refreshTokenIssuer, accessTokenIssuer
}

func getShelterUserAuthenticForLogin(expectId string, expectEmail string) *shelterUser.UserAuthentic {
	var userId uint = 1
	var userIdentifier, _ = pkgText.NewIdentifier(expectId)
	var emailId, _ = pkgText.NewEmail(expectEmail)
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

	var email, _ = pkgText.NewEmail(expectEmail)
	return shelterUser.NewUserAuthentic(userValue, companyRole, &email)
}

func getLoginEntry(expectId string, expectEmail string, expectPassword string) entryAuth.AuthLoginRequest {
	return entryAuth.AuthLoginRequest{
		AuthLogin: entryAuth.AuthLogin{
			AuthIdentifier: entryAuth.AuthIdentifier{
				Identifier:      &expectId,
				EmailIdentifier: &expectEmail,
			},
			Password: expectPassword,
		},
	}
}

func TestLogin(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectEmail = "test@example.com"
	var expectPassword = "password123"
	var expectUUID, _ = uuid.NewUUID()
	var refreshToken, _ = shelterText.CreateToken(expectUUID)
	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)
	var userAuthentic = getShelterUserAuthenticForLogin(expectIdentifier, expectEmail)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)

	var userGetter, passwordChecker, refreshTokenIssuer, accessTokenIssuer = getBehaviorForLogin(t, userAuthentic, refreshToken, accessToken)
	var control = controlAuth.NewLoginControl(
		db,
		userGetter,
		passwordChecker,
		refreshTokenIssuer,
		accessTokenIssuer,
	)

	var entry = getLoginEntry(expectIdentifier, expectEmail, expectPassword)

	var authLoginResponse, err = controlAuth.LoginExecute(control, entry, nil)

	assert.NoError(t, err)
	assert.Equal(t, expectIdentifier, authLoginResponse.User.Identifier)
	assert.Equal(t, expectUUID.String(), authLoginResponse.RefreshToken)
	assert.Equal(t, expectToken, authLoginResponse.AccessToken)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 1, transactionCalledCount.CommitCalled)
	assert.Equal(t, 0, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)

	t.Logf("response: %+v", authLoginResponse)
}

func TestLoginErrLogin(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectEmail = "test@example.com"
	var expectPassword = "password123"
	var expectUUID, _ = uuid.NewUUID()
	var refreshToken, _ = shelterText.CreateToken(expectUUID)
	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)
	var userAuthentic = getShelterUserAuthenticForLogin(expectIdentifier, expectEmail)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)

	var userGetter, passwordChecker, refreshTokenIssuer, accessTokenIssuer = getBehaviorForLogin(t, userAuthentic, refreshToken, accessToken)
	passwordChecker.FakeExecute = func(entry entryAuth.AuthLoginner) (pkgText.Identifier, error) {
		return pkgText.Identifier(""), errors.New("login error")
	}
	var control = controlAuth.NewLoginControl(
		db,
		userGetter,
		passwordChecker,
		refreshTokenIssuer,
		accessTokenIssuer,
	)

	var entry = getLoginEntry(expectIdentifier, expectEmail, expectPassword)

	var _, err = controlAuth.LoginExecute(control, entry, nil)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestLoginErrGet(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectEmail = "test@example.com"
	var expectPassword = "password123"
	var expectUUID, _ = uuid.NewUUID()
	var refreshToken, _ = shelterText.CreateToken(expectUUID)
	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)
	var userAuthentic = getShelterUserAuthenticForLogin(expectIdentifier, expectEmail)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)

	var userGetter, passwordChecker, refreshTokenIssuer, accessTokenIssuer = getBehaviorForLogin(t, userAuthentic, refreshToken, accessToken)
	userGetter.FakeExecute = func(identifier pkgText.Identifier) (*shelterUser.UserAuthentic, error) {
		return nil, errors.New("get user error")
	}
	var control = controlAuth.NewLoginControl(
		db,
		userGetter,
		passwordChecker,
		refreshTokenIssuer,
		accessTokenIssuer,
	)

	var entry = getLoginEntry(expectIdentifier, expectEmail, expectPassword)

	var _, err = controlAuth.LoginExecute(control, entry, nil)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestLoginErrRefToken(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectEmail = "test@example.com"
	var expectPassword = "password123"
	var expectUUID, _ = uuid.NewUUID()
	var refreshToken, _ = shelterText.CreateToken(expectUUID)
	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)
	var userAuthentic = getShelterUserAuthenticForLogin(expectIdentifier, expectEmail)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)

	var userGetter, passwordChecker, refreshTokenIssuer, accessTokenIssuer = getBehaviorForLogin(t, userAuthentic, refreshToken, accessToken)
	refreshTokenIssuer.FakeExecute = func(user *shelterUser.UserAuthentic) (shelterText.Token, error) {
		return shelterText.Token(""), errors.New("refresh token error")
	}
	var control = controlAuth.NewLoginControl(
		db,
		userGetter,
		passwordChecker,
		refreshTokenIssuer,
		accessTokenIssuer,
	)

	var entry = getLoginEntry(expectIdentifier, expectEmail, expectPassword)

	var _, err = controlAuth.LoginExecute(control, entry, nil)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestLoginErrAccToken(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectEmail = "test@example.com"
	var expectPassword = "password123"
	var expectUUID, _ = uuid.NewUUID()
	var refreshToken, _ = shelterText.CreateToken(expectUUID)
	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)
	var userAuthentic = getShelterUserAuthenticForLogin(expectIdentifier, expectEmail)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)

	var userGetter, passwordChecker, refreshTokenIssuer, accessTokenIssuer = getBehaviorForLogin(t, userAuthentic, refreshToken, accessToken)
	accessTokenIssuer.FakeExecute = func(user *shelterUser.UserAuthentic) (pkgText.JwtToken, error) {
		return pkgText.JwtToken(""), errors.New("access token error")
	}
	var control = controlAuth.NewLoginControl(
		db,
		userGetter,
		passwordChecker,
		refreshTokenIssuer,
		accessTokenIssuer,
	)

	var entry = getLoginEntry(expectIdentifier, expectEmail, expectPassword)

	var _, err = controlAuth.LoginExecute(control, entry, nil)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}
