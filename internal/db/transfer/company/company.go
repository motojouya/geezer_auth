package company

import (
	"github.com/doug-martin/goqu/v9"
	core "github.com/motojouya/geezer_auth/internal/core/company"
	"github.com/motojouya/geezer_auth/internal/db/sql"
	text "github.com/motojouya/geezer_auth/pkg/core/text"
	"time"
)

type Company struct {
	PersistKey     uint      `db:"persist_key"`
	Identifier     string    `db:"identifier"`
	Name           string    `db:"name"`
	RegisteredDate time.Time `db:"register_date"`
}

var SelectCompany = sql.Dialect.From("company").As("c").Select(
	goqu.C("c.persist_key").As("persist_key"),
	goqu.C("c.identifier").As("identifier"),
	goqu.C("c.name").As("name"),
	goqu.C("c.register_date").As("register_date"),
)

func FromCoreCompany(coreCompany core.UnsavedCompany) Company {
	// PersistKey is zero value
	return Company{
		Identifier:     string(coreCompany.Identifier),
		Name:           string(coreCompany.Name),
		RegisteredDate: coreCompany.RegisteredDate,
	}
}

func (c Company) ToCoreCompany() (core.Company, error) {
	var identifier, idErr = text.NewIdentifier(c.Identifier)
	if idErr != nil {
		return core.Company{}, idErr
	}

	var name, nameErr = text.NewName(c.Name)
	if nameErr != nil {
		return core.Company{}, nameErr
	}

	return core.NewCompany(
		c.PersistKey,
		identifier,
		name,
		c.RegisteredDate,
	), nil
}
