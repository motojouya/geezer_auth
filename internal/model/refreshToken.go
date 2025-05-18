package model

import (
	"github.com/google/uuid"
)

// FIXME use type RefreshToken string
type RefreshToken = string

// TODO uuidの発行は、副作用があるので、このまま直接呼び出す形にすると、procedureのテストが難しい。
// company invite tokenもuuidを使うので、この手のtoken発行はなんかまとめてもいいかも
func GenerateRefreshToken() (RefreshToken, error) {
	token, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	return token.String(), nil
}
