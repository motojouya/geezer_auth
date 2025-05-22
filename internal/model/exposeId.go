package model

import (
	"time"
	"fmt"
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
 */
type ExposeId string

func CreateExposeId(prefix string, randoms string) (ExposeId, error) {
	return NewExposeId(prefix + randoms)
}

func NewExposeId(exposeId string) (ExposeId, error) {
	if exposeId == "" {
		return ExposeId(""), fmt.Error("exposeId cannot be empty")
	}

	var length = len([]rune(exposeId))
	if length != 9 {
		return ExposeId(""), fmt.Errorf("randoms must be 9 characters")
	}

	// TODO 正規表現あってる？
	re, err := regexp.Compile(`^[A-Z_]{2}-[A-Z_]{6}$`)
	if err != nil {
		// 固定値なのでエラーにはならないはず
		panic(err)
	}

	var result = re.MatchString(text, -1)
	if !result {
		return ExposeId(""), fmt.Errorf("exposeId must contain only uppercase letters and underscores")
	}

	return ExposeId(exposeId), nil
}
