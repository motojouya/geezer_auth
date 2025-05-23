package jwt_test

import (
	"github.com/motojouya/geezer_auth/pkg/accessToken"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"os"
)

func TestGenerateAccessToken(t *testing.T) {
	var companyExposeId = "CP-TESTES"
	var companyName = "TestCompany"
	var companyRole = "TestRole"
	var companyRoleName = "TestRoleName"
	var company = accessToken.CreateCompany(exposeId, name, role, roleName)

	var userExposeId = "TestExposeId"
	var emailId = "test@gmail.com"
	var email = "test_2@gmail.com"
	var userName = "TestName"
	var botFlag = false
	var updateDate = time.Now()
	var user = accessToken.NewUser(userExposeId, emailId, email, userName, botFlag, company, updateDate)

	var issuedAt = time.Now()
	var id = "TestId"

	var issuer = "TestIssuer"
	var application = "TestAudience"
	var audience = []string{issuer, application}
	var latestSecret = "TestSecretKeyId"
	var secretMap = map[string]string{latestSecret:"TestSecret"}
	var validityPeriodMinutes = 60
	var getId = func() (string, error) {
		return id, nil
	}
	var jwtParser = accessToken.NewJwtParser(issuer, application, latestSecret, secretMap)
	var jwtHandker = accessToken.NewJwtHandler(audience, jwtParser, validityPeriodMinutes, getId)

	var tokenString, err = jwtHandker.GenerateAccessToken(user, issuedAt)
	if err != nil {
		t.Errorf("failed to generate token: %v", err)
		return
	}

	var token, err = jwtHandker.GetUserFromAccessToken(tokenString)
	if err != nil {
		t.Errorf("failed to create token: %v", err)
		return
	}

	assert.Equal(t, issuer, token.Issuer)
	assert.Equal(t, userExposeId, token.Subject)
	assert.Equal(t, len(audience), len(token.Audience))
	assert.Equal(t, issuer, token.Audience[0])
	assert.Equal(t, application, token.Audience[1])
	assert.Equal(t, expiresAt, token.ExpiresAt)
	assert.Equal(t, issuedAt, token.NotBefore)
	assert.Equal(t, issuedAt, token.IssuedAt)
	assert.Equal(t, id, token.ID)

	assert.Equal(t, userExposeId, token.User.ExposeId)
	assert.Equal(t, emailId, token.User.ExposeEmailId)
	assert.Equal(t, email, *token.User.Email)
	assert.Equal(t, userName, token.User.Name)
	assert.Equal(t, botFlag, token.User.BotFlag)
	assert.Equal(t, updateDate, token.User.UpdateDate)

	assert.Equal(t, companyExposeId, token.User.Company.ExposeId)
	assert.Equal(t, companyName, token.User.Company.Name)
	assert.Equal(t, companyRole, token.User.Company.Role)
	assert.Equal(t, companyRoleName, token.User.Company.RoleName)

	t.Logf("token: %+v", token)
	t.Logf("token.Issuer: %s", token.Issuer)
	t.Logf("token.Subject: %s", token.Subject)
	t.Logf("token.Audience[0]: %s", token.Audience[0])
	t.Logf("token.Audience[1]: %s", token.Audience[1])
	t.Logf("token.ExpiresAt: %s", token.ExpiresAt)
	t.Logf("token.NotBefore: %s", token.NotBefore)
	t.Logf("token.IssuedAt: %s", token.IssuedAt)
	t.Logf("token.ID: %s", token.ID)

	t.Logf("user: %+v", token.User)
	t.Logf("user.ExposeId: %s", token.User.ExposeId)
	t.Logf("user.ExposeEmailId: %s", token.User.ExposeEmailId)
	t.Logf("user.Email: %s", *token.User.Email)
	t.Logf("user.Name: %s", token.User.Name)
	t.Logf("user.BotFlag: %t", token.User.BotFlag)
	t.Logf("user.UpdateDate: %t", token.User.UpdateDate)

	t.Logf("company: %+v", token.User.Company)
	t.Logf("company.ExposeId: %s", token.User.Company.ExposeId)
	t.Logf("company.Name: %s", token.User.Company.Name)
	t.Logf("company.Role: %s", token.User.Company.Role)
	t.Logf("company.RoleName: %s", token.User.Company.RoleName)
}

func TestGenerateAccessTokenFailureId(t *testing.T) {
	var companyExposeId = "CP-TESTES"
	var companyName = "TestCompany"
	var companyRole = "TestRole"
	var companyRoleName = "TestRoleName"
	var company = accessToken.CreateCompany(exposeId, name, role, roleName)

	var userExposeId = "TestExposeId"
	var emailId = "test@gmail.com"
	var email = "test_2@gmail.com"
	var userName = "TestName"
	var botFlag = false
	var updateDate = time.Now()
	var user = accessToken.NewUser(userExposeId, emailId, email, userName, botFlag, company, updateDate)

	var issuedAt = time.Now()
	var id = "TestId"

	var issuer = "TestIssuer"
	var application = "TestAudience"
	var audience = []string{issuer, application}
	var latestSecret = "TestSecretKeyId"
	var secretMap = map[string]string{latestSecret:"TestSecret"}
	var validityPeriodMinutes = 60
	var getId = func() (string, error) {
		return "", fmt.Error("failed to get id")
	}
	var jwtParser = accessToken.NewJwtParser(issuer, application, latestSecret, secretMap)
	var jwtHandker = accessToken.NewJwtHandler(audience, jwtParser, validityPeriodMinutes, getId)

	var tokenString, err = jwtHandker.GenerateAccessToken(user, issuedAt)
	if err == nil {
		t.Errorf("failed to generate token: %v", err)
	}
}

func TestGenerateAccessTokenFailureIssuer(t *testing.T) {
	var companyExposeId = "CP-TESTES"
	var companyName = "TestCompany"
	var companyRole = "TestRole"
	var companyRoleName = "TestRoleName"
	var company = accessToken.CreateCompany(exposeId, name, role, roleName)

	var userExposeId = "TestExposeId"
	var emailId = "test@gmail.com"
	var email = "test_2@gmail.com"
	var userName = "TestName"
	var botFlag = false
	var updateDate = time.Now()
	var user = accessToken.NewUser(userExposeId, emailId, email, userName, botFlag, company, updateDate)

	var issuedAt = time.Now()
	var id = "TestId"

	var issuer = "TestIssuer"
	var application = "TestAudience"
	var audience = []string{issuer, application}
	var latestSecret = "TestSecretKeyId"
	var secretMap = map[string]string{latestSecret:"TestSecret"}
	var validityPeriodMinutes = 60
	var getId = func() (string, error) {
		return id, nil
	}
	var parserServer = accessToken.NewJwtParser(issuer, application, latestSecret, secretMap)
	var jwtHandker = accessToken.NewJwtHandler(audience, parserServer, validityPeriodMinutes, getId)

	var tokenString, err = jwtHandker.GenerateAccessToken(user, issuedAt)
	if err != nil {
		t.Errorf("failed to generate token: %v", err)
		return
	}

	var wrongIssuer = "WrongIssuer"
	var parserClient = accessToken.NewJwtParser(wrongIssuer, application, latestSecret, secretMap)

	var token, err = parserClient.GetUserFromAccessToken(tokenString)
	if err == nil {
		t.Errorf("failed to generate token: %v", err)
	}
}

func TestGenerateAccessTokenFailureSecret(t *testing.T) {
	var companyExposeId = "CP-TESTES"
	var companyName = "TestCompany"
	var companyRole = "TestRole"
	var companyRoleName = "TestRoleName"
	var company = accessToken.CreateCompany(exposeId, name, role, roleName)

	var userExposeId = "TestExposeId"
	var emailId = "test@gmail.com"
	var email = "test_2@gmail.com"
	var userName = "TestName"
	var botFlag = false
	var updateDate = time.Now()
	var user = accessToken.NewUser(userExposeId, emailId, email, userName, botFlag, company, updateDate)

	var issuedAt = time.Now()
	var id = "TestId"

	var issuer = "TestIssuer"
	var application = "TestAudience"
	var audience = []string{issuer, application}
	var latestSecret = "TestSecretKeyId"
	var secretMap = map[string]string{latestSecret:"TestSecret"}
	var validityPeriodMinutes = 60
	var getId = func() (string, error) {
		return id, nil
	}
	var parserServer = accessToken.NewJwtParser(issuer, application, latestSecret, secretMap)
	var jwtHandker = accessToken.NewJwtHandler(audience, parserServer, validityPeriodMinutes, getId)

	var tokenString, err = jwtHandker.GenerateAccessToken(user, issuedAt)
	if err != nil {
		t.Errorf("failed to generate token: %v", err)
		return
	}

	var wrongSecretMap = map[string]string{"wrongKey":"TestSecret"}
	var parserClient = accessToken.NewJwtParser(issuer, application, latestSecret, wrongSecretMap)

	var token, err = parserClient.GetUserFromAccessToken(tokenString)
	if err == nil {
		t.Errorf("failed to generate token: %v", err)
	}
}
