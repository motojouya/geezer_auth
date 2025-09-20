package companyUser_test

import (
	"errors"
	"github.com/google/uuid"
	companyTestUtility "github.com/motojouya/geezer_auth/internal/behavior/company/testUtility"
	userTestUtility "github.com/motojouya/geezer_auth/internal/behavior/user/testUtility"
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

func getBehaviorForAccept(userAuthentic *shelterUser.UserAuthentic, company *shelterCompany.Company, role shelterRole.Role, accessToken pkgText.JwtToken) (*userTestUtility.UserGetterMock, *companyTestUtility.CompanyGetterMock, *companyTestUtility.InviteTokenCheckerMock, *companyTestUtility.RoleAssignerMock, *userTestUtility.AccessTokenIssuerMock) {
	var userGetter = &userTestUtility.UserGetterMock{
		FakeExecute: func(identifier pkgText.Identifier) (*shelterUser.UserAuthentic, error) {
			return userAuthentic, nil
		},
	}

	var companyGetter = &companyTestUtility.CompanyGetterMock{
		FakeExecute: func(entry entryCompany.CompanyGetter) (*shelterCompany.Company, error) {
			return company, nil
		},
	}

	var inviteTokenChecker = &companyTestUtility.InviteTokenCheckerMock{
		FakeExecute: func(entry entryCompanyUser.InviteTokenGetter, company shelterCompany.Company) (shelterRole.Role, error) {
			return role, nil
		},
	}

	var roleAssigner = &companyTestUtility.RoleAssignerMock{
		FakeExecute: func(company shelterCompany.Company, userAuthentic *shelterUser.UserAuthentic, role shelterRole.Role) (*shelterUser.UserAuthentic, error) {
			return userAuthentic, nil
		},
	}

	var accessTokenIssuer = &userTestUtility.AccessTokenIssuerMock{
		FakeExecute: func(user *shelterUser.UserAuthentic) (pkgText.JwtToken, error) {
			return accessToken, nil
		},
	}

	return userGetter, companyGetter, inviteTokenChecker, roleAssigner, accessTokenIssuer
}

func getShelterUserAuthenticForAccept(expectId string) *shelterUser.UserAuthentic {
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

func getInviteAcceptEntry(expectId string, expectToken string) entryCompanyUser.CompanyUserAcceptRequest {
	return entryCompanyUser.CompanyUserAcceptRequest{
		CompanyGetRequest: entryCompany.CompanyGetRequest{
			CompanyGet: entryCompany.CompanyGet{
				Identifier: expectId,
			},
		},
		CompanyAccept: entryCompanyUser.CompanyAccept{
			Token: expectToken,
		},
	}
}

func getAuthorizationForAccept() *shelterAuth.Authorization {
	return shelterAuth.NewAuthorization([]shelterRole.RolePermission{
		shelterRole.AnonymousPermission,
		shelterRole.RoleLessPermission,
		shelterRole.NewRolePermission("EMPLOYEE", true, true, false, false, 5),
		shelterRole.NewRolePermission("MANAGER", true, true, true, true, 9),
	})
}

func getRoleForAccept(expectLabel string) shelterRole.Role {
	var label, _ = pkgText.NewLabel(expectLabel)
	var name, _ = pkgText.NewName("TestRole")
	var description, _ = shelterText.NewText("Role for testing")
	var registeredDate = time.Now()

	return shelterRole.NewRole(name, label, description, registeredDate)
}

func getPkgAuthenticForAccept(expectId string) *pkgUser.Authentic {
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

func getShelterCompanyForAccept(expectId string) shelterCompany.Company {
	var companyIdentifier, _ = pkgText.NewIdentifier(expectId)
	var companyId uint = 1
	var companyName, _ = pkgText.NewName("TestCompany")
	var companyRegisteredDate = time.Now()

	return shelterCompany.NewCompany(companyId, companyIdentifier, companyName, companyRegisteredDate)
}

func TestAccept(t *testing.T) {
	var expectUserId = "US-TESTES"
	var userAuthentic = getShelterUserAuthenticForAccept(expectUserId)

	var expectCompanyId = "CP-TESTES"
	var company = getShelterCompanyForAccept(expectCompanyId)

	var expectRoleLabel = "EMPLOYEE"
	var role = getRoleForAccept(expectRoleLabel)

	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForAccept()

	var userGetter, companyGetter, inviteTokenChecker, roleAssigner, accessTokenIssuer = getBehaviorForAccept(userAuthentic, &company, role, accessToken)
	var control = controlCompanyUser.NewAcceptControl(
		db,
		authorization,
		userGetter,
		companyGetter,
		inviteTokenChecker,
		roleAssigner,
		accessTokenIssuer,
	)

	var entry = getInviteAcceptEntry(expectCompanyId, "TestToken")
	var pkgAuthentic = getPkgAuthenticForAccept(expectUserId)

	var companyUserAcceptResponse, err = controlCompanyUser.AcceptExecute(control, entry, pkgAuthentic)

	assert.NoError(t, err)
	assert.Equal(t, expectUserId, companyUserAcceptResponse.User.Identifier)
	assert.Equal(t, expectToken, companyUserAcceptResponse.AccessToken)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 1, transactionCalledCount.CommitCalled)
	assert.Equal(t, 0, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)

	t.Logf("company: %+v", companyUserAcceptResponse)
}

func TestAcceptErrAuth(t *testing.T) {
	var expectUserId = "US-TESTES"
	var userAuthentic = getShelterUserAuthenticForAccept(expectUserId)

	var expectCompanyId = "CP-TESTES"
	var company = getShelterCompanyForAccept(expectCompanyId)

	var expectRoleLabel = "EMPLOYEE"
	var role = getRoleForAccept(expectRoleLabel)

	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForAccept()

	var userGetter, companyGetter, inviteTokenChecker, roleAssigner, accessTokenIssuer = getBehaviorForAccept(userAuthentic, &company, role, accessToken)
	var control = controlCompanyUser.NewAcceptControl(
		db,
		authorization,
		userGetter,
		companyGetter,
		inviteTokenChecker,
		roleAssigner,
		accessTokenIssuer,
	)

	var entry = getInviteAcceptEntry(expectCompanyId, "TestToken")

	var _, err = controlCompanyUser.AcceptExecute(control, entry, nil)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestAcceptErrGetUser(t *testing.T) {
	var expectUserId = "US-TESTES"
	var userAuthentic = getShelterUserAuthenticForAccept(expectUserId)

	var expectCompanyId = "CP-TESTES"
	var company = getShelterCompanyForAccept(expectCompanyId)

	var expectRoleLabel = "EMPLOYEE"
	var role = getRoleForAccept(expectRoleLabel)

	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForAccept()

	var userGetter, companyGetter, inviteTokenChecker, roleAssigner, accessTokenIssuer = getBehaviorForAccept(userAuthentic, &company, role, accessToken)
	userGetter.FakeExecute = func(identifier pkgText.Identifier) (*shelterUser.UserAuthentic, error) {
		return nil, errors.New("db error")
	}
	var control = controlCompanyUser.NewAcceptControl(
		db,
		authorization,
		userGetter,
		companyGetter,
		inviteTokenChecker,
		roleAssigner,
		accessTokenIssuer,
	)

	var entry = getInviteAcceptEntry(expectCompanyId, "TestToken")
	var pkgAuthentic = getPkgAuthenticForAccept(expectUserId)

	var _, err = controlCompanyUser.AcceptExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestAcceptErrGetUserNil(t *testing.T) {
	var expectUserId = "US-TESTES"

	var expectCompanyId = "CP-TESTES"
	var company = getShelterCompanyForAccept(expectCompanyId)

	var expectRoleLabel = "EMPLOYEE"
	var role = getRoleForAccept(expectRoleLabel)

	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForAccept()

	var userGetter, companyGetter, inviteTokenChecker, roleAssigner, accessTokenIssuer = getBehaviorForAccept(nil, &company, role, accessToken)
	var control = controlCompanyUser.NewAcceptControl(
		db,
		authorization,
		userGetter,
		companyGetter,
		inviteTokenChecker,
		roleAssigner,
		accessTokenIssuer,
	)

	var entry = getInviteAcceptEntry(expectCompanyId, "TestToken")
	var pkgAuthentic = getPkgAuthenticForAccept(expectUserId)

	var _, err = controlCompanyUser.AcceptExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestAcceptErrGetCompany(t *testing.T) {
	var expectUserId = "US-TESTES"
	var userAuthentic = getShelterUserAuthenticForAccept(expectUserId)

	var expectCompanyId = "CP-TESTES"
	var company = getShelterCompanyForAccept(expectCompanyId)

	var expectRoleLabel = "EMPLOYEE"
	var role = getRoleForAccept(expectRoleLabel)

	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForAccept()

	var userGetter, companyGetter, inviteTokenChecker, roleAssigner, accessTokenIssuer = getBehaviorForAccept(userAuthentic, &company, role, accessToken)
	companyGetter.FakeExecute = func(entry entryCompany.CompanyGetter) (*shelterCompany.Company, error) {
		return nil, errors.New("db error")
	}
	var control = controlCompanyUser.NewAcceptControl(
		db,
		authorization,
		userGetter,
		companyGetter,
		inviteTokenChecker,
		roleAssigner,
		accessTokenIssuer,
	)

	var entry = getInviteAcceptEntry(expectCompanyId, "TestToken")
	var pkgAuthentic = getPkgAuthenticForAccept(expectUserId)

	var _, err = controlCompanyUser.AcceptExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestAcceptErrGetCompanyNil(t *testing.T) {
	var expectUserId = "US-TESTES"
	var userAuthentic = getShelterUserAuthenticForAccept(expectUserId)

	var expectCompanyId = "CP-TESTES"

	var expectRoleLabel = "EMPLOYEE"
	var role = getRoleForAccept(expectRoleLabel)

	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForAccept()

	var userGetter, companyGetter, inviteTokenChecker, roleAssigner, accessTokenIssuer = getBehaviorForAccept(userAuthentic, nil, role, accessToken)
	var control = controlCompanyUser.NewAcceptControl(
		db,
		authorization,
		userGetter,
		companyGetter,
		inviteTokenChecker,
		roleAssigner,
		accessTokenIssuer,
	)

	var entry = getInviteAcceptEntry(expectCompanyId, "TestToken")
	var pkgAuthentic = getPkgAuthenticForAccept(expectUserId)

	var _, err = controlCompanyUser.AcceptExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestAcceptErrCheck(t *testing.T) {
	var expectUserId = "US-TESTES"
	var userAuthentic = getShelterUserAuthenticForAccept(expectUserId)

	var expectCompanyId = "CP-TESTES"
	var company = getShelterCompanyForAccept(expectCompanyId)

	var expectRoleLabel = "EMPLOYEE"
	var role = getRoleForAccept(expectRoleLabel)

	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForAccept()

	var userGetter, companyGetter, inviteTokenChecker, roleAssigner, accessTokenIssuer = getBehaviorForAccept(userAuthentic, &company, role, accessToken)
	inviteTokenChecker.FakeExecute = func(entry entryCompanyUser.InviteTokenGetter, company shelterCompany.Company) (shelterRole.Role, error) {
		return shelterRole.Role{}, errors.New("db error")
	}
	var control = controlCompanyUser.NewAcceptControl(
		db,
		authorization,
		userGetter,
		companyGetter,
		inviteTokenChecker,
		roleAssigner,
		accessTokenIssuer,
	)

	var entry = getInviteAcceptEntry(expectCompanyId, "TestToken")
	var pkgAuthentic = getPkgAuthenticForAccept(expectUserId)

	var _, err = controlCompanyUser.AcceptExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestAcceptErrAssign(t *testing.T) {
	var expectUserId = "US-TESTES"
	var userAuthentic = getShelterUserAuthenticForAccept(expectUserId)

	var expectCompanyId = "CP-TESTES"
	var company = getShelterCompanyForAccept(expectCompanyId)

	var expectRoleLabel = "EMPLOYEE"
	var role = getRoleForAccept(expectRoleLabel)

	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForAccept()

	var userGetter, companyGetter, inviteTokenChecker, roleAssigner, accessTokenIssuer = getBehaviorForAccept(userAuthentic, &company, role, accessToken)
	roleAssigner.FakeExecute = func(company shelterCompany.Company, userAuthentic *shelterUser.UserAuthentic, role shelterRole.Role) (*shelterUser.UserAuthentic, error) {
		return nil, errors.New("db error")
	}
	var control = controlCompanyUser.NewAcceptControl(
		db,
		authorization,
		userGetter,
		companyGetter,
		inviteTokenChecker,
		roleAssigner,
		accessTokenIssuer,
	)

	var entry = getInviteAcceptEntry(expectCompanyId, "TestToken")
	var pkgAuthentic = getPkgAuthenticForAccept(expectUserId)

	var _, err = controlCompanyUser.AcceptExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestAcceptErrIssue(t *testing.T) {
	var expectUserId = "US-TESTES"
	var userAuthentic = getShelterUserAuthenticForAccept(expectUserId)

	var expectCompanyId = "CP-TESTES"
	var company = getShelterCompanyForAccept(expectCompanyId)

	var expectRoleLabel = "EMPLOYEE"
	var role = getRoleForAccept(expectRoleLabel)

	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForAccept()

	var userGetter, companyGetter, inviteTokenChecker, roleAssigner, accessTokenIssuer = getBehaviorForAccept(userAuthentic, &company, role, accessToken)
	accessTokenIssuer.FakeExecute = func(user *shelterUser.UserAuthentic) (pkgText.JwtToken, error) {
		return "", errors.New("db error")
	}
	var control = controlCompanyUser.NewAcceptControl(
		db,
		authorization,
		userGetter,
		companyGetter,
		inviteTokenChecker,
		roleAssigner,
		accessTokenIssuer,
	)

	var entry = getInviteAcceptEntry(expectCompanyId, "TestToken")
	var pkgAuthentic = getPkgAuthenticForAccept(expectUserId)

	var _, err = controlCompanyUser.AcceptExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}
