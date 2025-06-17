package jwt

import (
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/motojouya/geezer_auth/pkg/core/user"
	"time"
)

type JwtHandling struct {
	Audience              []string `env:"JWT_AUDIENCE,notEmpty"`
	ValidityPeriodMinutes uint     `env:"JWT_VALIDITY_PERIOD_MINUTES,notEmpty"`
	JwtParsing
}

func NewJwtHandling(
	audience []string,
	jwtParsing JwtParsing,
	validityPeriodMinutes uint,
) JwtHandling {
	return JwtHandling{
		Audience:              audience,
		ValidityPeriodMinutes: validityPeriodMinutes,
		JwtParsing:            jwtParsing,
	}
}

func (jwtHandling *JwtHandling) getToken(claims *GeezerClaims) (text.JwtToken, error) {

	var token = gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	token.Header["kid"] = jwtHandling.LatestKeyId

	tokenString, err := token.SignedString([]byte(jwtHandling.LatestSecret))
	if err != nil {
		return text.JwtToken(""), err
	}

	return text.NewJwtToken(tokenString), nil
}

// idはuuidを想定
func (jwtHandling *JwtHandling) Generate(userValue *user.User, issueDate time.Time, id string) (*user.Authentic, text.JwtToken, error) {
	var authentic = user.CreateAuthentic(
		jwtHandling.Issuer,
		jwtHandling.Audience,
		issueDate,
		jwtHandling.ValidityPeriodMinutes,
		id,
		userValue,
	)

	var claims = FromAuthentic(authentic)

	var token, err = jwtHandling.getToken(claims)
	if err != nil {
		return nil, text.JwtToken(""), err
	}

	return authentic, token, nil
}
