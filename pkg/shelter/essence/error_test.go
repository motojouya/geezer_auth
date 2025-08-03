package essence_test

import (
	"github.com/motojouya/geezer_auth/pkg/shelter/essence"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestNewNilError(t *testing.T) {
	var name = "TestNilError"
	var message = "This is a test nil error"
	var httpStatus uint = 400

	var err = essence.NewNilError(name, message)

	assert.Equal(t, name, err.Name)
	assert.Equal(t, message, err.Unwrap().Error())
	assert.Equal(t, message+", name: "+name, err.Error())
	assert.Equal(t, httpStatus, err.HttpStatus())

	t.Logf("error: %s", err.Error())
	t.Logf("error.Name: %s", err.Name)
	t.Logf("error.HttpStatus: %d", err.HttpStatus())
}

func TestNewSystemConfigError(t *testing.T) {
	var name = "TestSystemConfigError"
	var message = "This is a test system config error"
	var httpStatus uint = 500

	var err = essence.NewSystemConfigError(name, message)

	assert.Equal(t, name, err.Name)
	assert.Equal(t, message, err.Unwrap().Error())
	assert.Equal(t, message+", name: "+name, err.Error())
	assert.Equal(t, httpStatus, err.HttpStatus())

	t.Logf("error: %s", err.Error())
	t.Logf("error.Name: %s", err.Name)
}

func TestNewPropertyError(t *testing.T) {
	var name = "TestNilError"
	var message = "This is a test nil error"
	var err = essence.NewNilError(name, message)

	var prop = "TestPath"
	var httpStatus uint = 210

	var propertyError = essence.NewPropertyError(prop, httpStatus, err)

	assert.Equal(t, prop, propertyError.Property)
	assert.Equal(t, httpStatus, propertyError.HttpStatusCode)
	assert.Equal(t, message+", name: "+name, propertyError.Unwrap().Error())
	assert.Equal(t, message+", name: "+name+", property: "+prop+", httpStatus: "+strconv.Itoa(int(httpStatus)), propertyError.Error())

	t.Logf("error: %s", propertyError.Error())
	t.Logf("error.Property: %s", propertyError.Property)
}

func TestCreatePropertyError(t *testing.T) {
	var name = "TestNilError"
	var message = "This is a test nil error"
	var err = essence.NewNilError(name, message)

	var prop = "TestPath"
	var httpStatus uint = 400

	var propertyError = essence.CreatePropertyError(prop, err)

	assert.Equal(t, prop, propertyError.Property)
	assert.Equal(t, httpStatus, propertyError.HttpStatusCode)
	assert.Equal(t, message+", name: "+name, propertyError.Unwrap().Error())
	assert.Equal(t, message+", name: "+name+", property: "+prop+", httpStatus: "+strconv.Itoa(int(httpStatus)), propertyError.Error())

	t.Logf("error: %s", propertyError.Error())
	t.Logf("error.Property: %s", propertyError.Property)
}

func TestPropertyErrorAdd(t *testing.T) {
	var name = "TestNilError"
	var message = "This is a test nil error"
	var err = essence.NewNilError(name, message)

	var prop = "TestPath"
	var httpStatus uint = 210

	var propertyError = essence.NewPropertyError(prop, httpStatus, err)
	var path = "additional"
	var added = propertyError.Add(path)

	assert.Equal(t, path+"."+prop, added.Property)
	assert.Equal(t, httpStatus, added.HttpStatusCode)
	assert.Equal(t, message+", name: "+name, added.Unwrap().Error())
	assert.Equal(t, message+", name: "+name+", property: "+path+"."+prop+", httpStatus: "+strconv.Itoa(int(httpStatus)), added.Error())

	t.Logf("error: %s", added.Error())
	t.Logf("error.Property: %s", added.Property)
}

func TestPropertyErrorChange(t *testing.T) {
	var name = "TestNilError"
	var message = "This is a test nil error"
	var err = essence.NewNilError(name, message)

	var prop = "TestPath"
	var httpStatus uint = 210

	var propertyError = essence.NewPropertyError(prop, httpStatus, err)
	var path = "additional"
	var changedStatus uint = 220
	var added = propertyError.Change(path, changedStatus)

	assert.Equal(t, path+"."+prop, added.Property)
	assert.Equal(t, changedStatus, added.HttpStatusCode)
	assert.Equal(t, message+", name: "+name, added.Unwrap().Error())
	assert.Equal(t, message+", name: "+name+", property: "+path+"."+prop+", httpStatus: "+strconv.Itoa(int(changedStatus)), added.Error())

	t.Logf("error: %s", added.Error())
	t.Logf("error.Property: %s", added.Property)
}

func TestAddPropertyError(t *testing.T) {
	var name = "TestNilError"
	var message = "This is a test nil error"
	var err = essence.NewNilError(name, message)

	var prop = "TestPath"
	var propertyError = essence.AddPropertyError(prop, err)

	var wrapPath = "additional"
	var httpStatus uint = 400
	var wrappedPropertyError = essence.AddPropertyError(wrapPath, propertyError)

	assert.Equal(t, wrapPath+"."+prop, wrappedPropertyError.Property)
	assert.Equal(t, httpStatus, wrappedPropertyError.HttpStatusCode)
	assert.Equal(t, message+", name: "+name, wrappedPropertyError.Unwrap().Error())
	assert.Equal(t, message+", name: "+name+", property: "+wrapPath+"."+prop+", httpStatus: 400", wrappedPropertyError.Error())

	t.Logf("error: %s", wrappedPropertyError.Error())
	t.Logf("error.Property: %s", wrappedPropertyError.Property)
}

func TestChangePropertyError(t *testing.T) {
	var name = "TestNilError"
	var message = "This is a test nil error"
	var err = essence.NewNilError(name, message)

	var prop = "TestPath"
	var httpStatus uint = 210
	var propertyError = essence.ChangePropertyError(prop, err, httpStatus)

	var wrapPath = "additional"
	var wraphttpStatus uint = 210
	var wrappedPropertyError = essence.ChangePropertyError(wrapPath, propertyError, wraphttpStatus)

	assert.Equal(t, wrapPath+"."+prop, wrappedPropertyError.Property)
	assert.Equal(t, wraphttpStatus, wrappedPropertyError.HttpStatusCode)
	assert.Equal(t, message+", name: "+name, wrappedPropertyError.Unwrap().Error())
	assert.Equal(t, message+", name: "+name+", property: "+wrapPath+"."+prop+", httpStatus: "+strconv.Itoa(int(wraphttpStatus)), wrappedPropertyError.Error())

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
	var _ = essence.AddPropertyError(prop, nil)

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
	var _ = essence.ChangePropertyError(prop, nil, httpStatus)

	t.Error("Expected panic for nil source error, but did not panic")
}
