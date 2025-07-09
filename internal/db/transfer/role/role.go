package role

import (
	"github.com/doug-martin/goqu/v9"
	core "github.com/motojouya/geezer_auth/internal/core/role"
	coreText "github.com/motojouya/geezer_auth/internal/core/text"
	"github.com/motojouya/geezer_auth/internal/db/sql"
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
	"time"
)

type Role struct {
	Label          string    `db:"label"`
	Name           string    `db:"name"`
	Description    string    `db:"description"`
	RegisteredDate time.Time `db:"register_date"`
}

var SelectRole = sql.Dialect.From("role").As("r").Select(
	goqu.C("r.label").As("label"),
	goqu.C("r.name").As("name"),
	goqu.C("r.description").As("description"),
	goqu.C("r.register_date").As("register_date"),
)

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
