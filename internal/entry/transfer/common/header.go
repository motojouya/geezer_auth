package common

import (
	"errors"
	"strings"
)

var tokenPrefix = "Bearer "

type RequestHeader struct {
	token string `header:"Authorization"`
}

func (r *RequestHeader) GetBearerToken() (string, error) {
	if !strings.HasPrefix(r.token, tokenPrefix) {
		return "", errors.New("invalid token") // TODO 独自エラーにする
	}
	return strings.TrimPrefix(r.token, tokenPrefix), nil
}

type BearerTokenGetter struct {
	GetBearerToken func() (string, error)
}
