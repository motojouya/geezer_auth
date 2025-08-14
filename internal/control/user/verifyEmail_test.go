package user_test

import (
	"errors"
	"github.com/google/uuid"
	userTestUtility "github.com/motojouya/geezer_auth/internal/behavior/user/testUtility"
	controlUser "github.com/motojouya/geezer_auth/internal/control/user"
	dbTestUtility "github.com/motojouya/geezer_auth/internal/db/testUtility"
	entryUser "github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	shelterAuth "github.com/motojouya/geezer_auth/internal/shelter/authorization"
	shelterCompany "github.com/motojouya/geezer_auth/internal/shelter/company"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
	shelterText "github.com/motojouya/geezer_auth/internal/shelter/text"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	pkgUser "github.com/motojouya/geezer_auth/pkg/shelter/user"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getBehaviorForEmailVerify(t *testing.T, userAuthentic *shelterUser.UserAuthentic, email pkgText.Email, accessToken pkgText.JwtToken) (*userTestUtility.UserGetterMock, *userTestUtility.EmailVerifierMock, *userTestUtility.AccessTokenIssuerMock) {
	var userGetter = &userTestUtility.UserGetterMock{
		FakeExecute: func(identifier pkgText.Identifier) (*shelterUser.UserAuthentic, error) {
			assert.Equal(t, identifier, userAuthentic.Identifier)
			return userAuthentic, nil
		},
	}

	var emailVerifier = &userTestUtility.EmailVerifierMock{
		FakeExecute: func(entry entryUser.EmailVerifier, user *shelterUser.UserAuthentic) (*shelterUser.UserAuthentic, error) {
			emailArg, err := entry.GetEmail()
			assert.NoError(t, err)
			assert.Equal(t, emailArg, email)
			userAuthentic.Email = &email
			return userAuthentic, nil
		},
	}

	var accessTokenIssuer = &userTestUtility.AccessTokenIssuerMock{
		FakeExecute: func(user *shelterUser.UserAuthentic) (pkgText.JwtToken, error) {
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

func getEmailVerifyEntry(expectEmail string) entryUser.UserVerifyEmailRequest {
	return entryUser.UserVerifyEmailRequest{
		UserVerifyEmail: entryUser.UserVerifyEmail{
			Email:       expectEmail,
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

func getPkgAuthenticForEmailVerify(expectId string, expectEmail string) *pkgUser.Authentic {
	var userIdentifier, _ = pkgText.NewIdentifier(expectId)
	var emailId, _ = pkgText.NewEmail("test@gmail.com")
	var email, _ = pkgText.NewEmail(expectEmail)
	var userName, _ = pkgText.NewName("TestName")
	var botFlag = false
	var updateDate = time.Now()

	var userValue = pkgUser.NewUser(userIdentifier, emailId, &email, userName, botFlag, nil, updateDate)

	var issuer = "issuer_id"
	var subject = "subject_id"
	var aud01 = "aud1"
	var aud02 = "aud2"
	var audience = []string{aud01, aud02}
	var expiresAt = time.Now()
	var notBefore = time.Now()
	var issuedAt = time.Now()
	var id, _ = uuid.NewUUID()

	return pkgUser.NewAuthentic(issuer, subject, audience, expiresAt, notBefore, issuedAt, id.String(), userValue)
}

func TestEmailVerifier(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectOldEmail = "test01@example.com"
	var expectNewEmail = "test02@example.com"
	var email, _ = pkgText.NewEmail(expectNewEmail)
	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)
	var userAuthentic = getShelterUserAuthenticForEmailVerify(expectIdentifier)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)

	var authorization = getAuthorizationForEmailVerify()

	var userGetter, emailVerifier, accessTokenIssuer = getBehaviorForEmailVerify(t, userAuthentic, email, accessToken)
	var control = controlUser.NewVerifyEmailControl(
		db,
		authorization,
		userGetter,
		emailVerifier,
		accessTokenIssuer,
	)

	var pkgAuthentic = getPkgAuthenticForEmailVerify(expectIdentifier, expectOldEmail)
	var entry = getEmailVerifyEntry(expectNewEmail)

	var emailVerifyResponse, err = controlUser.EmailVerifyExecute(control, entry, pkgAuthentic)

	assert.NoError(t, err)
	assert.Equal(t, expectIdentifier, emailVerifyResponse.User.Identifier)
	assert.NotNil(t, emailVerifyResponse.User.Email)
	assert.Equal(t, expectNewEmail, *emailVerifyResponse.User.Email)
	assert.Equal(t, expectToken, emailVerifyResponse.AccessToken)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 1, transactionCalledCount.CommitCalled)
	assert.Equal(t, 0, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)

	t.Logf("User Identifier: %+v", emailVerifyResponse)
}

func TestEmailVerifierErrRole(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectNewEmail = "test02@example.com"
	var email, _ = pkgText.NewEmail(expectNewEmail)
	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)
	var userAuthentic = getShelterUserAuthenticForEmailVerify(expectIdentifier)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)

	var authorization = getAuthorizationForEmailVerify()

	var userGetter, emailVerifier, accessTokenIssuer = getBehaviorForEmailVerify(t, userAuthentic, email, accessToken)
	var control = controlUser.NewVerifyEmailControl(
		db,
		authorization,
		userGetter,
		emailVerifier,
		accessTokenIssuer,
	)

	var entry = getEmailVerifyEntry(expectNewEmail)

	var _, err = controlUser.EmailVerifyExecute(control, entry, nil)

	assert.Error(t, err)
	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestEmailVerifierErrGet(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectOldEmail = "test01@example.com"
	var expectNewEmail = "test02@example.com"
	var email, _ = pkgText.NewEmail(expectNewEmail)
	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)
	var userAuthentic = getShelterUserAuthenticForEmailVerify(expectIdentifier)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)

	var authorization = getAuthorizationForEmailVerify()

	var userGetter, emailVerifier, accessTokenIssuer = getBehaviorForEmailVerify(t, userAuthentic, email, accessToken)
	userGetter.FakeExecute = func(identifier pkgText.Identifier) (*shelterUser.UserAuthentic, error) {
		return nil, errors.New("get user error")
	}

	var control = controlUser.NewVerifyEmailControl(
		db,
		authorization,
		userGetter,
		emailVerifier,
		accessTokenIssuer,
	)

	var pkgAuthentic = getPkgAuthenticForEmailVerify(expectIdentifier, expectOldEmail)
	var entry = getEmailVerifyEntry(expectNewEmail)

	var _, err = controlUser.EmailVerifyExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)
	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestEmailVerifierErrVerify(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectOldEmail = "test01@example.com"
	var expectNewEmail = "test02@example.com"
	var email, _ = pkgText.NewEmail(expectNewEmail)
	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)
	var userAuthentic = getShelterUserAuthenticForEmailVerify(expectIdentifier)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)

	var authorization = getAuthorizationForEmailVerify()

	var userGetter, emailVerifier, accessTokenIssuer = getBehaviorForEmailVerify(t, userAuthentic, email, accessToken)
	emailVerifier.FakeExecute = func(entry entryUser.EmailVerifier, user *shelterUser.UserAuthentic) (*shelterUser.UserAuthentic, error) {
		return nil, errors.New("verify email error")
	}

	var control = controlUser.NewVerifyEmailControl(
		db,
		authorization,
		userGetter,
		emailVerifier,
		accessTokenIssuer,
	)

	var pkgAuthentic = getPkgAuthenticForEmailVerify(expectIdentifier, expectOldEmail)
	var entry = getEmailVerifyEntry(expectNewEmail)

	var _, err = controlUser.EmailVerifyExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)
	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestEmailVerifierErrIssue(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectOldEmail = "test01@example.com"
	var expectNewEmail = "test02@example.com"
	var email, _ = pkgText.NewEmail(expectNewEmail)
	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)
	var userAuthentic = getShelterUserAuthenticForEmailVerify(expectIdentifier)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)

	var authorization = getAuthorizationForEmailVerify()

	var userGetter, emailVerifier, accessTokenIssuer = getBehaviorForEmailVerify(t, userAuthentic, email, accessToken)
	accessTokenIssuer.FakeExecute = func(user *shelterUser.UserAuthentic) (pkgText.JwtToken, error) {
		return pkgText.JwtToken(""), errors.New("issue access token error")
	}

	var control = controlUser.NewVerifyEmailControl(
		db,
		authorization,
		userGetter,
		emailVerifier,
		accessTokenIssuer,
	)

	var pkgAuthentic = getPkgAuthenticForEmailVerify(expectIdentifier, expectOldEmail)
	var entry = getEmailVerifyEntry(expectNewEmail)

	var _, err = controlUser.EmailVerifyExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)
	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}
