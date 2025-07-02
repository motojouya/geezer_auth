package transfer

import (
	core "github.com/motojouya/geezer_auth/internal/core/user"
	"github.com/motojouya/geezer_auth/internal/utility"
	"time"
)

type CompanyRole struct {
	Company Company
	Roles   []Role
}

type User struct {
	Identifier      string
	IdentifierEmail string
	Name            string
	BotFlag         bool
	UpdateDate      time.Time
	CompanyRole     *CompanyRole
	Email           *string
}

func FromCoreUser(u *core.UserAuthentic) *User {
	var companyRole *CompanyRole = nil
	if u.CompanyRole != nil {
		companyRole = &CompanyRole{
			Company: FromCoreCompany(u.CompanyRole.Company),
			Roles:   utility.Map(u.CompanyRole.Roles, FromCoreRole),
		}
	}

	var email *string = nil
	if u.Email != nil {
		var emailStr = string(*u.Email)
		email = &emailStr
	}

	return &User{
		Identifier:      string(u.Identifier),
		IdentifierEmail: string(u.ExposeEmailId),
		Name:            string(u.Name),
		BotFlag:         u.BotFlag,
		UpdateDate:      u.UpdateDate,
		CompanyRole:     companyRole,
		Email:           email,
	}
}
