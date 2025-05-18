package accessToken

import (
	"time"
	"github.com/motojouya/geezer_auth/internal/model"
	ext "github.com/motojouya/geezer_auth/pkg/accessToken"
)

// FIXME ここに実装するでいいのか
func ModelToAccessTokenUser(user *model.User, company *model.Company) *ext.User {
	var company *ext.Company = nil
	if user.CompanyRole != nil {
		// TODO role nameいりそう
		company = ext.NewCompany(
			user.CompanyRole.Company.ExposeId,
			user.CompanyRole.Company.Name,
			user.CompanyRole.Role.Label,
			user.CompanyRole.Role.Name,
		)
	}

	return ext.NewUser(
		user.ExposeId,
		user.ExposeEmailId,
		user.Email,
		user.Name,
		user.BotFlag,
		company,
		user.UpdateDate,
	)
}
