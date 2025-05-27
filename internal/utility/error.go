package utility

import (
	"errors"
	pkg "github.com/motojouya/geezer_auth/pkg/model"
)

/*
 * RangeError
 */
type RangeError struct {
	Name string
	Value uint
	min uint
	max uint
	error
}

func NewRangeError(name string, value uint, min uint, max uint, message string) *LengthError {
	return &LengthError{
		Name:   name,
		Value:  value,
		min:    min,
		max:    max,
		error:  errors.New(message),
	}
}

func (e RangeError) Error() string {
	return e.error.Error() + " (name: " + e.Name + ", value: " + string(e.Value) + ", min: " + string(e.min) + ", max: " + string(e.max) + ")"
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
	UserExposeId pkg.ExposeId
	error
}

func NewAuthenticationError(userExposeId pkg.ExposeId, message string) *AuthenticationError {
	return &AuthenticationError{
		UserExposeId: userExposeId,
		error:        errors.New(message),
	}
}

func (e AuthenticationError) Error() string {
	return e.error.Error() + " (userExposeId: " + e.UserExposeId.String() + ")"
}

func (e AuthenticationError) Unwrap() error {
	return e.error
}

func (e AuthenticationError) HttpStatus() uint {
	return 401
}
