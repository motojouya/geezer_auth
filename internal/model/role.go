package model

import (
	"time"
	pkg "github.com/motojouya/geezer_auth/pkg/model"
)

// TODO CompanyがRoleを持つ際には、propertyのCompanyは不要だし、UserがCompanyRoleを持つ際にも不要。
// 再利用性高く、かつそういうデータ型がほしい場合に何を定義すべきかな

type RoleWithoutCompany struct {
	pkg.Role
	Description    Text
	RegisteredDate time.Time
}

type UnsavedRole struct {
	RoleWithoutCompany
	Company Company
}

type Role struct {
	RoleId uint
	UnsavedRole
}

func NewRoleWithoutCompany(name pkg.Name, label pkg.Label, description Text, registeredDate time.Time) RoleWithoutCompany {
	return RoleWithoutCompany{
		Role:           pkg.NewRole(label, name),
		Description:    description,
		RegisteredDate: registeredDate,
	}
}

func CreateRole(company Company, name pkg.Name, label pkg.Label, description Text, registeredDate time.Time) UnsavedRole {
	return UnsavedRole{
		Company:     company,
		RoleWithoutCompany: NewRoleWithoutCompany(
			Name:        name,
			Label:       label,
			Description: description,
			RegisteredDate: registeredDate,
		)
	}
}

func NewRole(roleId uint, name pkg.Name, label pkg.Label, description Text, registeredDate time.Time) Role {
	return Role{
		RoleId:         roleId,
		UnsavedRole:    CreateRole(name, label, description, registeredDate),
	}
}
