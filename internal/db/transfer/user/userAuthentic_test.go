package user_test

import (
	"github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getUserCompanyRoleFull(persistKey uint, userId string, companyId string, roleLabel string) *user.UserCompanyRoleFull {
	var now = time.Now()
	var expireDate = now.Add(1 * time.Hour)
	return &user.UserCompanyRoleFull{
		UserCompanyRole: user.UserCompanyRole{
			PersistKey:        persistKey,
			UserPersistKey:    persistKey + 1,
			CompanyPersistKey: persistKey + 2,
			RoleLabel:         roleLabel,
			RegisterDate:      now,
			ExpireDate:        &expireDate,
		},
		UserIdentifier:        userId,
		UserExposeEmailId:     "test02@example.com",
		UserName:              "TestUserName",
		UserBotFlag:           false,
		UserRegisteredDate:    now.Add(2 * time.Hour),
		UserUpdateDate:        now.Add(3 * time.Hour),
		CompanyIdentifier:     companyId,
		CompanyName:           "TestCompanyName",
		CompanyRegisteredDate: now.Add(4 * time.Hour),
		RoleName:              "TestRoleName",
		RoleDescription:       "TestRoleDescription",
		RoleRegisteredDate:    now.Add(5 * time.Hour),
	}
}

func TestToCoreUserAuthentic(t *testing.T) {

	var userId = "US-TESTES"
	var companyId = "CP-TESTES"
	var userCompanyRole1 = getUserCompanyRoleFull(1, userId, companyId, "TEST_ROLE")
	var userCompanyRole2 = getUserCompanyRoleFull(2, userId, companyId, "TOST_ROLE")
	var userCompanyRoles = []*user.UserCompanyRoleFull{userCompanyRole1, userCompanyRole2}

	var now = time.Now()
	var email = "test01@example.com"
	var userAuthentic = &user.UserAuthentic{
		UserPersistKey:     1,
		UserIdentifier:     userId,
		UserExposeEmailId:  "test02@example.com",
		UserName:           "TestUserName",
		UserBotFlag:        false,
		UserRegisteredDate: now,
		UserUpdateDate:     now.Add(1 * time.Hour),
		Email:              &email,
		UserCompanyRole:    userCompanyRoles,
	}

	var coreUserAuthentic, err = userAuthentic.ToCoreUserAuthentic()

	assert.Nil(t, err)

	assert.Equal(t, uint(1), coreUserAuthentic.PersistKey)
	assert.Equal(t, userId, string(coreUserAuthentic.Identifier))
	assert.Equal(t, "test02@example.com", string(coreUserAuthentic.ExposeEmailId))
	assert.Equal(t, "TestUserName", string(coreUserAuthentic.Name))
	assert.Equal(t, false, coreUserAuthentic.BotFlag)
	assert.Equal(t, now, coreUserAuthentic.RegisteredDate)
	assert.Equal(t, now.Add(1*time.Hour), coreUserAuthentic.UpdateDate)

	assert.Equal(t, email, string(*coreUserAuthentic.Email))
	assert.Equal(t, companyId, string(coreUserAuthentic.CompanyRole.Company.Identifier))
	assert.Equal(t, 2, len(coreUserAuthentic.CompanyRole.Roles))
	assert.Equal(t, "TEST_ROLE", coreUserAuthentic.CompanyRole.Roles[0].Label)
	assert.Equal(t, "TOST_ROLE", coreUserAuthentic.CompanyRole.Roles[1].Label)

	t.Logf("coreUserAuthentic: %+v", coreUserAuthentic)
	t.Logf("coreUserAuthentic.Email: %s", *coreUserAuthentic.Email)
	t.Logf("coreUserAuthentic.CompanyRole: %v", *coreUserAuthentic.CompanyRole)
}

func TestToCoreUserAuthenticError(t *testing.T) {

	var userId = "US-TESTES"
	var companyId = "CP-TESTES"
	var userCompanyRole1 = getUserCompanyRoleFull(1, userId, companyId, "TEST_ROLE")
	var userCompanyRole2 = getUserCompanyRoleFull(2, userId, companyId, "TOST_ROLE")
	var userCompanyRoles = []*user.UserCompanyRoleFull{userCompanyRole1, userCompanyRole2}

	var now = time.Now()
	var email = "test01@example.com"
	var userAuthentic = &user.UserAuthentic{
		UserPersistKey:     1,
		UserIdentifier:     "invalid-idenfier",
		UserExposeEmailId:  "test02@example.com",
		UserName:           "TestUserName",
		UserBotFlag:        false,
		UserRegisteredDate: now,
		UserUpdateDate:     now.Add(1 * time.Hour),
		Email:              &email,
		UserCompanyRole:    userCompanyRoles,
	}

	var _, err = userAuthentic.ToCoreUserAuthentic()

	assert.Error(t, err)

	t.Logf("Error: %v", err)
}
