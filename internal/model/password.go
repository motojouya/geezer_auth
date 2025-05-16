package model

import (
	"golang.org/x/crypto/bcrypt"
)

// FIXME use type Password string
type Password = string

func GetPassword(password string) (Password, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 9)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func VerifyPassword(storedPassword Password, argumentPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(argumentPassword))
	if err != nil {
		// should be return bcrypt.ErrMismatchedHashAndPassword? or fmt.Println(err)?
		return false
	}
	return true
}
