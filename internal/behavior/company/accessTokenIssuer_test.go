package user_test

import (
	"errors"
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/motojouya/geezer_auth/internal/behavior/user"
	dbUtility "github.com/motojouya/geezer_auth/internal/db/testUtility"
	dbUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	localUtility "github.com/motojouya/geezer_auth/internal/local/testUtility"
	shelterCompany "github.com/motojouya/geezer_auth/internal/shelter/company"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
	shelterText "github.com/motojouya/geezer_auth/internal/shelter/text"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
	jwtUtility "github.com/motojouya/geezer_auth/pkg/shelter/jwt/testUtility"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	pkgUser "github.com/motojouya/geezer_auth/pkg/shelter/user"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type accessTokenIssuerDBMock struct {
	getUserAccessToken func(identifier string, now time.Time) ([]dbUser.UserAccessTokenFull, error)
	dbUtility.SqlExecutorMock
}

func (mock accessTokenIssuerDBMock) GetUserAccessToken(identifier string, now time.Time) ([]dbUser.UserAccessTokenFull, error) {
	return mock.getUserAccessToken(identifier, now)
}

func getShelterUserAuthenticForAccToken(expectId string) *shelterUser.UserAuthentic {
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

func getLocalerMockForAccToken(t *testing.T, expectUUID uuid.UUID, now time.Time) *localUtility.LocalerMock {
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

func getAccessTokenIssueDbMock(t *testing.T, expectId string, expectToken string, firstNow time.Time) accessTokenIssuerDBMock {
	var insert = func(userAccessTokens ...interface{}) error {
		userAccessToken, ok := userAccessTokens[0].(dbUser.UserAccessToken)
		if !ok {
			t.Errorf("Expected userAccessTokens[0] to be of type dbUser.UserAccessToken, got %T", userAccessTokens[0])
			return nil
		}
		assert.Equal(t, userAccessToken.AccessToken, expectToken, "Expected token to match")
		assert.WithinDuration(t, userAccessToken.ExpireDate, firstNow.AddDate(0, 0, 7), time.Second, "Expected expiration date to be 7 days from now")
		return nil
	}
	var getUserAccessToken = func(identifier string, now time.Time) ([]dbUser.UserAccessTokenFull, error) {
		assert.Equal(t, expectId, identifier, "Expected identifier to match")
		assert.WithinDuration(t, now, firstNow, time.Second, "Expected 'now' to be within 1 second of current time")
		return []dbUser.UserAccessTokenFull{}, nil
	}
	return accessTokenIssuerDBMock{
		SqlExecutorMock: dbUtility.SqlExecutorMock{
			FakeInsert: insert,
		},
		getUserAccessToken: getUserAccessToken,
	}
}

func getJwtHandlerMock(t *testing.T, expectId string, expectToken string, expectUUID string, firstNow time.Time) *jwtUtility.JwtHandlerMock {
	var generate = func(user *pkgUser.User, now time.Time, tokenId string) (*pkgUser.Authentic, pkgText.JwtToken, error) {
		assert.Equal(t, expectId, string(user.Identifier), "Expected token ID to match")
		assert.Equal(t, expectUUID, tokenId, "Expected token ID to match")
		assert.WithinDuration(t, now, firstNow, time.Second, "Expected 'now' to be within 1 second of current time")
		return &pkgUser.Authentic{
			RegisteredClaims: gojwt.RegisteredClaims{
				ExpiresAt: gojwt.NewNumericDate(firstNow.AddDate(0, 0, 7)),
			},
		}, pkgText.JwtToken(expectToken), nil
	}
	return &jwtUtility.JwtHandlerMock{
		FakeGenerate: generate,
	}
}

func getDbUserAccessTokenFull(persistKey uint, expectId string, expectToken string) dbUser.UserAccessTokenFull {
	return dbUser.UserAccessTokenFull{
		UserAccessToken: dbUser.UserAccessToken{
			PersistKey:       persistKey,
			UserPersistKey:   1,
			AccessToken:      expectToken,
			SourceUpdateDate: time.Now(),
			RegisterDate:     time.Now(),
			ExpireDate:       time.Now(),
		},
		UserIdentifier:     expectId,
		UserExposeEmailId:  "test@gmail.com",
		UserName:           "TestName",
		UserBotFlag:        false,
		UserRegisteredDate: time.Now(),
		UserUpdateDate:     time.Now(),
	}
}

func TestAccessTokenIssuer(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectUUID, _ = uuid.NewUUID()
	var expectToken = "test-access-token"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForAccToken(expectIdentifier)

	var localerMock = getLocalerMockForAccToken(t, expectUUID, firstNow)
	var dbMock = getAccessTokenIssueDbMock(t, expectIdentifier, expectToken, firstNow)
	var jwtHandlerMock = getJwtHandlerMock(t, expectIdentifier, expectToken, expectUUID.String(), firstNow)

	issuer := user.NewAccessTokenIssue(localerMock, dbMock, jwtHandlerMock)
	accessToken, err := issuer.Execute(userAuthentic)

	assert.NoError(t, err)
	assert.Equal(t, expectToken, string(accessToken), "Expected refresh token to match generated UUID")
}

func TestAccessTokenIssuerExistOne(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectUUID, _ = uuid.NewUUID()
	var expectToken = "test-access-token"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForAccToken(expectIdentifier)

	var localerMock = getLocalerMockForAccToken(t, expectUUID, firstNow)
	var dbMock = getAccessTokenIssueDbMock(t, expectIdentifier, expectToken, firstNow)
	var jwtHandlerMock = getJwtHandlerMock(t, expectIdentifier, expectToken, expectUUID.String(), firstNow)

	dbMock.getUserAccessToken = func(identifier string, now time.Time) ([]dbUser.UserAccessTokenFull, error) {
		return []dbUser.UserAccessTokenFull{
			dbUser.UserAccessTokenFull{},
		}, nil
	}

	issuer := user.NewAccessTokenIssue(localerMock, dbMock, jwtHandlerMock)
	accessToken, err := issuer.Execute(userAuthentic)

	assert.NoError(t, err)
	assert.Equal(t, expectToken, string(accessToken), "Expected refresh token to match generated UUID")
}

func TestAccessTokenIssuerExistOverOne(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectUUID, _ = uuid.NewUUID()
	var expectToken = "test-access-token"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForAccToken(expectIdentifier)

	var localerMock = getLocalerMockForAccToken(t, expectUUID, firstNow)
	var dbMock = getAccessTokenIssueDbMock(t, expectIdentifier, expectToken, firstNow)
	var jwtHandlerMock = getJwtHandlerMock(t, expectIdentifier, expectToken, expectUUID.String(), firstNow)

	var userAccessTokenFull1 = getDbUserAccessTokenFull(1, expectIdentifier, expectToken)
	var userAccessTokenFull2 = getDbUserAccessTokenFull(2, expectIdentifier, "some-string")
	dbMock.getUserAccessToken = func(identifier string, now time.Time) ([]dbUser.UserAccessTokenFull, error) {
		return []dbUser.UserAccessTokenFull{userAccessTokenFull1, userAccessTokenFull2}, nil
	}

	localerMock.FakeGenerateUUID = func() (uuid.UUID, error) {
		t.Error("GenerateUUID should not be called when there is an existing access token")
		return uuid.UUID{}, nil
	}

	issuer := user.NewAccessTokenIssue(localerMock, dbMock, jwtHandlerMock)
	accessToken, err := issuer.Execute(userAuthentic)

	assert.NoError(t, err)
	assert.Equal(t, expectToken, string(accessToken), "Expected refresh token to match generated UUID")
}

func TestAccessTokenIssuerExistOverOneErr(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectUUID, _ = uuid.NewUUID()
	var expectToken = "test-access-token"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForAccToken(expectIdentifier)

	var localerMock = getLocalerMockForAccToken(t, expectUUID, firstNow)
	var dbMock = getAccessTokenIssueDbMock(t, expectIdentifier, expectToken, firstNow)
	var jwtHandlerMock = getJwtHandlerMock(t, expectIdentifier, expectToken, expectUUID.String(), firstNow)

	// var userAccessTokenFull1 = getDbUserAccessTokenFull(1, expectIdentifier, expectToken)
	var userAccessTokenFull2 = getDbUserAccessTokenFull(2, expectIdentifier, "some-string")
	dbMock.getUserAccessToken = func(identifier string, now time.Time) ([]dbUser.UserAccessTokenFull, error) {
		return []dbUser.UserAccessTokenFull{dbUser.UserAccessTokenFull{}, userAccessTokenFull2}, nil
	}

	localerMock.FakeGenerateUUID = func() (uuid.UUID, error) {
		t.Error("GenerateUUID should not be called when there is an existing access token")
		return uuid.UUID{}, nil
	}

	issuer := user.NewAccessTokenIssue(localerMock, dbMock, jwtHandlerMock)
	_, err := issuer.Execute(userAuthentic)

	assert.Error(t, err)
}

func TestAccessTokenIssuerErrGetUserAccessToken(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectUUID, _ = uuid.NewUUID()
	var expectToken = "test-access-token"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForAccToken(expectIdentifier)

	var localerMock = getLocalerMockForAccToken(t, expectUUID, firstNow)
	var dbMock = getAccessTokenIssueDbMock(t, expectIdentifier, expectToken, firstNow)
	var jwtHandlerMock = getJwtHandlerMock(t, expectIdentifier, expectToken, expectUUID.String(), firstNow)

	dbMock.getUserAccessToken = func(identifier string, now time.Time) ([]dbUser.UserAccessTokenFull, error) {
		return []dbUser.UserAccessTokenFull{}, errors.New("test error")
	}

	issuer := user.NewAccessTokenIssue(localerMock, dbMock, jwtHandlerMock)
	_, err := issuer.Execute(userAuthentic)

	assert.Error(t, err)
}

func TestAccessTokenIssuerErrGenerateUUID(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectUUID, _ = uuid.NewUUID()
	var expectToken = "test-access-token"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForAccToken(expectIdentifier)

	var localerMock = getLocalerMockForAccToken(t, expectUUID, firstNow)
	var dbMock = getAccessTokenIssueDbMock(t, expectIdentifier, expectToken, firstNow)
	var jwtHandlerMock = getJwtHandlerMock(t, expectIdentifier, expectToken, expectUUID.String(), firstNow)

	localerMock.FakeGenerateUUID = func() (uuid.UUID, error) {
		return uuid.UUID{}, errors.New("test error")
	}

	issuer := user.NewAccessTokenIssue(localerMock, dbMock, jwtHandlerMock)
	_, err := issuer.Execute(userAuthentic)

	assert.Error(t, err)
}

func TestAccessTokenIssuerErrInsert(t *testing.T) {
	var expectIdentifier = "US-TESTES"
	var expectUUID, _ = uuid.NewUUID()
	var expectToken = "test-access-token"
	var firstNow = time.Now()
	var userAuthentic = getShelterUserAuthenticForAccToken(expectIdentifier)

	var localerMock = getLocalerMockForAccToken(t, expectUUID, firstNow)
	var dbMock = getAccessTokenIssueDbMock(t, expectIdentifier, expectToken, firstNow)
	var jwtHandlerMock = getJwtHandlerMock(t, expectIdentifier, expectToken, expectUUID.String(), firstNow)

	dbMock.FakeInsert = func(userAccessTokens ...interface{}) error {
		return errors.New("test error")
	}

	issuer := user.NewAccessTokenIssue(localerMock, dbMock, jwtHandlerMock)
	_, err := issuer.Execute(userAuthentic)

	assert.Error(t, err)
}
