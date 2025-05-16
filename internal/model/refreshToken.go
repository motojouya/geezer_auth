package model

import (
	"github.com/google/uuid"
)

// FIXME use type RefreshToken string
type RefreshToken = string

func GenerateRefreshToken() (RefreshToken, error) {
	token, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	return token.String(), nil
}
