package role_test

import (
	"errors"
	"github.com/google/uuid"
	roleTestUtility "github.com/motojouya/geezer_auth/internal/behavior/role/testUtility"
	controlRole "github.com/motojouya/geezer_auth/internal/control/role"
	dbTestUtility "github.com/motojouya/geezer_auth/internal/db/testUtility"
	entryCommon "github.com/motojouya/geezer_auth/internal/entry/transfer/common"
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

func getBehaviorForGetRole(roles []shelterRole.Role) *roleTestUtility.AllRoleGetterMock {
	var roleGetter = &roleTestUtility.AllRoleGetterMock{
		FakeExecute: func() ([]shelterRole.Role, error) {
			return roles, nil
		},
	}

	return roleGetter
}

func getShelterUserAuthenticForGetUser(expectId string) *shelterUser.UserAuthentic {
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

	var email, _ = pkgText.NewEmail("test_2@gmail.com")
	return shelterUser.NewUserAuthentic(userValue, companyRole, &email)
}

func getAuthorizationForGetRole() *shelterAuth.Authorization {
	return shelterAuth.NewAuthorization([]shelterRole.RolePermission{
		shelterRole.AnonymousPermission,
		shelterRole.RoleLessPermission,
		shelterRole.NewRolePermission("EMPLOYEE", true, true, false, false, 5),
		shelterRole.NewRolePermission("MANAGER", true, true, true, true, 9),
	})
}

func getPkgAuthenticForGetRole() *pkgUser.Authentic {
	var userIdentifier, _ = pkgText.NewIdentifier("US-TESTES")
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

func getRole(expectLabel string) shelterRole.Role {
	var name, _ = pkgText.NewName("TestRole")
	var label, _ = pkgText.NewLabel(expectLabel)
	var description, _ = shelterText.NewText("Role for testing")
	var registeredDate = time.Now()

	return shelterRole.NewRole(name, label, description, registeredDate)
}

func TestGetRole(t *testing.T) {
	var expectLabel01 = "ROLE_ONE"
	var role01 = getRole(expectLabel01)
	var expectLabel02 = "ROLE_TWO"
	var role02 = getRole(expectLabel02)
	var roles = []shelterRole.Role{role01, role02}

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForGetRole()

	var roleGetter = getBehaviorForGetRole(roles)
	var control = controlRole.NewGetRoleControl(
		db,
		authorization,
		roleGetter,
	)

	var entry = entryCommon.Empty{}
	var pkgAuthentic = getPkgAuthenticForGetRole()

	var getRoleResponse, err = controlRole.GetRoleExecute(control, entry, pkgAuthentic)

	assert.NoError(t, err)
	assert.NotNil(t, getRoleResponse)
	assert.Equal(t, 2, len(getRoleResponse.Roles))
	assert.Equal(t, expectLabel01, getRoleResponse.Roles[0].Label)
	assert.Equal(t, expectLabel02, getRoleResponse.Roles[1].Label)

	t.Logf("get roles: %+v", getRoleResponse)
}

func TestGetRoleErrAuth(t *testing.T) {
	var expectLabel01 = "ROLE_ONE"
	var role01 = getRole(expectLabel01)
	var expectLabel02 = "ROLE_TWO"
	var role02 = getRole(expectLabel02)
	var roles = []shelterRole.Role{role01, role02}

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForGetRole()

	var roleGetter = getBehaviorForGetRole(roles)
	var control = controlRole.NewGetRoleControl(
		db,
		authorization,
		roleGetter,
	)

	var entry = entryCommon.Empty{}

	var _, err = controlRole.GetRoleExecute(control, entry, nil)

	assert.Error(t, err)
}

func TestGetRoleErrGet(t *testing.T) {
	var expectLabel01 = "ROLE_ONE"
	var role01 = getRole(expectLabel01)
	var expectLabel02 = "ROLE_TWO"
	var role02 = getRole(expectLabel02)
	var roles = []shelterRole.Role{role01, role02}

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForGetRole()

	var roleGetter = getBehaviorForGetRole(roles)
	roleGetter.FakeExecute = func() ([]shelterRole.Role, error) {
		return nil, errors.New("failed to get roles")
	}
	var control = controlRole.NewGetRoleControl(
		db,
		authorization,
		roleGetter,
	)

	var entry = entryCommon.Empty{}
	var pkgAuthentic = getPkgAuthenticForGetRole()

	var _, err = controlRole.GetRoleExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)
}
