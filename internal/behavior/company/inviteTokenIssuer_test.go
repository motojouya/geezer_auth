package company_test

import (
	//"errors"
	dbUtility "github.com/motojouya/geezer_auth/internal/db/testUtility"
	"github.com/google/uuid"
	"github.com/motojouya/geezer_auth/internal/behavior/company"
	dbCompany "github.com/motojouya/geezer_auth/internal/db/transfer/company"
	localUtility "github.com/motojouya/geezer_auth/internal/local/testUtility"
	shelterCompany "github.com/motojouya/geezer_auth/internal/shelter/company"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
	shelterText "github.com/motojouya/geezer_auth/internal/shelter/text"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getLocalerMockForInviteToken(expectUUID uuid.UUID, now time.Time) *localUtility.LocalerMock {
	var getNow = func() time.Time {
		return now
	}
	var generateUUID = func() (uuid.UUID, error) {
		return expectUUID, nil
	}
	return &localUtility.LocalerMock{
		FakeGetNow:       getNow,
		FakeGenerateUUID: generateUUID,
	}
}

func getInviteTokenIssuerDBMock(t *testing.T, expectId string, expectRole string) *dbUtility.SqlExecutorMock {
	var insert = func(args ...interface{}) error {
		assert.Equal(t, 1, len(args), "Expected 1 argument")

		companyInvite, ok := args[0].(*dbCompany.CompanyInvite)
		if !ok {
			t.Errorf("Expected first argument to be of type *dbCompany.Company, got %T", args[0])
		}

		assert.NotNil(t, companyInvite)
		assert.Equal(t, expectRole, companyInvite.RoleLabel)

		return nil
	}
	return &dbUtility.SqlExecutorMock{
		FakeInsert: insert,
	}
}

func getShelterCompanyForInviteIssue(expectId string) shelterCompany.Company {
	var companyIdentifier, _ = pkgText.NewIdentifier(expectId)
	var companyId uint = 1
	var companyName, _ = pkgText.NewName("TestCompany")
	var companyRegisteredDate = time.Now()

	return shelterCompany.NewCompany(companyId, companyIdentifier, companyName, companyRegisteredDate)
}

func getShelterRoleForInviteIssue(expectLabel string) shelterRole.Role {
	var label, _ = pkgText.NewLabel(expectLabel)
	var roleName, _ = pkgText.NewName("TestRole")
	var description, _ = shelterText.NewText("Role for testing")
	var roleRegisteredDate = time.Now()

	return shelterRole.NewRole(roleName, label, description, roleRegisteredDate)
}

func TestRefreshTokenIssuer(t *testing.T) {
	var expectId = "CP-TESTES"
	var expectRole = "ROLE_TEST"
	var firstNow = time.Now()
	var expectUUID, _ = uuid.NewUUID()

	var shelterCompany = getShelterCompanyForInviteIssue(expectId)
	var role = getShelterRoleForInviteIssue(expectRole)

	var localerMock = getLocalerMockForInviteToken(expectUUID, firstNow)
	var dbMock = getInviteTokenIssuerDBMock(t, expectId, expectRole)

	issuer := company.NewInviteTokenIssue(localerMock, dbMock)
	inviteToken, err := issuer.Execute(shelterCompany, role)

	assert.NoError(t, err)
	assert.Equal(t, expectUUID.String(), string(inviteToken), "Expected refresh token to match generated UUID")
}
