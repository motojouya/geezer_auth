package utility_test

import (
	"github.com/motojouya/geezer_auth/pkg/utility"
	"github.com/stretchr/testify/assert"
	"testing"
	"strconv"
)

func TestNewNilError(t *testing.T) {
	var name = "TestNilError"
	var message = "This is a test nil error"

	var err = utility.NewNilError(name, message)

	assert.Equal(t, name, err.Name)
	assert.Equal(t, message, err.Unwrap().Error())
	assert.Equal(t, message+" (name: "+name+")", err.Error())
	assert.Equal(t, 400, err.HttpStatus())

	t.Logf("error: %s", err.Error())
	t.Logf("error.Name: %s", err.Name)
}

func TestNewSystemConfigError(t *testing.T) {
	var name = "TestSystemConfigError"
	var message = "This is a test system config error"

	var err = utility.NewSystemConfigError(name, message)

	assert.Equal(t, name, err.Name)
	assert.Equal(t, message, err.Unwrap().Error())
	assert.Equal(t, message+" (name: "+name+")", err.Error())
	assert.Equal(t, 500, err.HttpStatus())

	t.Logf("error: %s", err.Error())
	t.Logf("error.Name: %s", err.Name)
}

func TestNewPropertyError(t *testing.T) {
	var name = "TestNilError"
	var message = "This is a test nil error"
	var err = utility.NewNilError(name, message)

	var prop = "TestPath"
	var httpStatus uint = 210

	var propertyError = utility.NewPropertyError(prop, httpStatus, err)

	assert.Equal(t, prop, propertyError.Property)
	assert.Equal(t, httpStatus, propertyError.HttpStatus)
	assert.Equal(t, message, propertyError.Unwrap().Error())
	assert.Equal(t, message+" (property: "+prop+", httpStatus: "+strconv.Itoa(int(httpStatus))+")", propertyError.Error())

	t.Logf("error: %s", propertyError.Error())
	t.Logf("error.Property: %s", propertyError.Property)
}

func TestCreatePropertyError(t *testing.T) {
	var name = "TestNilError"
	var message = "This is a test nil error"
	var err = utility.NewNilError(name, message)

	var prop = "TestPath"
	var httpStatus = 400

	var propertyError = utility.CreatePropertyError(prop, err)

	assert.Equal(t, prop, propertyError.Property)
	assert.Equal(t, httpStatus, propertyError.HttpStatus)
	assert.Equal(t, message, propertyError.Unwrap().Error())
	assert.Equal(t, message+" (property: "+prop+", httpStatus: "+strconv.Itoa(httpStatus)+")", propertyError.Error())

	t.Logf("error: %s", propertyError.Error())
	t.Logf("error.Property: %s", propertyError.Property)
}

func TestPropertyErrorAdd(t *testing.T) {
	var name = "TestNilError"
	var message = "This is a test nil error"
	var err = utility.NewNilError(name, message)

	var prop = "TestPath"
	var httpStatus uint = 210

	var propertyError = utility.NewPropertyError(prop, httpStatus, err)
	var path = "additional"
	var added = propertyError.Add(path)

	assert.Equal(t, path+"."+prop, added.Property)
	assert.Equal(t, httpStatus, added.HttpStatus)
	assert.Equal(t, message, added.Unwrap().Error())
	assert.Equal(t, message+" (property: "+prop+", httpStatus: "+strconv.Itoa(int(httpStatus))+")", added.Error())

	t.Logf("error: %s", added.Error())
	t.Logf("error.Property: %s", added.Property)
}

func TestPropertyErrorChange(t *testing.T) {
	var name = "TestNilError"
	var message = "This is a test nil error"
	var err = utility.NewNilError(name, message)

	var prop = "TestPath"
	var httpStatus uint = 210

	var propertyError = utility.NewPropertyError(prop, httpStatus, err)
	var path = "additional"
	var changedStatus uint = 220
	var added = propertyError.Change(path, changedStatus)

	assert.Equal(t, path+"."+prop, added.Property)
	assert.Equal(t, changedStatus, added.HttpStatus)
	assert.Equal(t, message, added.Unwrap().Error())
	assert.Equal(t, message+" (property: "+prop+", httpStatus: "+strconv.Itoa(int(httpStatus))+")", added.Error())

	t.Logf("error: %s", added.Error())
	t.Logf("error.Property: %s", added.Property)
}

func TestAddPropertyError(t *testing.T) {
	var name = "TestNilError"
	var message = "This is a test nil error"
	var err = utility.NewNilError(name, message)

	var prop = "TestPath"
	var propertyError = utility.AddPropertyError(prop, err)

	var wrapPath = "additional"
	var wrappedPropertyError = utility.AddPropertyError(wrapPath, propertyError)

	assert.Equal(t, wrapPath+"."+prop+"."+name, wrappedPropertyError.Property)
	assert.Equal(t, 400, wrappedPropertyError.HttpStatus)
	assert.Equal(t, message, wrappedPropertyError.Unwrap().Error())
	assert.Equal(t, message+" (property: "+prop+", httpStatus: 400)", wrappedPropertyError.Error())

	t.Logf("error: %s", wrappedPropertyError.Error())
	t.Logf("error.Property: %s", wrappedPropertyError.Property)
}

func TestChangePropertyError(t *testing.T) {
	var name = "TestNilError"
	var message = "This is a test nil error"
	var err = utility.NewNilError(name, message)

	var prop = "TestPath"
	var httpStatus uint = 210
	var propertyError = utility.ChangePropertyError(prop, err, httpStatus)

	var wrapPath = "additional"
	var wraphttpStatus uint = 210
	var wrappedPropertyError = utility.ChangePropertyError(wrapPath, propertyError, wraphttpStatus)

	assert.Equal(t, wrapPath+"."+prop+"."+name, wrappedPropertyError.Property)
	assert.Equal(t, wraphttpStatus, wrappedPropertyError.HttpStatus)
	assert.Equal(t, message, wrappedPropertyError.Unwrap().Error())
	assert.Equal(t, message+" (property: "+prop+", httpStatus: "+strconv.Itoa(int(wraphttpStatus))+")", wrappedPropertyError.Error())

	t.Logf("error: %s", wrappedPropertyError.Error())
	t.Logf("error.Property: %s", wrappedPropertyError.Property)
}

func TestAddPropertyErrorNil(t *testing.T) {
	defer func() {
		if rec := recover(); rec != nil {
			t.Logf("Recovered from panic: %v", rec)
		}
	}()

	var prop = "TestPath"
	var _ = utility.AddPropertyError(prop, nil)

	t.Error("Expected panic for nil source error, but did not panic")
}

func TestChangePropertyErrorNil(t *testing.T) {
	defer func() {
		if rec := recover(); rec != nil {
			t.Logf("Recovered from panic: %v", rec)
		}
	}()

	var prop = "TestPath"
	var httpStatus uint = 210
	var _ = utility.ChangePropertyError(prop, nil, httpStatus)

	t.Error("Expected panic for nil source error, but did not panic")
}
