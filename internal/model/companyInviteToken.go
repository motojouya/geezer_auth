package model

import (
	"github.com/google/uuid"
)

// FIXME use type CompanyInviteToken string
type CompanyInviteToken = string

func GenerateCompanyInviteToken() (CompanyInviteToken, error) {
	token, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	return token.String(), nil
}
