package jwt_test

import (
	"github.com/motojouya/geezer_auth/pkg/accessToken"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"os"
)

func TestCreateGeezerTokenSuccess(t *testing.T) {
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

	var issuer = "TestIssuer"
	var application = "TestAudience"
	var audience = []string{issuer, application}
	var expiresAt = time.Now()
	var issuedAt = time.Now()
	var id = "TestId"

	var claims = accessToken.CreateClaims(user User, issuer, audience, expiresAt, issuedAt, id)

	var latestSecret = "TestSecretKeyId"
	var secretMap = map[string]string{latestSecret:"TestSecret"}
	var jwtParser = accessToken.NewJwtParser(issuer, application, latestSecret, secretMap)

	var token, err = jwtParser.CreateGeezerToken(claims)
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

func TestCreateGeezerTokenSuccessNilCompany(t *testing.T) {
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

	var issuer = "TestIssuer"
	var application = "TestAudience"
	var audience = []string{issuer, application}
	var expiresAt = time.Now()
	var issuedAt = time.Now()
	var id = "TestId"

	var claims = accessToken.CreateClaims(user User, issuer, audience, expiresAt, issuedAt, id)
	claims.CompanyExposeId = nil
	claims.CompanyName = nil
	claims.CompanyRole = nil
	claims.CompanyRoleName = nil

	var latestSecret = "TestSecretKeyId"
	var secretMap = map[string]string{latestSecret:"TestSecret"}
	var jwtParser = accessToken.NewJwtParser(issuer, application, latestSecret, secretMap)

	var token, err = jwtParser.CreateGeezerToken(claims)
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
	assert.Equal(t, nil, token.User.Company)

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
}

func TestCreateGeezerTokenFailureCompanyId(t *testing.T) {
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

	var issuer = "TestIssuer"
	var application = "TestAudience"
	var audience = []string{issuer, application}
	var expiresAt = time.Now()
	var issuedAt = time.Now()
	var id = "TestId"

	var claims = accessToken.CreateClaims(user User, issuer, audience, expiresAt, issuedAt, id)
	claims.CompanyExposeId = nil

	var latestSecret = "TestSecretKeyId"
	var secretMap = map[string]string{latestSecret:"TestSecret"}
	var jwtParser = accessToken.NewJwtParser(issuer, application, latestSecret, secretMap)

	var token, err = jwtParser.CreateGeezerToken(claims)
	if err == nil {
		// TODO error messageを確認する
		t.Errorf("expected error, but got nil")
	}
}

func TestCreateGeezerTokenFailureCompanyName(t *testing.T) {
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

	var issuer = "TestIssuer"
	var application = "TestAudience"
	var audience = []string{issuer, application}
	var expiresAt = time.Now()
	var issuedAt = time.Now()
	var id = "TestId"

	var claims = accessToken.CreateClaims(user User, issuer, audience, expiresAt, issuedAt, id)
	claims.CompanyName = nil

	var latestSecret = "TestSecretKeyId"
	var secretMap = map[string]string{latestSecret:"TestSecret"}
	var jwtParser = accessToken.NewJwtParser(issuer, application, latestSecret, secretMap)

	var token, err = jwtParser.CreateGeezerToken(claims)
	if err == nil {
		// TODO error messageを確認する
		t.Errorf("expected error, but got nil")
	}
}

func TestCreateGeezerTokenFailureCompanyRole(t *testing.T) {
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

	var issuer = "TestIssuer"
	var application = "TestAudience"
	var audience = []string{issuer, application}
	var expiresAt = time.Now()
	var issuedAt = time.Now()
	var id = "TestId"

	var claims = accessToken.CreateClaims(user User, issuer, audience, expiresAt, issuedAt, id)
	claims.CompanyRole = nil

	var latestSecret = "TestSecretKeyId"
	var secretMap = map[string]string{latestSecret:"TestSecret"}
	var jwtParser = accessToken.NewJwtParser(issuer, application, latestSecret, secretMap)

	var token, err = jwtParser.CreateGeezerToken(claims)
	if err == nil {
		// TODO error messageを確認する
		t.Errorf("expected error, but got nil")
	}
}

func TestCreateGeezerTokenFailureCompanyRoleName(t *testing.T) {
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

	var issuer = "TestIssuer"
	var application = "TestAudience"
	var audience = []string{issuer, application}
	var expiresAt = time.Now()
	var issuedAt = time.Now()
	var id = "TestId"

	var claims = accessToken.CreateClaims(user User, issuer, audience, expiresAt, issuedAt, id)
	claims.CompanyRoleName = nil

	var latestSecret = "TestSecretKeyId"
	var secretMap = map[string]string{latestSecret:"TestSecret"}
	var jwtParser = accessToken.NewJwtParser(issuer, application, latestSecret, secretMap)

	var token, err = jwtParser.CreateGeezerToken(claims)
	if err == nil {
		// TODO error messageを確認する
		t.Errorf("expected error, but got nil")
	}
}

func TestCreateGeezerTokenFailureIssuer(t *testing.T) {
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

	var wrongIssuer = "WrongIssuer"
	var application = "TestAudience"
	var audience = []string{wrongIssuer, application}
	var expiresAt = time.Now()
	var issuedAt = time.Now()
	var id = "TestId"

	var claims = accessToken.CreateClaims(user User, wrongIssuer, audience, expiresAt, issuedAt, id)

	var issuer = "TestIssuer"
	var latestSecret = "TestSecretKeyId"
	var secretMap = map[string]string{latestSecret:"TestSecret"}
	var jwtParser = accessToken.NewJwtParser(issuer, application, latestSecret, secretMap)

	var token, err = jwtParser.CreateGeezerToken(claims)
	if err == nil {
		// TODO error messageを確認する
		t.Errorf("expected error, but got nil")
	}
}

func TestCreateGeezerTokenFailureAudience(t *testing.T) {
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

	var issuer = "TestIssuer"
	var application = "TestAudience"
	var audience = []string{issuer, application}
	var expiresAt = time.Now()
	var issuedAt = time.Now()
	var id = "TestId"

	var claims = accessToken.CreateClaims(user User, issuer, audience, expiresAt, issuedAt, id)

	var wrongAudience = "WrongAudience"
	var latestSecret = "TestSecretKeyId"
	var secretMap = map[string]string{latestSecret:"TestSecret"}
	var jwtParser = accessToken.NewJwtParser(issuer, wrongAudience, latestSecret, secretMap)

	var token, err = jwtParser.CreateGeezerToken(claims)
	if err == nil {
		// TODO error messageを確認する
		t.Errorf("expected error, but got nil")
	}
}
