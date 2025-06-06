package text

import (
	"time"
	"strings"
	"regexp"
)

const (
	IdentifierChar = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	IdentifierLength = 6
)

/*
 * Identifierはlarge alphabet2文字と、large alphabet6文字をハイフンで繋いだもの。以下が例
 * `AB-ABCDEF`
 * 前2文字は、なんのIdentifierかを識別するための接頭語で、例としてCompanyならCP、UserならUSなど
 * 後ろ6文字はランダムな値
 *
 * 単に、ハイフンを付けずにprefixをつけてもよかったが、ランダム文字列が6文字でたりなくなって拡張する際に区別がつかなくなるのでハイフンを入れている
 */
type Identifier string

func CreateIdentifier(prefix string, randoms string) (Identifier, error) {
	return NewIdentifier(prefix + randoms)
}

func NewIdentifier(identifier string) (Identifier, error) {
	var trimmed = strings.TrimSpace(identifier)

	if trimmed == "" {
		return Identifier(""), NewLengthError("identifier", identifier, 9, 9, "exposeId should not be empty")
	}

	var length = len([]rune(trimmed))
	if length != 9 {
		return Identifier(""), NewLengthError("identifier", identifier, 9, 9, "exposeId must be exactly 9 characters")
	}

	// TODO 正規表現あってる？
	re, err := regexp.Compile(`^[A-Z_]{2}-[A-Z_]{6}$`)
	if err != nil {
		// 固定値なのでエラーにはならないはず
		panic(err)
	}

	var result = re.MatchString(text, -1)
	if !result {
		return Identifier(""), NewFormatError("identifier", "identifier", exposeId, "exposeId must be in the format of XX-XXXXXX where X is an uppercase letter")
	}

	return Identifier(trimmed), nil
}
