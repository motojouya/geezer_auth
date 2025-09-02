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

func getBehaviorForGetUser(users []shelterUser.UserAuthentic) *companyTestUtility.AllUserGetterMock {
	var allUserGetter = &companyTestUtility.AllUserGetterMock{
		FakeExecute: func(entry entryCompany.CompanyGetter) ([]shelterUser.UserAuthentic, error) {
			return users, nil
		},
	}

	return allUserGetter
}

func getShelterCompanyForGetUser(expectId string) shelterCompany.Company {
	var companyIdentifier, _ = pkgText.NewIdentifier(expectId)
	var companyId uint = 1
	var companyName, _ = pkgText.NewName("TestCompany")
	var companyRegisteredDate = time.Now()

	return shelterCompany.NewCompany(companyId, companyIdentifier, companyName, companyRegisteredDate)
}

func getShelterUserAuthenticForGetAllUser(expectId string) *shelterUser.UserAuthentic {
	var userId uint = 1
	var userIdentifier, _ = pkgText.NewIdentifier("US-TESTES")
	var emailId, _ = pkgText.NewEmail("test@example.com")
	var userName, _ = pkgText.NewName("Test User")
	var botFlag = false
	var userRegisteredDate = time.Now()
	var updateDate = time.Now()
	var userValue = shelterUser.NewUser(userId, userIdentifier, userName, emailId, botFlag, userRegisteredDate, updateDate)

	var companyIdentifier, _ = pkgText.NewIdentifier(expectId)
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

func getAuthorizationForGetAllUser() *shelterAuth.Authorization {
	return shelterAuth.NewAuthorization([]shelterRole.RolePermission{
		shelterRole.AnonymousPermission,
		shelterRole.RoleLessPermission,
		shelterRole.NewRolePermission("EMPLOYEE", true, true, false, false, 5),
		shelterRole.NewRolePermission("MANAGER", true, true, true, true, 9),
	})
}

func getCompanyRole(companyIdentifierStr string, roleLabelStr string) *pkgUser.CompanyRole {
	var companyIdentifier, _ = pkgText.NewIdentifier(companyIdentifierStr)
	var companyName, _ = pkgText.NewName("TestCompany")
	var company = pkgUser.NewCompany(companyIdentifier, companyName)

	var roleLabel, _ = pkgText.NewLabel(roleLabelStr)
	var roleName, _ = pkgText.NewName("TestRoleName")
	var role = pkgUser.NewRole(roleLabel, roleName)
	var roles = []pkgUser.Role{role}

	return pkgUser.NewCompanyRole(company, roles)
}

func getPkgAuthenticForGetAllUser(companyRole *pkgUser.CompanyRole) *pkgUser.Authentic {

	var userIdentifier, _ = pkgText.NewIdentifier("US-TESTES")
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

func getGetCompanyEntryForGetUser(expectId string) entryCompany.CompanyGetRequest {
	return entryCompany.CompanyGetRequest{
		CompanyGet: entryCompany.CompanyGet{
			Identifier: expectId,
		},
	}
}

func TestGetUser(t *testing.T) {
	var expectIdentifier = "CP-TESTES"
	var expectLabel = "EMPLOYEE"
	var userAuthentic = getShelterUserAuthenticForGetAllUser(expectIdentifier)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForGetAllUser()

	var userGetter = getBehaviorForGetUser([]shelterUser.UserAuthentic{*userAuthentic})
	var control = controlCompany.NewGetUserControl(
		db,
		authorization,
		userGetter,
	)

	var req = getGetCompanyEntryForGetUser(expectIdentifier)
	var companyRole = getCompanyRole(expectIdentifier, expectLabel)
	var pkgAuthentic = getPkgAuthenticForGetAllUser(companyRole)

	var companyUserResponse, err = controlCompany.GetUserExecute(control, req, pkgAuthentic)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(companyUserResponse.Users))
	assert.Equal(t, expectIdentifier, companyUserResponse.Users[0].CompanyRole.Company.Identifier)

	t.Logf("Company: %+v", companyUserResponse)
}

func TestGetUserErrAuth(t *testing.T) {
	var expectIdentifier = "CP-TESTES"
	var userAuthentic = getShelterUserAuthenticForGetAllUser(expectIdentifier)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForGetAllUser()

	var userGetter = getBehaviorForGetUser([]shelterUser.UserAuthentic{*userAuthentic})
	var control = controlCompany.NewGetUserControl(
		db,
		authorization,
		userGetter,
	)

	var req = getGetCompanyEntryForGetUser(expectIdentifier)
	var pkgAuthentic = getPkgAuthenticForGetAllUser(nil)

	var _, err = controlCompany.GetUserExecute(control, req, pkgAuthentic)

	assert.Error(t, err)
}

func TestGetUserErrGet(t *testing.T) {
	var expectIdentifier = "CP-TESTES"
	var expectLabel = "EMPLOYEE"
	var userAuthentic = getShelterUserAuthenticForGetAllUser(expectIdentifier)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForGetAllUser()

	var userGetter = getBehaviorForGetUser([]shelterUser.UserAuthentic{*userAuthentic})
	userGetter.FakeExecute = func(entry entryCompany.CompanyGetter) ([]shelterUser.UserAuthentic, error) {
		return nil, errors.New("db error")
	}
	var control = controlCompany.NewGetUserControl(
		db,
		authorization,
		userGetter,
	)

	var req = getGetCompanyEntryForGetUser(expectIdentifier)
	var companyRole = getCompanyRole(expectIdentifier, expectLabel)
	var pkgAuthentic = getPkgAuthenticForGetAllUser(companyRole)

	var _, err = controlCompany.GetUserExecute(control, req, pkgAuthentic)

	assert.Error(t, err)
}
