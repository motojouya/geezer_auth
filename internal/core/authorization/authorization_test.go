package authorization_test

import (
	"github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/motojouya/geezer_auth/pkg/core/user"
	"github.com/motojouya/geezer_auth/internal/core/authorization"
	"github.com/motojouya/geezer_auth/internal/core/role"
	"github.com/motojouya/geezer_auth/pkg/utility"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getAuthentic(role user.Role) *user.Authentic {
	var companyIdentifier, _ = text.NewIdentifier("CP-TESTES")
	var companyName, _ = text.NewName("TestCompany")
	var company = user.NewCompany(companyIdentifier, companyName)

	var roles = []user.Roles{role}

	var companyRole = user.NewCompanyRole(company, roles)

	var userIdentifier = text.NewIdentifier("TestIdentifier")
	var emailId = text.NewEmail("test@gmail.com")
	var email = text.NewEmail("test_2@gmail.com")
	var userName = text.NewName("TestName")
	var botFlag = false
	var updateDate = time.Now()

	var user = user.NewUser(userIdentifier, emailId, email, userName, botFlag, companyRole, updateDate)

	var issuer = "issuer_id"
	var subject = "subject_id"
	var aud01 = "aud1"
	var aud02 = "aud2"
	var audience = []string{aud01, aud02}
	var expiresAt = time.Now()
	var notBefore = time.Now()
	var issuedAt = time.Now()
	var id, _ := uuid.NewUUID()

	return user.NewAuthentic(issuer, subject, audience, expiresAt, notBefore, issuedAt, id, user)
}

func getAuthenticRoleLess() *user.Authentic {
	var userIdentifier = text.NewIdentifier("TestIdentifier")
	var emailId = text.NewEmail("test@gmail.com")
	var email = text.NewEmail("test_2@gmail.com")
	var userName = text.NewName("TestName")
	var botFlag = false
	var updateDate = time.Now()

	var user = user.NewUser(userIdentifier, emailId, email, userName, botFlag, nil, updateDate)

	var issuer = "issuer_id"
	var subject = "subject_id"
	var aud01 = "aud1"
	var aud02 = "aud2"
	var audience = []string{aud01, aud02}
	var expiresAt = time.Now()
	var notBefore = time.Now()
	var issuedAt = time.Now()
	var id, _ := uuid.NewUUID()

	return user.NewAuthentic(issuer, subject, audience, expiresAt, notBefore, issuedAt, id, user)
}

func TestAuthorizeSuccess(t *testing.T) {
	var auth = authorization.NewAuthorization([]authorization.RolePermission{
		role.AnonymousPermission,
		role.RoleLessPermission,
		role.NewRolePermission("EMPLOYEE", true, true, false, false, 5),
		role.NewRolePermission("MANAGER", true, true, true, true, 9),
	})

	var requirePermission = role.NewRequirePermission(true, true, false, false)

	var roleLabel, _ = text.NewLabel("EMPLOYEE")
	var roleName, _ = text.NewName("作業者")
	var role = user.NewRole(roleLabel, roleName)
	var authentic = getAuthentic(role)

	err := auth.Authorize(requirePermission, authentic)

	assert.NoError(t, err)

	t.Logf("Authorization successful for user: %s", authentic.User.Identifier)
}

func TestAuthorizeFailure(t *testing.T) {
	var auth = authorization.NewAuthorization([]authorization.RolePermission{
		role.AnonymousPermission,
		role.RoleLessPermission,
		role.NewRolePermission("EMPLOYEE", true, true, false, false, 5),
		role.NewRolePermission("MANAGER", true, true, true, true, 9),
	})

	var requirePermission = role.NewRequirePermission(true, true, true, false)

	var roleLabel, _ = text.NewLabel("EMPLOYEE")
	var roleName, _ = text.NewName("作業者")
	var role = user.NewRole(roleLabel, roleName)
	var authentic = getAuthentic(role)

	err := auth.Authorize(requirePermission, authentic)

	assert.Error(t, err)
	if !errors.As(err, &authorization.AuthorizationError{}) {
		t.Errorf("Expected AuthorizationError, got: %T", err)
		return
	}

	t.Logf("Authorization failed for user: %s, error: %s", authentic.User.Identifier, err.Error())
}

func TestAuthorizeError(t *testing.T) {
	var auth = authorization.NewAuthorization([]authorization.RolePermission{
		role.AnonymousPermission,
		role.RoleLessPermission,
		role.NewRolePermission("EMPLOYEE", true, true, false, false, 5),
		role.NewRolePermission("MANAGER", true, true, true, true, 9),
	})

	var requirePermission = role.NewRequirePermission(true, true, true, false)

	var roleLabel, _ = text.NewLabel("SUSPICIOUS_PERSON")
	var roleName, _ = text.NewName("不審者")
	var role = user.NewRole(roleLabel, roleName)
	var authentic = getAuthentic(role)

	err := auth.Authorize(requirePermission, authentic)

	assert.Error(t, err)
	if !errors.As(err, &utility.NilError{}) {
		t.Errorf("Expected NilError, got: %T", err)
		return
	}

	t.Logf("Authorization failed for user: %s, error: %s", authentic.User.Identifier, err.Error())
}

func TestGetPriorityRolePermission(t *testing.T) {
	var permissions = []authorization.RolePermission{
		authorization.AnonymousPermission,
		authorization.RoleLessPermission,
		authorization.NewRolePermission("EMPLOYEE", true, true, false, false, 5),
		authorization.NewRolePermission("MANAGER", true, true, true, true, 9),
	})

	var roleLabel, _ = text.NewLabel("MANAGER")
	var roleName, _ = text.NewName("管理者")
	var role = user.NewRole(roleLabel, roleName)
	var authentic = getAuthentic(role)

	permission, err := auth.GetPriorityRolePermission(permissions, authentic)

	assert.NoError(t, err)
	assert.Equal(t, "MANAGER", string(permission.Label))

	t.Logf("Priority role permission for user: %s is: %s", authentic.User.Identifier, permission.Label.Value)
}

func TestGetPriorityRolePermissionAnonymous(t *testing.T) {
	var permissions = []authorization.RolePermission{
		authorization.AnonymousPermission,
		authorization.RoleLessPermission,
		authorization.NewRolePermission("EMPLOYEE", true, true, false, false, 5),
		authorization.NewRolePermission("MANAGER", true, true, true, true, 9),
	})

	permission, err := auth.GetPriorityRolePermission(permissions, nil)

	assert.NoError(t, err)
	assert.Equal(t, "ANONYMOUS", string(permission.Label))

	t.Logf("Priority role permission for user: %s is: %s", authentic.User.Identifier, permission.Label.Value)
}

func TestGetPriorityRolePermissionRoleLess(t *testing.T) {
	var permissions = []authorization.RolePermission{
		authorization.AnonymousPermission,
		authorization.RoleLessPermission,
		authorization.NewRolePermission("EMPLOYEE", true, true, false, false, 5),
		authorization.NewRolePermission("MANAGER", true, true, true, true, 9),
	})

	var authentic = getAuthenticRoleLess()

	permission, err := auth.GetPriorityRolePermission(permissions, authentic)

	assert.NoError(t, err)
	assert.Equal(t, "ROLE_LESS", string(permission.Label))

	t.Logf("Priority role permission for user: %s is: %s", authentic.User.Identifier, permission.Label.Value)
}

func TestGetPriorityRolePermissionNil(t *testing.T) {
	var permissions = []authorization.RolePermission{
		authorization.AnonymousPermission,
		authorization.RoleLessPermission,
		authorization.NewRolePermission("EMPLOYEE", true, true, false, false, 5),
		authorization.NewRolePermission("MANAGER", true, true, true, true, 9),
	})

	var roleLabel, _ = text.NewLabel("SUSPICIOUS_PERSON")
	var roleName, _ = text.NewName("不審者")
	var role = user.NewRole(roleLabel, roleName)
	var authentic = getAuthentic(role)

	permission, err := auth.GetPriorityRolePermission(permissions, authentic)

	assert.Error(t, err)
	if !errors.As(err, &utility.NilError{}) {
		t.Errorf("Expected NilError, got: %T", err)
		return
	}

	t.Logf("Authorization failed for user: %s, error: %s", authentic.User.Identifier, err.Error())
}
