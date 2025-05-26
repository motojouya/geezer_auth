package model

import (
	"time"
	"fmt"
	"strings"
	"regexp"
)

type Email string

func NewEmail(email string) (Email, error) {
	if email == "" {
		return Email(""), fmt.Error("email cannot be empty")
	}

	var trimmed = strings.TrimSpace(email)

	var length = len([]rune(trimmed))
	if length < 3 || length > 320 {
		return Email(""), fmt.Errorf("email must be between 3 and 320 characters")
	}

	// TODO 正規表現あってる？
	re, err := regexp.Compile(`/^[a-z\d][\w.-]*@[\w.-]+\.[a-z\d]+$/i`)
	if err != nil {
		// 固定値なのでエラーにはならないはず
		panic(err)
	}

	var result = re.MatchString(text, -1)
	if !result {
		return Email(""), fmt.Errorf("email must contain only uppercase letters and underscores")
	}

	return Email(trimmed), nil
}
