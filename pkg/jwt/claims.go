package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/motojouya/geezer_auth/pkg/model/text"
	"github.com/motojouya/geezer_auth/pkg/model/user"
	"time"
)

// FIXME claimsのprivate keyが`github.com/motojouya/geezer_auth/`をprefixとしているが、本来は稼働するサーバのfqdnをprefixとして持つべき。
type GeezerClaims struct {
	jwt.RegisteredClaims
	UserEmail        *string   `json:"email"`
	UserName         string    `json:"name"`
	UpdateDate       time.Time `json:"update_at"`
	UserEmailId      string    `json:"github.com/motojouya/geezer_auth/email_id"`
	BotFlag          bool      `json:"github.com/motojouya/geezer_auth/bot_flag"`
	CompanyExposeId  *string   `json:"github.com/motojouya/geezer_auth/company_expose_id"`
	CompanyName      *string   `json:"github.com/motojouya/geezer_auth/company_name"`
	CompanyRoles     *[]string `json:"github.com/motojouya/geezer_auth/company_roles"`
	CompanyRoleNames *[]string `json:"github.com/motojouya/geezer_auth/company_role_names"`
}

func FromAuthentic(authentic *user.Authentic) *GeezerClaims {

	roleNames := make([]string, len(authentic.User.CompanyRole.Roles))
	roleLabels := make([]string, len(authentic.User.CompanyRole.Roles))
	for i, role := range authentic.User.CompanyRole.Roles {
		roleNames[i] = string(role.Name)
		roleLabels[i] = string(role.Label)
	}

	return &GeezerClaims{
		RegisteredClaims: authentic.RegisteredClaims,
		UserEmail:        string(authentic.User.Email),
		UserName:         string(authentic.User.Name),
		UpdateDate:       authentic.User.UpdateDate.Time,
		UserEmailId:      string(authentic.User.EmailId),
		BotFlag:          authentic.User.BotFlag,
		CompanyExposeId:  string(authentic.User.CompanyRole.Company.ExposeId),
		CompanyName:      string(authentic.User.CompanyRole.Company.Name),
		CompanyRoles:     &roleLabels,
		CompanyRoleNames: &roleNames,
	}
}

func getCompanyRole(claims GeezerClaims) (*CompanyRole, error) {
	if claims.CompanyExposeId != nil && claims.CompanyName != nil && claims.CompanyRoles != nil && claims.CompanyRoleNames != nil {
		if len(claims.CompanyRoles) != len(claims.CompanyRoleNames) {
			return nil, fmt.Error("CompanyRoles and CompanyRoleNames length is not equal")
		}

		var roles []Role = make([]Role, len(claims.CompanyRoles))
		for i := 0; i < len(claims.CompanyRoles); i++ {

			var label, err = text.NewLabel(claims.CompanyRoles[i])
			if err != nil {
				return nil, fmt.Error("CompanyRoles is not valid")
			}

			var name, err = text.NewName(claims.CompanyRoleNames[i])
			if err != nil {
				return nil, fmt.Error("CompanyRoleName is not valid")
			}

			var roles[i] = user.NewRole(label, name)
		}

		var companyExposeId, err = text.NewCompanyExposeId(*claims.CompanyExposeId)
		if err != nil {
			return nil, fmt.Error("CompanyExposeId is not valid")
		}
		var companyName, err = NewCompanyName(*claims.CompanyName)
		if err != nil {
			return nil, fmt.Error("CompanyName is not valid")
		}

		return user.NewCompany(companyExposeId, companyName, roles), nil
	} else {
		if claims.CompanyExposeId != nil {
			return nil, fmt.Error("CompanyExposeId is not nil")
		}
		if claims.CompanyName != nil {
			return nil, fmt.Error("CompanyName is not nil")
		}
		if claims.CompanyRoles != nil {
			return nil, fmt.Error("CompanyRole is not nil")
		}
		if claims.CompanyRoleNames != nil {
			return nil, fmt.Error("CompanyRoleName is not nil")
		}
		return nil, nil
	}
}

func (claims *GeezerClaims) ToAuthentic() (*user.Authentic, error) {
	var companyRole, err = getCompanyRole(claims)
	if err != nil {
		return nil, err
	}

	var userExposeId, err = text.NewUserExposeId(claims.Subject)
	if err != nil {
		return nil, fmt.Error("UserExposeId is not valid")
	}

	var userEmailId, err = text.NewEmail(claims.UserEmailId)
	if err != nil {
		return nil, fmt.Error("UserEmailId is not valid")
	}

	var userEmail = nil
	if claims.UserEmail != nil {
		var userEmail, err = text.NewEmail(claims.UserEmail)
		if err != nil {
			return nil, fmt.Error("UserEmail is not valid")
		}
	}

	var userName, err = text.NewName(claims.UserName)
	if err != nil {
		return nil, fmt.Error("UserName is not valid")
	}

	var user = user.NewUser(
		userExposeId,
		userEmailId,
		userEmail,
		userName,
		claims.BotFlag,
		companyRole,
		claims.UpdateDate.Time,
	)

	return user.NewAuthentic(
		claims.Issuer,
		claims.Subject,
		claims.Audience,
		claims.ExpiresAt.Time,
		claims.NotBefore.Time,
		claims.IssuedAt.Time,
		claims.ID,
		user,
	), nil
}
