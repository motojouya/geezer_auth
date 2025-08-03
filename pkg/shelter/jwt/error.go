package jwt

import (
	"errors"
)

/*
 * JwtError
 */
type JwtError struct {
	Claim string
	Value string
	error
}

func NewJwtError(claim string, value string, message string) *JwtError {
	return &JwtError{
		Claim: claim,
		Value: value,
		error: errors.New(message),
	}
}

func (e JwtError) Error() string {
	return e.error.Error() + ", claim: " + e.Claim + ", value: " + e.Value
}

func (e JwtError) Unwrap() error {
	return e.error
}

func (e JwtError) HttpStatus() uint {
	return 400
}
