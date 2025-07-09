package query

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-gorp/gorp"
	"github.com/motojouya/geezer_auth/internal/db/transfer/company"
	"github.com/motojouya/geezer_auth/internal/db"
)

type GetCompany interface {
	GetCompany(identifier string) (*company.Company, error)
}

func GetCompany(executer gorp.SqlExecutor, identifier string) (*company.Company, error) {
	var sql, args, sqlErr = company.SelectCompany.Where(goqu.C("c.identifier").Eq(identifier)).Prepared(true).ToSQL()
	if sqlErr != nil {
		return nil, sqlErr
	}

	var company company.Company
	var execErr = executer.SelectOne(&company, sql, args...)
	if execErr != nil {
		return nil, execErr
	}

	return &company, nil
}
