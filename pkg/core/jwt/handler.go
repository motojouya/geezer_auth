package jwt

import (
	gojwt "github.com/golang-jwt/jwt/v5"
	"time"
	"github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/motojouya/geezer_auth/pkg/core/user"
)

type JwtHandling struct {
	Audience              []string `env:"JWT_AUDIENCE,notEmpty"`
	ValidityPeriodMinutes uint     `env:"JWT_VALIDITY_PERIOD_MINUTES,notEmpty"`
	jwtParserConfig
}

func NewJwtHandling(
	audience []string,
	jwtParser jwtParserConfig,
	validityPeriodMinutes uint,
) JwtHandling {
	return &JwtHandling{
		Audience:              audience,
		ValidityPeriodMinutes: validityPeriodMinutes,
		JwtParser:             jwtParser
	}
}

func (jwtHandling *JwtHandling) getToken(claims *GeezerClaims) (text.JwtToken, error) {

	var token = gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	token.Header["kid"] = jwtHandling.LatestKeyId

	tokenString, err := token.SignedString(jwtHandling.LatestSecret)
	if err != nil {
        	return JwtToken(""), err
	}

	return text.NewJwtToken(tokenString), nil
}

// idはuuidを想定
func (jwtHandling *JwtHandling) Generate(user *user.User, issueDate time.Time, id string) (*user.Authentic, text.JwtToken, error) {
	var authentic = user.CreateAuthentic(
		jwtHandling.Issuer,
		jwtHandling.Audience,
		issueDate,
		jwtHandling.ValidityPeriodMinutes,
		id,
		user,
	)

	var claims = FromAuthentic(authentic)

	var token, err = jwtHandling.getClaims(claims)
	if err != nil {
		return nil, text.JwtToken(""), err
	}

	return authentic, token, nil
}
