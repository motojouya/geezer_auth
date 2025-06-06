package company

import (
	"time"
	text "github.com/motojouya/geezer_auth/pkg/core/text"
	user "github.com/motojouya/geezer_auth/pkg/core/user"
)

const CompanyExposeIdPrefix = "CP-"

type UnsavedCompany struct {
	user.Company
	RegisteredDate time.Time
}

type Company struct {
	PersistKey uint
	UnsavedCompany
}

func CreateCompanyExposeId(random string) (text.ExposeId, error) {
	return text.CreateExposeId(CompanyExposeIdPrefix, random)
}

func CreateCompany(
	exposeId text.ExposeId,
	name text.Name,
	registeredDate time.Time
) UnsavedCompany {
	return UnsavedCompany{
		user.Company:   user.NewCompany(exposeId, name)
		RegisteredDate: registeredDate,
	}
}

func NewCompany(
	persistKey uint,
	exposeId text.ExposeId,
	name text.Name,
	registeredDate time.Time
) Company {
	return Company{
		PersistKey:     persistKey,
		UnsavedCompany: UnsavedCompany{
			user.Company:    user.NewCompany(exposeId, name)
			RegisteredDate: registeredDate,
		},
	}
}
