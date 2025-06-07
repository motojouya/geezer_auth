package text

import (
	"errors"
)

/*
 * LengthError
 */
type LengthError struct {
	Name  string
	Value string
	min   uint
	max   uint
	error
}

func NewLengthError(name string, value string, min uint, max uint, message string) *LengthError {
	return &LengthError{
		Name:  name,
		Value: value,
		min:   min,
		max:   max,
		error: errors.New(message),
	}
}

func (e LengthError) Error() string {
	return e.error.Error() + " (name: " + e.Name + ", value: " + e.Value + ", min: " + string(e.min) + ", max: " + string(e.max) + ")"
}

func (e LengthError) Unwrap() error {
	return e.error
}

func (e LengthError) HttpStatus() uint {
	return 400
}

/*
 * CharacterError
 */
type CharacterError struct {
	Name  string
	Chars string
	Value string
	error
}

func NewCharacterError(name string, chars string, value string, message string) *CharacterError {
	return &CharacterError{
		Name:  name,
		Chars: chars,
		Value: value,
		error: errors.New(message),
	}
}

func (e CharacterError) Error() string {
	return e.error.Error() + " (name: " + e.Name + ", chars: " + e.Chars + ", value: " + *e.Value + ")"
}

func (e CharacterError) Unwrap() error {
	return e.error
}

func (e CharacterError) HttpStatus() uint {
	return 400
}

/*
 * FormatError
 */
type FormatError struct {
	Name   string
	Format string
	Value  string
	error
}

func NewFormatError(name string, format string, value string, message string) *FormatError {
	return &FormatError{
		Name:   name,
		Format: format,
		Value:  value,
		error:  errors.New(message),
	}
}

func (e FormatError) Error() string {
	return e.error.Error() + " (name: " + e.Name + ", format: " + e.Format + ", value: " + *e.Value + ")"
}

func (e FormatError) Unwrap() error {
	return e.error
}

func (e FormatError) HttpStatus() uint {
	return 400
}
