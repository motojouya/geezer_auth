package db

import (
	"errors"
)

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
	Table  string
	keys map[string]string
	error
}

func NewNotFoundError(table string, keys map[string]string, message string) *NotFoundError {
	return &NotFoundError{
		Table: table,
		keys:  keys,
		error: errors.New(message),
	}
}

func (e NotFoundError) Error() string {
	return e.error.Error() + " (table: " + e.Table + ", keys: " + formatKeys(e.keys) + ")"
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
	keys map[string]string
	error
}

func NewDuplicateError(table string, keys map[string]string, message string) *DuplicateError {
	return &DuplicateError{
		Table: table,
		keys:  keys,
		error: errors.New(message),
	}
}

func (e DuplicateError) Error() string {
	return e.error.Error() + " (table: " + e.Table + ", keys: " + formatKeys(e.keys) + ")"
}

func (e DuplicateError) Unwrap() error {
	return e.error
}

func (e DuplicateError) HttpStatus() uint {
	return 400
}
