package jwt

import (
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
	"os"
	"strconv"
	"github.com/motojouya/geezer_auth/pkg/model/text"
	"github.com/motojouya/geezer_auth/pkg/model/user"
	"github.com/motojouya/geezer_auth/pkg/utility"
)

type JwtHandler interface {
	GenerateAccessToken(user *user.User, issueDate time.Time) (*user.Authentic, text.JwtToken, error)
}

type jwtHandlerConfig struct {
	Audience              []string
	ValidityPeriodMinutes uint
	GetId                 func () (string, error)
	JwtParser
}

func NewJwtHandler(audience []string, jwtParser JwtParser, validityPeriodMinutes uint, getId (func() (string, error))) JwtHandler {
	return &jwtHandlerConfig{
		Issuer:                issuer,
		Audience:              audience,
		ValidityPeriodMinutes: validityPeriodMinutes,
		GetId:                 getId,
		JwtParser:             jwtParser
	}
}

func CreateJwtHandler() (JwtHandler, error) {
	var audience, audienceExist = os.LookupEnv("JWT_AUDIENCE");
	if !audienceExist {
		return nil, utility.NewSystemConfigError("JWT_AUDIENCE", "JWT_AUDIENCE is not set on env")
	}

	var validityPeriodMinutesStr, validityPeriodMinutesExist = os.LookupEnv("JWT_VALIDITY_PERIOD_MINUTES");
	if !validityPeriodMinutesExist {
		return nil, utility.NewSystemConfigError("JWT_VALIDITY_PERIOD_MINUTES", "JWT_VALIDITY_PERIOD_MINUTES is not set on env")
	}

	var validityPeriodMinutes, err = strconv.Atoi(validityPeriodMinutesStr)
	if err != nil {
		return nil, err
	}

	var GetId = func () (string, error) {
		token, err := uuid.NewUUID()
		if err != nil {
			return "", err
		}

		return token.String(), nil
	}

	var jwtParser, err = CreateJwtIssuerParser()
	if err != nil {
		return nil, err
	}
	if jwtParserConfig, ok := jwtParser.(jwtParserConfig); !ok {
		return nil, utility.NewNilError("jwtParser", "JwtParser is nil")
	}

	return NewJwtHandler(
		[]string{audience,jwtParser.Myself},
		validityPeriodMinutes,
		GetId,
		jwtParserConfig,
	), nil
}

(jwtHandler *jwtHandlerConfig) func GenerateAccessToken(user *user.User, issueDate time.Time) (*user.Authentic, text.JwtToken, error) {
	var id, err = jwtHandler.GetId()
	if err != nil {
		return JwtToken(""), err
	}

	var authentic = user.CreateAuthentic(
		jwtHandler.Issuer,
		jwtHandler.Audience,
		issueDate,
		jwtHandler.ValidityPeriodMinutes,
		id,
		user,
	)

	var claims = FromAuthentic(authentic)

	var token = gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	token.Header["kid"] = jwtHandler.LatestSecretKeyId

	var secret = jwtHandler.SecretMap[jwtHandler.LatestSecretKeyId]
	tokenString, err := token.SignedString(secret)
	if err != nil {
        	return JwtToken(""), err
	}

	return authentic, text.NewJwtToken(tokenString), nil
}
