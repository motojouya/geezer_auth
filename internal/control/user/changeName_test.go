package user_test

import (
	"errors"
	"github.com/google/uuid"
	pkgUser "github.com/motojouya/geezer_auth/pkg/shelter/user"
	shelterAuth "github.com/motojouya/geezer_auth/internal/shelter/authorization"
	userTestUtility "github.com/motojouya/geezer_auth/internal/behavior/user/testUtility"
	controlUser "github.com/motojouya/geezer_auth/internal/control/user"
	dbTestUtility "github.com/motojouya/geezer_auth/internal/db/testUtility"
	entryUser "github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	shelterCompany "github.com/motojouya/geezer_auth/internal/shelter/company"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
	shelterText "github.com/motojouya/geezer_auth/internal/shelter/text"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getBehaviorForChangeName(t *testing.T, userAuthentic *shelterUser.UserAuthentic, name pkgText.Name, accessToken pkgText.JwtToken) (*userTestUtility.UserGetterMock, *userTestUtility.NameChangerMock, *userTestUtility.AccessTokenIssuerMock) {
	var userGetter = &userTestUtility.UserGetterMock{
		FakeExecute: func(identifier pkgText.Identifier) (*shelterUser.UserAuthentic, error) {
			assert.Equal(t, identifier, userAuthentic.Identifier)
			return userAuthentic, nil
		},
	}

	var nameChanger = &userTestUtility.NameChangerMock{
		FakeExecute: func(entry entryUser.UserApplyer, user *shelterUser.UserAuthentic) (*shelterUser.UserAuthentic, error) {
			user.Name = name
			return user, nil
		},
	}

	var accessTokenIssuer = &userTestUtility.AccessTokenIssuerMock{
		FakeExecute: func(user *shelterUser.UserAuthentic) (pkgText.JwtToken, error) {
			return accessToken, nil
		},
	}

	return userGetter, nameChanger, accessTokenIssuer
}

func getShelterUserAuthenticForChangeName(expectId string, expectName string) *shelterUser.UserAuthentic {
	var userId uint = 1
	var userIdentifier, _ = pkgText.NewIdentifier(expectId)
	var emailId, _ = pkgText.NewEmail("test@example.com")
	var userName, _ = pkgText.NewName(expectName)
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

func getChangeNameEntry(expectedName string) entryUser.UserChangeNameRequest {
	return entryUser.UserChangeNameRequest{
		UserChangeName: entryUser.UserChangeName{
			Name:     expectedName,
		},
	}
}

func getAuthorizationForChangeName() *shelterAuth.Authorization {
	return shelterAuth.NewAuthorization([]shelterRole.RolePermission{
		shelterRole.AnonymousPermission,
		shelterRole.RoleLessPermission,
		shelterRole.NewRolePermission("EMPLOYEE", true, true, false, false, 5),
		shelterRole.NewRolePermission("MANAGER", true, true, true, true, 9),
	})
}

func getPkgAuthenticForChangeName(expectId string, expectName string) *pkgUser.Authentic {
	var userIdentifier, _ = pkgText.NewIdentifier(expectId)
	var emailId, _ = pkgText.NewEmail("test@example.com")
	var email, _ = pkgText.NewEmail("test@example.com")
	var userName, _ = pkgText.NewName(expectName)
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

func TestChangeName(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectOldName = "Test Old User"
	var expectNewName = "Test New User"
	var name, _ = pkgText.NewName(expectNewName)
	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)
	var userAuthentic = getShelterUserAuthenticForChangeName(expectIdentifier, expectOldName)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForChangeName()

	var userGetter, nameChanger, accessTokenIssuer = getBehaviorForChangeName(t, userAuthentic, name, accessToken)
	var control = controlUser.NewChangeNameControl(
		db,
		authorization,
		userGetter,
		nameChanger,
		accessTokenIssuer,
	)

	var entry = getChangeNameEntry(expectNewName)
	var pkgAuthentic = getPkgAuthenticForChangeName(expectIdentifier, expectOldName)

	var userUpdateResponse, err = controlUser.ChangeNameExecute(control, entry, pkgAuthentic)

	assert.NoError(t, err)
	assert.Equal(t, expectIdentifier, userUpdateResponse.User.Identifier)
	assert.Equal(t, expectNewName, userUpdateResponse.User.Name)
	assert.Equal(t, expectToken, userUpdateResponse.AccessToken)

	assert.Equal(t, 1, transactionCalledCount.BeginCalled)
	assert.Equal(t, 1, transactionCalledCount.CommitCalled)
	assert.Equal(t, 0, transactionCalledCount.RollbackCalled)
	assert.Equal(t, 0, transactionCalledCount.CloseCalled)

	t.Logf("User Identifier: %+v", userUpdateResponse)
}

func TestChangeNameErrAuth(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectOldName = "Test Old User"
	var expectNewName = "Test New User"
	var name, _ = pkgText.NewName(expectNewName)
	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)
	var userAuthentic = getShelterUserAuthenticForChangeName(expectIdentifier, expectOldName)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForChangeName()

	var userGetter, nameChanger, accessTokenIssuer = getBehaviorForChangeName(t, userAuthentic, name, accessToken)
	var control = controlUser.NewChangeNameControl(
		db,
		authorization,
		userGetter,
		nameChanger,
		accessTokenIssuer,
	)

	var entry = getChangeNameEntry(expectNewName)

	var _, err = controlUser.ChangeNameExecute(control, entry, nil)

	assert.Error(t, err)
}

func TestChangeNameErrGet(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectOldName = "Test Old User"
	var expectNewName = "Test New User"
	var name, _ = pkgText.NewName(expectNewName)
	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)
	var userAuthentic = getShelterUserAuthenticForChangeName(expectIdentifier, expectOldName)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForChangeName()

	var userGetter, nameChanger, accessTokenIssuer = getBehaviorForChangeName(t, userAuthentic, name, accessToken)
	userGetter.FakeExecute = func(identifier pkgText.Identifier) (*shelterUser.UserAuthentic, error) {
		return nil, errors.New("get user error")
	}
	var control = controlUser.NewChangeNameControl(
		db,
		authorization,
		userGetter,
		nameChanger,
		accessTokenIssuer,
	)

	var entry = getChangeNameEntry(expectNewName)
	var pkgAuthentic = getPkgAuthenticForChangeName(expectIdentifier, expectOldName)

	var _, err = controlUser.ChangeNameExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)
}

func TestChangeNameErrChange(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectOldName = "Test Old User"
	var expectNewName = "Test New User"
	var name, _ = pkgText.NewName(expectNewName)
	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)
	var userAuthentic = getShelterUserAuthenticForChangeName(expectIdentifier, expectOldName)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForChangeName()

	var userGetter, nameChanger, accessTokenIssuer = getBehaviorForChangeName(t, userAuthentic, name, accessToken)
	nameChanger.FakeExecute = func(entry entryUser.UserApplyer, user *shelterUser.UserAuthentic) (*shelterUser.UserAuthentic, error) {
		return nil, errors.New("change name error")
	}
	var control = controlUser.NewChangeNameControl(
		db,
		authorization,
		userGetter,
		nameChanger,
		accessTokenIssuer,
	)

	var entry = getChangeNameEntry(expectNewName)
	var pkgAuthentic = getPkgAuthenticForChangeName(expectIdentifier, expectOldName)

	var _, err = controlUser.ChangeNameExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)
}

func TestChangeNameErrIssue(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectOldName = "Test Old User"
	var expectNewName = "Test New User"
	var name, _ = pkgText.NewName(expectNewName)
	var expectToken = "test-access-token"
	var accessToken = pkgText.JwtToken(expectToken)
	var userAuthentic = getShelterUserAuthenticForChangeName(expectIdentifier, expectOldName)

	var transactionCalledCount = &dbTestUtility.TransactionCalledCount{}
	var db = dbTestUtility.GetTransactionalDatabaseMock(transactionCalledCount)
	var authorization = getAuthorizationForChangeName()

	var userGetter, nameChanger, accessTokenIssuer = getBehaviorForChangeName(t, userAuthentic, name, accessToken)
	accessTokenIssuer.FakeExecute = func(user *shelterUser.UserAuthentic) (pkgText.JwtToken, error) {
		return pkgText.JwtToken(""), errors.New("issue access token error")
	}
	var control = controlUser.NewChangeNameControl(
		db,
		authorization,
		userGetter,
		nameChanger,
		accessTokenIssuer,
	)

	var entry = getChangeNameEntry(expectNewName)
	var pkgAuthentic = getPkgAuthenticForChangeName(expectIdentifier, expectOldName)

	var _, err = controlUser.ChangeNameExecute(control, entry, pkgAuthentic)

	assert.Error(t, err)
}
