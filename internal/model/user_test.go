package model_test

import (
	"github.com/stretchr/testify/assert"
	"model"
	"testing"
	"time"
)

func TestSum(t *testing.T) {
	actual := sum(3, 5)
	expected := 8

	assert.Equal(t, expected, actual)
}

func CreateUser(exposeId string, emailId string, name string, botFlag bool) UnsavedUser {
	return UnsavedUser{
		exposeId:      exposeId,
		exposeEmailId: emailId,
		name:          name,
		botFlag:       botFlag,
	}
}

func NewUser(userId uint, exposeId string, name string, emailId string, email *string, botFlag bool, registeredDate time.Time, companyRole *CompanyRole) User {
	return User{
		userId:         userId,
		companyRole:    companyRole,
		email:          email,
		registeredDate: registeredDate,
		UnsavedUser:    CreateUser(exposeId, name, emailId, botFlag),
	}
}
