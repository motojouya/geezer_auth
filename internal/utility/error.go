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
 * RangeError
 */
type RangeError struct {
	Name string
	Value uint
	min *uint
	max *uint
	error
}

func NewRangeError(name string, value uint, min *uint, max *uint, message string) *LengthError {
	return &LengthError{
		Name:   name,
		Value:  value,
		min:    min,
		max:    max,
		error:  errors.New(message),
	}
}

func (e *RangeError) Error() string {
	return e.error.Error() + " (name: " + e.Name + ", value: " + string(e.Value) + ", min: " + string(e.min) + ", max: " + string(e.max) + ")"
}

func (e *RangeError) Unwrap() error {
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

/*
 * PropertyError
 * 
 * なんらかの不整合で生じたエラーは、特定の名前空間で処理されるため、その名前空間にたどりつくための経路を知ることができない
 * したがってその経路は呼び出し側で補填する必要があり、その機能を備えたエラー型
 */
type PropertyError struct {
	Property string
	error
}

func NewPropertyError(property string, source error) *ParameterError {
	return &ParameterError{
		Property: property,
		error:   source,
	}
}

func (e *PropertyError) Error() string {
	return e.error.Error() + " (property: " + e.Property + ")"
}

func (e *PropertyError) Unwrap() error {
	return e.error
}

func (e *PropertyError) Add(path string) *PropertyError {
	return &PropertyError{
		Property: path + "." + e.Property,
		error:    e.error,
	}
}

const propertyError *PropertyError = nil

func CreatePropertyError(path string, source error) *PropertyError {
	if source == nil {
		return panic("source error cannot be nil")
	}

	if errors.As(source, propertyError) {
		return source.(*PropertyError).Add(path)
	} else {
		return NewPropertyError(path, source)
	}
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

func (e *AuthenticationError) Error() string {
	return e.error.Error() + " (userExposeId: " + e.UserExposeId.String() + ")"
}

func (e *AuthenticationError) Unwrap() error {
	return e.error
}
