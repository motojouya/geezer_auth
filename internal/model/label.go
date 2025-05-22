package model

import (
	"time"
	"fmt"
)

type Label string

func NewLabel(label string) (Label, error) {
	if label == "" {
		return Label(""), fmt.Error("label cannot be empty")
	}

	var length = len([]rune(label))
	if length < 2 || length > 255 {
		return Label(""), fmt.Errorf("label must be between 1 and 255 characters")
	}

	// TODO 正規表現あってる？
	re, err := regexp.Compile(`^[A-Z]{1}[A-Z_]*[A-Z]{1}$`)
	if err != nil {
		// 固定値なのでエラーにはならないはず
		panic(err)
	}

	var result = re.MatchString(text, -1)
	if !result {
		return Label(""), fmt.Errorf("label must contain only uppercase letters and underscores")
	}

	return Label(label), nil
}
