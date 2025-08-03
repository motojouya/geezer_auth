package text

import (
	"regexp"
	"strings"
)

type Label string

func NewLabel(label string) (Label, error) {
	var trimmed = strings.TrimSpace(label)

	if trimmed == "" {
		return Label(""), NewLengthError("label", label, 2, 255, "label should not be empty")
	}

	var length = len([]rune(trimmed))
	if length < 2 || length > 255 {
		return Label(""), NewLengthError("label", label, 2, 255, "label must be between 2 and 255 characters")
	}

	// TODO 正規表現あってる？
	re, err := regexp.Compile(`^[A-Z]{1}[A-Z_]{0,253}[A-Z]{1}$`)
	if err != nil {
		// 固定値なのでエラーにはならないはず
		panic(err)
	}

	var result = re.MatchString(trimmed)
	if !result {
		return Label(""), NewFormatError("label", "label", label, "label must start and end with an uppercase letter and can contain undersshelters in between")
	}

	return Label(trimmed), nil
}
