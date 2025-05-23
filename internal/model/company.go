package model

import (
	"time"
	pkg "github.com/motojouya/geezer_auth/pkg/model"
)

const CompanyExposeIdPrefix = "CP-"

type UnsavedCompany struct {
	pkg.Company
	RegisteredDate time.Time
}

type Company struct {
	CompanyId uint
	UnsavedCompany
}

func CreateCompanyExposeId(random string) (pkg.ExposeId, error) {
	return pkg.CreateExposeId(CompanyExposeIdPrefix, random)
}

func CreateCompany(exposeId ExposeId, name Name, registeredDate time.Time) UnsavedCompany {
	return UnsavedCompany{
		ExposeId:       exposeId,
		Name:           name,
		RegisteredDate: registeredDate,
	}
}

func NewCompany(companyId uint, exposeId pkg.ExposeId, name pkg.Name, registeredDate time.Time, roles []RoleWithoutCompany) Company {
	return Company{
		CompanyId:      companyId,
		pkg.Company: CreateCompany(exposeId, name, registeredDate),
	}
}
