package user

import (
	shelter "github.com/motojouya/geezer_auth/internal/shelter/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"time"
)

type UserChangeName struct {
	Name string `json:"name"`
}

type UserChangeNameRequest struct {
	UserChangeName UserChangeName `http:"body"`
}

func (u UserChangeNameRequest) ToShelterUser(user shelter.User, updateDate time.Time) (shelter.User, error) {
	var name, nameErr = pkgText.NewName(u.UserChangeName.Name)
	if nameErr != nil {
		return shelter.User{}, nameErr
	}

	return user.UpdateName(name, updateDate), nil
}
