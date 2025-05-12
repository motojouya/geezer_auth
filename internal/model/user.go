package model

import (
  "time"
)

type UnsavedUser struct {
  exposeId string
  exposeEmailId string
  name string
  botFlag bool
}

type User struct {
  userId uint
  companyRole *CompanyRole
  email *string
  registeredDate time.Time
  UnsavedUser
}

func CreateUser(exposeId string, emailId string, name string, botFlag bool) UnsavedUser {
  return UnsavedUser{
    exposeId: exposeId,
    exposeEmailId: emailId,
    name: name,
    botFlag: botFlag,
  }
}

func NewUser(userId: uint, exposeId string, name string, emailId string, email *string, botFlag bool, registeredDate time.Time, companyRole: *CompanyRole) User {
  return User{
    userId: userId,
    companyRole: companyRole,
    email: email,
    registeredDate: registeredDate,
    UnsavedUser: UnsavedUser{
      exposeId: exposeId,
      name: name,
      exposeEmailId: emailId,
      botFlag: botFlag,
    },
  }
}

package model

import (
	"crypto/rand"
	"errors"
)

func NewExposeId() (string, error) {
  digit := 6
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 乱数を生成
	b := make([]byte, digit)
	if _, err := rand.Read(b); err != nil {
		return "", errors.New("unexpected error...")
	}

	// letters からランダムに取り出して文字列を生成
	var result string
	for _, v := range b {
        // index が letters の長さに収まるように調整
		result += string(letters[int(v)%len(letters)])
	}
	return result, nil
}


