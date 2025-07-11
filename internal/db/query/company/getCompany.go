package company

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-gorp/gorp"
	transfer "github.com/motojouya/geezer_auth/internal/db/transfer/company"
	"github.com/motojouya/geezer_auth/internal/db/utility"
)

type GetCompanyQuery interface {
	GetCompany(identifier string) (*transfer.Company, error)
}

func GetCompany(executer gorp.SqlExecutor, identifier string) (*transfer.Company, error) {
	var sql, args, sqlErr = transfer.SelectCompany.Where(goqu.C("c.identifier").Eq(identifier)).Prepared(true).ToSQL()
	if sqlErr != nil {
		return nil, sqlErr
	}

	var company, execErr = utility.SelectSingle[transfer.Company](executer, "company", map[string]string{"identifier":identifier}, sql, args...)
	if execErr != nil {
		return nil, execErr
	}

	return company, nil
}
