package model

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// TODO middlewareも作ってしまいたい。イメージを掴んで置く

type AccessToken struct {
	Token string
	ExpireDate time.Time
}

type Company struct {
	ExposeId string
	Name string
	Role string
}

type User struct {
	ExposeId string
	EmailId string
	Email *string
	Name string
	BotFlag bool
	Company *Company
}

// TODO expire dateを入れなくて大丈夫？
type GeezerClaims struct {
	jwt.RegisteredClaims
	UserExposeId string
	UserEmailId string
	UserEmail *string
	UserName string
	BotFlag bool
	CompanyExposeId *string
	CompanyName *string
	CompanyRole *string
}

type JwtSecret struct {
	secret []byte
	version string
}

// expireの基準日がいるので、日付が必要。オプションで期間の調整ができてもいいかもしれない
func PublishAccessToken(jwtSecret JwtSecret, user User, tokens []AccessToken, expireDate time.Time) (AccessToken, error) {

	var claims = GeezerClaims{
		Issuer:    "issuer", // iss
		Subject:   "subject", // sub
		Audience:  []string{"audience"}, // aud
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 + time.Minute)), // exp
		NotBefore: jwt.NewNumericDate(time.Now()), // nbf
		IssuedAt:  jwt.NewNumericDate(time.Now()), // iat
		ID:        "id", // jti

		"foo": "bar",
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	}

	var token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret.secret)
	if err != nil {
        	return "", err
	}

	return AccessToken{
		token: tokenString,
		expireDate: expireDate,
	}, nil
}

func GetUserFromAccessToken(jwtSecret JwtSecret, tokenString string) User {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return jwtSecret.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

 	token, err := jwt.ParseWithClaims(
 		tokenString,
 		&GeezerClaims{},
 		func(token *jwt.Token) (interface{}, error) {
 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
 				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
 			}
 			return secret, nil
 		},
 	)

	if err != nil {
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(GeezerClaims); ok && token.Valid {
		return NewUser(claims.Issuer)
	} else {
		return nil
	}

}
