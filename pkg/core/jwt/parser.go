package jwt

import (
	gojwt "github.com/golang-jwt/jwt/v5"
	"strings"
	"github.com/motojouya/geezer_auth/pkg/core/user"
)

// TODO middlewareも作ってしまいたい。イメージを掴んで置く
/*
 * MyselfはAudienceと照合する自サーバの識別情報
 */
type JwtParsering struct {
	Issuer       string  `env:"JWT_ISSUER,notEmpty"`
	Myself       string  `env:"JWT_MYSELF,notEmpty"`
	LatestKeyId  string  `env:"JWT_LATEST_KEY_ID,notEmpty"`
	LatestSecret string  `env:"JWT_LATEST_SECRET,notEmpty"`
	OldKeyId     string `env:"JWT_OLD_KEY_ID"`
	OldSecret    string `env:"JWT_OLD_SECRET"`
}

func NewJwtParsing(
	issuer       string,
	myself       string,
	latestKeyId  string,
	latestSecret string,
	oldKeyId     string,
	oldSecret    string,
) JwtParsing {
	return &JwtParsering{
		Issuer:       issuer,
		Myself:       myself,
		LatestKeyId:  latestKeyId,
		LatestSecret: latestSecret,
		OldKeyId:     oldKeyId,
		OldSecret:    oldSecret,
	}
}

// 引数のtokenStringはJwtToken型としてもいいが、いずれにしろこの関数で制約がかかるので、事前にチェックされた値ではなくstringを受けるほうが自然
func (jwtParsering *JwtParsering) Parse(tokenString string) (*user.Authentic, error) {
 	token, err := gojwt.ParseWithClaims(
 		tokenString,
 		&GeezerClaims{},
 		func(token *gojwt.Token) (interface{}, error) {
			// gojwt.SigningMethodHMAC?
 			if _, ok := token.Method.(*gojwt.SigningMethodHS256); !ok {
 				return nil, NewJwtError("header.alg", token.Header["alg"], "Unexpected signing method")
 			}

			if token.Header["kid"] == jwtParsering.LatestKeyId {
				return []byte(jwtParsering.LatestSecret), nil
			}

			if token.Header["kid"] == jwtParsering.OldKeyId {
				return []byte(jwtParsering.OldSecret), nil
			}

			return nil, NewJwtError("header.kid", token.Header["kid"], "Secret not found for key")
 		},
 	)

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(GeezerClaims); !ok || !token.Valid {
		return nil, NewJwtError("hole token", tokenString, "Invalid token")
	}

	if jwtParsering.Issuer != claims.Issuer {
		return NewJwtError("Issuer", claims.Issuer, "Issuer is not valid")
	}

	if claims.Audience.Contains(jwtParsering.Myself) {
		return NewJwtError("Audience", strings.Join(claims.Audience, ","), "Audience is not valid")
	}

	return claims.ToAuthentic()
}
