package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strings"
	"github.com/motojouya/geezer_auth/pkg/model/user"
	"github.com/motojouya/geezer_auth/pkg/utility"
)

// TODO middlewareも作ってしまいたい。イメージを掴んで置く
type JwtParser struct {
	Issuer            string
	Myself            string // Audienceと付き合わせるための自分自身の情報
	LatestSecretKeyId string
	SecretMap         map[string]string
}

func NewJwtParser(issuer string, myself string, latestSecret string, secretMap map[string]string) *JwtParser {
	return &JwtParser{
		Issuer:       issuer,
		Myself:       myself,
		LatestSecret: latestSecret,
		SecretMap:    secretMap,
	}
}

func CreateJwtAudienceParser() (*JwtParser, error) {
	if issuer, issuerExist := os.LookupEnv("JWT_ISSUER"); !issuerExist {
		return nil, utility.NewSystemConfigError("JWT_ISSUER", "JWT_ISSUER is not set on env")
	}
	if secretKeyId, secretKeyIdExist := os.LookupEnv("JWT_SECRET_KEY_ID"); !secretKeyIdExist {
		return nil, utility.NewSystemConfigError("JWT_SECRET_KEY_ID", "JWT_SECRET_KEY_ID is not set on env")
	}
	if secret, secretExist := os.LookupEnv("JWT_SECRET"); !secretExist {
		return nil, utility.NewSystemConfigError("JWT_SECRET", "JWT_SECRET is not set on env")
	}
	if myself, myselfExist := os.LookupEnv("JWT_MYSELF"); !myselfExist {
		return nil, utility.NewSystemConfigError("JWT_ISSUER", "JWT_ISSUER is not set on env")
	}

	return NewJwtParser(
		issuer,
		myself,
		secretKeyId,
		map[string]string{secretKeyId:secret},
	), nil
}

func CreateJwtIssuerParser() (*JwtParser, error) {
	if issuer, issuerExist := os.LookupEnv("JWT_ISSUER"); !issuerExist {
		return nil, utility.NewSystemConfigError("JWT_ISSUER", "JWT_ISSUER is not set on env")
	}
	if secretKeyId, secretKeyIdExist := os.LookupEnv("JWT_SECRET_KEY_ID"); !secretKeyIdExist {
		return nil, utility.NewSystemConfigError("JWT_SECRET_KEY_ID", "JWT_SECRET_KEY_ID is not set on env")
	}
	if secret, secretExist := os.LookupEnv("JWT_SECRET"); !secretExist {
		return nil, utility.NewSystemConfigError("JWT_SECRET", "JWT_SECRET is not set on env")
	}

	return NewJwtParser(
		issuer,
		issuer,
		secretKeyId,
		map[string]string{secretKeyId:secret},
	), nil
}

(jwtParser *JwtParser) func Validate(claims *GeezerClaims) (error) {
	if jwtParser.Issuer != claims.Issuer {
		return NewJwtError("Issuer", claims.Issuer, "Issuer is not valid")
	}
	if claims.Audience.Contains(jwtParser.Myself) {
		return NewJwtError("Audience", strings.Join(claims.Audience, ","), "Audience is not valid")
	}
	return nil
}

// 引数のtokenStringはJwtToken型としてもいいが、いずれにしろこの関数で制約がかかるので、事前にチェックされた値ではなくstringを受けるほうが自然
(jwtParser *JwtParser) func GetUserFromAccessToken(tokenString string) (*user.Authentic, error) {
 	token, err := jwt.ParseWithClaims(
 		tokenString,
 		&GeezerClaims{},
 		func(token *jwt.Token) (interface{}, error) {
			// jwt.SigningMethodHMAC?
 			if _, ok := token.Method.(*jwt.SigningMethodHS256); !ok {
 				return nil, NewJwtError("header.alg", token.Header["alg"], "Unexpected signing method")
 			}

			var secret, exist = jwtParser.SecretMap[token.Header["kid"]]
			if !exist {
				return nil, NewJwtError("header.kid", token.Header["kid"], "Secret not found for key")
			}
			return secret, nil
 		},
 	)

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(GeezerClaims); !ok || !token.Valid {
		return nil, NewJwtError("hole token", tokenString, "Invalid token")
	}

	if err := jwtParser.Validate(claims); err != nil {
		return nil, err
	}

	return claims.ToAuthentic()
}
