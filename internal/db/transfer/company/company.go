package company

import (
	core "github.com/motojouya/geezer_auth/internal/core/company"
	text "github.com/motojouya/geezer_auth/pkg/core/text"
	"time"
)

type Company struct {
	PersistKey     uint
	Identifier     string
	Name           string
	RegisteredDate time.Time
}

func FromCoreCompany(coreCompany core.UnsavedCompany) Company {
	// PersistKey is zero value
	return Company{
		Identifier:     string(coreCompany.Identifier),
		Name:           string(coreCompany.Name),
		RegisteredDate: coreCompany.RegisteredDate,
	}
}

func (c Company) ToCoreCompany() (core.Company, error) {
	var identifier, idErr = text.NewIdentifier(c.Identifier)
	if idErr != nil {
		return core.Company{}, idErr
	}

	var name, nameErr = text.NewName(c.Name)
	if nameErr != nil {
		return core.Company{}, nameErr
	}

	return core.NewCompany(
		c.PersistKey,
		identifier,
		name,
		c.RegisteredDate,
	), nil
}
