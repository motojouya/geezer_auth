package user_test

import (
	"github.com/motojouya/geezer_auth/pkg/model/text"
	"github.com/motojouya/geezer_auth/pkg/model/user"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"github.com/google/uuid"
)

func getUser(userExposeIdStr string) *user.User {
	var companyExposeId, _ = text.NewExposeId("CP-TESTES")
	var companyName, _ = text.NewName("TestCompany")
	var company = user.NewCompany(companyExposeId, companyName)

	var roleLabel, _ = text.NewLabel("TestRole")
	var roleName, _ = text.NewName("TestRoleName")
	var role = user.NewRole(roleLabel, roleName)
	var roles = []user.Roles{role}

	var companyRole = user.NewCompanyRole(company, roles)

	var userExposeId = text.NewExposeId(userExposeIdStr)
	var emailId = text.NewEmail("test@gmail.com")
	var email = text.NewEmail("test_2@gmail.com")
	var userName = text.NewName("TestName")
	var botFlag = false
	var updateDate = time.Now()

	return user.NewUser(userExposeId, emailId, email, userName, botFlag, companyRole, updateDate)
}

func TestNewAuthentic(t *testing.T) {
	var userExposeId = "TestExposeId"
	var user = getUser(userExposeId)

	var issuer = "issuer_id"
	var subject = "subject_id"
	var aud01 = "aud1"
	var aud02 = "aud2"
	var audience = []string{aud01, aud02}
	var expiresAt = time.Now()
	var notBefore = time.Now()
	var issuedAt = time.Now()
	var id, _ := uuid.NewUUID()

	var authentic = user.NewAuthentic(issuer, subject, audience, expiresAt, notBefore, issuedAt, id, user)

	assert.Equal(t, issuer, authentic.Issuer)
	assert.Equal(t, subject, authentic.Subject)
	assert.Equal(t, len(audience), len(authentic.Audience))
	assert.Equal(t, aud01, authentic.Audience[0])
	assert.Equal(t, aud02, authentic.Audience[1])
	assert.Equal(t, expiresAt, authentic.ExpiresAt)
	assert.Equal(t, notBefore, authentic.NotBefore)
	assert.Equal(t, issuedAt, authentic.IssuedAt)
	assert.Equal(t, id.String(), authentic.ID)
	assert.Equal(t, userExposeId, string(authentic.User.ExposeId))

	t.Logf("authentic: %+v", authentic)
	t.Logf("authentic.Issuer: %s", authentic.Issuer)
	t.Logf("authentic.Audience[0]: %s", authentic.Audience[0])
	t.Logf("authentic.Audience[1]: %s", authentic.Audience[1])
	t.Logf("authentic.ExpiresAt: %t", authentic.ExpiresAt)
	t.Logf("authentic.NotBefore: %t", authentic.NotBefore)
	t.Logf("authentic.IssuedAt: %t", authentic.IssuedAt)
	t.Logf("authentic.ID: %s", authentic.ID)

	t.Logf("authentic.User: %+v", authentic.User)
	t.Logf("authentic.User.ExposeId: %s", authentic.User.ExposeId)
}

func TestCreateAuthentic(t *testing.T) {
	var userExposeId = "TestExposeId"
	var user = getUser(userExposeId)

	var issuer = "issuer_id"
	var aud01 = "aud1"
	var aud02 = "aud2"
	var audience = []string{aud01, aud02}
	var issuedAt = time.Now()
	var id, _ := uuid.NewUUID()
	var validityPeriodMinutes = 60

	var authentic = user.CreateAuthentic(issuer, audience, issuedAt, validityPeriodMinutes, id, user)

	var expiresAt = issuedAt.Add(validityPeriodMinutes * time.Minute)

	assert.Equal(t, issuer, authentic.Issuer)
	assert.Equal(t, userExposeId, authentic.Subject)
	assert.Equal(t, len(audience), len(authentic.Audience))
	assert.Equal(t, aud01, authentic.Audience[0])
	assert.Equal(t, aud02, authentic.Audience[1])
	assert.Equal(t, expiresAt, authentic.ExpiresAt)
	assert.Equal(t, issuedAt, authentic.NotBefore)
	assert.Equal(t, issuedAt, authentic.IssuedAt)
	assert.Equal(t, id.String(), authentic.ID)
	assert.Equal(t, userExposeId, string(authentic.User.ExposeId))

	t.Logf("authentic: %+v", authentic)
	t.Logf("authentic.Issuer: %s", authentic.Issuer)
	t.Logf("authentic.Audience[0]: %s", authentic.Audience[0])
	t.Logf("authentic.Audience[1]: %s", authentic.Audience[1])
	t.Logf("authentic.ExpiresAt: %t", authentic.ExpiresAt)
	t.Logf("authentic.NotBefore: %t", authentic.NotBefore)
	t.Logf("authentic.IssuedAt: %t", authentic.IssuedAt)
	t.Logf("authentic.ID: %s", authentic.ID)

	t.Logf("authentic.User: %+v", authentic.User)
	t.Logf("authentic.User.ExposeId: %s", authentic.User.ExposeId)
}
