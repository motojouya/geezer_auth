package company

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-gorp/gorp"
	core "github.com/motojouya/geezer_auth/internal/core/company"
	"github.com/motojouya/geezer_auth/internal/db/utility"
	text "github.com/motojouya/geezer_auth/pkg/core/text"
	"time"
)

type Company struct {
	PersistKey     uint      `db:"persist_key,primarykey,autoincrement"`
	Identifier     string    `db:"identifier"`
	Name           string    `db:"name"`
	RegisteredDate time.Time `db:"register_date"`
}

func AddCompanyTable(dbMap *gorp.DbMap) {
	dbMap.AddTableWithName(Company{}, "company").SetKeys(true, "PersistKey")
}

var SelectCompany = utility.Dialect.From(goqu.T("company").As("c")).Select(
	goqu.I("c.persist_key").As("persist_key"),
	goqu.I("c.identifier").As("identifier"),
	goqu.I("c.name").As("name"),
	goqu.I("c.register_date").As("register_date"),
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

// testdata投入時に楽するためのもの。アプリケーションからは利用を想定しない。
func NewCompany(persistKey uint, identifier string, name string, registerDate time.Time) Company {
	return Company{
		PersistKey:     persistKey,
		Identifier:     identifier,
		Name:           name,
		RegisteredDate: registerDate,
	}
}
