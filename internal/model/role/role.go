package role

import (
	"time"
	text "github.com/motojouya/geezer_auth/pkg/model/text"
	user "github.com/motojouya/geezer_auth/pkg/model/user"
)

/*
 * Roleは管理者が登録する想定なので、基本的には削除されない
 * また、一意な識別子はlabelであるため、RoleIdは必要ない
 * 他のサービスからも参照されるので、内部に閉じるRoleIdは意味がないため
 */
type Role struct {
	user.Role
	Description    Text
	RegisteredDate time.Time
}

func NewRole(name text.Name, label text.Label, description Text, registeredDate time.Time) Role {
	return Role{
		Role:           user.NewRole(label, name),
		Description:    description,
		RegisteredDate: registeredDate,
	}
}
