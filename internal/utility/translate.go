package utility

import (
	"slices"
	"reflect"
)

func Translate[F any, T any](from F): (T, error) {

	var fromValue = reflect.ValueOf(from)
	var toType = reflect.TypeOf(T{})

	var toValue = translateInner(fromValue, toType)
	if value, ok := toValue.Interface().(T); ok {
		return value, nil
	} else {
		return T{}, NewConvertError("", "Cannot Cast from "+fromValue.Type()+" to "+toType.Kind())
	}
}

func translateInner(fromValue reflect.Value, toType reflect.Type) (reflect.Value, error) {

	var toElementType = toType
	var isToPtr = false
	var kindTo = toType.Kind()
	if kindTo == reflect.Ptr {
		toElementType = toType.Elem()
		kindTo = toElementType.Kind()
		isToPtr = true
	}

	var fromElementValue = fromValue
	var isFromPtr = false
	var kindFrom = fromValue.Kind()
	if kindFrom == reflect.Ptr {
		fromElementValue = fromValue.Elem()
		kindFrom = fromElementValue.Kind()
		isFromPtr = true
	}

	switch kindTo {
	case reflect.Invalid, reflect.Uintptr, reflect.Complex64, reflect.Complex128, reflect.Chan, reflect.Interface, reflect.Func, reflect.Ptr, reflect.UnsafePointer, reflect.Map, reflect.Array:
		retrun T{}, NewConvertError("", kindTo+" is not supported.")

	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.String:
		// FIXME 文字は[]byteなので、intはstringに変換できたりするので、そのチェックが必要
		if fromElementValue.CanConvert(toElementType) {
			return fromElementValue.Convert(toElementType), nil
		}
		return T{}, NewConvertEror("", "Cannot convert from "+kindFrom+" to "+kindTo)

	case reflect.Struct:
		if kindFrom != reflect.Struct {
			return T{}, NewConvertError("", "Cannot convert from "+kindFrom+" to "+kindTo)
		}
		var toValue = reflect.ValueOf(T{})
		for i := 0; i < toElementType.NumField(); i++ {
			var toFieldType = rteTo.Field(i)
			var toFieldName = field.Name
			var fromFieldValue = fromElementValue.FieldByName(toFieldName)
			var toFieldValue = translateInner(fromFieldValue, toFieldType)
			toValue.Field(i).Set(toFieldValue)
		}
		return toValue, nil

	case reflect.Slice:
		if kindFrom != reflect.Slice {
			return T{}, NewConvertError("", "Cannot convert from "+kindFrom+" to "+kindTo)
		}
}

/*
 * ConvertError
 */
type ConvertError struct {
	Name  string
	error
}

func NewConvertError(name string, message string) ConvertError {
	return ConvertError{
		Name:  name,
		error: errors.New(message),
	}
}

func (e ConvertError) Error() string {
	return e.error.Error() + ", name: " + e.Name
}

func (e ConvertError) Unwrap() error {
	return e.error
}

func (e ConvertError) HttpStatus() uint {
	return 500
}

/*
 * PropertyNotFoundError
 */
type PropertyNotFoundError struct {
	Name  string
	error
}

func NewPropertyNotFoundError(name string, message string) PropertyNotFoundError {
	return PropertyNotFoundError{
		Name:  name,
		error: errors.New(message),
	}
}

func (e PropertyNotFoundError) Error() string {
	return e.error.Error() + ", name: " + e.Name
}

func (e PropertyNotFoundError) Unwrap() error {
	return e.error
}

func (e PropertyNotFoundError) HttpStatus() uint {
	return 500
}

/*
 * FactoryError
 */
type FactoryError struct {
	Name  string
	error
}

func NewFactoryError(name string, message string) FactoryError {
	return FactoryError{
		Name:  name,
		error: errors.New(message),
	}
}

func (e FactoryError) Error() string {
	return e.error.Error() + ", name: " + e.Name
}

func (e FactoryError) Unwrap() error {
	return e.error
}

func (e FactoryError) HttpStatus() uint {
	return 500
}
