package model

type CompanyRole struct {
	Company Company
	Role    Role
}

func NewCompanyRole(company Company, role Role) CompanyRole {
	return CompanyRole{
		Company: company,
		Role:    role,
	}
}
