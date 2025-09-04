package companyUser_test

import (
	"errors"
	"github.com/google/uuid"
	companyTestUtility "github.com/motojouya/geezer_auth/internal/behavior/company/testUtility"
	roleTestUtility "github.com/motojouya/geezer_auth/internal/behavior/role/testUtility"
	controlCompanyUser "github.com/motojouya/geezer_auth/internal/control/companyUser"
	dbTestUtility "github.com/motojouya/geezer_auth/internal/db/testUtility"
	entryCompany "github.com/motojouya/geezer_auth/internal/entry/transfer/company"
	entryCompanyUser "github.com/motojouya/geezer_auth/internal/entry/transfer/companyUser"
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

func getBehaviorForAssign(userAuthentic *shelterUser.UserAuthentic, company *shelterCompany.Company, role *shelterRole.Role) (*companyTestUtility.CompanyGetterMock, *companyTestUtility.UserGetterMock, *roleTestUtility.RoleGetterMock, *companyTestUtility.RoleAssignerMock) {
	var companyGetter = &companyTestUtility.CompanyGetterMock{
		FakeExecute: func(entry entryCompany.CompanyGetter) (*shelterCompany.Company, error) {
			return company, nil
		},
	}

	var userGetter = &companyTestUtility.UserGetterMock{
		FakeExecute: func(entry entryCompanyUser.CompanyUserGetter, company shelterCompany.Company) (*shelterUser.UserAuthentic, error) {
			return userAuthentic, nil
		},
	}

	var roleGetter = &roleTestUtility.RoleGetterMock{
		FakeExecute: func(entry entryCompanyUser.RoleGetter) (*shelterRole.Role, error) {
			return role, nil
		},
	}

	var roleAssigner = &companyTestUtility.RoleAssignerMock{
		FakeExecute: func(company shelterCompany.Company, userAuthentic *shelterUser.UserAuthentic, role shelterRole.Role) (*shelterUser.UserAuthentic, error) {
			return userAuthentic, nil
		},
	}

	return companyGetter, userGetter, roleGetter, roleAssigner
}

func getShelterUserAuthenticForAssign(expectId string) *shelterUser.UserAuthentic {
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

func getInviteAssignEntry(expectCompanyId string, expectUserId string, expectLabel string) entryCompanyUser.CompanyUserAssignRequest {
	return entryCompanyUser.CompanyUserAssignRequest{
		CompanyGetRequest: entryCompany.CompanyGetRequest{
			CompanyGet: entryCompany.CompanyGet{
				Identifier: expectCompanyId,
			},
		},
		RoleAssign: entryCompanyUser.RoleAssign{
			UserIdentifier: expectUserId,
			RoleInvite: entryCompanyUser.RoleInvite{
				RoleLabel: expectLabel,
			},
		},
	}
}

func getAuthorizationForAssign() *shelterAuth.Authorization {
	return shelterAuth.NewAuthorization([]shelterRole.RolePermission{
		shelterRole.AnonymousPermission,
		shelterRole.RoleLessPermission,
		shelterRole.NewRolePermission("EMPLOYEE", true, true, false, false, 5),
		shelterRole.NewRolePermission("MANAGER", true, true, true, true, 9),
	})
}

func getRoleForAssign(expectLabel string) shelterRole.Role {
	var label, _ = pkgText.NewLabel(expectLabel)
	var name, _ = pkgText.NewName("TestRole")
	var description, _ = shelterText.NewText("Role for testing")
	var registeredDate = time.Now()

	return shelterRole.NewRole(name, label, description, registeredDate)
}

func getCompanyRoleForAccept(companyIdentifierStr string, roleLabelStr string) *pkgUser.CompanyRole {
	var companyIdentifier, _ = pkgText.NewIdentifier(companyIdentifierStr)
	var companyName, _ = pkgText.NewName("TestCompany")
	var company = pkgUser.NewCompany(companyIdentifier, companyName)

	var roleLabel, _ = pkgText.NewLabel(roleLabelStr)
	var roleName, _ = pkgText.NewName("TestRoleName")
	var role = pkgUser.NewRole(roleLabel, roleName)
	var roles = []pkgUser.Role{role}

	return pkgUser.NewCompanyRole(company, roles)
}

func getPkgAuthenticForAssign(expectId string, companyRole *pkgUser.CompanyRole) *pkgUser.Authentic {
	var userIdentifier, _ = pkgText.NewIdentifier(expectId)
	var emailId, _ = pkgText.NewEmail("test@example.com")
	var email, _ = pkgText.NewEmail("test@example.com")
	var userName, _ = pkgText.NewName("Test User")
	var botFlag = false
	var updateDate = time.Now()

	var userValue = pkgUser.NewUser(userIdentifier, emailId, &email, userName, botFlag, companyRole, updateDate)

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

func getShelterCompanyForAssign(expectId string) shelterCompany.Company {
	var companyIdentifier, _ = pkgText.NewIdentifier(expectId)
	var companyId uint = 1
	var companyName, _ = pkgText.NewName("TestCompany")
	var companyRegisteredDate = time.Now()

	return shelterCompany.NewCompany(companyId, companyIdentifier, companyName, companyRegisteredDate)
}

func TestAssign(t *testing.T) {
	var expectUserId = "US-TESTES"
	var userAuthentic = getShelterUserAuthenticForAssign(expectUserId)

	var expectCompanyId = "CP-TESTES"
	var company = getShelterCompanyForAssign(expectCompanyId)

	var expectRoleLabel = "EMPLOYEE"
	var role = getRoleForAssign(expectRoleLabel)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForAssign()

	var companyGetter, userGetter, roleGetter, roleAssigner = getBehaviorForAssign(userAuthentic, &company, &role)
	var control = controlCompanyUser.NewAssignControl(
		db,
		authorization,
		companyGetter,
		userGetter,
		roleGetter,
		roleAssigner,
	)

	var entry = getInviteAssignEntry(expectCompanyId, expectUserId, expectRoleLabel)
	var companyRole = getCompanyRoleForAccept(expectCompanyId, "MANAGER")
	var pkgAuthentic = getPkgAuthenticForAssign(expectUserId, companyRole)

	var companyUserAssignResponse, err = controlCompanyUser.AssignExecute(control, entry, pkgAuthentic)

	assert.NoError(t, err)
	assert.Equal(t, expectUserId, companyUserAssignResponse.User.Identifier)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 1, transactionCalledCount.CommitCalled)
	assert.Equal(t, 0, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)

	t.Logf("company: %+v", companyUserAssignResponse)
}

func TestAssignErrAuth(t *testing.T) {
	var expectUserId = "US-TESTES"
	var userAuthentic = getShelterUserAuthenticForAssign(expectUserId)

	var expectCompanyId = "CP-TESTES"
	var company = getShelterCompanyForAssign(expectCompanyId)

	var expectRoleLabel = "EMPLOYEE"
	var role = getRoleForAssign(expectRoleLabel)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForAssign()

	var companyGetter, userGetter, roleGetter, roleAssigner = getBehaviorForAssign(userAuthentic, &company, &role)
	var control = controlCompanyUser.NewAssignControl(
		db,
		authorization,
		companyGetter,
		userGetter,
		roleGetter,
		roleAssigner,
	)

	var entry = getInviteAssignEntry(expectCompanyId, expectUserId, expectRoleLabel)
	var companyRole = getCompanyRoleForAccept(expectCompanyId, "EMPLOYEE")
	var pkgAuthentic = getPkgAuthenticForAssign(expectUserId, companyRole)

	var _, err = controlCompanyUser.AssignExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestAssignErrGetCompany(t *testing.T) {
	var expectUserId = "US-TESTES"
	var userAuthentic = getShelterUserAuthenticForAssign(expectUserId)

	var expectCompanyId = "CP-TESTES"
	var company = getShelterCompanyForAssign(expectCompanyId)

	var expectRoleLabel = "EMPLOYEE"
	var role = getRoleForAssign(expectRoleLabel)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForAssign()

	var companyGetter, userGetter, roleGetter, roleAssigner = getBehaviorForAssign(userAuthentic, &company, &role)
	companyGetter.FakeExecute = func(entry entryCompany.CompanyGetter) (*shelterCompany.Company, error) {
		return nil, errors.New("error get company")
	}
	var control = controlCompanyUser.NewAssignControl(
		db,
		authorization,
		companyGetter,
		userGetter,
		roleGetter,
		roleAssigner,
	)

	var entry = getInviteAssignEntry(expectCompanyId, expectUserId, expectRoleLabel)
	var companyRole = getCompanyRoleForAccept(expectCompanyId, "MANAGER")
	var pkgAuthentic = getPkgAuthenticForAssign(expectUserId, companyRole)

	var _, err = controlCompanyUser.AssignExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestAssignErrGetCompanyNil(t *testing.T) {
	var expectUserId = "US-TESTES"
	var userAuthentic = getShelterUserAuthenticForAssign(expectUserId)

	var expectCompanyId = "CP-TESTES"

	var expectRoleLabel = "EMPLOYEE"
	var role = getRoleForAssign(expectRoleLabel)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForAssign()

	var companyGetter, userGetter, roleGetter, roleAssigner = getBehaviorForAssign(userAuthentic, nil, &role)
	var control = controlCompanyUser.NewAssignControl(
		db,
		authorization,
		companyGetter,
		userGetter,
		roleGetter,
		roleAssigner,
	)

	var entry = getInviteAssignEntry(expectCompanyId, expectUserId, expectRoleLabel)
	var companyRole = getCompanyRoleForAccept(expectCompanyId, "MANAGER")
	var pkgAuthentic = getPkgAuthenticForAssign(expectUserId, companyRole)

	var _, err = controlCompanyUser.AssignExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestAssignErrGetUser(t *testing.T) {
	var expectUserId = "US-TESTES"
	var userAuthentic = getShelterUserAuthenticForAssign(expectUserId)

	var expectCompanyId = "CP-TESTES"
	var company = getShelterCompanyForAssign(expectCompanyId)

	var expectRoleLabel = "EMPLOYEE"
	var role = getRoleForAssign(expectRoleLabel)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForAssign()

	var companyGetter, userGetter, roleGetter, roleAssigner = getBehaviorForAssign(userAuthentic, &company, &role)
	userGetter.FakeExecute = func(entry entryCompanyUser.CompanyUserGetter, company shelterCompany.Company) (*shelterUser.UserAuthentic, error) {
		return nil, errors.New("error get user")
	}
	var control = controlCompanyUser.NewAssignControl(
		db,
		authorization,
		companyGetter,
		userGetter,
		roleGetter,
		roleAssigner,
	)

	var entry = getInviteAssignEntry(expectCompanyId, expectUserId, expectRoleLabel)
	var companyRole = getCompanyRoleForAccept(expectCompanyId, "MANAGER")
	var pkgAuthentic = getPkgAuthenticForAssign(expectUserId, companyRole)

	var _, err = controlCompanyUser.AssignExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestAssignErrGetUserNil(t *testing.T) {
	var expectUserId = "US-TESTES"

	var expectCompanyId = "CP-TESTES"
	var company = getShelterCompanyForAssign(expectCompanyId)

	var expectRoleLabel = "EMPLOYEE"
	var role = getRoleForAssign(expectRoleLabel)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForAssign()

	var companyGetter, userGetter, roleGetter, roleAssigner = getBehaviorForAssign(nil, &company, &role)
	var control = controlCompanyUser.NewAssignControl(
		db,
		authorization,
		companyGetter,
		userGetter,
		roleGetter,
		roleAssigner,
	)

	var entry = getInviteAssignEntry(expectCompanyId, expectUserId, expectRoleLabel)
	var companyRole = getCompanyRoleForAccept(expectCompanyId, "MANAGER")
	var pkgAuthentic = getPkgAuthenticForAssign(expectUserId, companyRole)

	var _, err = controlCompanyUser.AssignExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestAssignErrGetRole(t *testing.T) {
	var expectUserId = "US-TESTES"
	var userAuthentic = getShelterUserAuthenticForAssign(expectUserId)

	var expectCompanyId = "CP-TESTES"
	var company = getShelterCompanyForAssign(expectCompanyId)

	var expectRoleLabel = "EMPLOYEE"
	var role = getRoleForAssign(expectRoleLabel)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForAssign()

	var companyGetter, userGetter, roleGetter, roleAssigner = getBehaviorForAssign(userAuthentic, &company, &role)
	roleGetter.FakeExecute = func(entry entryCompanyUser.RoleGetter) (*shelterRole.Role, error) {
		return nil, errors.New("error get role")
	}
	var control = controlCompanyUser.NewAssignControl(
		db,
		authorization,
		companyGetter,
		userGetter,
		roleGetter,
		roleAssigner,
	)

	var entry = getInviteAssignEntry(expectCompanyId, expectUserId, expectRoleLabel)
	var companyRole = getCompanyRoleForAccept(expectCompanyId, "MANAGER")
	var pkgAuthentic = getPkgAuthenticForAssign(expectUserId, companyRole)

	var _, err = controlCompanyUser.AssignExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestAssignErrGetRoleNil(t *testing.T) {
	var expectUserId = "US-TESTES"
	var userAuthentic = getShelterUserAuthenticForAssign(expectUserId)

	var expectCompanyId = "CP-TESTES"
	var company = getShelterCompanyForAssign(expectCompanyId)

	var expectRoleLabel = "EMPLOYEE"

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForAssign()

	var companyGetter, userGetter, roleGetter, roleAssigner = getBehaviorForAssign(userAuthentic, &company, nil)
	var control = controlCompanyUser.NewAssignControl(
		db,
		authorization,
		companyGetter,
		userGetter,
		roleGetter,
		roleAssigner,
	)

	var entry = getInviteAssignEntry(expectCompanyId, expectUserId, expectRoleLabel)
	var companyRole = getCompanyRoleForAccept(expectCompanyId, "MANAGER")
	var pkgAuthentic = getPkgAuthenticForAssign(expectUserId, companyRole)

	var _, err = controlCompanyUser.AssignExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestAssignErrAssign(t *testing.T) {
	var expectUserId = "US-TESTES"
	var userAuthentic = getShelterUserAuthenticForAssign(expectUserId)

	var expectCompanyId = "CP-TESTES"
	var company = getShelterCompanyForAssign(expectCompanyId)

	var expectRoleLabel = "EMPLOYEE"
	var role = getRoleForAssign(expectRoleLabel)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForAssign()

	var companyGetter, userGetter, roleGetter, roleAssigner = getBehaviorForAssign(userAuthentic, &company, &role)
	roleAssigner.FakeExecute = func(company shelterCompany.Company, userAuthentic *shelterUser.UserAuthentic, role shelterRole.Role) (*shelterUser.UserAuthentic, error) {
		return nil, errors.New("error assign role")
	}
	var control = controlCompanyUser.NewAssignControl(
		db,
		authorization,
		companyGetter,
		userGetter,
		roleGetter,
		roleAssigner,
	)

	var entry = getInviteAssignEntry(expectCompanyId, expectUserId, expectRoleLabel)
	var companyRole = getCompanyRoleForAccept(expectCompanyId, "MANAGER")
	var pkgAuthentic = getPkgAuthenticForAssign(expectUserId, companyRole)

	var _, err = controlCompanyUser.AssignExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}
