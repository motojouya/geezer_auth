package role

import (
	roleQuery "github.com/motojouya/geezer_auth/internal/db/query/role"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
	entryCompanyUser "github.com/motojouya/geezer_auth/internal/entry/transfer/companyUser"
)

type RoleGetterDB interface {
	roleQuery.GetRoleQuery
}

type RoleGetter interface {
	Execute(entry entryCompanyUser.RoleGetter) (shelterRole.Role, error)
}

type RoleGet struct {
	db RoleGetterDB
}

func NewRoleGet(db RoleGetterDB) *RoleGet {
	return &RoleGet{
		db: db,
	}
}

func (getter RoleGet) Execute(entry entryCompanyUser.RoleGetter) (shelterRole.Role, error) {
	roleLabel, err := entry.GetRoleLabel()
	if err != nil {
		return shelterRole.Role{}, err
	}

	dbRole, err := getter.db.GetRole(string(roleLabel))
	if err != nil {
		return shelterRole.Role{}, err
	}

	if dbRole == nil {
		return shelterRole.Role{}, nil
	}

	role, err := dbRole.ToShelterRole()
	if err != nil {
		return shelterRole.Role{}, err
	}

	return role, nil
}
