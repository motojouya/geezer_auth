package model

import (
	"time"
)

const CompanyExposeIdPrefix = "CP-"

type UnsavedCompany struct {
	ExposeId string
	Name     string
}

type Company struct {
	CompanyId      uint
	RegisteredDate time.Time
	Roles          []Role
	UnsavedCompany
}

func CreateCompany(exposeId string, name string) UnsavedCompany {
	return UnsavedCompany{
		ExposeId: exposeId,
		Name:      name,
	}
}

func NewCompany(companyId uint, exposeId string, name string, registeredDate time.Time, roles []Role) Company {
	return Company{
		CompanyId:      companyId,
		RegisteredDate: registeredDate,
		Roles:          roles,
		UnsavedCompany: CreateCompany(exposeId, name),
	}
}
