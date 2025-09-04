package company_test

import (
	"errors"
	"github.com/google/uuid"
	companyTestUtility "github.com/motojouya/geezer_auth/internal/behavior/company/testUtility"
	roleTestUtility "github.com/motojouya/geezer_auth/internal/behavior/role/testUtility"
	userTestUtility "github.com/motojouya/geezer_auth/internal/behavior/user/testUtility"
	controlCompany "github.com/motojouya/geezer_auth/internal/control/company"
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

func getBehaviorForCreate(t *testing.T, userAuthentic *shelterUser.UserAuthentic, company shelterCompany.Company, role shelterRole.Role) (*companyTestUtility.CompanyCreatorMock, *roleTestUtility.RoleGetterMock, *userTestUtility.UserGetterMock, *companyTestUtility.RoleAssignerMock) {
	var companyCreator = &companyTestUtility.CompanyCreatorMock{
		FakeExecute: func(entry entryCompany.CompanyCreator) (shelterCompany.Company, error) {
			return company, nil
		},
	}

	var roleGetter = &roleTestUtility.RoleGetterMock{
		FakeExecute: func(entry entryCompanyUser.RoleGetter) (*shelterRole.Role, error) {
			return &role, nil
		},
	}

	var userGetter = &userTestUtility.UserGetterMock{
		FakeExecute: func(identififer pkgText.Identifier) (*shelterUser.UserAuthentic, error) {
			assert.Equal(t, string(userAuthentic.Identifier), string(identififer), "User Identifier should match")
			return userAuthentic, nil
		},
	}

	var RoleAssigner = &companyTestUtility.RoleAssignerMock{
		FakeExecute: func(company shelterCompany.Company, user *shelterUser.UserAuthentic, role shelterRole.Role) (*shelterUser.UserAuthentic, error) {
			return userAuthentic, nil
		},
	}

	return companyCreator, roleGetter, userGetter, RoleAssigner
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

func getCreateEntry(expectName string) entryCompany.CompanyCreateRequest {
	return entryCompany.CompanyCreateRequest{
		CompanyCreate: entryCompany.CompanyCreate{
			Name: expectName,
		},
	}
}

func getAuthorizationForCreate() *shelterAuth.Authorization {
	return shelterAuth.NewAuthorization([]shelterRole.RolePermission{
		shelterRole.AnonymousPermission,
		shelterRole.RoleLessPermission,
		shelterRole.NewRolePermission("EMPLOYEE", true, true, false, false, 5),
		shelterRole.NewRolePermission("MANAGER", true, true, true, true, 9),
	})
}

func getRoleForCreate(expectLabel pkgText.Label) shelterRole.Role {
	var name, _ = pkgText.NewName("TestRole")
	var description, _ = shelterText.NewText("Role for testing")
	var registeredDate = time.Now()

	return shelterRole.NewRole(name, expectLabel, description, registeredDate)
}

func getPkgAuthenticForCreate(expectId string) *pkgUser.Authentic {
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

func getShelterCompanyForCreate(expectId string, expectName string) shelterCompany.Company {
	var companyIdentifier, _ = pkgText.NewIdentifier(expectId)
	var companyId uint = 1
	var companyName, _ = pkgText.NewName(expectName)
	var companyRegisteredDate = time.Now()

	return shelterCompany.NewCompany(companyId, companyIdentifier, companyName, companyRegisteredDate)
}

func TestCreate(t *testing.T) {
	var expectUserId = "US-TESTES"
	var userAuthentic = getShelterUserAuthenticForCreate(expectUserId)
	var expectCompanyId = "CP-TESTES"
	var expectCompanyName = "TestCompany"
	var company = getShelterCompanyForCreate(expectCompanyId, expectCompanyName)
	var role = getRoleForCreate(shelterRole.RoleAdminLabel)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForCreate()

	var companyCreator, roleGetter, userGetter, roleAssigner = getBehaviorForCreate(t, userAuthentic, company, role)
	var control = controlCompany.NewCreateControl(
		db,
		authorization,
		companyCreator,
		roleGetter,
		userGetter,
		roleAssigner,
	)

	var entry = getCreateEntry(expectCompanyName)
	var pkgAuthentic = getPkgAuthenticForCreate(expectUserId)

	var companyCreateResponse, err = controlCompany.CreateExecute(control, entry, pkgAuthentic)

	assert.NoError(t, err)
	assert.Equal(t, expectCompanyId, companyCreateResponse.Company.Identifier)
	assert.Equal(t, expectCompanyName, companyCreateResponse.Company.Name)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 1, transactionCalledCount.CommitCalled)
	assert.Equal(t, 0, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)

	t.Logf("company: %+v", companyCreateResponse)
}

func TestCreateErrAuth(t *testing.T) {
	var expectUserId = "US-TESTES"
	var userAuthentic = getShelterUserAuthenticForCreate(expectUserId)
	var expectCompanyId = "CP-TESTES"
	var expectCompanyName = "TestCompany"
	var company = getShelterCompanyForCreate(expectCompanyId, expectCompanyName)
	var role = getRoleForCreate(shelterRole.RoleAdminLabel)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForCreate()

	var companyCreator, roleGetter, userGetter, roleAssigner = getBehaviorForCreate(t, userAuthentic, company, role)
	var control = controlCompany.NewCreateControl(
		db,
		authorization,
		companyCreator,
		roleGetter,
		userGetter,
		roleAssigner,
	)

	var entry = getCreateEntry(expectCompanyName)

	var _, err = controlCompany.CreateExecute(control, entry, nil)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestCreateErrCreate(t *testing.T) {
	var expectUserId = "US-TESTES"
	var userAuthentic = getShelterUserAuthenticForCreate(expectUserId)
	var expectCompanyId = "CP-TESTES"
	var expectCompanyName = "TestCompany"
	var company = getShelterCompanyForCreate(expectCompanyId, expectCompanyName)
	var role = getRoleForCreate(shelterRole.RoleAdminLabel)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForCreate()

	var companyCreator, roleGetter, userGetter, roleAssigner = getBehaviorForCreate(t, userAuthentic, company, role)
	companyCreator.FakeExecute = func(entry entryCompany.CompanyCreator) (shelterCompany.Company, error) {
		return shelterCompany.Company{}, errors.New("failed to create company")
	}
	var control = controlCompany.NewCreateControl(
		db,
		authorization,
		companyCreator,
		roleGetter,
		userGetter,
		roleAssigner,
	)

	var entry = getCreateEntry(expectCompanyName)
	var pkgAuthentic = getPkgAuthenticForCreate(expectUserId)

	var _, err = controlCompany.CreateExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestCreateErrGetRole(t *testing.T) {
	var expectUserId = "US-TESTES"
	var userAuthentic = getShelterUserAuthenticForCreate(expectUserId)
	var expectCompanyId = "CP-TESTES"
	var expectCompanyName = "TestCompany"
	var company = getShelterCompanyForCreate(expectCompanyId, expectCompanyName)
	var role = getRoleForCreate(shelterRole.RoleAdminLabel)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForCreate()

	var companyCreator, roleGetter, userGetter, roleAssigner = getBehaviorForCreate(t, userAuthentic, company, role)
	roleGetter.FakeExecute = func(entry entryCompanyUser.RoleGetter) (*shelterRole.Role, error) {
		return nil, errors.New("failed to get role")
	}
	var control = controlCompany.NewCreateControl(
		db,
		authorization,
		companyCreator,
		roleGetter,
		userGetter,
		roleAssigner,
	)

	var entry = getCreateEntry(expectCompanyName)
	var pkgAuthentic = getPkgAuthenticForCreate(expectUserId)

	var _, err = controlCompany.CreateExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestCreateErrGetRoleNil(t *testing.T) {
	var expectUserId = "US-TESTES"
	var userAuthentic = getShelterUserAuthenticForCreate(expectUserId)
	var expectCompanyId = "CP-TESTES"
	var expectCompanyName = "TestCompany"
	var company = getShelterCompanyForCreate(expectCompanyId, expectCompanyName)
	var role = getRoleForCreate(shelterRole.RoleAdminLabel)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForCreate()

	var companyCreator, roleGetter, userGetter, roleAssigner = getBehaviorForCreate(t, userAuthentic, company, role)
	roleGetter.FakeExecute = func(entry entryCompanyUser.RoleGetter) (*shelterRole.Role, error) {
		return nil, nil
	}
	var control = controlCompany.NewCreateControl(
		db,
		authorization,
		companyCreator,
		roleGetter,
		userGetter,
		roleAssigner,
	)

	var entry = getCreateEntry(expectCompanyName)
	var pkgAuthentic = getPkgAuthenticForCreate(expectUserId)

	var _, err = controlCompany.CreateExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestCreateErrGetUser(t *testing.T) {
	var expectUserId = "US-TESTES"
	var userAuthentic = getShelterUserAuthenticForCreate(expectUserId)
	var expectCompanyId = "CP-TESTES"
	var expectCompanyName = "TestCompany"
	var company = getShelterCompanyForCreate(expectCompanyId, expectCompanyName)
	var role = getRoleForCreate(shelterRole.RoleAdminLabel)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForCreate()

	var companyCreator, roleGetter, userGetter, roleAssigner = getBehaviorForCreate(t, userAuthentic, company, role)
	userGetter.FakeExecute = func(identififer pkgText.Identifier) (*shelterUser.UserAuthentic, error) {
		return nil, errors.New("failed to get user")
	}
	var control = controlCompany.NewCreateControl(
		db,
		authorization,
		companyCreator,
		roleGetter,
		userGetter,
		roleAssigner,
	)

	var entry = getCreateEntry(expectCompanyName)
	var pkgAuthentic = getPkgAuthenticForCreate(expectUserId)

	var _, err = controlCompany.CreateExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}

func TestCreateErrAssign(t *testing.T) {
	var expectUserId = "US-TESTES"
	var userAuthentic = getShelterUserAuthenticForCreate(expectUserId)
	var expectCompanyId = "CP-TESTES"
	var expectCompanyName = "TestCompany"
	var company = getShelterCompanyForCreate(expectCompanyId, expectCompanyName)
	var role = getRoleForCreate(shelterRole.RoleAdminLabel)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForCreate()

	var companyCreator, roleGetter, userGetter, roleAssigner = getBehaviorForCreate(t, userAuthentic, company, role)
	roleAssigner.FakeExecute = func(company shelterCompany.Company, user *shelterUser.UserAuthentic, role shelterRole.Role) (*shelterUser.UserAuthentic, error) {
		return nil, errors.New("failed to assign role")
	}
	var control = controlCompany.NewCreateControl(
		db,
		authorization,
		companyCreator,
		roleGetter,
		userGetter,
		roleAssigner,
	)

	var entry = getCreateEntry(expectCompanyName)
	var pkgAuthentic = getPkgAuthenticForCreate(expectUserId)

	var _, err = controlCompany.CreateExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 0, transactionCalledCount.CommitCalled)
	assert.Equal(t, 1, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)
}
