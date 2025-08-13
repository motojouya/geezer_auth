package user_test

import (
	"errors"
	"github.com/google/uuid"
	controlUser "github.com/motojouya/geezer_auth/internal/control/user"
	testUtility "github.com/motojouya/geezer_auth/internal/db/testUtility"
	entryUser "github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	shelterCompany "github.com/motojouya/geezer_auth/internal/shelter/company"
	shelterAuth "github.com/motojouya/geezer_auth/internal/shelter/authorization"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
	shelterText "github.com/motojouya/geezer_auth/internal/shelter/text"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getBehaviorForEmailVerify(userAuthentic *shelterUser.UserAuthentic, accessToken pkgText.JwtToken) (*userGetterMock, *emailVerifierMock, *accessTokenIssuerMock) {
	var userGetter = &userGetterMock{
		execute: func(entry entryUser.UserGetter) (*shelterUser.UserAuthentic, error) {
			return userAuthentic, nil
		},
	}

	var emailVerifier = &emailVerifierMock{
		execute: func(entry entryUser.EmailVeifier, user *shelterUser.UserAuthentic) (*shelterUser.UserAuthentic, error) {
			return userAuthentic, nil
		},
	}

	var accessTokenIssuer = &accessTokenIssuerMock{
		execute: func(user *shelterUser.UserAuthentic) (pkgText.JwtToken, error) {
			return accessToken, nil
		},
	}

	return userGetter, emailVerifier, accessTokenIssuer
}

func getShelterUserAuthenticForEmailVerify(expectId string) *shelterUser.UserAuthentic {
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

func getEmailVerifyEntry() entryUser.UserVerifyEmailRequest {
	return entryUser.UserVerifyEmailRequest{
		UserVerifyEmail: entryUser.UserVerifyEmail{
			Email:       "test@example.com",
			VerifyToken: "test-verify-token",
		},
	}
}

func getAuthorizationForEmailVerify() *shelterAuth.Authorization {
	return shelterAuth.NewAuthorization([]shelterRole.RolePermission{
		shelterRole.AnonymousPermission,
		shelterRole.RoleLessPermission,
		shelterRole.NewRolePermission("EMPLOYEE", true, true, false, false, 5),
		shelterRole.NewRolePermission("MANAGER", true, true, true, true, 9),
	})
}

func getPkgAuthenticForEmailVerify() *pkgUser.Authentic {
	var userIdentifierStr = "US-TESTES"
	var companyIdentifier, _ = text.NewIdentifier("CP-TESTES")
	var companyName, _ = text.NewName("TestCompany")
	var company = user.NewCompany(companyIdentifier, companyName)

	var roleLabel, _ = text.NewLabel("TestRole")
	var roleName, _ = text.NewName("TestRoleName")
	var role = user.NewRole(roleLabel, roleName)
	var roles = []user.Role{role}

	var companyRole = user.NewCompanyRole(company, roles)

	var userIdentifier, _ = text.NewIdentifier(userIdentifierStr)
	var emailId, _ = text.NewEmail("test@gmail.com")
	var email, _ = text.NewEmail("test_2@gmail.com")
	var userName, _ = text.NewName("TestName")
	var botFlag = false
	var updateDate = time.Now()

	var userValue = user.NewUser(userIdentifier, emailId, &email, userName, botFlag, companyRole, updateDate)

	var issuer = "issuer_id"
	var subject = "subject_id"
	var aud01 = "aud1"
	var aud02 = "aud2"
	var audience = []string{aud01, aud02}
	var expiresAt = time.Now()
	var notBefore = time.Now()
	var issuedAt = time.Now()
	var id, _ = uuid.NewUUID()

	return user.NewAuthentic(issuer, subject, audience, expiresAt, notBefore, issuedAt, id.String(), userValue)
}

func TestRegister(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectOldEmail = "test01@example.com"
	var expectNewEmail = "test01@example.com"
	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)
	var userAuthentic = getShelterUserAuthenticForEmailVerify(expectIdentifier)

	var transactionCalledCount = &testUtility.TransactionCalledCount{}
	var db = testUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForEmailVerify()

	var userGetter, emailVerifier, accessTokenIssuer = getBehaviorForEmailVerify(userAuthentic, accessToken)
	var control = controlUser.NewVerifyEmailControl(
		db,
		authorization,
		userGetter,
		emailVerifier,
		accessTokenIssuer,
	)

	var pkgAuthentic = getPkgAuthenticForEmailVerify()
	var entry = getEmailVerifyEntry()

	var emailVerifyResponse, err = controlUser.EmailVerifyExecute(control, entry, nil)

	assert.NoError(t, err)
	assert.Equal(t, expectIdentifier, emailVerifyResponse.User.Identifier)
	assert.Equal(t, expectEmail, emailVerifyResponse.User.Email)
	assert.Equal(t, expectToken, emailVerifyResponse.AccessToken)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 1, transactionCalledCount.CommitCalled)
	assert.Equal(t, 0, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)

	t.Logf("User Identifier: %+v", emailVerifyResponse)
}
