package company

import (
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
)

type CompanyGet struct {
    Identifier string `json:"identifier"`
}

type CompanyGetRequest struct {
	CompanyGet CompanyGet `http:"path"`
}

func (c CompanyGetRequest) GetCompanyIdentifier() (pkgText.Identifier, error) {
	return pkgText.NewIdentifier(c.CompanyGet.Identifier)
}
