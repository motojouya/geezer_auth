package company

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/common"
	shelter "github.com/motojouya/geezer_auth/internal/shelter/company"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"time"
)

type CompanyCreator interface {
	ToShelterCompany(identifier pkgText.Identifier, registerDate time.Time) (shelter.UnsavedCompany, error)
}

type CompanyCreate struct {
	Name string `json:"name"`
}

type CompanyCreateRequest struct {
	common.RequestHeader
	CompanyCreate
}

func (c CompanyCreateRequest) ToShelterCompany(identifier pkgText.Identifier, registerDate time.Time) (shelter.UnsavedCompany, error) {
	var name, nameErr = pkgText.NewName(c.Name)
	if nameErr != nil {
		return shelter.UnsavedCompany{}, nameErr
	}

	return shelter.CreateCompany(identifier, name, registerDate), nil
}
