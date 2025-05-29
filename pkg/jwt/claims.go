package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/motojouya/geezer_auth/pkg/model/text"
	"github.com/motojouya/geezer_auth/pkg/model/user"
	"github.com/motojouya/geezer_auth/pkg/utility"
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
	CompanyRoles     []string `json:"github.com/motojouya/geezer_auth/company_roles"`
	CompanyRoleNames []string `json:"github.com/motojouya/geezer_auth/company_role_names"`
}

func FromAuthentic(authentic *user.Authentic) *GeezerClaims {

	var companyExposeId *string = nil
	var companyName *string = nil
	var companyRoles []string = nil
	var companyRoleNames []string = nil

	if authentic.User.CompanyRole != nil {
		companyExposeId = &string(authentic.User.CompanyRole.Company.ExposeId)
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
		CompanyExposeId:  companyExposeId,
		CompanyName:      companyName,
		CompanyRoles:     companyRoles,
		CompanyRoleNames: companyRoleNames,
	}
}

func getCompanyRole(claims *GeezerClaims) (*CompanyRole, error) {

	if claims.CompanyExposeId != nil && claims.CompanyName != nil && claims.CompanyRoles != nil && claims.CompanyRoleNames != nil {
		if len(claims.CompanyRoles) != len(claims.CompanyRoleNames) {
			return nil, NewJwtError("len(CompanyRoles)", "len(CompanyRoleNames)", "CompanyRoles and CompanyRoleNames length is not equal")
		}

		var roles []Role = make([]Role, len(claims.CompanyRoles))
		for i := 0; i < len(claims.CompanyRoles); i++ {

			var label, err = text.NewLabel(claims.CompanyRoles[i])
			if err != nil {
				return nil, utility.CreatePropertyError("Company.Role[" + string(i) + "]", err)
			}

			var name, err = text.NewName(claims.CompanyRoleNames[i])
			if err != nil {
				return nil, utility.CreatePropertyError("Company.Role[" + string(i) + "]", err)
			}

			var roles[i] = user.NewRole(label, name)
		}

		var companyExposeId, err = text.NewCompanyExposeId(*claims.CompanyExposeId)
		if err != nil {
			return nil, utility.CreatePropertyError("company", err)
		}
		var companyName, err = NewCompanyName(*claims.CompanyName)
		if err != nil {
			return nil, utility.CreatePropertyError("company", err)
		}

		return user.NewCompany(companyExposeId, companyName, roles), nil
	} else {
		if claims.CompanyExposeId != nil {
			return nil, NewJwtError("CompanyExposeId", claims.CompanyExposeId, "CompanyExposeId is not nil")
		}
		if claims.CompanyName != nil {
			return nil, NewJwtError("CompanyName", claims.CompanyName, "CompanyExposeId is not nil")
		}
		if claims.CompanyRoles != nil {
			return nil, NewJwtError("CompanyRoles", claims.CompanyRoles, "CompanyExposeId is not nil")
		}
		if claims.CompanyRoleNames != nil {
			return nil, NewJwtError("CompanyRoleNames", claims.CompanyRoleNames, "CompanyExposeId is not nil")
		}
		return nil, nil
	}
}

func (claims *GeezerClaims) ToAuthentic() (*user.Authentic, error) {

	var userExposeId, err = text.NewUserExposeId(claims.Subject)
	if err != nil {
		return nil, utility.CreatePropertyError("claims", err)
	}

	var userEmailId, err = text.NewEmail(claims.UserEmailId)
	if err != nil {
		return nil, utility.CreatePropertyError("claims", err)
	}

	var userEmail = nil
	if claims.UserEmail != nil {
		var userEmail, err = text.NewEmail(claims.UserEmail)
		if err != nil {
			return nil, utility.CreatePropertyError("claims", err)
		}
	}

	var userName, err = text.NewName(claims.UserName)
	if err != nil {
		return nil, utility.CreatePropertyError("claims", err)
	}

	var companyRole, err = getCompanyRole(claims)
	if err != nil {
		return nil, utility.CreatePropertyError("claims", err)
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
