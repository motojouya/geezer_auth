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
	var sql, args, sqlErr = transfer.SelectCompany.Where(goqu.I("c.identifier").Eq(identifier)).Prepared(true).ToSQL()
	if sqlErr != nil {
		return nil, sqlErr
	}

	return utility.SelectSingle[transfer.Company](executer, "company", map[string]string{"identifier": identifier}, sql, args...)
}
