package accessToken

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
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

(jwtParser JwtParser) func CreateGeezerToken(claims GeezerClaims) (*GeezerToken, error) {
	var company *Company = nil
	if claims.CompanyExposeId != nil && claims.CompanyName != nil && claims.CompanyRole != nil && claims.CompanyRoleName != nil {
		company = NewCompany(*claims.CompanyExposeId, *claims.CompanyName, *claims.CompanyRole, *claims.CompanyRoleName)
	} else {
		if claims.CompanyExposeId != nil {
			return nil, fmt.Error("CompanyExposeId is not nil")
		}
		if claims.CompanyName != nil {
			return nil, fmt.Error("CompanyName is not nil")
		}
		if claims.CompanyRole != nil {
			return nil, fmt.Error("CompanyRole is not nil")
		}
		if claims.CompanyRoleName != nil {
			return nil, fmt.Error("CompanyRoleName is not nil")
		}
		// company = nil
	}

	var user = NewUser(
		claims.Subject,
		claims.UserEmailId,
		claims.UserEmail,
		claims.UserName,
		claims.BotFlag,
		company,
		claims.UpdateDate.Time,
	)

	// issuerチェックでerrorとしているが、Tokenの中にvalidというflagを入れる方法もある。audienceも同様。
	if jwtParser.Issuer == claims.Issuer {
		return nil, fmt.Error("Issuer is not valid")
	}
	if claims.Audience.Contains(jwtParser.Myself) {
		return nil, fmt.Error("Audience is not valid")
	}

	return NewGeezerToken(
		claims.Issuer,
		claims.Subject,
		claims.Audience,
		claims.ExpiresAt.Time,
		claims.NotBefore.Time,
		claims.IssuedAt.Time,
		claims.ID,
		user,
	), nil
}

(jwtParser *JwtParser) func GetUserFromAccessToken(tokenString string) (*GeezerToken, error) {
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

	return jwtParser.CreateGeezerToken(claims)
}
