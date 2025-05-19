package accessToken_test

import (
	"github.com/motojouya/geezer_auth/pkg/accessToken"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// TODO 未実装

func TestNewGeezerToken(t *testing.T) {
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
	var subject = "TestSubject"
	var audience = []string{"TestAudience1", "TestAudience2"}
	var expiresAt = time.Now()
	var notBefore = time.Now()
	var issuedAt = time.Now()
	var id = "TestId"

	var token = accessToken.NewGeezerToken(issuer, subject, audience, expiresAt, notBefore, issuedAt, id, user)

	assert.Equal(t, issuer, token.Issuer)
	assert.Equal(t, subject, token.Subject)
	assert.Equal(t, len(audience), len(token.Audience))
	assert.Equal(t, audience[0], token.Audience[0])
	assert.Equal(t, audience[1], token.Audience[1])
	assert.Equal(t, expiresAt, token.ExpiresAt)
	assert.Equal(t, notBefore, token.NotBefore)
	assert.Equal(t, issuedAt, token.IssuedAt)
	assert.Equal(t, id, token.ID)

	assert.Equal(t, userExposeId, token.User.ExposeId)
	assert.Equal(t, companyExposeId, token.User.Company.ExposeId)

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
	t.Logf("company: %+v", token.User.Company)
	t.Logf("company.ExposeId: %s", token.User.Company.ExposeId)
}

func TestCreateClaims(t *testing.T) {
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
	var audience = []string{"TestAudience1", "TestAudience2"}
	var expiresAt = time.Now()
	var issuedAt = time.Now()
	var id = "TestId"

	var claims = accessToken.CreateClaims(user User, issuer, audience, expiresAt, issuedAt, id)

	assert.Equal(t, issuer, claims.Issuer)
	assert.Equal(t, userExposeId, claims.Subject)
	assert.Equal(t, len(audience), len(claims.Audience))
	assert.Equal(t, audience[0], claims.Audience[0])
	assert.Equal(t, audience[1], claims.Audience[1])
	assert.Equal(t, expiresAt, claims.ExpiresAt) // TODO numericDateに変換する必要がある?
	assert.Equal(t, issuedAt, claims.NotBefore) // TODO numericDateに変換する必要がある?
	assert.Equal(t, issuedAt, claims.IssuedAt) // TODO numericDateに変換する必要がある?
	assert.Equal(t, id, claims.ID)

	assert.Equal(t, email, claims.UserEmail)
	assert.Equal(t, userName, claims.UserName)
	assert.Equal(t, updateDate, claims.UpdateDate) // TODO numericDateに変換する必要がある?
	assert.Equal(t, emailId, claims.UserEmailId)
	assert.Equal(t, botFlag, claims.BotFlag)
	assert.Equal(t, companyExposeId, claims.CompanyExposeId)
	assert.Equal(t, companyName, claims.CompanyName)
	assert.Equal(t, companyRole, claims.CompanyRole)
	assert.Equal(t, companyRoleName, claims.CompanyRoleName)

	t.Logf("claims: %+v", claims)
	t.Logf("claims.Issuer: %s", claims.Issuer)
	t.Logf("claims.Subject: %s", claims.Subject)
	t.Logf("claims.Audience.length: %d", len(claims.Audience))
	t.Logf("claims.Audience[0]: %s", claims.Audience[0])
	t.Logf("claims.Audience[1]: %s", claims.Audience[1])
	t.Logf("claims.ExpiresAt: %s", claims.ExpiresAt)
	t.Logf("claims.NotBefore: %s", claims.NotBefore)
	t.Logf("claims.IssuedAt: %s", claims.IssuedAt)
	t.Logf("claims.ID: %s", claims.ID)

	t.Logf("claims.UserEmail: %s", claims.UserEmail)
	t.Logf("claims.UserName: %s", claims.UserName)
	t.Logf("claims.UpdateDate: %s", claims.UpdateDate)
	t.Logf("claims.UserEmailId: %s", claims.UserEmailId)
	t.Logf("claims.BotFlag: %t", claims.BotFlag)
	t.Logf("claims.CompanyExposeId: %t", claims.CompanyExposeId)
	t.Logf("claims.CompanyName: %t", claims.CompanyName)
	t.Logf("claims.CompanyRole: %t", claims.CompanyRole)
	t.Logf("claims.CompanyRoleName: %t", claims.CompanyRoleName)
}
