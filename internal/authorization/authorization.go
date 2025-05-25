package authorization

import (
	"github.com/motojouya/geezer_auth/internal/model/role"
	pkg "github.com/motojouya/geezer_auth/pkg/model"
)

type Authorization struct {
	Permisions []role.RolePermission
}

func NewAuthorization(permissions []role.RolePermission) *Authorization {
	return &Authorization{
		Permisions: permissions,
	}
}

// TODO  DBアクセスしてロードする
func CreateAuthorization() *Authorization {
	var EmployeeLabel = pkg.NewLabel("EMPLOYEE")
	var EmployeePermission = role.NewRolePermission(EmployeeLabel, true, true, false, false, 5)

	var ManagerLabel = pkg.NewLabel("MANAGER")
	var ManagerPermission = role.NewRolePermission(ManagerLabel, true, true, true, true, 9)

	var permissions = []role.RolePermission{
		role.AnonymousPermission,
		role.RoleLessPermission,
		EmployeePermission,
		ManagerPermission,
	}

	return NewAuthorization(permissions)
}

// TODO working
func (auth *Authorization) getPriorityRolePermission(roles []role.Role) role.RolePermission {
	if len(roles) == 0 {
		return role.RoleLessPermission
	}

	priorityRole := roles[0]
	for _, r := range roles {
		if r.Priority > priorityRole.Priority {
			priorityRole = r
		}
	}

	return priorityRole
}

func (auth *Authorization) Authorize(authentic *pkg.Authentic) error {

}

func GenerateLargeCharactorString(length int, source string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = source[rand.Intn(len(source))]
	}
	return string(b)
}

func GenerateUUID() (UUID, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return UUID(""), err
	}

	return uuid, nil
}

func GetNow() time.Time {
	return time.Now()
}
