package utility

import (
	"errors"
)

/*
 * RangeError
 */
type RangeError struct {
	Name string
	Value uint
	Min uint
	Max uint
	error
}

func NewRangeError(name string, value uint, min uint, max uint, message string) *LengthError {
	return &LengthError{
		Name:   name,
		Value:  value,
		Min:    min,
		Max:    max,
		error:  errors.New(message),
	}
}

func (e RangeError) Error() string {
	return e.error.Error() + " (name: " + e.Name + ", value: " + string(e.Value) + ", min: " + string(e.Min) + ", max: " + string(e.Max) + ")"
}

func (e RangeError) Unwrap() error {
	return e.error
}

func (e RangeError) HttpStatus() uint {
	return 400
}

/*
 * AuthenticationError
 */
type AuthenticationError struct {
	UserIdentifier string
	error
}

func NewAuthenticationError(userIdentifier string, message string) *AuthenticationError {
	return &AuthenticationError{
		UserIdentifier: userIdentifier,
		error:        errors.New(message),
	}
}

func (e AuthenticationError) Error() string {
	return e.error.Error() + " (userIdentifier: " + e.UserIdentifier + ")"
}

func (e AuthenticationError) Unwrap() error {
	return e.error
}

func (e AuthenticationError) HttpStatus() uint {
	return 401
}
