package companyUser

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/company"
	"github.com/motojouya/geezer_auth/internal/shelter/text"
	"github.com/motojouya/geezer_auth/internal/entry/transfer/common"
)

type InviteTokenGetter interface {
	GetToken() (text.Token, error)
}

type CompanyAccept struct {
	Token string `json:"token"`
}

type CompanyUserAcceptRequest struct {
	common.RequestHeader
	company.CompanyGetRequest
	CompanyAccept
}

func (c CompanyUserAcceptRequest) GetToken() (text.Token, error) {
	return text.NewToken(c.Token)
}
