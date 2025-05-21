package model

import (
	"time"
	"fmt"
)

type Email string

func NewEmail(email string) (Email, error) {
	if email == "" {
		return Email(""), fmt.Error("email cannot be empty")
	}
	// TODO: validate email format
	// 面倒なので未実装

	return Email(email), nil
}
