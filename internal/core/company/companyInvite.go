package company

import (
	"time"
	"github.com/google/uuid"
)

type UnsavedCompanyInvite struct {
	Company        Company
	Token          uuid.UUID
	Role           Role
	RegisteredDate time.Time
	ExpireDate     time.Time
}

type CompanyInvite struct {
	CompanyInviteId uint
	UnsavedCompanyInvite
}

// FIXME 外から環境変数で設定できてもいいかも
const TokenValidityPeriodHours = 50

func CreateCompanyInvite(company Company, token uuid.UUID, role Role, registerDate time.Time) UnsavedCompanyInvite {
	var expireDate = registerDate.Add(TokenValidityPeriodHours * time.Hour)

	return UnsavedCompanyInvite{
		Company:      company,
		Token:        token,
		Role:         role,
		RegisterDate: registerDate,
		ExpireDate:   expireDate,
	}
}

func NewUserRefreshToken(companyInviteId uint, company Company, token uuid.UUID, role Role, registerDate time.Time, expireDate time.Time) CompanyInvite {
	return CompanyInvite{
		CompanyInviteId: companyInviteId,
		UnsavedCompanyInvite: UnsavedCompanyInvite{
			Company:      company,
			Token:        token,
			Role:         role,
			RegisterDate: registerDate,
			ExpireDate:   expireDate,
		},
	}
}
