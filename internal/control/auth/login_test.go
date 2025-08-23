package user_test

import (
	"errors"
	"github.com/google/uuid"
	userTestUtility "github.com/motojouya/geezer_auth/internal/behavior/user/testUtility"
	controlUser "github.com/motojouya/geezer_auth/internal/control/user"
	dbTestUtility "github.com/motojouya/geezer_auth/internal/db/testUtility"
	entryUser "github.com/motojouya/geezer_auth/internal/entry/transfer/user"
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

func getBehaviorForLogin(userAuthentic *shelterUser.UserAuthentic, refreshToken shelterText.Token, accessToken pkgText.JwtToken) (*userTestUtility.UserGetterMock, *userTestUtility.PasswordCheckerMock, *userTestUtility.RefreshTokenIssuerMock, *userTestUtility.AccessTokenIssuerMock) {
	var userGetter = &userTestUtility.UserGetterMock{
		FakeExecute: func(identifier pkgText.Identifier) (*shelterUser.UserAuthentic, error) {
			return userAuthentic, nil
		},
	}

	var passwordChecker = &userTestUtility.PasswordCheckerMock{
		FakeExecute: func(entry entryAuth.AuthLoginner) (pkgText.Identifier, error) {
			return pkgText.Identifier(""), nil
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
				Identifier:      expectId,
				EmailIdentifier: expectEmail,
			},
			Password: expectPassword,
		}
	}
}

func TestLogin(t *testing.T) {
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
