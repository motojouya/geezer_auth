package role

import (
	roleQuery "github.com/motojouya/geezer_auth/internal/db/query/role"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
)

type AllRoleGetterDB interface {
	roleQuery.GetAllRoleQuery
}

type AllRoleGetter interface {
	Execute() ([]shelterRole.Role, error)
}

type AllRoleGet struct {
	db AllRoleGetterDB
}

func NewAllRoleGet(db AllRoleGetterDB) *AllRoleGet {
	return &AllRoleGet{
		db: db,
	}
}

func (getter AllRoleGet) Execute() ([]shelterRole.Role, error) {
	dbRoles, err := getter.db.GetAllRole()
	if err != nil {
		return nil, err
	}

	var roles []shelterRole.Role
	for _, dbRole := range dbRoles {
		role, err := dbRole.ToShelterRole()
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}
