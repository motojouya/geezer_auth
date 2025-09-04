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

func getBehaviorForInviteIssue(company *shelterCompany.Company, role *shelterRole.Role, token shelterText.Token) (*companyTestUtility.CompanyGetterMock, *roleTestUtility.RoleGetterMock, *companyTestUtility.InviteTokenIssuerMock) {
	var companyGetter = &companyTestUtility.CompanyGetterMock{
		FakeExecute: func(entry entryCompany.CompanyGetter) (*shelterCompany.Company, error) {
			return company, nil
		},
	}

	var roleGetter = &roleTestUtility.RoleGetterMock{
		FakeExecute: func(entry entryCompanyUser.RoleGetter) (*shelterRole.Role, error) {
			return role, nil
		},
	}

	var inviteTokenIssuer = &companyTestUtility.InviteTokenIssuerMock{
		FakeExecute: func(company shelterCompany.Company, role shelterRole.Role) (shelterText.Token, error) {
			return token, nil
		},
	}

	return companyGetter, roleGetter, inviteTokenIssuer
}

func getShelterUserAuthenticForCreate(expectId string) *shelterUser.UserAuthentic {
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

func getInviteIssueEntry(expectId string, expectLabel string) entryCompanyUser.CompanyUserInviteRequest {
	return entryCompanyUser.CompanyUserInviteRequest{
		CompanyGetRequest: entryCompany.CompanyGetRequest{
			CompanyGet: entryCompany.CompanyGet{
				Identifier: expectId,
			},
		},
		RoleInvite: entryCompanyUser.RoleInvite{
			RoleLabel: expectLabel,
		},
	}
}

func getAuthorizationForInviteIssue() *shelterAuth.Authorization {
	return shelterAuth.NewAuthorization([]shelterRole.RolePermission{
		shelterRole.AnonymousPermission,
		shelterRole.RoleLessPermission,
		shelterRole.NewRolePermission("EMPLOYEE", true, true, false, false, 5),
		shelterRole.NewRolePermission("MANAGER", true, true, true, true, 9),
	})
}

func getRoleForInviteIssue(expectLabel string) shelterRole.Role {
	var label, _ = pkgText.NewLabel(expectLabel)
	var name, _ = pkgText.NewName("TestRole")
	var description, _ = shelterText.NewText("Role for testing")
	var registeredDate = time.Now()

	return shelterRole.NewRole(name, label, description, registeredDate)
}

func getCompanyRoleForInviteIssue(companyIdentifierStr string, roleLabelStr string) *pkgUser.CompanyRole {
	var companyIdentifier, _ = pkgText.NewIdentifier(companyIdentifierStr)
	var companyName, _ = pkgText.NewName("TestCompany")
	var company = pkgUser.NewCompany(companyIdentifier, companyName)

	var roleLabel, _ = pkgText.NewLabel(roleLabelStr)
	var roleName, _ = pkgText.NewName("TestRoleName")
	var role = pkgUser.NewRole(roleLabel, roleName)
	var roles = []pkgUser.Role{role}

	return pkgUser.NewCompanyRole(company, roles)
}

func getPkgAuthenticForInviteIssue(expectId string, companyRole *pkgUser.CompanyRole) *pkgUser.Authentic {
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

func getShelterCompanyForInviteIssue(expectId string) shelterCompany.Company {
	var companyIdentifier, _ = pkgText.NewIdentifier(expectId)
	var companyId uint = 1
	var companyName, _ = pkgText.NewName("TestCompany")
	var companyRegisteredDate = time.Now()

	return shelterCompany.NewCompany(companyId, companyIdentifier, companyName, companyRegisteredDate)
}

func TestInviteIssue(t *testing.T) {
	var expectUserId = "US-TESTES"

	var expectCompanyId = "CP-TESTES"
	var company = getShelterCompanyForInviteIssue(expectCompanyId)

	var expectRoleLabel = "EMPLOYEE"
	var role = getRoleForInviteIssue(expectRoleLabel)

	var expectUUID, _ = uuid.NewUUID()
	var inviteToken, _ = shelterText.CreateToken(expectUUID)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForInviteIssue()

	var companyGetter, roleGetter, inviteTokenIssuer = getBehaviorForInviteIssue(&company, &role, inviteToken)
	var control = controlCompanyUser.NewInviteControl(
		db,
		authorization,
		companyGetter,
		roleGetter,
		inviteTokenIssuer,
	)

	var entry = getInviteIssueEntry(expectCompanyId, expectRoleLabel)
	var companyRole = getCompanyRoleForInviteIssue(expectCompanyId, "MANAGER")
	var pkgAuthentic = getPkgAuthenticForInviteIssue(expectUserId, companyRole)

	var companyUserInviteResponse, err = controlCompanyUser.InviteExecute(control, entry, pkgAuthentic)

	assert.NoError(t, err)
	assert.Equal(t, string(inviteToken), companyUserInviteResponse.Token)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 1, transactionCalledCount.CommitCalled)
	assert.Equal(t, 0, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)

	t.Logf("company: %+v", companyUserInviteResponse)
}

func TestInviteIssueErrAuth(t *testing.T) {
	var expectUserId = "US-TESTES"

	var expectCompanyId = "CP-TESTES"
	var company = getShelterCompanyForInviteIssue(expectCompanyId)

	var expectRoleLabel = "EMPLOYEE"
	var role = getRoleForInviteIssue(expectRoleLabel)

	var expectUUID, _ = uuid.NewUUID()
	var inviteToken, _ = shelterText.CreateToken(expectUUID)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForInviteIssue()

	var companyGetter, roleGetter, inviteTokenIssuer = getBehaviorForInviteIssue(&company, &role, inviteToken)
	var control = controlCompanyUser.NewInviteControl(
		db,
		authorization,
		companyGetter,
		roleGetter,
		inviteTokenIssuer,
	)

	var entry = getInviteIssueEntry(expectCompanyId, expectRoleLabel)
	var companyRole = getCompanyRoleForInviteIssue(expectCompanyId, "EMPLOYEE")
	var pkgAuthentic = getPkgAuthenticForInviteIssue(expectUserId, companyRole)

	var _, err = controlCompanyUser.InviteExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestInviteIssueErrGetCompany(t *testing.T) {
	var expectUserId = "US-TESTES"

	var expectCompanyId = "CP-TESTES"
	var company = getShelterCompanyForInviteIssue(expectCompanyId)

	var expectRoleLabel = "EMPLOYEE"
	var role = getRoleForInviteIssue(expectRoleLabel)

	var expectUUID, _ = uuid.NewUUID()
	var inviteToken, _ = shelterText.CreateToken(expectUUID)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForInviteIssue()

	var companyGetter, roleGetter, inviteTokenIssuer = getBehaviorForInviteIssue(&company, &role, inviteToken)
	companyGetter.FakeExecute = func(entry entryCompany.CompanyGetter) (*shelterCompany.Company, error) {
		return nil, errors.New("error in get company")
	}
	var control = controlCompanyUser.NewInviteControl(
		db,
		authorization,
		companyGetter,
		roleGetter,
		inviteTokenIssuer,
	)

	var entry = getInviteIssueEntry(expectCompanyId, expectRoleLabel)
	var companyRole = getCompanyRoleForInviteIssue(expectCompanyId, "MANAGER")
	var pkgAuthentic = getPkgAuthenticForInviteIssue(expectUserId, companyRole)

	var _, err = controlCompanyUser.InviteExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestInviteIssueErrCompanyNil(t *testing.T) {
	var expectUserId = "US-TESTES"

	var expectCompanyId = "CP-TESTES"

	var expectRoleLabel = "EMPLOYEE"
	var role = getRoleForInviteIssue(expectRoleLabel)

	var expectUUID, _ = uuid.NewUUID()
	var inviteToken, _ = shelterText.CreateToken(expectUUID)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForInviteIssue()

	var companyGetter, roleGetter, inviteTokenIssuer = getBehaviorForInviteIssue(nil, &role, inviteToken)
	var control = controlCompanyUser.NewInviteControl(
		db,
		authorization,
		companyGetter,
		roleGetter,
		inviteTokenIssuer,
	)

	var entry = getInviteIssueEntry(expectCompanyId, expectRoleLabel)
	var companyRole = getCompanyRoleForInviteIssue(expectCompanyId, "MANAGER")
	var pkgAuthentic = getPkgAuthenticForInviteIssue(expectUserId, companyRole)

	var _, err = controlCompanyUser.InviteExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestInviteIssueErrGetRole(t *testing.T) {
	var expectUserId = "US-TESTES"

	var expectCompanyId = "CP-TESTES"
	var company = getShelterCompanyForInviteIssue(expectCompanyId)

	var expectRoleLabel = "EMPLOYEE"
	var role = getRoleForInviteIssue(expectRoleLabel)

	var expectUUID, _ = uuid.NewUUID()
	var inviteToken, _ = shelterText.CreateToken(expectUUID)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForInviteIssue()

	var companyGetter, roleGetter, inviteTokenIssuer = getBehaviorForInviteIssue(&company, &role, inviteToken)
	roleGetter.FakeExecute = func(entry entryCompanyUser.RoleGetter) (*shelterRole.Role, error) {
		return nil, errors.New("error in get role")
	}
	var control = controlCompanyUser.NewInviteControl(
		db,
		authorization,
		companyGetter,
		roleGetter,
		inviteTokenIssuer,
	)

	var entry = getInviteIssueEntry(expectCompanyId, expectRoleLabel)
	var companyRole = getCompanyRoleForInviteIssue(expectCompanyId, "MANAGER")
	var pkgAuthentic = getPkgAuthenticForInviteIssue(expectUserId, companyRole)

	var _, err = controlCompanyUser.InviteExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestInviteIssueErrRoleNil(t *testing.T) {
	var expectUserId = "US-TESTES"

	var expectCompanyId = "CP-TESTES"
	var company = getShelterCompanyForInviteIssue(expectCompanyId)

	var expectRoleLabel = "EMPLOYEE"

	var expectUUID, _ = uuid.NewUUID()
	var inviteToken, _ = shelterText.CreateToken(expectUUID)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForInviteIssue()

	var companyGetter, roleGetter, inviteTokenIssuer = getBehaviorForInviteIssue(&company, nil, inviteToken)
	var control = controlCompanyUser.NewInviteControl(
		db,
		authorization,
		companyGetter,
		roleGetter,
		inviteTokenIssuer,
	)

	var entry = getInviteIssueEntry(expectCompanyId, expectRoleLabel)
	var companyRole = getCompanyRoleForInviteIssue(expectCompanyId, "MANAGER")
	var pkgAuthentic = getPkgAuthenticForInviteIssue(expectUserId, companyRole)

	var _, err = controlCompanyUser.InviteExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestInviteIssueErrIssue(t *testing.T) {
	var expectUserId = "US-TESTES"

	var expectCompanyId = "CP-TESTES"
	var company = getShelterCompanyForInviteIssue(expectCompanyId)

	var expectRoleLabel = "EMPLOYEE"
	var role = getRoleForInviteIssue(expectRoleLabel)

	var expectUUID, _ = uuid.NewUUID()
	var inviteToken, _ = shelterText.CreateToken(expectUUID)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForInviteIssue()

	var companyGetter, roleGetter, inviteTokenIssuer = getBehaviorForInviteIssue(&company, &role, inviteToken)
	inviteTokenIssuer.FakeExecute = func(company shelterCompany.Company, role shelterRole.Role) (shelterText.Token, error) {
		return shelterText.Token(""), errors.New("error in issue token")
	}
	var control = controlCompanyUser.NewInviteControl(
		db,
		authorization,
		companyGetter,
		roleGetter,
		inviteTokenIssuer,
	)

	var entry = getInviteIssueEntry(expectCompanyId, expectRoleLabel)
	var companyRole = getCompanyRoleForInviteIssue(expectCompanyId, "MANAGER")
	var pkgAuthentic = getPkgAuthenticForInviteIssue(expectUserId, companyRole)

	var _, err = controlCompanyUser.InviteExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}
