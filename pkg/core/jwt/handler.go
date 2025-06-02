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
	GetId                 func () (string, error)
	jwtParserConfig
}

func NewJwtHandling(
	audience []string,
	jwtParser jwtParserConfig,
	validityPeriodMinutes uint,
	getId (func() (string, error)),
) JwtHandling {
	return &JwtHandling{
		Audience:              audience,
		ValidityPeriodMinutes: validityPeriodMinutes,
		GetId:                 getId,
		JwtParser:             jwtParser
	}
}

func (jwtHandling *JwtHandling) Generate(user *user.User, issueDate time.Time) (*user.Authentic, text.JwtToken, error) {
	var id, err = jwtHandling.GetId()
	if err != nil {
		return JwtToken(""), err
	}

	var authentic = user.CreateAuthentic(
		jwtHandling.Issuer,
		jwtHandling.Audience,
		issueDate,
		jwtHandling.ValidityPeriodMinutes,
		id,
		user,
	)

	var claims = FromAuthentic(authentic)

	var token = gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	token.Header["kid"] = jwtHandling.LatestKeyId

	tokenString, err := token.SignedString(jwtHandling.LatestSecret)
	if err != nil {
        	return JwtToken(""), err
	}

	return authentic, text.NewJwtToken(tokenString), nil
}
