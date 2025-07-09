package query

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-gorp/gorp"
	"github.com/motojouya/geezer_auth/internal/db/transfer/company"
)

type GetCompanyInviteQuery interface {
	GetCompanyInvite(identifier string, verifyToken string) (*company.CompanyInviteFull, error)
}

func GetCompanyInvite(executer gorp.SqlExecutor, companyIdentifier string, verifyToken string) (*company.CompanyInviteFull, error) {
	var sql, args, sqlErr = company.SelectCompanyInvite.Where(
		goqu.C("c.identifier").Eq(companyIdentifier),
		goqu.C("ci.verify_token").Eq(verifyToken),
	).Prepared(true).ToSQL()
	if sqlErr != nil {
		return nil, sqlErr
	}

	var companyInvite company.CompanyInviteFull
	var execErr = executer.SelectOne(&companyInvite, sql, args...)
	if execErr != nil {
		return nil, execErr
	}

	return &companyInvite, nil
}
