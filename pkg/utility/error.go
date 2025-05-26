package utility

import (
	"errors"
	pkg "github.com/motojouya/geezer_auth/pkg/model"
)

/*
 * NilError
 */
type NilError struct {
	Name  string
	error
}

func NewNilError(name string, message string) *NilError {
	return &NilError{
		Name:  name,
		error: errors.New(message),
	}
}

func (e *NilError) Error() string {
	return e.error.Error() + " (name: " + e.Name + ")"
}

func (e *NilError) Unwrap() error {
	return e.error
}

/*
 * SystemConfigError
 */
type SystemConfigError struct {
	Name string
	error
}

func NewSystemConfigError(name string, message string) *SystemConfigError {
	return &SystemConfigError{
		Name:  name,
		error: errors.New(message),
	}
}

func (e *SystemConfigError) Error() string {
	return e.error.Error() + " (name: " + e.Name + ")"
}

func (e *SystemConfigError) Unwrap() error {
	return e.error
}
