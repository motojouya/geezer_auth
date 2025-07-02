package common

import (
	core "github.com/motojouya/geezer_auth/internal/core/company"
)

type Company struct {
	Identifier string
	Name       string
}

func FromCoreCompany(c core.Company) Company {
	return Company{
		Identifier: string(c.Identifier),
		Name:       string(c.Name),
	}
}
