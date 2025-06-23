package role

import (
	coreText "github.com/motojouya/geezer_auth/internal/core/text"
	core "github.com/motojouya/geezer_auth/internal/core/role"
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
	"time"
)

type Role struct {
	Label       string
	Name        string
	Description string
	RegisteredDate time.Time
}

func FromCoreRole(role core.Role) Role {
	return Role{
		Label:          string(role.Label),
		Name:           string(role.Name),
		Description:    string(role.Description),
		RegisteredDate: role.RegisteredDate,
	}
}

func (r Role) ToCoreRole() (core.Role, error) {
	var label, labelErr = pkgText.NewLabel(r.Label)
	if labelErr != nil {
		return core.Role{}, labelErr
	}

	var name, nameErr = pkgText.NewName(r.Name)
	if nameErr != nil {
		return core.Role{}, nameErr
	}

	var description, descErr = coreText.NewText(r.Description)
	if descErr != nil {
		return core.Role{}, descErr
	}

	return core.NewRole(
		name,
		label,
		description,
		r.RegisteredDate,
	), nil
}
