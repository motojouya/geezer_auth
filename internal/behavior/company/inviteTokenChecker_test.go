package company_test

import (
	//"errors"
	"github.com/motojouya/geezer_auth/internal/behavior/company"
	dbUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	localUtility "github.com/motojouya/geezer_auth/internal/local/testUtility"
	shelterText "github.com/motojouya/geezer_auth/internal/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type inviteTokenCheckerDBMock struct {
	getCompanyInvite func(identifier string, token string) (*dbCompany.CompanyInviteFull, error)
}

func (mock inviteTokenCheckerDBMock) GetUserRefreshToken(token string, now time.Time) (*dbCompany.CompanyInviteFull, error) {
	return mock.getCompanyInvite(identifier, token)
}

func getLocalerMockForInviteTokenCheck(t *testing.T, now time.Time) *localUtility.LocalerMock {
	var getNow = func() time.Time {
		return now
	}
	return &localUtility.LocalerMock{
		FakeGetNow:       getNow,
	}
}

func getCompanyInviteFull(expectId string, expectLabel string, expectToken string) dbCompany.CompanyInviteFull {
	var registerDate = time.Now()
	var expireDate = registerDate.Add(50 * time.Hour)

	return company.CompanyInviteFull{
		CompanyInvite: company.CompanyInvite{
			PersistKey:        1,
			CompanyPersistKey: 2,
			Token:             expectToken,
			RoleLabel:         expectLabel,
			RegisterDate:      registerDate,
			ExpireDate:        expireDate,
		},
		CompanyIdentifier:     expectId,
		CompanyName:           "Test Company",
		CompanyRegisteredDate: companyValue.RegisteredDate,
		RoleName:              "TestRole",
		RoleDescription:       "Role for testing",
		RoleRegisteredDate:    role.RegisteredDate,
	}
}

func getInviteTokenCheckerDbMock(t *testing.T, expectId string, expectLabel string, expectToken string, firstNow time.Time) inviteTokenCheckerDBMock {
	var companyInviteFull = getCompanyInviteFull(expectId, expectLabel, expectToken)
	var getCompanyInvite = func(identifier string, token string) (*dbCompany.CompanyInviteFull, error) {
		assert.Equal(t, identifier, expectId)
		assert.Equal(t, token, expectToken)
		return &companyInviteFull, nil
	}
	return inviteTokenCheckerDBMock{
		getCompanyInvite: getCompanyInvite,
	}
}

type userInviteTokenGetterMock struct {
	getToken func() (shelterText.Token, error)
}

func (getter userInviteTokenGetterMock) GetToken() (shelterText.Token, error) {
	return getter.getToken()
}

func getUserInviteTokenGetterMock(expectToken string) userInviteTokenGetterMock {
	var getToken = func() (shelterText.Token, error) {
		return shelterText.NewToken(expectToken)
	}
	return userInviteTokenGetterMock{
		getToken: getToken,
	}
}

func getShelterCompanyForInviteCheck(expectId string) shelterCompany.Company {
	var companyIdentifier, _ = pkgText.NewIdentifier(expectId)
	var companyId uint = 1
	var companyName, _ = pkgText.NewName("TestCompany")
	var companyRegisteredDate = time.Now()

	return shelterCompany.NewCompany(companyId, companyIdentifier, companyName, companyRegisteredDate)
}

func TestRefreshTokenChecker(t *testing.T) {
	var firstNow = time.Now()
	var expectToken = "refresh_token01"
	var expectId = "CP-TESTES"
	var expectLabel = "ROLE_LABEL"
	var company = getShelterCompanyForInviteCheck(expectId)

	var localerMock = getLocalerMockForInviteTokenCheck(t, firstNow)
	var dbMock = getInviteTokenCheckerDbMock(t, expectId, expectLabel, expectToken, firstNow)
	var entryMock = getUserInviteTokenGetterMock(expectToken)

	checker := company.NewInviteTokenCheck(localerMock, dbMock)
	role, err := checker.Execute(entryMock, company)

	assert.NoError(t, err)
	assert.NotNil(t, role)
	assert.Equal(t, expectLabel, string(role.Label))
}
