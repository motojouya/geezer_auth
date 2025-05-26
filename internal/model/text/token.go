package model

import (
	"fmt"
)

/*
 * Tokenは、認証トークンを表す型
 * 長さの仕様は特に制限がないという前提にはなるのでText型と同等の形で定義する。
 */
type Token string

func NewToken(token *string) (*Token, error) {
	if token == nil {
		return nil, fmt.Errorf("token cannot be nil")
	}

	var length = len(*[]rune(token))
	if length > 10000 {
		return Token(""), fmt.Errorf("token must be less then 10000 characters")
	}

	return &Token(token), nil
}
