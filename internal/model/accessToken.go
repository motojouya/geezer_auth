package model

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// TODO middlewareも作ってしまいたい。イメージを掴んで置く

type AccessToken struct {
  token string
  expireDate time.Time
}

type CustomClaims struct {
	// TODO 
}

// expireの基準日がいるので、日付が必要。オプションで期間の調整ができてもいいかもしれない
func PublishAccessToken(accessTokenSecret string, user User, tokens []AccessToken, expireDate time.Time) AccessToken {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo": "bar",
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	tokenString, err := token.SignedString([]byte(accessTokenSecret))
	if err != nil {
		// TODO: handle error
		panic(err)
	}

	return AccessToken{
		token: tokenString,
		expireDate: expireDate,
	}
}

func GetUserFromAccessToken(accesstokensecret string, tokenString string) User {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(accesstokensecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return NewUser(claims["foo"])
	} else {
		return nil
	}

}
