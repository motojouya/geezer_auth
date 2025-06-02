package user

import (
	"github.com/motojouya/geezer_auth/pkg/core/text"
)

type Company struct {
	ExposeId text.ExposeId
	Name     text.Name
}

func NewCompany(exposeId text.ExposeId, name text.Name) Company {
	return Company{
		ExposeId: exposeId,
		Name:     name,
	}
}
