package model

import (
	"time"
	"fmt"
)

const (
	ExposeIdChar = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	ExposeIdLength = 6
)

type ExposeId string

func NewExposeId(prefix string, randoms string) (ExposeId, error) {
	var length = len([]rune(randoms))
	if length != 6 {
		return ExposeId(""), fmt.Errorf("randoms must be 6 characters")
	}

	// TODO 本当はExposeIdCharのみであることを保証したい

	return ExposeId(prefix + randoms), nil
}
