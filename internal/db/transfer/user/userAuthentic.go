package user

import (
	text "github.com/motojouya/geezer_auth/pkg/core/text"
	core "github.com/motojouya/geezer_auth/internal/core/user"
	"time"
)

type UserAuthentic struct {
	UserPersistKey     uint
	UserIdentifier     string
	UserExposeEmailId  string
	UserName           string
	UserBotFlag        bool
	UserRegisteredDate time.Time
	UserUpdateDate     time.Time
	Email              *string
	UserCompanyRole    []*UserCompanyRoleFull
}

func (ua UserAuthentic) ToCoreUserAuthentic() (*core.UserAuthentic, error) {
	var user, userErr = (User{
		PersistKey:     ua.UserPersistKey,
		Identifier:     ua.UserIdentifier,
		ExposeEmailId:  ua.UserExposeEmailId,
		Name:           ua.UserName,
		BotFlag:        ua.UserBotFlag,
		RegisteredDate: ua.UserRegisteredDate,
		UpdateDate:     ua.UserUpdateDate,
	}).ToCoreUser()
	if userErr != nil {
		return &core.UserAuthentic{}, userErr
	}

	var email *text.Email = nil
	if ua.Email != nil {
		var emailResult, emailErr = text.NewEmail(*ua.Email)
		if emailErr != nil {
			return &core.UserAuthentic{}, emailErr
		}
		email = &emailResult
	}

	var coreUserCompanyRoles = make([]*core.UserCompanyRole, 0, len(ua.UserCompanyRole))
	for _, ucr := range ua.UserCompanyRole {
		var coreUserCompanyRole, companyRoleErr = ucr.ToCoreUserCompanyRole()
		if companyRoleErr != nil {
			return &core.UserAuthentic{}, companyRoleErr
		}
		coreUserCompanyRoles = append(coreUserCompanyRoles, coreUserCompanyRole)
	}

	var companyRole, companyRoleErr = core.ListToCompanyRole(user, coreUserCompanyRoles)
	if companyRoleErr != nil {
		return &core.UserAuthentic{}, companyRoleErr
	}

	return core.NewUserAuthentic(
		user,
		companyRole,
		email,
	), nil
}
