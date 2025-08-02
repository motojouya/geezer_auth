package jwt

import (
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/motojouya/geezer_auth/pkg/core/user"
	"time"
)

type JwtHandler interface {
	Generate(user *user.User, issueDate time.Time, id string) (*user.Authentic, text.JwtToken, error)
}

type JwtHandle struct {
	Audience              []string `env:"JWT_AUDIENCE,notEmpty"`
	ValidityPeriodMinutes uint     `env:"JWT_VALIDITY_PERIOD_MINUTES,notEmpty"`
	JwtParse
}

func NewJwtHandle(
	audience []string,
	jwtParse JwtParse,
	validityPeriodMinutes uint,
) JwtHandle {
	return JwtHandle{
		Audience:              audience,
		ValidityPeriodMinutes: validityPeriodMinutes,
		JwtParse:              jwtParse,
	}
}

func (jwtHandle *JwtHandle) getToken(claims *GeezerClaims) (text.JwtToken, error) {

	var token = gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	token.Header["kid"] = jwtHandle.LatestKeyId

	tokenString, err := token.SignedString([]byte(jwtHandle.LatestSecret))
	if err != nil {
		return text.JwtToken(""), err
	}

	return text.NewJwtToken(tokenString), nil
}

// idはuuidを想定
func (jwtHandle *JwtHandle) Generate(userValue *user.User, issueDate time.Time, id string) (*user.Authentic, text.JwtToken, error) {
	var authentic = user.CreateAuthentic(
		jwtHandle.Issuer,
		jwtHandle.Audience,
		issueDate,
		jwtHandle.ValidityPeriodMinutes,
		id,
		userValue,
	)

	var claims = FromAuthentic(authentic)

	var token, err = jwtHandle.getToken(claims)
	if err != nil {
		return nil, text.JwtToken(""), err
	}

	return authentic, token, nil
}
