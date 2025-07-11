package company

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-gorp/gorp"
	transfer "github.com/motojouya/geezer_auth/internal/db/transfer/company"
	"github.com/motojouya/geezer_auth/internal/db/utility"
)

type GetCompanyInviteQuery interface {
	GetCompanyInvite(identifier string, verifyToken string) (*transfer.CompanyInviteFull, error)
}

func GetCompanyInvite(executer gorp.SqlExecutor, companyIdentifier string, verifyToken string) (*transfer.CompanyInviteFull, error) {
	var sql, args, sqlErr = transfer.SelectCompanyInvite.Where(
		goqu.C("c.identifier").Eq(companyIdentifier),
		goqu.C("ci.verify_token").Eq(verifyToken),
	).Prepared(true).ToSQL()
	if sqlErr != nil {
		return nil, sqlErr
	}

	var keys = map[string]string{
		"identifier":   companyIdentifier,
		"verify_token": verifyToken,
	}
	return utility.SelectSingle[transfer.CompanyInviteFull](executer, "company_invite", keys, sql, args...)
}
