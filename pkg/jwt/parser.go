package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"github.com/motojouya/geezer_auth/pkg/model"
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
		return nil, fmt.Error("JWT_ISSUER is not set on env")
	}
	if secretKeyId, secretKeyIdExist := os.LookupEnv("JWT_SECRET_KEY_ID"); !secretKeyIdExist {
		return nil, fmt.Error("JWT_SECRET_KEY_ID is not set on env")
	}
	if secret, secretExist := os.LookupEnv("JWT_SECRET"); !secretExist {
		return nil, fmt.Error("JWT_SECRET is not set on env")
	}
	if myself, myselfExist := os.LookupEnv("JWT_MYSELF"); !myselfExist {
		return nil, fmt.Error("JWT_ISSUER is not set on env")
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
		return nil, fmt.Error("JWT_ISSUER is not set on env")
	}
	if secretKeyId, secretKeyIdExist := os.LookupEnv("JWT_SECRET_KEY_ID"); !secretKeyIdExist {
		return nil, fmt.Error("JWT_SECRET_KEY_ID is not set on env")
	}
	if secret, secretExist := os.LookupEnv("JWT_SECRET"); !secretExist {
		return nil, fmt.Error("JWT_SECRET is not set on env")
	}

	return NewJwtParser(
		issuer,
		issuer,
		secretKeyId,
		map[string]string{secretKeyId:secret},
	), nil
}

// 引数のtokenStringはJwtToken型としてもいいが、いずれにしろこの関数で制約がかかるので、事前にチェックされた値ではなくstringを受けるほうが自然
(jwtParser *JwtParser) func GetUserFromAccessToken(tokenString string) (*model.Authentic, error) {
 	token, err := jwt.ParseWithClaims(
 		tokenString,
 		&GeezerClaims{},
 		func(token *jwt.Token) (interface{}, error) {
			// jwt.SigningMethodHMAC?
 			if _, ok := token.Method.(*jwt.SigningMethodHS256); !ok {
 				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
 			}

			var secret, exist = jwtParser.SecretMap[token.Header["kid"]]
			if !exist {
				return nil, fmt.Errorf("Secret not found for key: %v", token.Header["kid"])
			}
			return secret, nil
 		},
 	)

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(GeezerClaims); !ok || !token.Valid {
		return nil, fmt.Error("Invalid token")
	}

	return claims.ToAuthentic(jwtParser.Issuer, jwtParser.Myself)
}
