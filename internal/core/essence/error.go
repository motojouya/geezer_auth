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

func formatKeys(keys map[string]string) string {
	if len(keys) == 0 {
		return "{}"
	}
	result := "{"
	for k, v := range keys {
		result += k + ": " + v + ", "
	}
	return result + "}"
}

/*
 * NotFoundError
 */
type NotFoundError struct {
	Table string
	Keys  map[string]string
	error
}

func NewNotFoundError(table string, keys map[string]string, message string) *NotFoundError {
	return &NotFoundError{
		Table: table,
		Keys:  keys,
		error: errors.New(message),
	}
}

func (e NotFoundError) Error() string {
	return e.error.Error() + ", table: " + e.Table + ", keys: " + formatKeys(e.Keys)
}

func (e NotFoundError) Unwrap() error {
	return e.error
}

func (e NotFoundError) HttpStatus() uint {
	return 400
}

/*
 * DuplicateError
 */
type DuplicateError struct {
	Table string
	Keys  map[string]string
	error
}

func NewDuplicateError(table string, keys map[string]string, message string) *DuplicateError {
	return &DuplicateError{
		Table: table,
		Keys:  keys,
		error: errors.New(message),
	}
}

func (e DuplicateError) Error() string {
	return e.error.Error() + ", table: " + e.Table + ", keys: " + formatKeys(e.Keys)
}

func (e DuplicateError) Unwrap() error {
	return e.error
}

func (e DuplicateError) HttpStatus() uint {
	return 400
}
