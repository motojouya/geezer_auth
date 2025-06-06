package jwt

import (
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/motojouya/geezer_auth/pkg/core/user"
	"github.com/motojouya/geezer_auth/pkg/utility"
	"time"
)

// FIXME claimsのprivate keyが`github.com/motojouya/geezer_auth/`をprefixとしているが、本来は稼働するサーバのfqdnをprefixとして持つべき。
type GeezerClaims struct {
	gojwt.RegisteredClaims
	UserEmail        *string   `json:"email"`
	UserName         string    `json:"name"`
	UpdateDate       time.Time `json:"update_at"`
	UserEmailId      string    `json:"github.com/motojouya/geezer_auth/email_id"`
	BotFlag          bool      `json:"github.com/motojouya/geezer_auth/bot_flag"`
	CompanyIdentifier  *string   `json:"github.com/motojouya/geezer_auth/company_expose_id"`
	CompanyName      *string   `json:"github.com/motojouya/geezer_auth/company_name"`
	CompanyRoles     []string `json:"github.com/motojouya/geezer_auth/company_roles"`
	CompanyRoleNames []string `json:"github.com/motojouya/geezer_auth/company_role_names"`
}

func FromAuthentic(authentic *user.Authentic) *GeezerClaims {

	var companyIdentifier *string = nil
	var companyName *string = nil
	var companyRoles []string = nil
	var companyRoleNames []string = nil

	if authentic.User.CompanyRole != nil {
		companyIdentifier = &string(authentic.User.CompanyRole.Company.Identifier)
		companyName = &string(authentic.User.CompanyRole.Company.Name)
		companyRoles = make([]string, 0, len(authentic.User.CompanyRole.Roles))
		companyRoleNames = make([]string, 0, len(authentic.User.CompanyRole.Roles))

		for i, role := range authentic.User.CompanyRole.Roles {
			companyRoles[i] = string(role.Label)
			companyRoleNames[i] = string(role.Name)
		}
	}

	var userEmail *string = nil
	if authentic.User.Email != nil {
		userEmail = &string(authentic.User.Email)
	}

	return &GeezerClaims{
		RegisteredClaims: authentic.RegisteredClaims,
		UserEmail:        userEmail,
		UserName:         string(authentic.User.Name),
		UpdateDate:       authentic.User.UpdateDate.Time,
		UserEmailId:      string(authentic.User.EmailId),
		BotFlag:          authentic.User.BotFlag,
		CompanyIdentifier:  companyIdentifier,
		CompanyName:      companyName,
		CompanyRoles:     companyRoles,
		CompanyRoleNames: companyRoleNames,
	}
}

func getCompanyRole(claims *GeezerClaims) (*CompanyRole, error) {

	if claims.CompanyIdentifier != nil && claims.CompanyName != nil && claims.CompanyRoles != nil && claims.CompanyRoleNames != nil {
		if len(claims.CompanyRoles) != len(claims.CompanyRoleNames) {
			return nil, NewJwtError("len(CompanyRoles)", "len(CompanyRoleNames)", "CompanyRoles and CompanyRoleNames length is not equal")
		}

		var roles []Role = make([]Role, len(claims.CompanyRoles))
		for i := 0; i < len(claims.CompanyRoles); i++ {

			var label, err = text.NewLabel(claims.CompanyRoles[i])
			if err != nil {
				return nil, utility.AddPropertyError("Company.Role[" + string(i) + "]", err)
			}

			var name, err = text.NewName(claims.CompanyRoleNames[i])
			if err != nil {
				return nil, utility.AddPropertyError("Company.Role[" + string(i) + "]", err)
			}

			var roles[i] = user.NewRole(label, name)
		}

		var companyIdentifier, err = text.NewCompanyIdentifier(*claims.CompanyIdentifier)
		if err != nil {
			return nil, utility.AddPropertyError("company", err)
		}
		var companyName, err = NewCompanyName(*claims.CompanyName)
		if err != nil {
			return nil, utility.AddPropertyError("company", err)
		}

		return user.NewCompany(companyIdentifier, companyName, roles), nil
	} else {
		if claims.CompanyIdentifier != nil {
			return nil, NewJwtError("CompanyIdentifier", claims.CompanyIdentifier, "CompanyIdentifier is not nil")
		}
		if claims.CompanyName != nil {
			return nil, NewJwtError("CompanyName", claims.CompanyName, "CompanyIdentifier is not nil")
		}
		if claims.CompanyRoles != nil {
			return nil, NewJwtError("CompanyRoles", claims.CompanyRoles, "CompanyIdentifier is not nil")
		}
		if claims.CompanyRoleNames != nil {
			return nil, NewJwtError("CompanyRoleNames", claims.CompanyRoleNames, "CompanyIdentifier is not nil")
		}
		return nil, nil
	}
}

func (claims *GeezerClaims) ToAuthentic() (*user.Authentic, error) {

	var userIdentifier, err = text.NewUserIdentifier(claims.Subject)
	if err != nil {
		return nil, utility.AddPropertyError("claims", err)
	}

	var userEmailId, err = text.NewEmail(claims.UserEmailId)
	if err != nil {
		return nil, utility.AddPropertyError("claims", err)
	}

	var userEmail = nil
	if claims.UserEmail != nil {
		var userEmail, err = text.NewEmail(claims.UserEmail)
		if err != nil {
			return nil, utility.AddPropertyError("claims", err)
		}
	}

	var userName, err = text.NewName(claims.UserName)
	if err != nil {
		return nil, utility.AddPropertyError("claims", err)
	}

	var companyRole, err = getCompanyRole(claims)
	if err != nil {
		return nil, utility.AddPropertyError("claims", err)
	}

	var user = user.NewUser(
		userIdentifier,
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
