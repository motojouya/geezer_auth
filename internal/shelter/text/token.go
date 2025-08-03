package text

import (
	"github.com/google/uuid"
	pkg "github.com/motojouya/geezer_auth/pkg/core/text"
)

/*
 * Tokenは、認証トークンを表す型
 * 長さの仕様は特に制限がないという前提にはなるのでText型と同等の形で定義する。
 */
type Token string

func NewToken(token string) (Token, error) {
	var length = len([]rune(token))
	if length > 10000 {
		return Token(""), pkg.NewLengthError("token", token, 1, 10000, "token must be between 1 and 10000 characters")
	}

	return Token(token), nil
}

func CreateToken(token uuid.UUID) (Token, error) {
	var tokenStr = token.String()
	return NewToken(tokenStr)
}
