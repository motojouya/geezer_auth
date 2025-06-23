package company

import (
	"github.com/motojouya/geezer_auth/internal/core/role"
	"github.com/motojouya/geezer_auth/internal/core/text"
	"time"
)

type UnsavedCompanyInvite struct {
	Company      Company
	Token        text.Token
	Role         role.Role
	RegisterDate time.Time
	ExpireDate   time.Time
}

type CompanyInvite struct {
	PersistKey uint
	UnsavedCompanyInvite
}

// FIXME 外から環境変数で設定できてもいいかも
const TokenValidityPeriodHours = 50

func CreateCompanyInvite(
	company Company,
	token text.Token,
	role role.Role,
	registerDate time.Time,
) UnsavedCompanyInvite {
	var expireDate = registerDate.Add(TokenValidityPeriodHours * time.Hour)

	return UnsavedCompanyInvite{
		Company:      company,
		Token:        token,
		Role:         role,
		RegisterDate: registerDate,
		ExpireDate:   expireDate,
	}
}

func NewCompanyInvite(
	persistKey uint,
	company Company,
	token text.Token,
	role role.Role,
	registerDate time.Time,
	expireDate time.Time,
) CompanyInvite {
	return CompanyInvite{
		PersistKey: persistKey,
		UnsavedCompanyInvite: UnsavedCompanyInvite{
			Company:      company,
			Token:        token,
			Role:         role,
			RegisterDate: registerDate,
			ExpireDate:   expireDate,
		},
	}
}
