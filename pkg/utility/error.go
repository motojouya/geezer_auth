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

func (e NilError) Error() string {
	return e.error.Error() + " (name: " + e.Name + ")"
}

func (e NilError) Unwrap() error {
	return e.error
}

func (e NilError) HttpStatus() uint {
	return 400
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

func (e SystemConfigError) Error() string {
	return e.error.Error() + " (name: " + e.Name + ")"
}

func (e SystemConfigError) Unwrap() error {
	return e.error
}

func (e SystemConfigError) HttpStatus() uint {
	return 500
}

/*
 * PropertyError
 * 
 * なんらかの不整合で生じたエラーは、特定の名前空間で処理されるため、その名前空間にたどりつくための経路を知ることができない
 * したがってその経路は呼び出し側で補填する必要があり、その機能を備えたエラー型
 */
type PropertyError struct {
	Property string
	HttpStatus uint
	error
}

func CreatePropertyError(property string, source error) *PropertyError {
	return NewPropertyError(property, source.HttpStatus(), source)
}

func NewPropertyError(property string, httpStatus uint, source error) *PropertyError {
	return &PropertyError{
		Property:   property,
		HttpStatus: httpStatus,
		error:      source,
	}
}

func (e PropertyError) Error() string {
	return e.error.Error() + " (property: " + e.Property + " httpStatus: " + string(e.HttpStatus) + ")"
}

func (e PropertyError) Unwrap() error {
	return e.error
}

func (e PropertyError) HttpStatus() uint {
	if e.HttpStatus != 0 {
		return e.HttpStatus
	}
	return 400
}

func (e PropertyError) Add(path string) *PropertyError {
	return &PropertyError{
		Property:   path + "." + e.Property,
		HttpStatus: e.HttpStatus,
		error:      e.error,
	}
}

func (e PropertyError) Change(path string, httpStatus uint) *PropertyError {
	return &PropertyError{
		Property:   path + "." + e.Property,
		HttpStatus: httpStatus,
		error:      e.error,
	}
}

const propertyError *PropertyError = nil

func AddPropertyError(path string, source error) *PropertyError {
	if source == nil {
		return panic("source error cannot be nil")
	}

	if errors.As(source, propertyError) {
		return source.(*PropertyError).Add(path)
	} else {
		return CreatePropertyError(path, source)
	}
}

func ChangePropertyError(path string, source error, httpStatus uint) *PropertyError {
	if source == nil {
		return panic("source error cannot be nil")
	}

	if errors.As(source, propertyError) {
		return source.(*PropertyError).Change(path, httpStatus)
	} else {
		return NewPropertyError(path, httpStatus, source)
	}
}
