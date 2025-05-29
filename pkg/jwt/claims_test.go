package jwt_test

import (
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/motojouya/geezer_auth/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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

func getRegisteredClaims(issuedAt time.Time, id uuid.UUID) jwt.RegisteredClaims {
	var issuer = "issuer_id"
	var subject = "subject_id"
	var aud01 = "aud1"
	var aud02 = "aud2"
	var audience = []string{aud01, aud02}
	return jwt.RegisteredClaims{
		Issuer:    issuer,
		Subject:   subject,
		Audience:  audience,
		ExpiresAt: jwt.NewNumericDate(issuedAt),
		NotBefore: jwt.NewNumericDate(issuedAt),
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ID:        id.String(),
	}
}

func TestFromAuthentic(t *testing.T) {
	var companyExposeId, _ = text.NewExposeId("CP-TESTES")
	var companyName, _ = text.NewName("TestCompany")
	var company = user.NewCompany(companyExposeId, companyName)

	var roleLabel, _ = text.NewLabel("TestRole")
	var roleName, _ = text.NewName("TestRoleName")
	var role = user.NewRole(roleLabel, roleName)
	var roles = []user.Roles{role}

	var companyRole = user.NewCompanyRole(company, roles)

	var userExposeId, _ = text.NewExposeId("TestExposeId")
	var emailId, _ = text.NewEmail("test@gmail.com")
	var email, _ = text.NewEmail("test_2@gmail.com")
	var userName, _ = text.NewName("TestName")
	var botFlag = false
	var updateDate = time.Now()

	var user = user.NewUser(userExposeId, emailId, email, userName, botFlag, companyRole, updateDate)

	var issuer = "issuer_id"
	var aud01 = "aud1"
	var aud02 = "aud2"
	var audience = []string{aud01, aud02}
	var issuedAt = time.Now()
	var id, _ := uuid.NewUUID()
	var validityPeriodMinutes = 60

	var authentic = user.CreateAuthentic(issuer, audience, issuedAt, validityPeriodMinutes, id, user)

	var claims = jwt.FromAuthentic(authentic)

	assert.Equal(t, issuer, claims.RegisteredClaims.Issuer)
	assert.Equal(t, string(userExposeId), claims.RegisteredClaims.Subject)
	assert.Equal(t, len(audience), len(claims.RegisteredClaims.Audience))
	assert.Equal(t, aud01, claims.RegisteredClaims.Audience[0])
	assert.Equal(t, aud02, claims.RegisteredClaims.Audience[1])
	assert.Equal(t, expiresAt, claims.RegisteredClaims.ExpiresAt)
	assert.Equal(t, notBefore, claims.RegisteredClaims.NotBefore)
	assert.Equal(t, issuedAt, claims.RegisteredClaims.IssuedAt)
	assert.Equal(t, id.String(), claims.RegisteredClaims.ID)

	assert.Equal(t, string(email), *claims.UserEmail)
	assert.Equal(t, string(userName), claims.UserName)
	assert.Equal(t, botFlag, claims.BotFlag)
	assert.Equal(t, updateDate, claims.UpdateDate)
	assert.Equal(t, string(emailId), claims.UserEmailId)
	assert.Equal(t, string(companyExposeId), claims.CompanyExposeId)
	assert.Equal(t, string(companyName), claims.CompanyName)
	assert.Equal(t, len(roles), len(claims.CompanyRoles))
	assert.Equal(t, len(roles), len(claims.CompanyRoleNames))
	assert.Equal(t, string(role), claims.CompanyRoles[0])
	assert.Equal(t, string(roleName), claims.CompanyRoleNames[0])

	t.Logf("claims: %+v", claims)
	t.Logf("claims.RegisteredClaims: %+v", claims.RegisteredClaims)
	t.Logf("claims.Issuer: %s", claims.RegisteredClaims.Issuer)
	t.Logf("claims.Audience[0]: %s", claims.RegisteredClaims.Audience[0])
	t.Logf("claims.Audience[1]: %s", claims.RegisteredClaims.Audience[1])
	t.Logf("claims.ExpiresAt: %t", claims.RegisteredClaims.ExpiresAt)
	t.Logf("claims.NotBefore: %t", claims.RegisteredClaims.NotBefore)
	t.Logf("claims.IssuedAt: %t", claims.RegisteredClaims.IssuedAt)
	t.Logf("claims.ID: %s", claims.RegisteredClaims.ID)

	t.Logf("claims.User: %+v", user)
	t.Logf("claims.UserEmail: %s", *claims.Email)
	t.Logf("claims.UserName: %s", claims.Name)
	t.Logf("claims.UserEmailId: %s", claims.UserEmailId)
	t.Logf("claims.BotFlag: %t", claims.BotFlag)
	t.Logf("claims.UpdateDate: %t", claims.UpdateDate)
	t.Logf("claims.CompanyExposeId: %s", claims.CompanyExposeId)
	t.Logf("claims.CompanyName: %s", claims.CompanyName)
	t.Logf("claims.CompanyRoles[0]: %s", claims.CompanyRoles[0])
	t.Logf("claims.CompanyRoleNames[0]: %s", claims.CompanyRoleNames[0])
}

func TestFromAuthenticNil(t *testing.T) {
	var userExposeId, _ = text.NewExposeId("TestExposeId")
	var emailId, _ = text.NewEmail("test@gmail.com")
	var userName, _ = text.NewName("TestName")
	var botFlag = false
	var updateDate = time.Now()

	var user = user.NewUser(userExposeId, emailId, nil, userName, botFlag, nil, updateDate)

	var issuer = "issuer_id"
	var aud01 = "aud1"
	var aud02 = "aud2"
	var audience = []string{aud01, aud02}
	var issuedAt = time.Now()
	var id, _ := uuid.NewUUID()
	var validityPeriodMinutes = 60

	var authentic = user.CreateAuthentic(issuer, audience, issuedAt, validityPeriodMinutes, id, user)

	var claims = jwt.FromAuthentic(authentic)

	assert.Equal(t, issuer, claims.RegisteredClaims.Issuer)
	assert.Equal(t, string(userExposeId), claims.RegisteredClaims.Subject)
	assert.Equal(t, len(audience), len(claims.RegisteredClaims.Audience))
	assert.Equal(t, aud01, claims.RegisteredClaims.Audience[0])
	assert.Equal(t, aud02, claims.RegisteredClaims.Audience[1])
	assert.Equal(t, expiresAt, claims.RegisteredClaims.ExpiresAt)
	assert.Equal(t, notBefore, claims.RegisteredClaims.NotBefore)
	assert.Equal(t, issuedAt, claims.RegisteredClaims.IssuedAt)
	assert.Equal(t, id.String(), claims.RegisteredClaims.ID)

	assert.Equal(t, nil, *claims.UserEmail)
	assert.Equal(t, string(userName), claims.UserName)
	assert.Equal(t, botFlag, claims.BotFlag)
	assert.Equal(t, updateDate, claims.UpdateDate)
	assert.Equal(t, string(emailId), claims.UserEmailId)
	assert.Equal(t, nil, claims.CompanyExposeId)
	assert.Equal(t, nil, claims.CompanyName)
	assert.Equal(t, nil, claims.CompanyRoles)
	assert.Equal(t, nil, claims.CompanyRoleNames)

	t.Logf("claims: %+v", claims)
	t.Logf("claims.RegisteredClaims: %+v", claims.RegisteredClaims)
	t.Logf("claims.Issuer: %s", claims.RegisteredClaims.Issuer)
	t.Logf("claims.Audience[0]: %s", claims.RegisteredClaims.Audience[0])
	t.Logf("claims.Audience[1]: %s", claims.RegisteredClaims.Audience[1])
	t.Logf("claims.ExpiresAt: %t", claims.RegisteredClaims.ExpiresAt)
	t.Logf("claims.NotBefore: %t", claims.RegisteredClaims.NotBefore)
	t.Logf("claims.IssuedAt: %t", claims.RegisteredClaims.IssuedAt)
	t.Logf("claims.ID: %s", claims.RegisteredClaims.ID)

	t.Logf("claims.User: %+v", user)
	t.Logf("claims.UserEmail: %s", *claims.Email)
	t.Logf("claims.UserName: %s", claims.Name)
	t.Logf("claims.UserEmailId: %s", claims.UserEmailId)
	t.Logf("claims.BotFlag: %t", claims.BotFlag)
	t.Logf("claims.UpdateDate: %t", claims.UpdateDate)
	t.Logf("claims.CompanyExposeId: %s", claims.CompanyExposeId)
	t.Logf("claims.CompanyName: %s", claims.CompanyName)
	t.Logf("claims.CompanyRoles: %s", claims.CompanyRoles)
	t.Logf("claims.CompanyRoleNames: %s", claims.CompanyRoleNames)
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

func TestToAuthentic(t *testing.T) {
	var companyExposeId, _ = text.NewExposeId("CP-TESTES")
	var companyName, _ = text.NewName("TestCompany")
	var roleLabel, _ = text.NewLabel("TestRole")
	var roleName, _ = text.NewName("TestRoleName")

	var userExposeId, _ = text.NewExposeId("TestExposeId")
	var emailId, _ = text.NewEmail("test@gmail.com")
	var email, _ = text.NewEmail("test_2@gmail.com")
	var userName, _ = text.NewName("TestName")
	var botFlag = false
	var updateDate = time.Now()

	var issuer = "issuer_id"
	var aud01 = "aud1"
	var aud02 = "aud2"
	var audience = []string{aud01, aud02}
	var issuedAt = time.Now()
	var expiresAt = issuedAt.Add(validityPeriodMinutes * time.Minute)
	var id, _ := uuid.NewUUID()
	var validityPeriodMinutes = 60

	var registeredClaims = gojwt.RegisteredClaims{
		Issuer:    issuer,
		Subject:   string(userExposeId),
		Audience:  audience,
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		NotBefore: jwt.NewNumericDate(issuedAt),
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ID:        id.String(),
	}

	var claims = &jwt.GeezerClaims{
		RegisteredClaims: registeredClaims
		UserEmail:        email,
		UserName:         userName,
		BotFlag:          botFlag,
		UpdateDate:       jwt.NewNumericDate(updateDate),
		UserEmailId:      string(emailId),
		CompanyExposeId:  string(companyExposeId),
		CompanyName:      string(companyName),
		CompanyRoles:     []string{string(roleLabel)},
		CompanyRoleNames: []string{string(roleName)},
	}

	var authentic = claims.ToAuthentic()

	assert.Equal(t, issuer, authentic.RegisteredClaims.Issuer)
	assert.Equal(t, string(userExposeId), authentic.RegisteredClaims.Subject)
	assert.Equal(t, len(audience), len(authentic.RegisteredClaims.Audience))
	assert.Equal(t, aud01, authentic.RegisteredClaims.Audience[0])
	assert.Equal(t, aud02, authentic.RegisteredClaims.Audience[1])
	assert.Equal(t, expiresAt, authentic.RegisteredClaims.ExpiresAt)
	assert.Equal(t, notBefore, authentic.RegisteredClaims.NotBefore)
	assert.Equal(t, issuedAt, authentic.RegisteredClaims.IssuedAt)
	assert.Equal(t, id.String(), authentic.RegisteredClaims.ID)

	assert.Equal(t, string(email), *string(authentic.User.UserEmail))
	assert.Equal(t, string(userName), string(authentic.User.UserName))
	assert.Equal(t, botFlag, authentic.User.BotFlag)
	assert.Equal(t, updateDate, authentic.User.UpdateDate)
	assert.Equal(t, string(emailId), string(authentic.User.UserEmailId))
	assert.Equal(t, string(companyExposeId), string(authentic.User.CompanyRole.Company.ExposeId))
	assert.Equal(t, string(companyName), string(authentic.User.CompanyRole.Company.Name))
	assert.Equal(t, 1, len(authentic.User.CompanyRole.Roles))
	assert.Equal(t, string(role), string(authentic.User.CompanyRole.Roles[0].Label))
	assert.Equal(t, string(roleName), string(authentic.User.CompanyRole.Roles[0].Name))

	t.Logf("authentic: %+v", authentic)
	t.Logf("authentic.RegisteredClaims: %+v", authentic.RegisteredClaims)
	t.Logf("authentic.Issuer: %s", authentic.RegisteredClaims.Issuer)
	t.Logf("authentic.Audience[0]: %s", authentic.RegisteredClaims.Audience[0])
	t.Logf("authentic.Audience[1]: %s", authentic.RegisteredClaims.Audience[1])
	t.Logf("authentic.ExpiresAt: %t", authentic.RegisteredClaims.ExpiresAt)
	t.Logf("authentic.NotBefore: %t", authentic.RegisteredClaims.NotBefore)
	t.Logf("authentic.IssuedAt: %t", authentic.RegisteredClaims.IssuedAt)
	t.Logf("authentic.ID: %s", authentic.RegisteredClaims.ID)

	t.Logf("authentic.User: %+v", user)
	t.Logf("authentic.User.Email: %s", string(*authentic.User.Email))
	t.Logf("authentic.User.Name: %s", string(authentic.User.Name))
	t.Logf("authentic.User.EmailId: %s", string(authentic.User.EmailId))
	t.Logf("authentic.User.BotFlag: %t", authentic.User.BotFlag)
	t.Logf("authentic.User.UpdateDate: %t", authentic.User.UpdateDate)
	t.Logf("authentic.User.CompanyRole.Company.ExposeId: %s", string(authentic.User.CompanyRole.Company.ExposeId))
	t.Logf("authentic.User.CompanyRole.Company.Name: %s", string(authentic.User.CompanyRole.Company.Name))
	t.Logf("authentic.User.CompanyRole.Roles[0].Label: %s", string(authentic.User.CompanyRole.Roles[0].Label))
	t.Logf("authentic.User.CompanyRole.Roles[0].Name: %s", string(authentic.User.CompanyRole.Roles[0].Name))
}

func TestToAuthenticNil(t *testing.T) {
	var userExposeId, _ = text.NewExposeId("TestExposeId")
	var emailId, _ = text.NewEmail("test@gmail.com")
	var userName, _ = text.NewName("TestName")
	var botFlag = false
	var updateDate = time.Now()

	var issuer = "issuer_id"
	var aud01 = "aud1"
	var aud02 = "aud2"
	var audience = []string{aud01, aud02}
	var issuedAt = time.Now()
	var expiresAt = issuedAt.Add(validityPeriodMinutes * time.Minute)
	var id, _ := uuid.NewUUID()
	var validityPeriodMinutes = 60

	var registeredClaims = gojwt.RegisteredClaims{
		Issuer:    issuer,
		Subject:   string(userExposeId),
		Audience:  audience,
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		NotBefore: jwt.NewNumericDate(issuedAt),
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ID:        id.String(),
	}

	var claims = &jwt.GeezerClaims{
		RegisteredClaims: registeredClaims
		UserEmail:        nil,
		UserName:         userName,
		BotFlag:          botFlag,
		UpdateDate:       jwt.NewNumericDate(updateDate),
		UserEmailId:      string(emailId),
		CompanyExposeId:  nil,
		CompanyName:      nil,
		CompanyRoles:     []string{},
		CompanyRoleNames: []string{},
	}

	var authentic = claims.ToAuthentic()

	assert.Equal(t, issuer, authentic.RegisteredClaims.Issuer)
	assert.Equal(t, string(userExposeId), authentic.RegisteredClaims.Subject)
	assert.Equal(t, len(audience), len(authentic.RegisteredClaims.Audience))
	assert.Equal(t, aud01, authentic.RegisteredClaims.Audience[0])
	assert.Equal(t, aud02, authentic.RegisteredClaims.Audience[1])
	assert.Equal(t, expiresAt, authentic.RegisteredClaims.ExpiresAt)
	assert.Equal(t, notBefore, authentic.RegisteredClaims.NotBefore)
	assert.Equal(t, issuedAt, authentic.RegisteredClaims.IssuedAt)
	assert.Equal(t, id.String(), authentic.RegisteredClaims.ID)

	assert.Equal(t, nil, authentic.User.UserEmail)
	assert.Equal(t, string(userName), string(authentic.User.UserName))
	assert.Equal(t, botFlag, authentic.User.BotFlag)
	assert.Equal(t, updateDate, authentic.User.UpdateDate)
	assert.Equal(t, string(emailId), string(authentic.User.UserEmailId))
	assert.Equal(t, nil, authentic.User.CompanyRole)

	t.Logf("authentic: %+v", authentic)
	t.Logf("authentic.RegisteredClaims: %+v", authentic.RegisteredClaims)
	t.Logf("authentic.Issuer: %s", authentic.RegisteredClaims.Issuer)
	t.Logf("authentic.Audience[0]: %s", authentic.RegisteredClaims.Audience[0])
	t.Logf("authentic.Audience[1]: %s", authentic.RegisteredClaims.Audience[1])
	t.Logf("authentic.ExpiresAt: %t", authentic.RegisteredClaims.ExpiresAt)
	t.Logf("authentic.NotBefore: %t", authentic.RegisteredClaims.NotBefore)
	t.Logf("authentic.IssuedAt: %t", authentic.RegisteredClaims.IssuedAt)
	t.Logf("authentic.ID: %s", authentic.RegisteredClaims.ID)

	t.Logf("authentic.User: %+v", user)
	t.Logf("authentic.User.Name: %s", string(authentic.User.Name))
	t.Logf("authentic.User.EmailId: %s", string(authentic.User.EmailId))
	t.Logf("authentic.User.BotFlag: %t", authentic.User.BotFlag)
	t.Logf("authentic.User.UpdateDate: %t", authentic.User.UpdateDate)
}
