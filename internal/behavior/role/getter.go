package role

import (
	roleQuery "github.com/motojouya/geezer_auth/internal/db/query/role"
	entryCompanyUser "github.com/motojouya/geezer_auth/internal/entry/transfer/companyUser"
	shelterRole "github.com/motojouya/geezer_auth/internal/shelter/role"
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

func (getter RoleGet) Execute(entry entryCompanyUser.RoleGetter) (*shelterRole.Role, error) {
	roleLabel, err := entry.GetRoleLabel()
	if err != nil {
		return nil, err
	}

	dbRole, err := getter.db.GetRole(string(roleLabel))
	if err != nil {
		return nil, err
	}

	if dbRole == nil {
		return nil, nil
	}

	role, err := dbRole.ToShelterRole()
	if err != nil {
		return nil, err
	}

	return &role, nil
}
