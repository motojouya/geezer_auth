package model

import (
	"time"
	"fmt"
)

// Nameは長さ1-255
// 255に特に強い意味はない。dbのcharの最大長が255なので、255にしているだけ。
type Name string

func NewName(name string) (Name, error) {
	if name == "" {
		return Name(""), fmt.Error("name cannot be empty")
	}

	var length = len([]rune(name))
	if length < 1 || length > 255 {
		return Name(""), fmt.Errorf("name must be between 1 and 255 characters")
	}

	return Name(name), nil
}
