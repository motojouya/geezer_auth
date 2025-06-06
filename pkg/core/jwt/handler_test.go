package jwt_test

import (
	"github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/motojouya/geezer_auth/pkg/core/user"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getUser() *user.User {
	var companyRole, _ = text.NewLabel("TestRole")
	var companyRoleName, _ = text.NewName("TestRoleName")
	var role = user.NewRole(companyRole, companyRoleName)
	var roles = []user.Role{role}

	var companyIdentifier, _ = text.NewIdentifier("CP-TESTES")
	var companyName, _ = text.NewName("TestCompany")
	var company = user.CreateCompany(identifier, name, role, roleName)

	var companyRole = user.NewCompanyRole(company, roles)

	var userIdentifier, _ = text.NewIdentifier("TestIdentifier")
	var emailId, _ = text.NewEmail("test@gmail.com")
	var email, _ = text.NewEmail("test_2@gmail.com")
	var userName, _ = text.NewName("TestName")
	var botFlag = false
	var updateDate = time.Now()

	return user.NewUser(userIdentifier, emailId, email, userName, botFlag, company, updateDate)
}

func TestHandleJwt(t *testing.T) {
	var user = getUser()

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
	var jwtParser = user.NewJwtParser(issuer, application, latestSecret, secretMap)
	var jwtHandker = user.NewJwtHandler(audience, jwtParser, validityPeriodMinutes, getId)

	var tokenString, err = jwtHandker.Generate(user, issuedAt)
	if err != nil {
		t.Errorf("failed to generate token: %v", err)
		return
	}

	var authentic, err = jwtHandker.Parse(tokenString)
	if err != nil {
		t.Errorf("failed to create token: %v", err)
		return
	}

	assert.Equal(t, issuer, authentic.Issuer)
	assert.Equal(t, userIdentifier, authentic.Subject)
	assert.Equal(t, len(audience), len(authentic.Audience))
	assert.Equal(t, issuer, authentic.Audience[0])
	assert.Equal(t, application, authentic.Audience[1])
	assert.Equal(t, expiresAt, authentic.ExpiresAt)
	assert.Equal(t, issuedAt, authentic.NotBefore)
	assert.Equal(t, issuedAt, authentic.IssuedAt)
	assert.Equal(t, id, authentic.ID)

	assert.Equal(t, string(userIdentifier), authentic.User.Identifier)
	assert.Equal(t, string(emailId), authentic.User.ExposeEmailId)
	assert.Equal(t, string(email), *authentic.User.Email)
	assert.Equal(t, string(userName), authentic.User.Name)
	assert.Equal(t, botFlag, authentic.User.BotFlag)
	assert.Equal(t, updateDate, authentic.User.UpdateDate)

	assert.Equal(t, string(companyIdentifier), authentic.User.CompanyRole.Company.Identifier)
	assert.Equal(t, string(companyName), authentic.User.CompanyRole.Company.Name)
	assert.Equal(t, string(companyRole), authentic.User.CompanyRole.Roles[0].Label)
	assert.Equal(t, string(companyRoleName), authentic.User.CompanyRole.Roles[0].Name)

	t.Logf("authentic: %+v", authentic)
	t.Logf("authentic.Issuer: %s", authentic.Issuer)
	t.Logf("authentic.Subject: %s", authentic.Subject)
	t.Logf("authentic.Audience[0]: %s", authentic.Audience[0])
	t.Logf("authentic.Audience[1]: %s", authentic.Audience[1])
	t.Logf("authentic.ExpiresAt: %s", authentic.ExpiresAt)
	t.Logf("authentic.NotBefore: %s", authentic.NotBefore)
	t.Logf("authentic.IssuedAt: %s", authentic.IssuedAt)
	t.Logf("authentic.ID: %s", authentic.ID)

	t.Logf("authentic.User: %+v", authentic.User)
	t.Logf("authentic.User.Identifier: %s", authentic.User.Identifier)
	t.Logf("authentic.User.ExposeEmailId: %s", authentic.User.ExposeEmailId)
	t.Logf("authentic.User.Email: %s", *authentic.User.Email)
	t.Logf("authentic.User.Name: %s", authentic.User.Name)
	t.Logf("authentic.User.BotFlag: %t", authentic.User.BotFlag)
	t.Logf("authentic.User.UpdateDate: %t", authentic.User.UpdateDate)

	t.Logf("authentic.User.CompanyRole.Company: %+v", authentic.User.CompanyRole)
	t.Logf("authentic.User.CompanyRole.Company.Identifier: %s", authentic.User.CompanyRole.Company.Identifier)
	t.Logf("authentic.User.CompanyRole.Company.Name: %s", authentic.User.CompanyRole.Company.Name)
	t.Logf("authentic.User.CompanyRole.Roles[0].Label: %s", authentic.User.CompanyRole.Roles[0].Label)
	t.Logf("authentic.User.CompanyRole.Roles[0].Name: %s", authentic.User.CompanyRole.Roles[0].Name)
}

func TestHandleJwtFailureId(t *testing.T) {
	var user = getUser()

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
	var jwtParser = user.NewJwtParser(issuer, application, latestSecret, secretMap)
	var jwtHandker = user.NewJwtHandler(audience, jwtParser, validityPeriodMinutes, getId)

	var tokenString, err = jwtHandker.Generate(user, issuedAt)
	if err == nil {
		t.Errorf("failed to generate token: %v", err)
	}
}

func TestHandleJwtFailureIssuer(t *testing.T) {
	var user = getUser()

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
	var parserServer = user.NewJwtParser(issuer, application, latestSecret, secretMap)
	var jwtHandker = user.NewJwtHandler(audience, parserServer, validityPeriodMinutes, getId)

	var tokenString, err = jwtHandker.Generate(user, issuedAt)
	if err != nil {
		t.Errorf("failed to generate token: %v", err)
		return
	}

	var wrongIssuer = "WrongIssuer"
	var parserClient = user.NewJwtParser(wrongIssuer, application, latestSecret, secretMap)

	var token, err = parserClient.Parse(tokenString)
	if err == nil {
		t.Errorf("failed to generate token: %v", err)
	}
}

func TestHandleJwtFailureAudience(t *testing.T) {
	var user = getUser()

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
	var parserServer = user.NewJwtParser(issuer, application, latestSecret, secretMap)
	var jwtHandker = user.NewJwtHandler(audience, parserServer, validityPeriodMinutes, getId)

	var tokenString, err = jwtHandker.Generate(user, issuedAt)
	if err != nil {
		t.Errorf("failed to generate token: %v", err)
		return
	}

	var wrongApplication = "WrongApplication"
	var parserClient = user.NewJwtParser(issuer, wrongApplication, latestSecret, secretMap)

	var token, err = parserClient.Parse(tokenString)
	if err == nil {
		t.Errorf("failed to generate token: %v", err)
	}
}

func TestHandleJwtFailureSecret(t *testing.T) {
	var user = getUser()

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
	var parserServer = user.NewJwtParser(issuer, application, latestSecret, secretMap)
	var jwtHandker = user.NewJwtHandler(audience, parserServer, validityPeriodMinutes, getId)

	var tokenString, err = jwtHandker.Generate(user, issuedAt)
	if err != nil {
		t.Errorf("failed to generate token: %v", err)
		return
	}

	var wrongSecretMap = map[string]string{"wrongKey":"TestSecret"}
	var parserClient = user.NewJwtParser(issuer, application, latestSecret, wrongSecretMap)

	var token, err = parserClient.Parse(tokenString)
	if err == nil {
		t.Errorf("failed to generate token: %v", err)
	}
}
