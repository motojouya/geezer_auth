package role

import (
	text "github.com/motojouya/geezer_auth/internal/shelter/text"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	user "github.com/motojouya/geezer_auth/pkg/shelter/user"
	"time"
)

/*
 * Roleは管理者が登録する想定なので、基本的には削除されない
 * また、一意な識別子はlabelであるため、RoleIdは必要ない
 * 他のサービスからも参照されるので、内部に閉じるRoleIdは意味がないため
 */
type Role struct {
	user.Role
	Description    text.Text
	RegisteredDate time.Time
}

func NewRole(
	name pkgText.Name,
	label pkgText.Label,
	description text.Text,
	registeredDate time.Time,
) Role {
	return Role{
		Role:           user.NewRole(label, name),
		Description:    description,
		RegisteredDate: registeredDate,
	}
}

var RoleAdminLabel = pkgText.Label("MANAGER")
