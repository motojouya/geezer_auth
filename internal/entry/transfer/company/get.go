package company

import (
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
)

type CompanyGetter interface {
	GetCompanyIdentifier() (pkgText.Identifier, error)
}

type CompanyGet struct {
	Identifier string `param:"identifier"`
}

type CompanyGetRequest struct {
	CompanyGet
}

func (c CompanyGetRequest) GetCompanyIdentifier() (pkgText.Identifier, error) {
	return pkgText.NewIdentifier(c.Identifier)
}
