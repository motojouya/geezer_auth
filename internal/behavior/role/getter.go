package user

import (
	roleQuery "github.com/motojouya/geezer_auth/internal/db/query/role"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
)

type RoleGetterDB interface {
	roleQuery.GetRoleQuery
}

type RoleGetter interface {
	Execute() ([]shelterRole.Role, error)
}

type RoleGet struct {
	db RoleGetterDB
}

func NewRoleGet(db RoleGetterDB) *RoleGet {
	return &RoleGet{
		db: db,
	}
}

func (getter RoleGet) Execute() ([]shelterRole.Role, error) {
	dbRoles, err := getter.db.GetRole()
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
