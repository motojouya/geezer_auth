package text

import (
	"time"
	"strings"
	"regexp"
)

const (
	ExposeIdChar = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	ExposeIdLength = 6
)

/*
 * ExposeIdはlarge alphabet2文字と、large alphabet6文字をハイフンで繋いだもの。以下が例
 * `AB-ABCDEF`
 * 前2文字は、なんのExposeIdかを識別するための接頭語で、例としてCompanyならCP、UserならUSなど
 * 後ろ6文字はランダムな値
 *
 * 単に、ハイフンを付けずにprefixをつけてもよかったが、ランダム文字列が6文字でたりなくなって拡張する際に区別がつかなくなるのでハイフンを入れている
 */
type ExposeId string

func CreateExposeId(prefix string, randoms string) (ExposeId, error) {
	return NewExposeId(prefix + randoms)
}

func NewExposeId(exposeId string) (ExposeId, error) {
	var trimmed = strings.TrimSpace(exposeId)

	if trimmed == "" {
		return ExposeId(""), NewLengthError("exposeId", &exposeId, 9, 9, "exposeId should not be empty")
	}

	var length = len([]rune(trimmed))
	if length != 9 {
		return ExposeId(""), NewLengthError("exposeId", &exposeId, 9, 9, "exposeId must be exactly 9 characters")
	}

	// TODO 正規表現あってる？
	re, err := regexp.Compile(`^[A-Z_]{2}-[A-Z_]{6}$`)
	if err != nil {
		// 固定値なのでエラーにはならないはず
		panic(err)
	}

	var result = re.MatchString(text, -1)
	if !result {
		return ExposeId(""), NewFormatError("exposeId", "exposeId", &exposeId, "exposeId must be in the format of XX-XXXXXX where X is an uppercase letter")
	}

	return ExposeId(trimmed), nil
}
