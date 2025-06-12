package company

import (
	"github.com/google/uuid"
	"time"
	"github.com/motojouya/geezer_auth/pkg/core/user"
)

type UnsavedCompanyInvite struct {
	Company      Company
	Token        uuid.UUID
	Role         user.Role
	RegisterDate time.Time
	ExpireDate   time.Time
}

type CompanyInvite struct {
	PersistKey uint
	UnsavedCompanyInvite
}

// FIXME 外から環境変数で設定できてもいいかも
const TokenValidityPeriodHours = 50

func CreateCompanyInvite(company Company, token uuid.UUID, role user.Role, registerDate time.Time) UnsavedCompanyInvite {
	var expireDate = registerDate.Add(TokenValidityPeriodHours * time.Hour)

	return UnsavedCompanyInvite{
		Company:      company,
		Token:        token,
		Role:         role,
		RegisterDate: registerDate,
		ExpireDate:   expireDate,
	}
}

func NewUserRefreshToken(
	persistKey uint,
	company Company,
	token uuid.UUID,
	role user.Role,
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
