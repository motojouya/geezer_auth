package model

import (
	"time"
)

const CompanyExposeIdPrefix = "CP-"

type UnsavedCompany struct {
	ExposeId ExposeId
	Name     Name
}

type Company struct {
	CompanyId      uint
	RegisteredDate time.Time
	Roles          []Role
	UnsavedCompany
}

type CompanyWithRole struct {
	Company
	Roles          []Role
}

func NewCompanyExposeId(random string) (ExposeId, error) {
	return NewExposeId(CompanyExposeIdPrefix, random)
}

func CreateCompany(exposeId ExposeId, name Name) UnsavedCompany {
	return UnsavedCompany{
		ExposeId: exposeId,
		Name:     name,
	}
}

func NewCompany(companyId uint, exposeId ExposeId, name Name, registeredDate time.Time, roles []Role) Company {
	return Company{
		CompanyId:      companyId,
		RegisteredDate: registeredDate,
		UnsavedCompany: CreateCompany(exposeId, name),
	}
}

func NewCompanyWithRole(companyId uint, exposeId ExposeId, name Name, registeredDate time.Time, roles []Role) CompanyWithRole {
	return CompanyWithRole{
		Company: NewCompany(companyId, exposeId, name, registeredDate),
		Roles:   roles,
	}
}
