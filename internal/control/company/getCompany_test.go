package company_test

import (
	"errors"
	"github.com/google/uuid"
	companyTestUtility "github.com/motojouya/geezer_auth/internal/behavior/company/testUtility"
	controlCompany "github.com/motojouya/geezer_auth/internal/control/company"
	dbTestUtility "github.com/motojouya/geezer_auth/internal/db/testUtility"
	entryCompany "github.com/motojouya/geezer_auth/internal/entry/transfer/company"
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

func getBehaviorForGetCompany(company shelterCompany.Company) *companyTestUtility.CompanyGetterMock {
	var companyGetter = &companyTestUtility.CompanyGetterMock{
		FakeExecute: func(entry entryCompany.CompanyGetter) (*shelterCompany.Company, error) {
			return &company, nil
		},
	}

	return companyGetter
}

func getShelterCompanyForGet(expectId string) shelterCompany.Company {
	var companyIdentifier, _ = pkgText.NewIdentifier(expectId)
	var companyId uint = 1
	var companyName, _ = pkgText.NewName("TestCompany")
	var companyRegisteredDate = time.Now()

	return shelterCompany.NewCompany(companyId, companyIdentifier, companyName, companyRegisteredDate)
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

func getAuthorizationForGetUser() *shelterAuth.Authorization {
	return shelterAuth.NewAuthorization([]shelterRole.RolePermission{
		shelterRole.AnonymousPermission,
		shelterRole.RoleLessPermission,
		shelterRole.NewRolePermission("EMPLOYEE", true, true, false, false, 5),
		shelterRole.NewRolePermission("MANAGER", true, true, true, true, 9),
	})
}

func getCompanyRoleForGetCompany(companyIdentifierStr string, roleLabelStr string) *pkgUser.CompanyRole {
	var companyIdentifier, _ = pkgText.NewIdentifier(companyIdentifierStr)
	var companyName, _ = pkgText.NewName("TestCompany")
	var company = pkgUser.NewCompany(companyIdentifier, companyName)

	var roleLabel, _ = pkgText.NewLabel(roleLabelStr)
	var roleName, _ = pkgText.NewName("TestRoleName")
	var role = pkgUser.NewRole(roleLabel, roleName)
	var roles = []pkgUser.Role{role}

	return pkgUser.NewCompanyRole(company, roles)
}

func getPkgAuthenticForGetUser(expectId string, companyRole *pkgUser.CompanyRole) *pkgUser.Authentic {
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

type getCompanyEntryMock struct {
	getCompanyIdentifier func() (pkgText.Identifier, error)
}

func (mock getCompanyEntryMock) GetCompanyIdentifier() (pkgText.Identifier, error) {
	return mock.getCompanyIdentifier()
}

func getGetCompanyEntry(expectId string) entryCompany.CompanyGetRequest {
	return entryCompany.CompanyGetRequest{
		CompanyGet: entryCompany.CompanyGet{
			Identifier: expectId,
		},
	}
}

func TestGetCompany(t *testing.T) {
	var expectIdentifier = "CP-TESTES"
	var expectLabel = "EMPLOYEE"
	var company = getShelterCompanyForGet(expectIdentifier)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForGetUser()

	var companyGetter = getBehaviorForGetCompany(company)
	var control = controlCompany.NewGetCompanyControl(
		db,
		authorization,
		companyGetter,
	)

	var req = getGetCompanyEntry(expectIdentifier)
	var companyRole = getCompanyRoleForGetCompany(expectIdentifier, expectLabel)
	var pkgAuthentic = getPkgAuthenticForGetUser(expectIdentifier, companyRole)

	var companyGetResponse, err = controlCompany.GetCompanyExecute(control, req, pkgAuthentic)

	assert.NoError(t, err)
	assert.Equal(t, expectIdentifier, companyGetResponse.Company.Identifier)

	t.Logf("Company: %+v", companyGetResponse)
}

func TestGetCompanyErrAuth(t *testing.T) {
	var expectIdentifier = "CP-TESTES"
	var company = getShelterCompanyForGet(expectIdentifier)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForGetUser()

	var companyGetter = getBehaviorForGetCompany(company)
	var control = controlCompany.NewGetCompanyControl(
		db,
		authorization,
		companyGetter,
	)

	var req = getGetCompanyEntry(expectIdentifier)

	var _, err = controlCompany.GetCompanyExecute(control, req, nil)

	assert.Error(t, err)
}

func TestGetCompanyErrGet(t *testing.T) {
	var expectIdentifier = "CP-TESTES"
	var expectLabel = "EMPLOYEE"
	var company = getShelterCompanyForGet(expectIdentifier)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForGetUser()

	var companyGetter = getBehaviorForGetCompany(company)
	companyGetter.FakeExecute = func(entry entryCompany.CompanyGetter) (*shelterCompany.Company, error) {
		return nil, errors.New("test error")
	}
	var control = controlCompany.NewGetCompanyControl(
		db,
		authorization,
		companyGetter,
	)

	var req = getGetCompanyEntry(expectIdentifier)
	var companyRole = getCompanyRoleForGetCompany(expectIdentifier, expectLabel)
	var pkgAuthentic = getPkgAuthenticForGetUser(expectIdentifier, companyRole)

	var _, err = controlCompany.GetCompanyExecute(control, req, pkgAuthentic)

	assert.Error(t, err)
}

func TestGetCompanyErrGetNil(t *testing.T) {
	var expectIdentifier = "CP-TESTES"
	var expectLabel = "EMPLOYEE"
	var company = getShelterCompanyForGet(expectIdentifier)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForGetUser()

	var companyGetter = getBehaviorForGetCompany(company)
	companyGetter.FakeExecute = func(entry entryCompany.CompanyGetter) (*shelterCompany.Company, error) {
		return nil, nil
	}
	var control = controlCompany.NewGetCompanyControl(
		db,
		authorization,
		companyGetter,
	)

	var req = getGetCompanyEntry(expectIdentifier)
	var companyRole = getCompanyRoleForGetCompany(expectIdentifier, expectLabel)
	var pkgAuthentic = getPkgAuthenticForGetUser(expectIdentifier, companyRole)

	var _, err = controlCompany.GetCompanyExecute(control, req, pkgAuthentic)

	assert.Error(t, err)
}
