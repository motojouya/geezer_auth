package user

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/common"
	shelter "github.com/motojouya/geezer_auth/internal/shelter/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"time"
)

type UserApplyer interface {
	ApplyShelterUser(user shelter.User, updateDate time.Time) (shelter.User, error)
}

type UserChangeName struct {
	Name string `json:"name"`
}

type UserChangeNameRequest struct {
	common.RequestHeader
	UserChangeName
}

func (u UserChangeNameRequest) ApplyShelterUser(user shelter.User, updateDate time.Time) (shelter.User, error) {
	var name, nameErr = pkgText.NewName(u.Name)
	if nameErr != nil {
		return shelter.User{}, nameErr
	}

	return user.UpdateName(name, updateDate), nil
}
