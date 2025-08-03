package company

import (
	text "github.com/motojouya/geezer_auth/pkg/core/text"
	user "github.com/motojouya/geezer_auth/pkg/core/user"
	"time"
)

const CompanyIdentifierPrefix = "CP-"

type UnsavedCompany struct {
	user.Company
	RegisteredDate time.Time
}

type Company struct {
	PersistKey uint
	UnsavedCompany
}

func CreateCompanyIdentifier(random string) (text.Identifier, error) {
	return text.CreateIdentifier(CompanyIdentifierPrefix, random)
}

func CreateCompany(
	identifier text.Identifier,
	name text.Name,
	registeredDate time.Time,
) UnsavedCompany {
	return UnsavedCompany{
		Company:        user.NewCompany(identifier, name),
		RegisteredDate: registeredDate,
	}
}

func NewCompany(
	persistKey uint,
	identifier text.Identifier,
	name text.Name,
	registeredDate time.Time,
) Company {
	return Company{
		PersistKey: persistKey,
		UnsavedCompany: UnsavedCompany{
			Company:        user.NewCompany(identifier, name),
			RegisteredDate: registeredDate,
		},
	}
}
