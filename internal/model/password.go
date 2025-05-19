package model

import (
	"golang.org/x/crypto/bcrypt"
)

type Password string

func GetPassword(password string) (Password, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 9)
	if err != nil {
		return "", err
	}

	return Password(hashed), nil
}

func VerifyPassword(storedPassword Password, argumentPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(string(storedPassword)), []byte(argumentPassword))
	if err != nil {
		// should be return bcrypt.ErrMismatchedHashAndPassword? or fmt.Println(err)?
		return false
	}
	return true
}
