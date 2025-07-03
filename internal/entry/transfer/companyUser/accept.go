package companyUser

import (
	"github.com/motojouya/geezer_auth/internal/core/text"
	"github.com/motojouya/geezer_auth/internal/entry/transfer/company"
)

type CompanyAccept struct {
	Token string `json:"token"`
}

type CompanyUserAcceptRequest struct {
	company.CompanyGetRequest
	CompanyAccept CompanyAccept `http:"body"`
}

func (c CompanyUserAcceptRequest) GetToken() (text.Token, error) {
	return text.NewToken(c.CompanyAccept.Token)
}
