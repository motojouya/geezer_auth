package jwt

import (
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/motojouya/geezer_auth/pkg/core/user"
	"github.com/motojouya/geezer_auth/pkg/utility"
	"strconv"
	"strings"
	"time"
)

// FIXME claimsのprivate keyが`github.com/motojouya/geezer_auth/`をprefixとしているが、本来は稼働するサーバのfqdnをprefixとして持つべき。
type GeezerClaims struct {
	gojwt.RegisteredClaims
	UserEmail         *string   `json:"email"`
	UserName          string    `json:"name"`
	UpdateDate        time.Time `json:"update_at"`
	UserEmailId       string    `json:"github.com/motojouya/geezer_auth/email_id"`
	BotFlag           bool      `json:"github.com/motojouya/geezer_auth/bot_flag"`
	CompanyIdentifier *string   `json:"github.com/motojouya/geezer_auth/company_expose_id"`
	CompanyName       *string   `json:"github.com/motojouya/geezer_auth/company_name"`
	CompanyRoles      []string  `json:"github.com/motojouya/geezer_auth/company_roles"`
	CompanyRoleNames  []string  `json:"github.com/motojouya/geezer_auth/company_role_names"`
}

func FromAuthentic(authentic *user.Authentic) *GeezerClaims {

	var companyIdentifier *string = nil
	var companyName *string = nil
	var companyRoles []string = nil
	var companyRoleNames []string = nil

	if authentic.User.CompanyRole != nil {
		var companyIdentifierValue = string(authentic.User.CompanyRole.Company.Identifier)
		companyIdentifier = &companyIdentifierValue
		var companyNameValue = string(authentic.User.CompanyRole.Company.Name)
		companyName = &companyNameValue
		companyRoles = make([]string, 0, len(authentic.User.CompanyRole.Roles))
		companyRoleNames = make([]string, 0, len(authentic.User.CompanyRole.Roles))

		for i, role := range authentic.User.CompanyRole.Roles {
			companyRoles[i] = string(role.Label)
			companyRoleNames[i] = string(role.Name)
		}
	}

	var userEmail *string = nil
	if authentic.User.Email != nil {
		var userEmailValue = string(*authentic.User.Email)
		userEmail = &userEmailValue
	}

	return &GeezerClaims{
		RegisteredClaims:  authentic.RegisteredClaims,
		UserEmail:         userEmail,
		UserName:          string(authentic.User.Name),
		UpdateDate:        authentic.User.UpdateDate,
		UserEmailId:       string(authentic.User.EmailId),
		BotFlag:           authentic.User.BotFlag,
		CompanyIdentifier: companyIdentifier,
		CompanyName:       companyName,
		CompanyRoles:      companyRoles,
		CompanyRoleNames:  companyRoleNames,
	}
}

func getCompanyRole(claims *GeezerClaims) (*user.CompanyRole, error) {

	if claims.CompanyIdentifier != nil && claims.CompanyName != nil && claims.CompanyRoles != nil && claims.CompanyRoleNames != nil {
		if len(claims.CompanyRoles) != len(claims.CompanyRoleNames) {
			return nil, NewJwtError("len(CompanyRoles)", "len(CompanyRoleNames)", "CompanyRoles and CompanyRoleNames length is not equal")
		}

		var roles = make([]user.Role, len(claims.CompanyRoles))
		for i := 0; i < len(claims.CompanyRoles); i++ {

			var label, labelErr = text.NewLabel(claims.CompanyRoles[i])
			if labelErr != nil {
				return nil, utility.AddPropertyError("Company.Role["+strconv.Itoa(i)+"]", labelErr)
			}

			var name, nameErr = text.NewName(claims.CompanyRoleNames[i])
			if nameErr != nil {
				return nil, utility.AddPropertyError("Company.Role["+strconv.Itoa(i)+"]", nameErr)
			}

			roles[i] = user.NewRole(label, name)
		}

		var companyIdentifier, idErr = text.NewIdentifier(*claims.CompanyIdentifier)
		if idErr != nil {
			return nil, utility.AddPropertyError("company", idErr)
		}
		var companyName, nameErr = text.NewName(*claims.CompanyName)
		if nameErr != nil {
			return nil, utility.AddPropertyError("company", nameErr)
		}

		var company = user.NewCompany(companyIdentifier, companyName)

		return user.NewCompanyRole(company, roles), nil

	} else {
		if claims.CompanyIdentifier != nil {
			return nil, NewJwtError("CompanyIdentifier", *claims.CompanyIdentifier, "CompanyIdentifier is not nil")
		}
		if claims.CompanyName != nil {
			return nil, NewJwtError("CompanyName", *claims.CompanyName, "CompanyIdentifier is not nil")
		}
		if claims.CompanyRoles != nil {
			return nil, NewJwtError("CompanyRoles", strings.Join(claims.CompanyRoles, ","), "CompanyIdentifier is not nil")
		}
		if claims.CompanyRoleNames != nil {
			return nil, NewJwtError("CompanyRoleNames", strings.Join(claims.CompanyRoleNames, ","), "CompanyIdentifier is not nil")
		}
		return nil, nil
	}
}

func (claims *GeezerClaims) ToAuthentic() (*user.Authentic, error) {

	var userIdentifier, idErr = text.NewIdentifier(claims.Subject)
	if idErr != nil {
		return nil, utility.AddPropertyError("claims", idErr)
	}

	var userEmailId, emailErr = text.NewEmail(claims.UserEmailId)
	if emailErr != nil {
		return nil, utility.AddPropertyError("claims", emailErr)
	}

	var userEmail *text.Email = nil
	if claims.UserEmail != nil {
		var userEmailValue, err = text.NewEmail(*claims.UserEmail)
		if err != nil {
			return nil, utility.AddPropertyError("claims", err)
		} else {
			userEmail = &userEmailValue
		}
	}

	var userName, nameErr = text.NewName(claims.UserName)
	if nameErr != nil {
		return nil, utility.AddPropertyError("claims", nameErr)
	}

	var companyRole, crErr = getCompanyRole(claims)
	if crErr != nil {
		return nil, utility.AddPropertyError("claims", crErr)
	}

	var userValue = user.NewUser(
		userIdentifier,
		userEmailId,
		userEmail,
		userName,
		claims.BotFlag,
		companyRole,
		claims.UpdateDate,
	)

	return user.NewAuthentic(
		claims.Issuer,
		claims.Subject,
		claims.Audience,
		claims.ExpiresAt.Time,
		claims.NotBefore.Time,
		claims.IssuedAt.Time,
		claims.ID,
		userValue,
	), nil
}
