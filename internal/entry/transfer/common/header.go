package common

import (
	"errors"
	"strings"
)

var tokenPrefix = "Bearer "

type RequestHeader struct {
	Token string `header:"Authorization"`
}

func (r *RequestHeader) GetBearerToken() (string, error) {
	if r.Token == "" {
		return "", nil
	}

	if !strings.HasPrefix(r.Token, tokenPrefix) {
		return "", errors.New("invalid token") // TODO 独自エラーにする
	}

	return strings.TrimPrefix(r.Token, tokenPrefix), nil
}

type BearerTokenGetter struct {
	GetBearerToken func() (string, error)
}
