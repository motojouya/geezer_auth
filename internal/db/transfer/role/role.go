package role

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-gorp/gorp"
	core "github.com/motojouya/geezer_auth/internal/core/role"
	coreText "github.com/motojouya/geezer_auth/internal/core/text"
	"github.com/motojouya/geezer_auth/internal/db/utility"
	pkgText "github.com/motojouya/geezer_auth/pkg/core/text"
	"time"
)

type Role struct {
	Label          string    `db:"label,primarykey"`
	Name           string    `db:"name"`
	Description    string    `db:"description"`
	RegisteredDate time.Time `db:"register_date"`
}

func AddRoleTable(dbMap *gorp.DbMap) {
	dbMap.AddTableWithName(Role{}, "role")
}

var SelectRole = utility.Dialect.From(goqu.T("role").As("r")).Select(
	goqu.I("r.label").As("label"),
	goqu.I("r.name").As("name"),
	goqu.I("r.description").As("description"),
	goqu.I("r.register_date").As("register_date"),
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

// testdata投入時に楽するためのもの。アプリケーションからは利用を想定しない。
func NewRole(label string, name string, description string, registeredDate time.Time) Role {
	return Role{
		Label:          label,
		Name:           name,
		Description:    description,
		RegisteredDate: registeredDate,
	}
}
