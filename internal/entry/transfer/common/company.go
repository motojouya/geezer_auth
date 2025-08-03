package common

import (
	shelter "github.com/motojouya/geezer_auth/internal/shelter/company"
)

type Company struct {
	Identifier string
	Name       string
}

func FromCoreCompany(c shelter.Company) Company {
	return Company{
		Identifier: string(c.Identifier),
		Name:       string(c.Name),
	}
}
