package text

import (
	"time"
	"strings"
)

/*
 * Nameは長さ1-255
 * 255に特に強い意味はない。dbのcharの最大長が255なので、255にしているだけ。
 */
type Name string

func NewName(name string) (Name, error) {
	if name == "" {
		return Name(""), NewLengthError("name", &name, 1, 255, "name should not be empty")
	}

	var trimmed = strings.TrimSpace(name)

	var length = len([]rune(trimmed))
	if length < 1 || length > 255 {
		return Name(""), NewLengthError("name", &name, 1, 255, "name must be between 1 and 255 characters")
	}

	return Name(trimmed), nil
}
