package company_test

import (
	//"errors"
	dbUtility "github.com/motojouya/geezer_auth/internal/db/testUtility"
	"github.com/google/uuid"
	"github.com/motojouya/geezer_auth/internal/behavior/user"
	dbUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	localUtility "github.com/motojouya/geezer_auth/internal/local/testUtility"
	shelterCompany "github.com/motojouya/geezer_auth/internal/shelter/company"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
	shelterText "github.com/motojouya/geezer_auth/internal/shelter/text"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type roleAssignerDBMock struct {
	getUserAuthentic func(identifier string, now time.Time) (*dbUser.UserAuthentic, error)
	dbUtility.SqlExecutorMock
}

func (mock roleAssignerDBMock) GetUserAuthentic(identifier string, now time.Time) (*dbUser.UserAuthentic, error) {
	return mock.getUserAuthentic(identifier, now)
}

func getShelterUserAuthenticForRoleAssign(expectId string) *shelterUser.UserAuthentic {

	var userId uint = 1
	var userIdentifier, _ = pkgText.NewIdentifier(expectId)
	var emailId, _ = pkgText.NewEmail("test@exaple.com")
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

func getLocalerMockForRoleAssign(t *testing.T, now time.Time) *localUtility.LocalerMock {
	var getNow = func() time.Time {
		return now
	}
	var generateUUID = func() (uuid.UUID, error) {
		return uuid.NewUUID()
	}
	return &localUtility.LocalerMock{
		FakeGenerateUUID: generateUUID,
		FakeGetNow:       getNow,
	}
}

func getRoleAssignDbMock(t *testing.T, expectId string, firstNow time.Time) roleAssignerDBMock {
	var getUserAuthentic = func(identifier string, now time.Time) (*dbUser.UserAuthentic, error) {
		assert.Equal(t, expectId, identifier)
		return []dbUser.UserEmailFull{}, nil
	}
	var insert = func(args ...interface{}) error {
		assert.Equal(t, 1, len(args), "Expected 1 argument")

		userCompanyRole, ok := args[0].(*dbUser.UnsavedUserCompanyRole)
		if !ok {
			t.Errorf("Expected first argument to be of type *dbUser.User, got %T", args[0])
		}

		assert.NotNil(t, userCompanyRole)
		assert.Equal(t, expectId, userCompanyRole.User.Identifier)

		return nil
	}
	return roleAssignerDBMock{
		getUserEmail: getUserEmail,
		addEmail:     addEmail,
		SqlExecutorMock: dbUtility.SqlExecutorMock{
			FakeInsert: insert,
		},
	}
}

func TestRoleAssigner(t *testing.T) {
	var expectEmail = "test@example.com"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForEmail(expectEmail)

	var localerMock = getLocalerMockForEmail(t, firstNow)
	var dbMock = getEmailSetDbMock(t, expectEmail, firstNow)
	var entryMock = getGetEmailEntryMock(t, expectEmail, firstNow)

	setter := user.NewEmailSet(localerMock, dbMock)
	err := setter.Execute(entryMock, userAuthentic)

	assert.NoError(t, err)
}
