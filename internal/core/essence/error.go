package essence

import (
	"errors"
	"strconv"
)

/*
 * InvalidArgumentError
 */
type InvalidArgumentError struct {
	Name  string
	Value string
	error
}

func NewInvalidArgumentError(name string, value string, message string) InvalidArgumentError {
	return InvalidArgumentError{
		Name:  name,
		Value: value,
		error: errors.New(message),
	}
}

func (e InvalidArgumentError) Error() string {
	return e.error.Error() + ", name: " + e.Name + ", value: " + e.Value
}

func (e InvalidArgumentError) Unwrap() error {
	return e.error
}

func (e InvalidArgumentError) HttpStatus() uint {
	return 400
}

/*
 * RangeError
 */
type RangeError struct {
	Name  string
	Value int
	Min   int
	Max   int
	error
}

func NewRangeError(name string, value int, min int, max int, message string) RangeError {
	return RangeError{
		Name:  name,
		Value: value,
		Min:   min,
		Max:   max,
		error: errors.New(message),
	}
}

func (e RangeError) Error() string {
	return e.error.Error() + ", name: " + e.Name + ", value: " + strconv.Itoa(e.Value) + ", min: " + strconv.Itoa(e.Min) + ", max: " + strconv.Itoa(e.Max)
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
		error:          errors.New(message),
	}
}

func (e AuthenticationError) Error() string {
	return e.error.Error() + ", userIdentifier: " + e.UserIdentifier
}

func (e AuthenticationError) Unwrap() error {
	return e.error
}

func (e AuthenticationError) HttpStatus() uint {
	return 401
}
