package model

import (
	"time"
)

const CompanyExposeIdPrefix = "CP-"

type UnsavedCompany struct {
	exposeId string
	name      string
}

type Company struct {
	companyId     uint
	registeredDate time.Time
	roles          []Role
	UnsavedCompany
}

func CreateCompany(exposeId string, name string): UnsavedCompany {
	return UnsavedCompany{
		expose_id: exposeId,
		name:      name,
	}
}

func NewCompany(companyId uint, exposeId string, name string, registeredDate time.Time, roles []Role): Company {
	return Company{
		companyId:     companyId,
		registeredDate: registeredDate,
		roles:          roles,
		UnsavedCompany: CreateCompany(exposeId, name),
	}
}
