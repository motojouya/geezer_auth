package company

import (
	core "github.com/motojouya/geezer_auth/internal/core/company"
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
	"time"
)

type CompanyCreate struct {
	Name string `json:"name"`
}

type CompanyCreateRequest struct {
	CompanyCreate CompanyCreate `http:"body"`
}

func (c CompanyCreateRequest) ToCoreCompany(identifier pkgText.Identifier, registerDate time.Time) (core.UnsavedCompany, error) {
	var name, nameErr = pkgText.NewName(c.CompanyCreate.Name)
	if nameErr != nil {
		return core.UnsavedCompany{}, nameErr
	}

	return core.CreateCompany(identifier, name, registerDate), nil
}
