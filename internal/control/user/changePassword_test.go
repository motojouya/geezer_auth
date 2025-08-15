package user_test

import (
	"errors"
	"github.com/google/uuid"
	pkgUser "github.com/motojouya/geezer_auth/pkg/shelter/user"
	shelterAuth "github.com/motojouya/geezer_auth/internal/shelter/authorization"
	userTestUtility "github.com/motojouya/geezer_auth/internal/behavior/user/testUtility"
	controlUser "github.com/motojouya/geezer_auth/internal/control/user"
	dbTestUtility "github.com/motojouya/geezer_auth/internal/db/testUtility"
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

func getBehaviorForChangePassword(t *testing.T, userAuthentic *shelterUser.UserAuthentic) (*userTestUtility.UserGetterMock, *userTestUtility.PasswordSetterMock) {
	var userGetter = &userTestUtility.UserGetterMock{
		FakeExecute: func(identifier pkgText.Identifier) (*shelterUser.UserAuthentic, error) {
			assert.Equal(t, identifier, userAuthentic.Identifier)
			return userAuthentic, nil
		},
	}

	var passwordSetter = &userTestUtility.PasswordSetterMock{
		FakeExecute: func(entry entryUser.PasswordGetter, user *shelterUser.UserAuthentic) error {
			return nil
		},
	}

	return userGetter, passwordSetter
}

func getShelterUserAuthenticForChangePassword(expectId string) *shelterUser.UserAuthentic {
	var userId uint = 1
	var userIdentifier, _ = pkgText.NewIdentifier(expectId)
	var emailId, _ = pkgText.NewEmail("test@example.com")
	var userName, _ = pkgText.NewName("Test User")
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

func getChangePasswordEntry(expectedPassword string) entryUser.UserChangePasswordRequest {
	return entryUser.UserChangePasswordRequest{
		UserChangePassword: entryUser.UserChangePassword{
			Password:     expectedPassword,
		},
	}
}

func getAuthorizationForChangePassword() *shelterAuth.Authorization {
	return shelterAuth.NewAuthorization([]shelterRole.RolePermission{
		shelterRole.AnonymousPermission,
		shelterRole.RoleLessPermission,
		shelterRole.NewRolePermission("EMPLOYEE", true, true, false, false, 5),
		shelterRole.NewRolePermission("MANAGER", true, true, true, true, 9),
	})
}

func getPkgAuthenticForChangePassword(expectId string) *pkgUser.Authentic {
	var userIdentifier, _ = pkgText.NewIdentifier(expectId)
	var emailId, _ = pkgText.NewEmail("test@example.com")
	var email, _ = pkgText.NewEmail("test@example.com")
	var userName, _ = pkgText.NewName("Test User")
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

func TestChangePassword(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectPassword = "password01"
	var userAuthentic = getShelterUserAuthenticForChangePassword(expectIdentifier)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForChangePassword()

	var userGetter, passwordSetter = getBehaviorForChangePassword(t, userAuthentic)
	var control = controlUser.NewChangePasswordControl(
		db,
		authorization,
		userGetter,
		passwordSetter,
	)

	var entry = getChangePasswordEntry(expectPassword)
	var pkgAuthentic = getPkgAuthenticForChangePassword(expectIdentifier)

	var userUpdateResponse, err = controlUser.ChangePasswordExecute(control, entry, pkgAuthentic)

	assert.NoError(t, err)
	assert.Equal(t, expectIdentifier, userUpdateResponse.User.Identifier)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 1, transactionCalledCount.CommitCalled)
	assert.Equal(t, 0, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)

	t.Logf("User Identifier: %+v", userUpdateResponse)
}

func TestChangePasswordErrAuth(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectPassword = "password01"
	var userAuthentic = getShelterUserAuthenticForChangePassword(expectIdentifier)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForChangePassword()

	var userGetter, passwordSetter = getBehaviorForChangePassword(t, userAuthentic)
	var control = controlUser.NewChangePasswordControl(
		db,
		authorization,
		userGetter,
		passwordSetter,
	)

	var entry = getChangePasswordEntry(expectPassword)

	var _, err = controlUser.ChangePasswordExecute(control, entry, nil)

	assert.Error(t, err)
}

func TestChangePasswordErrGet(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectPassword = "password01"
	var userAuthentic = getShelterUserAuthenticForChangePassword(expectIdentifier)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForChangePassword()

	var userGetter, passwordSetter = getBehaviorForChangePassword(t, userAuthentic)
	userGetter.FakeExecute = func(identifier pkgText.Identifier) (*shelterUser.UserAuthentic, error) {
		return nil, errors.New("user not found")
	}
	var control = controlUser.NewChangePasswordControl(
		db,
		authorization,
		userGetter,
		passwordSetter,
	)

	var entry = getChangePasswordEntry(expectPassword)
	var pkgAuthentic = getPkgAuthenticForChangePassword(expectIdentifier)

	var _, err = controlUser.ChangePasswordExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)
}

func TestChangePasswordErrChange(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectPassword = "password01"
	var userAuthentic = getShelterUserAuthenticForChangePassword(expectIdentifier)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForChangePassword()

	var userGetter, passwordSetter = getBehaviorForChangePassword(t, userAuthentic)
	passwordSetter.FakeExecute = func(entry entryUser.PasswordGetter, user *shelterUser.UserAuthentic) error {
		return errors.New("password change failed")
	}
	var control = controlUser.NewChangePasswordControl(
		db,
		authorization,
		userGetter,
		passwordSetter,
	)

	var entry = getChangePasswordEntry(expectPassword)
	var pkgAuthentic = getPkgAuthenticForChangePassword(expectIdentifier)

	var _, err = controlUser.ChangePasswordExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)
}
