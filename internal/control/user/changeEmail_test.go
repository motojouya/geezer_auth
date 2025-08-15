package user_test

import (
	//"errors"
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

func getBehaviorForChangeEmail(t *testing.T, userAuthentic *shelterUser.UserAuthentic) (*userTestUtility.UserGetterMock, *userTestUtility.EmailSetterMock) {
	var userGetter = &userTestUtility.UserGetterMock{
		FakeExecute: func(identifier pkgText.Identifier) (*shelterUser.UserAuthentic, error) {
			assert.Equal(t, identifier, userAuthentic.Identifier)
			return userAuthentic, nil
		},
	}

	var emailSetter = &userTestUtility.EmailSetterMock{
		FakeExecute: func(entry entryUser.EmailGetter, user *shelterUser.UserAuthentic) error {
			return nil
		},
	}

	return userGetter, emailSetter
}

func getShelterUserAuthenticForChangeEmail(expectId string, expectEmail string) *shelterUser.UserAuthentic {
	var userId uint = 1
	var userIdentifier, _ = pkgText.NewIdentifier(expectId)
	var emailId, _ = pkgText.NewEmail(expectEmail)
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

	var email, _ = pkgText.NewEmail(expectEmail)
	return shelterUser.NewUserAuthentic(userValue, companyRole, &email)
}

func getChangeEmailEntry(expectedEmail string) entryUser.UserChangeEmailRequest {
	return entryUser.UserChangeEmailRequest{
		UserChangeEmail: entryUser.UserChangeEmail{
			Email:     expectedEmail,
		},
	}
}

func getAuthorizationForChangeEmail() *shelterAuth.Authorization {
	return shelterAuth.NewAuthorization([]shelterRole.RolePermission{
		shelterRole.AnonymousPermission,
		shelterRole.RoleLessPermission,
		shelterRole.NewRolePermission("EMPLOYEE", true, true, false, false, 5),
		shelterRole.NewRolePermission("MANAGER", true, true, true, true, 9),
	})
}

func getPkgAuthenticForChangeEmail(expectId string, expectEmail string) *pkgUser.Authentic {
	var userIdentifier, _ = pkgText.NewIdentifier(expectId)
	var emailId, _ = pkgText.NewEmail(expectEmail)
	var email, _ = pkgText.NewEmail(expectEmail)
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

func TestChangeEmail(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectOldEmail = "test01@example.com"
	var expectNewEmail = "test02@example.com"
	var userAuthentic = getShelterUserAuthenticForChangeEmail(expectIdentifier, expectOldEmail)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForChangeEmail()

	var userGetter, emailSetter = getBehaviorForChangeEmail(t, userAuthentic)
	var control = controlUser.NewChangeEmailControl(
		db,
		authorization,
		userGetter,
		emailSetter,
	)

	var entry = getChangeEmailEntry(expectNewEmail)
	var pkgAuthentic = getPkgAuthenticForChangeEmail(expectIdentifier, expectOldEmail)

	var userUpdateResponse, err = controlUser.ChangeEmailExecute(control, entry, pkgAuthentic)

	assert.NoError(t, err)
	assert.Equal(t, expectIdentifier, userUpdateResponse.User.Identifier)
	assert.Equal(t, expectOldEmail, *userUpdateResponse.User.Email)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 1, transactionCalledCount.CommitCalled)
	assert.Equal(t, 0, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)

	t.Logf("User Identifier: %+v", userUpdateResponse)
}

// TODO working err case
