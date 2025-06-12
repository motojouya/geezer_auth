package jwt

import (
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/motojouya/geezer_auth/pkg/core/user"
	"strings"
	"slices"
)

// TODO middlewareも作ってしまいたい。イメージを掴んで置く
/*
 * MyselfはAudienceと照合する自サーバの識別情報
 */
type JwtParsering struct {
	Issuer       string `env:"JWT_ISSUER,notEmpty"`
	Myself       string `env:"JWT_MYSELF,notEmpty"`
	LatestKeyId  string `env:"JWT_LATEST_KEY_ID,notEmpty"`
	LatestSecret string `env:"JWT_LATEST_SECRET,notEmpty"`
	OldKeyId     string `env:"JWT_OLD_KEY_ID"`
	OldSecret    string `env:"JWT_OLD_SECRET"`
}

func NewJwtParsering(
	issuer string,
	myself string,
	latestKeyId string,
	latestSecret string,
	oldKeyId string,
	oldSecret string,
) JwtParsering {
	return JwtParsering{
		Issuer:       issuer,
		Myself:       myself,
		LatestKeyId:  latestKeyId,
		LatestSecret: latestSecret,
		OldKeyId:     oldKeyId,
		OldSecret:    oldSecret,
	}
}

func (jwtParsering *JwtParsering) getClaims(tokenString string) (*GeezerClaims, error) {
	token, err := gojwt.ParseWithClaims(
		tokenString,
		&GeezerClaims{},
		func(token *gojwt.Token) (interface{}, error) {
			// gojwt.SigningMethodHS256?
			if _, ok := token.Method.(*gojwt.SigningMethodHMAC); !ok {
				var alg, ok = token.Header["alg"].(string)
				if !ok {
					return nil, NewJwtError("header.alg", "", "Unexpected signing method")
				}
				return nil, NewJwtError("header.alg", alg, "Unexpected signing method")
			}

			if token.Header["kid"] == jwtParsering.LatestKeyId {
				return []byte(jwtParsering.LatestSecret), nil
			}

			if token.Header["kid"] == jwtParsering.OldKeyId {
				return []byte(jwtParsering.OldSecret), nil
			}

			var kid, ok = token.Header["kid"].(string)
			if !ok {
				return nil, NewJwtError("header.kid", "", "Secret not found for key")
			}
			return nil, NewJwtError("header.kid", kid, "Secret not found for key")
		},
	)

	if err != nil {
		return nil, err
	}

	var claims, ok = token.Claims.(GeezerClaims);
	if  !ok || !token.Valid {
		return nil, NewJwtError("hole token", tokenString, "Invalid token")
	}

	if jwtParsering.Issuer != claims.Issuer {
		return nil, NewJwtError("Issuer", claims.Issuer, "Issuer is not valid")
	}

	if slices.Contains(claims.Audience, jwtParsering.Myself) {
		return nil, NewJwtError("Audience", strings.Join(claims.Audience, ","), "Audience is not valid")
	}

	return &claims, nil
}

// 引数のtokenStringはJwtToken型としてもいいが、いずれにしろこの関数で制約がかかるので、事前にチェックされた値ではなくstringを受けるほうが自然
func (jwtParsering *JwtParsering) Parse(tokenString string) (*user.Authentic, error) {

	var claims, err = jwtParsering.getClaims(tokenString)
	if err != nil {
		return nil, err
	}

	return claims.ToAuthentic()
}
