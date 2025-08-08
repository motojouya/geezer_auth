package user_test

import (
	localUtility "github.com/motojouya/geezer_auth/internal/local/testUtility"
	"github.com/motojouya/geezer_auth/internal/behavior/user"
	dbUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
	shelterCompany "github.com/motojouya/geezer_auth/internal/shelter/company"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	shelterText "github.com/motojouya/geezer_auth/internal/shelter/text"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"errors"
)

type emailSetterDBMock struct {
	getUserEmail func(email string) ([]dbUser.UserEmailFull, error)
	addEmail func(userEmail *dbUser.UserEmail, now time.Time) (*dbUser.UserEmail, error)
}

func (mock emailSetterDBMock) GetUserEmail(email string) ([]dbUser.UserEmailFull, error) {
	return mock.getUserEmail(email)
}

func (mock emailSetterDBMock) AddEmail(userEmail *dbUser.UserEmail, now time.Time) (*dbUser.UserEmail, error) {
	return mock.addEmail(userEmail, now)
}

type emailGetEntryMock struct {
	getEmail func() (pkgText.Email, error)
}

func (mock emailGetEntryMock) GetEmail() (pkgText.Email, error) {
	return mock.getEmail()
}

func getShelterUserAuthenticForEmail(expectEmail string) *shelterUser.UserAuthentic {

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

func getLocalerMockForEmail(t *testing.T, now time.Time) *localUtility.LocalerMock {
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

func getEmailSetDbMock(t *testing.T, expectEmail string, firstNow time.Time) emailSetterDBMock {
	var getUserEmail = func(email string) ([]dbUser.UserEmailFull, error) {
		assert.Equal(t, expectEmail, email)
		return []dbUser.UserEmailFull{}, nil
	}
	var addEmail = func(userEmail *dbUser.UserEmail, now time.Time) (*dbUser.UserEmail, error) {
		assert.Equal(t, expectEmail, userEmail.Email)
		assert.WithinDuration(t, now, firstNow, time.Second, "Expected 'now' to be within 1 second of current time")
		return userEmail, nil
	}
	return emailSetterDBMock{
		getUserEmail: getUserEmail,
		addEmail:     addEmail,
	}
}

func getGetEmailEntryMock(t *testing.T, expectEmail string, firstNow time.Time) emailGetEntryMock {
	var email, _ = pkgText.NewEmail(expectEmail)
	var getEmail = func() (pkgText.Email, error) {
		return email, nil
	}
	return emailGetEntryMock{
		getEmail: getEmail,
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
