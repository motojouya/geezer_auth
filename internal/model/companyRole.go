package model

type CompanyRole struct {
	Company Company
	Roles   []Role
}

func NewCompanyRole(company Company, roles []Role) CompanyRole {
	return CompanyRole{
		Company: company,
		Roles:   roles,
	}
}
