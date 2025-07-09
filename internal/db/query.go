package db

import (
	"github.com/motojouya/geezer_auth/internal/db/query"
	"github.com/motojouya/geezer_auth/internal/db/transfer/company"
)

type Query interface {
	query.GetCompanyQuery
	query.GetCompanyInviteQuery
}

/** GetCompany */
func (orp ORPImpl) GetCompany(identifier string) (*company.Company, error) {
	return query.GetCompany(orp, identifier)
}

func (orp ORPTransactionImpl) GetCompany(identifier string) (*company.Company, error) {
	return query.GetCompany(orp, identifier)
}

/** GetCompanyInvite */
func (orp ORPImpl) GetCompanyInvite(identifier string, verifyToken string) (*company.CompanyInviteFull, error) {
	return query.GetCompanyInvite(orp, identifier, verifyToken)
}

func (orp ORPTransactionImpl) GetCompanyInvite(identifier string, verifyToken string) (*company.CompanyInviteFull, error) {
	return query.GetCompanyInvite(orp, identifier, verifyToken)
}
