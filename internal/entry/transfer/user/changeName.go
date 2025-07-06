package user

import (
	core "github.com/motojouya/geezer_auth/internal/core/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
	"time"
)

type UserChangeName struct {
	Name string `json:"name"`
}

type UserChangeNameRequest struct {
	UserChangeName UserChangeName `http:"body"`
}

func (u UserChangeNameRequest) ToCoreUser(user core.User, updateDate time.Time) (core.User, error) {
	var name, nameErr = pkgText.NewName(u.UserChangeName.Name)
	if nameErr != nil {
		return core.User{}, nameErr
	}

	return user.UpdateName(name, updateDate), nil
}
