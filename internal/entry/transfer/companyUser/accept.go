package companyUser

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/company"
	"github.com/motojouya/geezer_auth/internal/shelter/text"
)

type InviteTokenGetter interface {
	GetToken() (text.Token, error)
}

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
