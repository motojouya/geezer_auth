package user_test

import (
	"errors"
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

type emailVerifierDBMock struct {
	getUserEmailOfToken func(identifier string, email string) (*dbUser.UserEmailFull, error)
	verifyEmail         func(userEmail *dbUser.UserEmail, now time.Time) (*dbUser.UserEmail, error)
	getUserAuthentic    func(identifier string) (*dbUser.UserAuthentic, error)
	dbUtility.SqlExecutorMock
}

func (mock emailVerifierDBMock) GetUserEmailOfToken(identifier string, email string) (*dbUser.UserEmailFull, error) {
	return mock.getUserEmailOfToken(identifier, email)
}

func (mock emailVerifierDBMock) VerifyEmail(userEmail *dbUser.UserEmail, now time.Time) (*dbUser.UserEmail, error) {
	return mock.verifyEmail(userEmail, now)
}

func (mock emailVerifierDBMock) GetUserAuthentic(identifier string) (*dbUser.UserAuthentic, error) {
	return mock.getUserAuthentic(identifier)
}

type emailVerifierEntryMock struct {
	getEmail func() (pkgText.Email, error)
	getVerifyToken func() (shelterText.Token, error)
}

func (mock emailVerifierEntryMock) GetEmail() (pkgText.Email, error) {
	return mock.getEmail()
}

func (mock emailVerifierEntryMock) GetVerifyToken() (shelterText.Token, error) {
	return mock.getVerifyToken()
}

func getShelterUserAuthenticForEmailVerify(expectEmail string) *shelterUser.UserAuthentic {

	var userId uint = 1
	var userIdentifier, _ = pkgText.NewIdentifier("TestIdentifier")
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
		FakeGetNow:       getNow,
	}
}

func getEmailVerifyDbMock(t *testing.T, expectId string, expectEmail string, exectToken string, firstNow time.Time) emailVerifierDBMock {
	var getUserEmailOfToken = func(identifier string, email string) (*dbUser.UserEmailFull, error) {
		assert.Equal(t, expectId, identifier)
		assert.Equal(t, expectEmail, email)
		return &dbUser.UserEmailFull{
			Email: expectEmail,
			VerifyToken: exectToken,
		}, nil
	}
	var verifyEmail = func(userEmail *dbUser.UserEmail, now time.Time) (*dbUser.UserEmail, error) {
		assert.Equal(t, expectEmail, userEmail.Email)
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

		return nil
	}
	var getUserAuthentic = func(identifier string) (*dbUser.UserAuthentic, error) {
		assert.Equal(t, expectId, identifier)
		return []dbUser.UserEmailFull{}, nil
	}
	return emailVerifierDBMock{
		getUserEmailOfToken: getUserEmailOfToken,
		verifyEmail:         verifyEmail,
		getUserAuthentic:    getUserAuthentic,
		update:              update,
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

func TestEmailSetter(t *testing.T) {
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

func TestEmailSetterErrGetEmail(t *testing.T) {
	var expectEmail = "test@example.com"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForEmail(expectEmail)

	var localerMock = getLocalerMockForEmail(t, firstNow)
	var dbMock = getEmailSetDbMock(t, expectEmail, firstNow)
	var entryMock = getGetEmailEntryMock(t, expectEmail, firstNow)

	entryMock.getEmail = func() (pkgText.Email, error) {
		return pkgText.Email(""), errors.New("failed to get email")
	}

	setter := user.NewEmailSet(localerMock, dbMock)
	err := setter.Execute(entryMock, userAuthentic)

	assert.Error(t, err)
}

func TestEmailSetterErrGetUserEmail(t *testing.T) {
	var expectEmail = "test@example.com"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForEmail(expectEmail)

	var localerMock = getLocalerMockForEmail(t, firstNow)
	var dbMock = getEmailSetDbMock(t, expectEmail, firstNow)
	var entryMock = getGetEmailEntryMock(t, expectEmail, firstNow)

	dbMock.getUserEmail = func(email string) ([]dbUser.UserEmailFull, error) {
		return nil, errors.New("failed to get user email")
	}

	setter := user.NewEmailSet(localerMock, dbMock)
	err := setter.Execute(entryMock, userAuthentic)

	assert.Error(t, err)
}

func TestEmailSetterErrGetUserEmailMany(t *testing.T) {
	var expectEmail = "test@example.com"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForEmail(expectEmail)

	var localerMock = getLocalerMockForEmail(t, firstNow)
	var dbMock = getEmailSetDbMock(t, expectEmail, firstNow)
	var entryMock = getGetEmailEntryMock(t, expectEmail, firstNow)

	var verifiyDate = firstNow.Add(1 * time.Hour)
	var expireDate = firstNow.Add(1 * time.Hour)
	var userEmail = dbUser.UserEmailFull{
		UserEmail: dbUser.UserEmail{
			PersistKey:     1,
			UserPersistKey: 2,
			Email:          "test01@example.com",
			VerifyToken:    "TestVerifyToken",
			RegisterDate:   firstNow,
			VerifyDate:     &verifiyDate,
			ExpireDate:     &expireDate,
		},
		UserIdentifier:     "invalid-identifier",
		UserExposeEmailId:  "test02@example.com",
		UserName:           "TestUserName",
		UserBotFlag:        false,
		UserRegisteredDate: firstNow.Add(2 * time.Hour),
		UserUpdateDate:     firstNow.Add(3 * time.Hour),
	}

	dbMock.getUserEmail = func(email string) ([]dbUser.UserEmailFull, error) {
		return []dbUser.UserEmailFull{userEmail}, nil
	}

	setter := user.NewEmailSet(localerMock, dbMock)
	err := setter.Execute(entryMock, userAuthentic)

	assert.Error(t, err)
}

func TestEmailSetterErrGenerateUUID(t *testing.T) {
	var expectEmail = "test@example.com"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForEmail(expectEmail)

	var localerMock = getLocalerMockForEmail(t, firstNow)
	var dbMock = getEmailSetDbMock(t, expectEmail, firstNow)
	var entryMock = getGetEmailEntryMock(t, expectEmail, firstNow)

	localerMock.FakeGenerateUUID = func() (uuid.UUID, error) {
		return uuid.UUID{}, errors.New("failed to generate UUID")
	}

	setter := user.NewEmailSet(localerMock, dbMock)
	err := setter.Execute(entryMock, userAuthentic)

	assert.Error(t, err)
}

func TestEmailSetterErrAddEmail(t *testing.T) {
	var expectEmail = "test@example.com"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForEmail(expectEmail)

	var localerMock = getLocalerMockForEmail(t, firstNow)
	var dbMock = getEmailSetDbMock(t, expectEmail, firstNow)
	var entryMock = getGetEmailEntryMock(t, expectEmail, firstNow)

	dbMock.addEmail = func(userEmail *dbUser.UserEmail, now time.Time) (*dbUser.UserEmail, error) {
		return nil, errors.New("failed to add email")
	}

	setter := user.NewEmailSet(localerMock, dbMock)
	err := setter.Execute(entryMock, userAuthentic)

	assert.Error(t, err)
}
