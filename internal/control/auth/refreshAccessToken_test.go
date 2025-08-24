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

func getBehaviorForRefresh(t *testing.T, userAuthentic *shelterUser.UserAuthentic, accessToken pkgText.JwtToken) (*userTestUtility.RefreshTokenCheckerMock, *userTestUtility.AccessTokenIssuerMock) {
	var refreshTokenChecker = &userTestUtility.RefreshTokenCheckerMock{
		FakeExecute: func(entry entryAuth.RefreshTokenGetter) (*shelterUser.UserAuthentic, error) {
			return userAuthentic, nil
		},
	}

	var accessTokenIssuer = &userTestUtility.AccessTokenIssuerMock{
		FakeExecute: func(user *shelterUser.UserAuthentic) (pkgText.JwtToken, error) {
			return accessToken, nil
		},
	}

	return refreshTokenChecker, accessTokenIssuer
}

func getShelterUserAuthenticForRefresh(expectId string) *shelterUser.UserAuthentic {
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

	var email, _ = pkgText.NewEmail("test@example.com")
	return shelterUser.NewUserAuthentic(userValue, companyRole, &email)
}

func getRefreshEntry(expectToken string) entryAuth.AuthRefreshRequest {
	return entryAuth.AuthRefreshRequest{
		AuthRefresh: entryAuth.AuthRefresh{
			RefreshToken: expectToken,
		},
	}
}

func TestRefresh(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectUUID, _ = uuid.NewUUID()
	var refreshToken, _ = shelterText.CreateToken(expectUUID)
	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)
	var userAuthentic = getShelterUserAuthenticForRefresh(expectIdentifier)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)

	var refreshTokenChecker, accessTokenIssuer = getBehaviorForRefresh(t, userAuthentic, accessToken)
	var control = controlAuth.NewRefreshAccessTokenControl(
		db,
		refreshTokenChecker,
		accessTokenIssuer,
	)

	var entry = getRefreshEntry(string(refreshToken))

	var authRefreshResponse, err = controlAuth.RefreshAccessTokenExecute(control, entry, nil)

	assert.NoError(t, err)
	assert.Equal(t, expectIdentifier, authRefreshResponse.User.Identifier)
	assert.Equal(t, expectToken, authRefreshResponse.AccessToken)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 1, transactionCalledCount.CommitCalled)
	assert.Equal(t, 0, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)

	t.Logf("response: %+v", authRefreshResponse)
}

func TestRefreshErrRefresh(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectUUID, _ = uuid.NewUUID()
	var refreshToken, _ = shelterText.CreateToken(expectUUID)
	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)
	var userAuthentic = getShelterUserAuthenticForRefresh(expectIdentifier)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)

	var refreshTokenChecker, accessTokenIssuer = getBehaviorForRefresh(t, userAuthentic, accessToken)
	refreshTokenChecker.FakeExecute = func(entry entryAuth.RefreshTokenGetter) (*shelterUser.UserAuthentic, error) {
		return nil, errors.New("refresh token error")
	}
	var control = controlAuth.NewRefreshAccessTokenControl(
		db,
		refreshTokenChecker,
		accessTokenIssuer,
	)

	var entry = getRefreshEntry(string(refreshToken))

	var _, err = controlAuth.RefreshAccessTokenExecute(control, entry, nil)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestRefreshErrIssue(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectUUID, _ = uuid.NewUUID()
	var refreshToken, _ = shelterText.CreateToken(expectUUID)
	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)
	var userAuthentic = getShelterUserAuthenticForRefresh(expectIdentifier)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)

	var refreshTokenChecker, accessTokenIssuer = getBehaviorForRefresh(t, userAuthentic, accessToken)
	accessTokenIssuer.FakeExecute = func(user *shelterUser.UserAuthentic) (pkgText.JwtToken, error) {
		return pkgText.JwtToken(""), errors.New("access token error")
	}
	var control = controlAuth.NewRefreshAccessTokenControl(
		db,
		refreshTokenChecker,
		accessTokenIssuer,
	)

	var entry = getRefreshEntry(string(refreshToken))

	var _, err = controlAuth.RefreshAccessTokenExecute(control, entry, nil)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}
