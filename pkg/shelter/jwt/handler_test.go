package jwt_test

import (
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/motojouya/geezer_auth/pkg/core/jwt"
	"github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/motojouya/geezer_auth/pkg/core/user"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getUserForHandler() *user.User {
	var roleLabel, _ = text.NewLabel("TestRole")
	var roleName, _ = text.NewName("TestRoleName")
	var role = user.NewRole(roleLabel, roleName)
	var roles = []user.Role{role}

	var companyIdentifier, _ = text.NewIdentifier("CP-TESTES")
	var companyName, _ = text.NewName("TestCompany")
	var company = user.NewCompany(companyIdentifier, companyName)

	var companyRole = user.NewCompanyRole(company, roles)

	var userIdentifier, _ = text.NewIdentifier("TestIdentifier")
	var emailId, _ = text.NewEmail("test@gmail.com")
	var email, _ = text.NewEmail("test_2@gmail.com")
	var userName, _ = text.NewName("TestName")
	var botFlag = false
	var updateDate = time.Now()

	return user.NewUser(userIdentifier, emailId, &email, userName, botFlag, companyRole, updateDate)
}

func TestHandleJwt(t *testing.T) {
	var roleLabel, _ = text.NewLabel("TEST_ROLE")
	var roleName, _ = text.NewName("TestRoleName")
	var role = user.NewRole(roleLabel, roleName)
	var roles = []user.Role{role}

	var companyIdentifier, _ = text.NewIdentifier("CP-TESTES")
	var companyName, _ = text.NewName("TestCompany")
	var company = user.NewCompany(companyIdentifier, companyName)

	var companyRole = user.NewCompanyRole(company, roles)

	var userIdentifier, _ = text.NewIdentifier("US-TESTES")
	var emailId, _ = text.NewEmail("test@gmail.com")
	var email, _ = text.NewEmail("test_2@gmail.com")
	var userName, _ = text.NewName("TestName")
	var botFlag = false
	var updateDate = time.Now()

	var userValue = user.NewUser(userIdentifier, emailId, &email, userName, botFlag, companyRole, updateDate)

	var issuedAt = time.Now()
	var id = "TestId"

	var issuer = "TestIssuer"
	var application = "TestAudience"
	var audience = []string{issuer, application}

	var latestKeyId = "TestLatestKeyId"
	var latestSecret = "TestLatestSecret"
	var oldKeyId = "TestOldKeyId"
	var oldSecret = "TestOldSecret"
	var validityPeriodMinutes uint = 60
	var expiresAt = issuedAt.Add(time.Duration(validityPeriodMinutes) * time.Minute)

	var jwtParsing = jwt.NewJwtParse(issuer, application, latestKeyId, latestSecret, oldKeyId, oldSecret)
	var jwtHandle = jwt.NewJwtHandle(audience, jwtParsing, validityPeriodMinutes)

	var _, tokenString, generateErr = jwtHandle.Generate(userValue, issuedAt, id)
	if generateErr != nil {
		t.Errorf("failed to generate token: %v", generateErr)
		return
	}

	var authentic, parseErr = jwtHandle.Parse(string(tokenString))
	if parseErr != nil {
		t.Errorf("failed to create token: %v", parseErr)
		return
	}

	assert.Equal(t, issuer, authentic.Issuer)
	assert.Equal(t, string(userIdentifier), authentic.Subject)
	assert.Equal(t, len(audience), len(authentic.Audience))
	assert.Equal(t, issuer, authentic.Audience[0])
	assert.Equal(t, application, authentic.Audience[1])
	assert.Equal(t, gojwt.NewNumericDate(expiresAt), authentic.ExpiresAt)
	assert.Equal(t, gojwt.NewNumericDate(issuedAt), authentic.NotBefore)
	assert.Equal(t, gojwt.NewNumericDate(issuedAt), authentic.IssuedAt)
	assert.Equal(t, id, authentic.ID)

	assert.Equal(t, string(userIdentifier), string(authentic.User.Identifier))
	assert.Equal(t, string(emailId), string(authentic.User.EmailId))
	assert.Equal(t, string(email), string(*authentic.User.Email))
	assert.Equal(t, string(userName), string(authentic.User.Name))
	assert.Equal(t, botFlag, authentic.User.BotFlag)
	assert.WithinDuration(t, updateDate, authentic.User.UpdateDate, time.Second)

	assert.Equal(t, string(companyIdentifier), string(authentic.User.CompanyRole.Company.Identifier))
	assert.Equal(t, string(companyName), string(authentic.User.CompanyRole.Company.Name))
	assert.Equal(t, string(roleLabel), string(authentic.User.CompanyRole.Roles[0].Label))
	assert.Equal(t, string(roleName), string(authentic.User.CompanyRole.Roles[0].Name))

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
	t.Logf("authentic.User.ExposeEmailId: %s", authentic.User.EmailId)
	t.Logf("authentic.User.Email: %s", *authentic.User.Email)
	t.Logf("authentic.User.Name: %s", authentic.User.Name)
	t.Logf("authentic.User.BotFlag: %t", authentic.User.BotFlag)
	t.Logf("authentic.User.UpdateDate: %s", authentic.User.UpdateDate)

	t.Logf("authentic.User.CompanyRole.Company: %+v", authentic.User.CompanyRole)
	t.Logf("authentic.User.CompanyRole.Company.Identifier: %s", authentic.User.CompanyRole.Company.Identifier)
	t.Logf("authentic.User.CompanyRole.Company.Name: %s", authentic.User.CompanyRole.Company.Name)
	t.Logf("authentic.User.CompanyRole.Roles[0].Label: %s", authentic.User.CompanyRole.Roles[0].Label)
	t.Logf("authentic.User.CompanyRole.Roles[0].Name: %s", authentic.User.CompanyRole.Roles[0].Name)
}

func TestHandleJwtFailureIssuer(t *testing.T) {
	var user = getUserForHandler()

	var issuedAt = time.Now()
	var id = "TestId"

	var issuer = "TestIssuer"
	var application = "TestAudience"
	var audience = []string{issuer, application}

	var latestKeyId = "TestLatestKeyId"
	var latestSecret = "TestLatestSecret"
	var oldKeyId = "TestOldKeyId"
	var oldSecret = "TestOldSecret"
	var validityPeriodMinutes uint = 60

	var parserServer = jwt.NewJwtParse(issuer, application, latestKeyId, latestSecret, oldKeyId, oldSecret)
	var jwtHandle = jwt.NewJwtHandle(audience, parserServer, validityPeriodMinutes)

	var _, tokenString, generateErr = jwtHandle.Generate(user, issuedAt, id)
	if generateErr != nil {
		t.Errorf("failed to generate token: %v", generateErr)
		return
	}

	var wrongIssuer = "WrongIssuer"
	var parserClient = jwt.NewJwtParse(wrongIssuer, application, latestKeyId, latestSecret, oldKeyId, oldSecret)

	var _, parseErr = parserClient.Parse(string(tokenString))
	if parseErr == nil {
		t.Errorf("failed to generate token: %v", parseErr)
	}
}

func TestHandleJwtFailureAudience(t *testing.T) {
	var user = getUserForHandler()

	var issuedAt = time.Now()
	var id = "TestId"

	var issuer = "TestIssuer"
	var application = "TestAudience"
	var audience = []string{issuer, application}

	var latestKeyId = "TestLatestKeyId"
	var latestSecret = "TestLatestSecret"
	var oldKeyId = "TestOldKeyId"
	var oldSecret = "TestOldSecret"
	var validityPeriodMinutes uint = 60

	var parserServer = jwt.NewJwtParse(issuer, application, latestKeyId, latestSecret, oldKeyId, oldSecret)
	var jwtHandle = jwt.NewJwtHandle(audience, parserServer, validityPeriodMinutes)

	var _, tokenString, generateErr = jwtHandle.Generate(user, issuedAt, id)
	if generateErr != nil {
		t.Errorf("failed to generate token: %v", generateErr)
		return
	}

	var wrongApplication = "WrongApplication"
	var parserClient = jwt.NewJwtParse(issuer, wrongApplication, latestKeyId, latestSecret, oldKeyId, oldSecret)

	var _, parseErr = parserClient.Parse(string(tokenString))
	if parseErr == nil {
		t.Errorf("failed to generate token: %v", parseErr)
	}
}

func TestHandleJwtFailureSecret(t *testing.T) {
	var user = getUserForHandler()

	var issuedAt = time.Now()
	var id = "TestId"

	var issuer = "TestIssuer"
	var application = "TestAudience"
	var audience = []string{issuer, application}

	var latestKeyId = "TestLatestKeyId"
	var latestSecret = "TestLatestSecret"
	var oldKeyId = "TestOldKeyId"
	var oldSecret = "TestOldSecret"
	var validityPeriodMinutes uint = 60

	var parserServer = jwt.NewJwtParse(issuer, application, latestKeyId, latestSecret, oldKeyId, oldSecret)
	var jwtHandle = jwt.NewJwtHandle(audience, parserServer, validityPeriodMinutes)

	var _, tokenString, generateErr = jwtHandle.Generate(user, issuedAt, id)
	if generateErr != nil {
		t.Errorf("failed to generate token: %v", generateErr)
		return
	}

	var parserClient = jwt.NewJwtParse(issuer, application, "WrongKey", latestSecret, "", "")

	var _, parseErr = parserClient.Parse(string(tokenString))
	if parseErr == nil {
		t.Errorf("failed to generate token: %v", parseErr)
	}
}
