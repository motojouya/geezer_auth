package user

import (
	"github.com/motojouya/geezer_auth/pkg/shelter/text"
)

type Company struct {
	Identifier text.Identifier
	Name       text.Name
}

func NewCompany(identifier text.Identifier, name text.Name) Company {
	return Company{
		Identifier: identifier,
		Name:       name,
	}
}
