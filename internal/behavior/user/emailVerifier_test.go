package user_test

import (
	"errors"
	"github.com/motojouya/geezer_auth/internal/behavior/user"
	dbUtility "github.com/motojouya/geezer_auth/internal/db/testUtility"
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

type emailVerifierDBMock struct {
	getUserEmailOfToken func(identifier string, email string) (*dbUser.UserEmailFull, error)
	verifyEmail         func(userEmail *dbUser.UserEmail, now time.Time) (*dbUser.UserEmail, error)
	getUserAuthentic    func(identifier string, now time.Time) (*dbUser.UserAuthentic, error)
	dbUtility.SqlExecutorMock
}

func (mock emailVerifierDBMock) GetUserEmailOfToken(identifier string, email string) (*dbUser.UserEmailFull, error) {
	return mock.getUserEmailOfToken(identifier, email)
}

func (mock emailVerifierDBMock) VerifyEmail(userEmail *dbUser.UserEmail, now time.Time) (*dbUser.UserEmail, error) {
	return mock.verifyEmail(userEmail, now)
}

func (mock emailVerifierDBMock) GetUserAuthentic(identifier string, now time.Time) (*dbUser.UserAuthentic, error) {
	return mock.getUserAuthentic(identifier, now)
}

type emailVerifierEntryMock struct {
	getEmail       func() (pkgText.Email, error)
	getVerifyToken func() (shelterText.Token, error)
}

func (mock emailVerifierEntryMock) GetEmail() (pkgText.Email, error) {
	return mock.getEmail()
}

func (mock emailVerifierEntryMock) GetVerifyToken() (shelterText.Token, error) {
	return mock.getVerifyToken()
}

func getShelterUserAuthenticForEmailVerify(expectId string, expectEmail string) *shelterUser.UserAuthentic {
	var userId uint = 1
	var userIdentifier, _ = pkgText.NewIdentifier(expectId)
	var emailId, _ = pkgText.NewEmail(expectEmail)
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

func getLocalerMockForEmailVerify(t *testing.T, now time.Time) *localUtility.LocalerMock {
	var getNow = func() time.Time {
		return now
	}
	return &localUtility.LocalerMock{
		FakeGetNow: getNow,
	}
}

func getDbUserAuthentic(expectId string, expectOldEmail string, expectNewEmail string) *dbUser.UserAuthentic {
	var companyId = "CP-TESTES"
	var now = time.Now()
	var expireDate = now.Add(1 * time.Hour)
	var userCompanyRole1 = dbUser.UserCompanyRoleFull{
		UserCompanyRole: dbUser.UserCompanyRole{
			PersistKey:        1,
			UserPersistKey:    2,
			CompanyPersistKey: 3,
			RoleLabel:         "TEST_ROLE",
			RegisterDate:      now,
			ExpireDate:        &expireDate,
		},
		UserIdentifier:        expectId,
		UserExposeEmailId:     expectOldEmail,
		UserName:              "TestUserName",
		UserBotFlag:           false,
		UserRegisteredDate:    now.Add(2 * time.Hour),
		UserUpdateDate:        now.Add(3 * time.Hour),
		CompanyIdentifier:     companyId,
		CompanyName:           "TestCompanyName",
		CompanyRegisteredDate: now.Add(4 * time.Hour),
		RoleName:              "TestRoleName",
		RoleDescription:       "TestRoleDescription",
		RoleRegisteredDate:    now.Add(5 * time.Hour),
	}
	var userCompanyRoles = []dbUser.UserCompanyRoleFull{userCompanyRole1}

	return &dbUser.UserAuthentic{
		UserPersistKey:     1,
		UserIdentifier:     expectId,
		UserExposeEmailId:  expectOldEmail,
		UserName:           "TestUserName",
		UserBotFlag:        false,
		UserRegisteredDate: now,
		UserUpdateDate:     now.Add(1 * time.Hour),
		Email:              &expectNewEmail,
		UserCompanyRole:    userCompanyRoles,
	}
}

func getUserEmailFull(expectId string, expectOldEmail string, expectNewEmail string, now time.Time) *dbUser.UserEmailFull {
	var verifiyDate = now.Add(1 * time.Hour)
	var expireDate = now.Add(1 * time.Hour)
	return &dbUser.UserEmailFull{
		UserEmail: dbUser.UserEmail{
			PersistKey:     1,
			UserPersistKey: 2,
			Email:          expectNewEmail,
			VerifyToken:    "TestVerifyToken",
			RegisterDate:   now,
			VerifyDate:     &verifiyDate,
			ExpireDate:     &expireDate,
		},
		UserIdentifier:     expectId,
		UserExposeEmailId:  expectOldEmail,
		UserName:           "TestUserName",
		UserBotFlag:        false,
		UserRegisteredDate: now.Add(2 * time.Hour),
		UserUpdateDate:     now.Add(3 * time.Hour),
	}
}

func getEmailVerifyDbMock(t *testing.T, expectId string, expectOldEmail string, expectNewEmail string, exectToken string, firstNow time.Time) emailVerifierDBMock {
	var getUserEmailOfToken = func(identifier string, email string) (*dbUser.UserEmailFull, error) {
		assert.Equal(t, expectId, identifier)
		assert.Equal(t, expectNewEmail, email)
		return getUserEmailFull(expectId, expectOldEmail, expectNewEmail, firstNow), nil
	}
	var verifyEmail = func(userEmail *dbUser.UserEmail, now time.Time) (*dbUser.UserEmail, error) {
		assert.Equal(t, expectNewEmail, userEmail.Email)
		assert.WithinDuration(t, now, firstNow, time.Second, "Expected 'now' to be within 1 second of current time")
		return userEmail, nil
	}
	var update = func(args ...interface{}) (int64, error) {
		assert.Equal(t, 1, len(args), "Expected 1 argument")

		user, ok := args[0].(*dbUser.User)
		if !ok {
			t.Errorf("Expected first argument to be of type *dbUser.User, got %T", args[0])
		}

		assert.NotNil(t, user, "Expected user to be not nil")
		assert.Equal(t, expectId, user.Identifier, "Expected user identifier 'US-TESTES'")

		return 0, nil
	}
	var getUserAuthentic = func(identifier string, now time.Time) (*dbUser.UserAuthentic, error) {
		assert.Equal(t, expectId, identifier)
		assert.WithinDuration(t, now, firstNow, time.Second, "Expected 'now' to be within 1 second of current time")
		return getDbUserAuthentic(expectId, expectOldEmail, expectNewEmail), nil
	}
	return emailVerifierDBMock{
		getUserEmailOfToken: getUserEmailOfToken,
		verifyEmail:         verifyEmail,
		getUserAuthentic:    getUserAuthentic,
		SqlExecutorMock: dbUtility.SqlExecutorMock{
			FakeUpdate: update,
		},
	}
}

func getVerifyEmailEntryMock(t *testing.T, expectEmail string, expectToken string) emailVerifierEntryMock {
	var email, _ = pkgText.NewEmail(expectEmail)
	var getEmail = func() (pkgText.Email, error) {
		return email, nil
	}
	var token, _ = shelterText.NewToken(expectToken)
	var getVerifyToken = func() (shelterText.Token, error) {
		return token, nil
	}
	return emailVerifierEntryMock{
		getEmail:       getEmail,
		getVerifyToken: getVerifyToken,
	}
}

func TestEmailVerifier(t *testing.T) {
	var expectId = "US-TESTES"
	var expectOldEmail = "test01@example.com"
	var expectNewEmail = "test02@example.com"
	var expectToken = "TestVerifyToken"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForEmailVerify(expectId, expectOldEmail)

	var localerMock = getLocalerMockForEmailVerify(t, firstNow)
	var dbMock = getEmailVerifyDbMock(t, expectId, expectOldEmail, expectNewEmail, expectToken, firstNow)
	var entryMock = getVerifyEmailEntryMock(t, expectNewEmail, expectToken)

	verifier := user.NewEmailVerify(localerMock, dbMock)
	resultAuthentic, err := verifier.Execute(entryMock, userAuthentic)

	assert.NoError(t, err)
	assert.NotNil(t, resultAuthentic, "Expected resultAuthentic to be not nil")
	assert.Equal(t, expectId, string(resultAuthentic.Identifier), "Expected user identifier to match")
	assert.NotNil(t, resultAuthentic.Email, "Expected user email to be not nil")
	assert.Equal(t, expectNewEmail, string(*resultAuthentic.Email), "Expected user email to match")
}

func TestEmailVerifierErrNilAuthentic(t *testing.T) {
	var expectId = "US-TESTES"
	var expectOldEmail = "test01@example.com"
	var expectNewEmail = "test02@example.com"
	var expectToken = "TestVerifyToken"
	var firstNow = time.Now()
	// var userAuthentic = getShelterUserAuthenticForEmailVerify(expectId, expectOldEmail)

	var localerMock = getLocalerMockForEmailVerify(t, firstNow)
	var dbMock = getEmailVerifyDbMock(t, expectId, expectOldEmail, expectNewEmail, expectToken, firstNow)
	var entryMock = getVerifyEmailEntryMock(t, expectNewEmail, expectToken)

	verifier := user.NewEmailVerify(localerMock, dbMock)
	_, err := verifier.Execute(entryMock, nil)

	assert.Error(t, err)
}

func TestEmailVerifierErrEntryGetEmail(t *testing.T) {
	var expectId = "US-TESTES"
	var expectOldEmail = "test01@example.com"
	var expectNewEmail = "test02@example.com"
	var expectToken = "TestVerifyToken"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForEmailVerify(expectId, expectOldEmail)

	var localerMock = getLocalerMockForEmailVerify(t, firstNow)
	var dbMock = getEmailVerifyDbMock(t, expectId, expectOldEmail, expectNewEmail, expectToken, firstNow)
	var entryMock = getVerifyEmailEntryMock(t, expectNewEmail, expectToken)

	entryMock.getEmail = func() (pkgText.Email, error) {
		return pkgText.Email(""), errors.New("test error")
	}

	verifier := user.NewEmailVerify(localerMock, dbMock)
	_, err := verifier.Execute(entryMock, userAuthentic)

	assert.Error(t, err)
}

func TestEmailVerifierErrEntryGetToken(t *testing.T) {
	var expectId = "US-TESTES"
	var expectOldEmail = "test01@example.com"
	var expectNewEmail = "test02@example.com"
	var expectToken = "TestVerifyToken"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForEmailVerify(expectId, expectOldEmail)

	var localerMock = getLocalerMockForEmailVerify(t, firstNow)
	var dbMock = getEmailVerifyDbMock(t, expectId, expectOldEmail, expectNewEmail, expectToken, firstNow)
	var entryMock = getVerifyEmailEntryMock(t, expectNewEmail, expectToken)

	entryMock.getVerifyToken = func() (shelterText.Token, error) {
		return shelterText.Token(""), errors.New("test error")
	}

	verifier := user.NewEmailVerify(localerMock, dbMock)
	_, err := verifier.Execute(entryMock, userAuthentic)

	assert.Error(t, err)
}

func TestEmailVerifierErrDbGetEmail(t *testing.T) {
	var expectId = "US-TESTES"
	var expectOldEmail = "test01@example.com"
	var expectNewEmail = "test02@example.com"
	var expectToken = "TestVerifyToken"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForEmailVerify(expectId, expectOldEmail)

	var localerMock = getLocalerMockForEmailVerify(t, firstNow)
	var dbMock = getEmailVerifyDbMock(t, expectId, expectOldEmail, expectNewEmail, expectToken, firstNow)
	var entryMock = getVerifyEmailEntryMock(t, expectNewEmail, expectToken)

	dbMock.getUserEmailOfToken = func(identifier string, email string) (*dbUser.UserEmailFull, error) {
		return nil, errors.New("test error")
	}

	verifier := user.NewEmailVerify(localerMock, dbMock)
	_, err := verifier.Execute(entryMock, userAuthentic)

	assert.Error(t, err)
}

func TestEmailVerifierErrDbGetEmailNil(t *testing.T) {
	var expectId = "US-TESTES"
	var expectOldEmail = "test01@example.com"
	var expectNewEmail = "test02@example.com"
	var expectToken = "TestVerifyToken"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForEmailVerify(expectId, expectOldEmail)

	var localerMock = getLocalerMockForEmailVerify(t, firstNow)
	var dbMock = getEmailVerifyDbMock(t, expectId, expectOldEmail, expectNewEmail, expectToken, firstNow)
	var entryMock = getVerifyEmailEntryMock(t, expectNewEmail, expectToken)

	dbMock.getUserEmailOfToken = func(identifier string, email string) (*dbUser.UserEmailFull, error) {
		return nil, nil
	}

	verifier := user.NewEmailVerify(localerMock, dbMock)
	_, err := verifier.Execute(entryMock, userAuthentic)

	assert.Error(t, err)
}

func TestEmailVerifierErrDbGetEmailInvalid(t *testing.T) {
	var expectId = "US-TESTES"
	var expectOldEmail = "test01@example.com"
	var expectNewEmail = "test02@example.com"
	var expectToken = "TestVerifyToken"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForEmailVerify(expectId, expectOldEmail)

	var localerMock = getLocalerMockForEmailVerify(t, firstNow)
	var dbMock = getEmailVerifyDbMock(t, expectId, expectOldEmail, expectNewEmail, expectToken, firstNow)
	var entryMock = getVerifyEmailEntryMock(t, expectNewEmail, expectToken)

	dbMock.getUserEmailOfToken = func(identifier string, email string) (*dbUser.UserEmailFull, error) {
		var emailFull = getUserEmailFull(expectId, expectOldEmail, expectNewEmail, firstNow)
		emailFull.UserIdentifier = "INVALID_ID"
		return emailFull, nil
	}

	verifier := user.NewEmailVerify(localerMock, dbMock)
	_, err := verifier.Execute(entryMock, userAuthentic)

	assert.Error(t, err)
}

func TestEmailVerifierErrInvalidToken(t *testing.T) {
	var expectId = "US-TESTES"
	var expectOldEmail = "test01@example.com"
	var expectNewEmail = "test02@example.com"
	var expectToken = "TestVerifyToken"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForEmailVerify(expectId, expectOldEmail)

	var localerMock = getLocalerMockForEmailVerify(t, firstNow)
	var dbMock = getEmailVerifyDbMock(t, expectId, expectOldEmail, expectNewEmail, expectToken, firstNow)
	var entryMock = getVerifyEmailEntryMock(t, expectNewEmail, "InvalidToken")

	verifier := user.NewEmailVerify(localerMock, dbMock)
	_, err := verifier.Execute(entryMock, userAuthentic)

	assert.Error(t, err)
}

func TestEmailVerifierErrDbVerifyEmail(t *testing.T) {
	var expectId = "US-TESTES"
	var expectOldEmail = "test01@example.com"
	var expectNewEmail = "test02@example.com"
	var expectToken = "TestVerifyToken"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForEmailVerify(expectId, expectOldEmail)

	var localerMock = getLocalerMockForEmailVerify(t, firstNow)
	var dbMock = getEmailVerifyDbMock(t, expectId, expectOldEmail, expectNewEmail, expectToken, firstNow)
	var entryMock = getVerifyEmailEntryMock(t, expectNewEmail, "InvalidToken")

	dbMock.verifyEmail = func(userEmail *dbUser.UserEmail, now time.Time) (*dbUser.UserEmail, error) {
		return nil, errors.New("test error")
	}

	verifier := user.NewEmailVerify(localerMock, dbMock)
	_, err := verifier.Execute(entryMock, userAuthentic)

	assert.Error(t, err)
}

func TestEmailVerifierErrDbUpdate(t *testing.T) {
	var expectId = "US-TESTES"
	var expectOldEmail = "test01@example.com"
	var expectNewEmail = "test02@example.com"
	var expectToken = "TestVerifyToken"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForEmailVerify(expectId, expectOldEmail)

	var localerMock = getLocalerMockForEmailVerify(t, firstNow)
	var dbMock = getEmailVerifyDbMock(t, expectId, expectOldEmail, expectNewEmail, expectToken, firstNow)
	var entryMock = getVerifyEmailEntryMock(t, expectNewEmail, "InvalidToken")

	dbMock.SqlExecutorMock.FakeUpdate = func(args ...interface{}) (int64, error) {
		return 0, errors.New("test error")
	}

	verifier := user.NewEmailVerify(localerMock, dbMock)
	_, err := verifier.Execute(entryMock, userAuthentic)

	assert.Error(t, err)
}

func TestEmailVerifierErrDbGetAuthentic(t *testing.T) {
	var expectId = "US-TESTES"
	var expectOldEmail = "test01@example.com"
	var expectNewEmail = "test02@example.com"
	var expectToken = "TestVerifyToken"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForEmailVerify(expectId, expectOldEmail)

	var localerMock = getLocalerMockForEmailVerify(t, firstNow)
	var dbMock = getEmailVerifyDbMock(t, expectId, expectOldEmail, expectNewEmail, expectToken, firstNow)
	var entryMock = getVerifyEmailEntryMock(t, expectNewEmail, "InvalidToken")

	dbMock.getUserAuthentic = func(identifier string, now time.Time) (*dbUser.UserAuthentic, error) {
		return nil, errors.New("test error")
	}

	verifier := user.NewEmailVerify(localerMock, dbMock)
	_, err := verifier.Execute(entryMock, userAuthentic)

	assert.Error(t, err)
}

func TestEmailVerifierErrDbGetAuthenticNil(t *testing.T) {
	var expectId = "US-TESTES"
	var expectOldEmail = "test01@example.com"
	var expectNewEmail = "test02@example.com"
	var expectToken = "TestVerifyToken"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForEmailVerify(expectId, expectOldEmail)

	var localerMock = getLocalerMockForEmailVerify(t, firstNow)
	var dbMock = getEmailVerifyDbMock(t, expectId, expectOldEmail, expectNewEmail, expectToken, firstNow)
	var entryMock = getVerifyEmailEntryMock(t, expectNewEmail, "InvalidToken")

	dbMock.getUserAuthentic = func(identifier string, now time.Time) (*dbUser.UserAuthentic, error) {
		return nil, nil
	}

	verifier := user.NewEmailVerify(localerMock, dbMock)
	_, err := verifier.Execute(entryMock, userAuthentic)

	assert.Error(t, err)
}
