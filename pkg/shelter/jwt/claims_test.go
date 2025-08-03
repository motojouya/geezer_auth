package jwt_test

import (
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/motojouya/geezer_auth/pkg/shelter/jwt"
	"github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/motojouya/geezer_auth/pkg/shelter/user"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getUser(userIdentifierStr string) *user.User {
	var companyIdentifier, _ = text.NewIdentifier("CP-TESTES")
	var companyName, _ = text.NewName("TestCompany")
	var company = user.NewCompany(companyIdentifier, companyName)

	var roleLabel, _ = text.NewLabel("TestRole")
	var roleName, _ = text.NewName("TestRoleName")
	var role = user.NewRole(roleLabel, roleName)
	var roles = []user.Role{role}

	var companyRole = user.NewCompanyRole(company, roles)

	var userIdentifier, _ = text.NewIdentifier(userIdentifierStr)
	var emailId, _ = text.NewEmail("test@gmail.com")
	var email, _ = text.NewEmail("test_2@gmail.com")
	var userName, _ = text.NewName("TestName")
	var botFlag = false
	var updateDate = time.Now()

	return user.NewUser(userIdentifier, emailId, &email, userName, botFlag, companyRole, updateDate)
}

func getRegisteredClaims(issuedAt time.Time, id uuid.UUID) gojwt.RegisteredClaims {
	var issuer = "issuer_id"
	var subject = "subject_id"
	var aud01 = "aud1"
	var aud02 = "aud2"
	var audience = []string{aud01, aud02}
	return gojwt.RegisteredClaims{
		Issuer:    issuer,
		Subject:   subject,
		Audience:  audience,
		ExpiresAt: gojwt.NewNumericDate(issuedAt),
		NotBefore: gojwt.NewNumericDate(issuedAt),
		IssuedAt:  gojwt.NewNumericDate(issuedAt),
		ID:        id.String(),
	}
}

func TestFromAuthentic(t *testing.T) {
	var companyIdentifier, _ = text.NewIdentifier("CP-TESTES")
	var companyName, _ = text.NewName("TestCompany")
	var company = user.NewCompany(companyIdentifier, companyName)

	var roleLabel, _ = text.NewLabel("TestRole")
	var roleName, _ = text.NewName("TestRoleName")
	var role = user.NewRole(roleLabel, roleName)
	var roles = []user.Role{role}

	var companyRole = user.NewCompanyRole(company, roles)

	var userIdentifier, _ = text.NewIdentifier("TestIdentifier")
	var emailId, _ = text.NewEmail("test@gmail.com")
	var email, _ = text.NewEmail("test_2@gmail.com")
	var userName, _ = text.NewName("TestName")
	var botFlag = false
	var updateDate = time.Now()

	var userValue = user.NewUser(userIdentifier, emailId, &email, userName, botFlag, companyRole, updateDate)

	var issuer = "issuer_id"
	var aud01 = "aud1"
	var aud02 = "aud2"
	var audience = []string{aud01, aud02}
	var issuedAt = time.Now()
	var id, _ = uuid.NewUUID()
	var validityPeriodMinutes uint = 60
	var expiresAt = issuedAt.Add(time.Duration(validityPeriodMinutes) * time.Minute)

	var authentic = user.CreateAuthentic(issuer, audience, issuedAt, validityPeriodMinutes, id.String(), userValue)

	var claims = jwt.FromAuthentic(authentic)

	assert.Equal(t, issuer, claims.RegisteredClaims.Issuer)
	assert.Equal(t, string(userIdentifier), claims.RegisteredClaims.Subject)
	assert.Equal(t, len(audience), len(claims.RegisteredClaims.Audience))
	assert.Equal(t, aud01, claims.RegisteredClaims.Audience[0])
	assert.Equal(t, aud02, claims.RegisteredClaims.Audience[1])
	assert.Equal(t, gojwt.NewNumericDate(expiresAt), claims.RegisteredClaims.ExpiresAt)
	assert.Equal(t, gojwt.NewNumericDate(issuedAt), claims.RegisteredClaims.NotBefore)
	assert.Equal(t, gojwt.NewNumericDate(issuedAt), claims.RegisteredClaims.IssuedAt)
	assert.Equal(t, id.String(), claims.RegisteredClaims.ID)

	assert.Equal(t, string(email), *claims.UserEmail)
	assert.Equal(t, string(userName), claims.UserName)
	assert.Equal(t, botFlag, claims.BotFlag)
	assert.WithinDuration(t, updateDate, claims.UpdateDate.Time, time.Second)
	assert.Equal(t, string(emailId), claims.UserEmailId)
	assert.Equal(t, string(companyIdentifier), *claims.CompanyIdentifier)
	assert.Equal(t, string(companyName), *claims.CompanyName)
	assert.Equal(t, len(roles), len(claims.CompanyRoles))
	assert.Equal(t, len(roles), len(claims.CompanyRoleNames))
	assert.Equal(t, string(roleLabel), claims.CompanyRoles[0])
	assert.Equal(t, string(roleName), claims.CompanyRoleNames[0])

	t.Logf("claims: %+v", claims)
	t.Logf("claims.RegisteredClaims: %+v", claims.RegisteredClaims)
	t.Logf("claims.Issuer: %s", claims.RegisteredClaims.Issuer)
	t.Logf("claims.Audience[0]: %s", claims.RegisteredClaims.Audience[0])
	t.Logf("claims.Audience[1]: %s", claims.RegisteredClaims.Audience[1])
	t.Logf("claims.ExpiresAt: %s", claims.RegisteredClaims.ExpiresAt)
	t.Logf("claims.NotBefore: %s", claims.RegisteredClaims.NotBefore)
	t.Logf("claims.IssuedAt: %s", claims.RegisteredClaims.IssuedAt)
	t.Logf("claims.ID: %s", claims.RegisteredClaims.ID)

	t.Logf("claims.User: %+v", userValue)
	t.Logf("claims.UserEmail: %s", *claims.UserEmail)
	t.Logf("claims.UserName: %s", claims.UserName)
	t.Logf("claims.UserEmailId: %s", claims.UserEmailId)
	t.Logf("claims.BotFlag: %t", claims.BotFlag)
	t.Logf("claims.UpdateDate: %s", claims.UpdateDate)
	t.Logf("claims.CompanyIdentifier: %s", *claims.CompanyIdentifier)
	t.Logf("claims.CompanyName: %s", *claims.CompanyName)
	t.Logf("claims.CompanyRoles[0]: %s", claims.CompanyRoles[0])
	t.Logf("claims.CompanyRoleNames[0]: %s", claims.CompanyRoleNames[0])
}

func TestFromAuthenticNil(t *testing.T) {
	var userIdentifier, _ = text.NewIdentifier("TestIdentifier")
	var emailId, _ = text.NewEmail("test@gmail.com")
	var userName, _ = text.NewName("TestName")
	var botFlag = false
	var updateDate = time.Now()

	var userValue = user.NewUser(userIdentifier, emailId, nil, userName, botFlag, nil, updateDate)

	var issuer = "issuer_id"
	var aud01 = "aud1"
	var aud02 = "aud2"
	var audience = []string{aud01, aud02}
	var issuedAt = time.Now()
	var id, _ = uuid.NewUUID()
	var validityPeriodMinutes uint = 60
	var expiresAt = issuedAt.Add(time.Duration(validityPeriodMinutes) * time.Minute)

	var authentic = user.CreateAuthentic(issuer, audience, issuedAt, validityPeriodMinutes, id.String(), userValue)

	var claims = jwt.FromAuthentic(authentic)

	assert.Equal(t, issuer, claims.RegisteredClaims.Issuer)
	assert.Equal(t, string(userIdentifier), claims.RegisteredClaims.Subject)
	assert.Equal(t, len(audience), len(claims.RegisteredClaims.Audience))
	assert.Equal(t, aud01, claims.RegisteredClaims.Audience[0])
	assert.Equal(t, aud02, claims.RegisteredClaims.Audience[1])
	assert.Equal(t, gojwt.NewNumericDate(expiresAt), claims.RegisteredClaims.ExpiresAt)
	assert.Equal(t, gojwt.NewNumericDate(issuedAt), claims.RegisteredClaims.NotBefore)
	assert.Equal(t, gojwt.NewNumericDate(issuedAt), claims.RegisteredClaims.IssuedAt)
	assert.Equal(t, id.String(), claims.RegisteredClaims.ID)

	assert.Nil(t, claims.UserEmail)
	assert.Equal(t, string(userName), claims.UserName)
	assert.Equal(t, botFlag, claims.BotFlag)
	assert.WithinDuration(t, updateDate, claims.UpdateDate.Time, time.Second)
	assert.Equal(t, string(emailId), claims.UserEmailId)
	assert.Nil(t, claims.CompanyIdentifier)
	assert.Nil(t, claims.CompanyName)
	assert.Nil(t, claims.CompanyRoles)
	assert.Nil(t, claims.CompanyRoleNames)

	t.Logf("claims: %+v", claims)
	t.Logf("claims.RegisteredClaims: %+v", claims.RegisteredClaims)
	t.Logf("claims.Issuer: %s", claims.RegisteredClaims.Issuer)
	t.Logf("claims.Audience[0]: %s", claims.RegisteredClaims.Audience[0])
	t.Logf("claims.Audience[1]: %s", claims.RegisteredClaims.Audience[1])
	t.Logf("claims.ExpiresAt: %s", claims.RegisteredClaims.ExpiresAt)
	t.Logf("claims.NotBefore: %s", claims.RegisteredClaims.NotBefore)
	t.Logf("claims.IssuedAt: %s", claims.RegisteredClaims.IssuedAt)
	t.Logf("claims.ID: %s", claims.RegisteredClaims.ID)

	t.Logf("claims.User: %+v", userValue)
	t.Logf("claims.UserName: %s", claims.UserName)
	t.Logf("claims.UserEmailId: %s", claims.UserEmailId)
	t.Logf("claims.BotFlag: %t", claims.BotFlag)
	t.Logf("claims.UpdateDate: %s", claims.UpdateDate)
	t.Logf("claims.CompanyRoles: %s", claims.CompanyRoles)
	t.Logf("claims.CompanyRoleNames: %s", claims.CompanyRoleNames)
}

func TestToAuthentic(t *testing.T) {
	var companyIdentifier, _ = text.NewIdentifier("CP-TESTES")
	var companyIdentifierStr = string(companyIdentifier)
	var companyName, _ = text.NewName("TestCompany")
	var companyNameStr = string(companyName)
	var roleLabel, _ = text.NewLabel("TEST_ROLE")
	var roleName, _ = text.NewName("TestRoleName")

	var userIdentifier, _ = text.NewIdentifier("US-TESTES")
	var emailId, _ = text.NewEmail("test@gmail.com")
	var email, _ = text.NewEmail("test_2@gmail.com")
	var emailStr = string(email)
	var userName, _ = text.NewName("TestName")
	var botFlag = false
	var updateDate = time.Now()

	var issuer = "issuer_id"
	var aud01 = "aud1"
	var aud02 = "aud2"
	var audience = []string{aud01, aud02}
	var issuedAt = time.Now()
	var validityPeriodMinutes uint = 60
	var expiresAt = issuedAt.Add(time.Duration(validityPeriodMinutes) * time.Minute)
	var id, _ = uuid.NewUUID()

	var registeredClaims = gojwt.RegisteredClaims{
		Issuer:    issuer,
		Subject:   string(userIdentifier),
		Audience:  audience,
		ExpiresAt: gojwt.NewNumericDate(expiresAt),
		NotBefore: gojwt.NewNumericDate(issuedAt),
		IssuedAt:  gojwt.NewNumericDate(issuedAt),
		ID:        id.String(),
	}

	var claims = &jwt.GeezerClaims{
		RegisteredClaims:  registeredClaims,
		UserEmail:         &emailStr,
		UserName:          string(userName),
		BotFlag:           botFlag,
		UpdateDate:        gojwt.NewNumericDate(updateDate),
		UserEmailId:       string(emailId),
		CompanyIdentifier: &companyIdentifierStr,
		CompanyName:       &companyNameStr,
		CompanyRoles:      []string{string(roleLabel)},
		CompanyRoleNames:  []string{string(roleName)},
	}

	var authentic, err = claims.ToAuthentic()
	if err != nil {
		t.Fatalf("ToAuthentic failed: %v", err)
	}

	assert.Equal(t, issuer, authentic.RegisteredClaims.Issuer)
	assert.Equal(t, string(userIdentifier), authentic.RegisteredClaims.Subject)
	assert.Equal(t, len(audience), len(authentic.RegisteredClaims.Audience))
	assert.Equal(t, aud01, authentic.RegisteredClaims.Audience[0])
	assert.Equal(t, aud02, authentic.RegisteredClaims.Audience[1])
	assert.Equal(t, gojwt.NewNumericDate(expiresAt), authentic.RegisteredClaims.ExpiresAt)
	assert.Equal(t, gojwt.NewNumericDate(issuedAt), authentic.RegisteredClaims.NotBefore)
	assert.Equal(t, gojwt.NewNumericDate(issuedAt), authentic.RegisteredClaims.IssuedAt)
	assert.Equal(t, id.String(), authentic.RegisteredClaims.ID)

	assert.Equal(t, string(email), string(*authentic.User.Email))
	assert.Equal(t, string(userName), string(authentic.User.Name))
	assert.Equal(t, botFlag, authentic.User.BotFlag)
	assert.WithinDuration(t, updateDate, authentic.User.UpdateDate, time.Second)
	assert.Equal(t, string(emailId), string(authentic.User.EmailId))
	assert.Equal(t, string(companyIdentifier), string(authentic.User.CompanyRole.Company.Identifier))
	assert.Equal(t, string(companyName), string(authentic.User.CompanyRole.Company.Name))
	assert.Equal(t, 1, len(authentic.User.CompanyRole.Roles))
	assert.Equal(t, string(roleLabel), string(authentic.User.CompanyRole.Roles[0].Label))
	assert.Equal(t, string(roleName), string(authentic.User.CompanyRole.Roles[0].Name))

	t.Logf("authentic: %+v", authentic)
	t.Logf("authentic.RegisteredClaims: %+v", authentic.RegisteredClaims)
	t.Logf("authentic.Issuer: %s", authentic.RegisteredClaims.Issuer)
	t.Logf("authentic.Audience[0]: %s", authentic.RegisteredClaims.Audience[0])
	t.Logf("authentic.Audience[1]: %s", authentic.RegisteredClaims.Audience[1])
	t.Logf("authentic.ExpiresAt: %s", authentic.RegisteredClaims.ExpiresAt)
	t.Logf("authentic.NotBefore: %s", authentic.RegisteredClaims.NotBefore)
	t.Logf("authentic.IssuedAt: %s", authentic.RegisteredClaims.IssuedAt)
	t.Logf("authentic.ID: %s", authentic.RegisteredClaims.ID)

	t.Logf("authentic.User: %+v", authentic.User)
	t.Logf("authentic.User.Email: %s", string(*authentic.User.Email))
	t.Logf("authentic.User.Name: %s", string(authentic.User.Name))
	t.Logf("authentic.User.EmailId: %s", string(authentic.User.EmailId))
	t.Logf("authentic.User.BotFlag: %t", authentic.User.BotFlag)
	t.Logf("authentic.User.UpdateDate: %s", authentic.User.UpdateDate)
	t.Logf("authentic.User.CompanyRole.Company.Identifier: %s", string(authentic.User.CompanyRole.Company.Identifier))
	t.Logf("authentic.User.CompanyRole.Company.Name: %s", string(authentic.User.CompanyRole.Company.Name))
	t.Logf("authentic.User.CompanyRole.Roles[0].Label: %s", string(authentic.User.CompanyRole.Roles[0].Label))
	t.Logf("authentic.User.CompanyRole.Roles[0].Name: %s", string(authentic.User.CompanyRole.Roles[0].Name))
}

func TestToAuthenticNil(t *testing.T) {
	var userIdentifier, _ = text.NewIdentifier("US-TESTES")
	var emailId, _ = text.NewEmail("test@gmail.com")
	var userName, _ = text.NewName("TestName")
	var botFlag = false
	var updateDate = time.Now()

	var issuer = "issuer_id"
	var aud01 = "aud1"
	var aud02 = "aud2"
	var audience = []string{aud01, aud02}
	var issuedAt = time.Now()
	var validityPeriodMinutes uint = 60
	var expiresAt = issuedAt.Add(time.Duration(validityPeriodMinutes) * time.Minute)
	var id, _ = uuid.NewUUID()

	var registeredClaims = gojwt.RegisteredClaims{
		Issuer:    issuer,
		Subject:   string(userIdentifier),
		Audience:  audience,
		ExpiresAt: gojwt.NewNumericDate(expiresAt),
		NotBefore: gojwt.NewNumericDate(issuedAt),
		IssuedAt:  gojwt.NewNumericDate(issuedAt),
		ID:        id.String(),
	}

	var claims = &jwt.GeezerClaims{
		RegisteredClaims:  registeredClaims,
		UserEmail:         nil,
		UserName:          string(userName),
		BotFlag:           botFlag,
		UpdateDate:        gojwt.NewNumericDate(updateDate),
		UserEmailId:       string(emailId),
		CompanyIdentifier: nil,
		CompanyName:       nil,
		CompanyRoles:      []string{},
		CompanyRoleNames:  []string{},
	}

	var authentic, err = claims.ToAuthentic()
	if err != nil {
		t.Fatalf("ToAuthentic failed: %v", err)
	}

	assert.Equal(t, issuer, authentic.RegisteredClaims.Issuer)
	assert.Equal(t, string(userIdentifier), authentic.RegisteredClaims.Subject)
	assert.Equal(t, len(audience), len(authentic.RegisteredClaims.Audience))
	assert.Equal(t, aud01, authentic.RegisteredClaims.Audience[0])
	assert.Equal(t, aud02, authentic.RegisteredClaims.Audience[1])
	assert.Equal(t, gojwt.NewNumericDate(expiresAt), authentic.RegisteredClaims.ExpiresAt)
	assert.Equal(t, gojwt.NewNumericDate(issuedAt), authentic.RegisteredClaims.NotBefore)
	assert.Equal(t, gojwt.NewNumericDate(issuedAt), authentic.RegisteredClaims.IssuedAt)
	assert.Equal(t, id.String(), authentic.RegisteredClaims.ID)

	assert.Nil(t, authentic.User.Email)
	assert.Equal(t, string(userName), string(authentic.User.Name))
	assert.Equal(t, botFlag, authentic.User.BotFlag)
	assert.WithinDuration(t, updateDate, authentic.User.UpdateDate, time.Second)
	assert.Equal(t, string(emailId), string(authentic.User.EmailId))
	assert.Nil(t, authentic.User.CompanyRole)

	t.Logf("authentic: %+v", authentic)
	t.Logf("authentic.RegisteredClaims: %+v", authentic.RegisteredClaims)
	t.Logf("authentic.Issuer: %s", authentic.RegisteredClaims.Issuer)
	t.Logf("authentic.Audience[0]: %s", authentic.RegisteredClaims.Audience[0])
	t.Logf("authentic.Audience[1]: %s", authentic.RegisteredClaims.Audience[1])
	t.Logf("authentic.ExpiresAt: %s", authentic.RegisteredClaims.ExpiresAt)
	t.Logf("authentic.NotBefore: %s", authentic.RegisteredClaims.NotBefore)
	t.Logf("authentic.IssuedAt: %s", authentic.RegisteredClaims.IssuedAt)
	t.Logf("authentic.ID: %s", authentic.RegisteredClaims.ID)

	t.Logf("authentic.User: %+v", authentic.User)
	t.Logf("authentic.User.Name: %s", string(authentic.User.Name))
	t.Logf("authentic.User.EmailId: %s", string(authentic.User.EmailId))
	t.Logf("authentic.User.BotFlag: %t", authentic.User.BotFlag)
	t.Logf("authentic.User.UpdateDate: %s", authentic.User.UpdateDate)
}

func getClaims() *jwt.GeezerClaims {
	var companyIdentifier, _ = text.NewIdentifier("CP-TESTES")
	var companyIdentifierStr = string(companyIdentifier)
	var companyName, _ = text.NewName("TestCompany")
	var companyNameStr = string(companyName)
	var roleLabel, _ = text.NewLabel("TestRole")
	var roleName, _ = text.NewName("TestRoleName")

	var userIdentifier, _ = text.NewIdentifier("TestIdentifier")
	var emailId, _ = text.NewEmail("test@gmail.com")
	var email, _ = text.NewEmail("test_2@gmail.com")
	var emailStr = string(email)
	var userName, _ = text.NewName("TestName")
	var botFlag = false
	var updateDate = time.Now()

	var issuer = "issuer_id"
	var aud01 = "aud1"
	var aud02 = "aud2"
	var audience = []string{aud01, aud02}
	var issuedAt = time.Now()
	var validityPeriodMinutes uint = 60
	var expiresAt = issuedAt.Add(time.Duration(validityPeriodMinutes) * time.Minute)
	var id, _ = uuid.NewUUID()

	var registeredClaims = gojwt.RegisteredClaims{
		Issuer:    issuer,
		Subject:   string(userIdentifier),
		Audience:  audience,
		ExpiresAt: gojwt.NewNumericDate(expiresAt),
		NotBefore: gojwt.NewNumericDate(issuedAt),
		IssuedAt:  gojwt.NewNumericDate(issuedAt),
		ID:        id.String(),
	}

	return &jwt.GeezerClaims{
		RegisteredClaims:  registeredClaims,
		UserEmail:         &emailStr,
		UserName:          string(userName),
		BotFlag:           botFlag,
		UpdateDate:        gojwt.NewNumericDate(updateDate),
		UserEmailId:       string(emailId),
		CompanyIdentifier: &companyIdentifierStr,
		CompanyName:       &companyNameStr,
		CompanyRoles:      []string{string(roleLabel)},
		CompanyRoleNames:  []string{string(roleName)},
	}
}

func TestToAuthenticError(t *testing.T) {
	var testTable = []struct {
		name   string
		change func(*jwt.GeezerClaims)
	}{
		{
			name: "Subject",
			change: func(claims *jwt.GeezerClaims) {
				claims.Subject = "WrongSubject"
			},
		},
		{
			name: "UserEmail",
			change: func(claims *jwt.GeezerClaims) {
				var userEmail = "WrongUserEmail"
				claims.UserEmail = &userEmail
			},
		},
		{
			name: "UserName",
			change: func(claims *jwt.GeezerClaims) {
				claims.UserName = ""
			},
		},
		{
			name: "UserEmailId",
			change: func(claims *jwt.GeezerClaims) {
				claims.UserEmailId = "WrongUserEmailId"
			},
		},
		{
			name: "CompanyIdentifier",
			change: func(claims *jwt.GeezerClaims) {
				var companyIdentifier = "WrongCompanyIdentifier"
				claims.CompanyIdentifier = &companyIdentifier
			},
		},
		{
			name: "CompanyName",
			change: func(claims *jwt.GeezerClaims) {
				var companyName = ""
				claims.CompanyName = &companyName
			},
		},
		{
			name: "CompanyRoles",
			change: func(claims *jwt.GeezerClaims) {
				claims.CompanyRoles = []string{""}
			},
		},
		{
			name: "CompanyRoleNames",
			change: func(claims *jwt.GeezerClaims) {
				claims.CompanyRoleNames = []string{""}
			},
		},
		{
			name: "CompanyIdentifier-Nil",
			change: func(claims *jwt.GeezerClaims) {
				claims.CompanyIdentifier = nil
			},
		},
		{
			name: "CompanyName-Nil",
			change: func(claims *jwt.GeezerClaims) {
				claims.CompanyName = nil
			},
		},
		{
			name: "CompanyRoles-Nil",
			change: func(claims *jwt.GeezerClaims) {
				claims.CompanyRoles = nil
			},
		},
		{
			name: "CompanyRoleNames-Nil",
			change: func(claims *jwt.GeezerClaims) {
				claims.CompanyRoleNames = nil
			},
		},
		{
			name: "CompanyRoles-len",
			change: func(claims *jwt.GeezerClaims) {
				claims.CompanyRoles = []string{"Role1", "Role2"}
			},
		},
		{
			name: "CompanyRoleNames-len",
			change: func(claims *jwt.GeezerClaims) {
				claims.CompanyRoleNames = []string{"RoleName1", "RoleName2"}
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {

			var claims = getClaims()
			tt.change(claims)

			var _, err = claims.ToAuthentic()

			assert.NotNil(t, err, "Expected error for %s", tt.name)
			if err != nil {
				t.Logf("Error for %s: %v", tt.name, err)
			} else {
				t.Errorf("Expected error for %s but got nil", tt.name)
			}
		})
	}
}
