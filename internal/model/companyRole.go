package model

type CompanyRole struct {
	company Company
	role    Role
}

func NewCompanyRole(company Company, role Role) CompanyRole {
	return CompanyRole{
		company: company,
		role:    role,
	}
}
