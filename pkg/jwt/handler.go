package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
	"os"
	"strconv"
)

type JwtHandler struct {
	Audience              []string
	ValidityPeriodMinutes uint
	GetId                 func () (string, error)
	JwtParser
}

func NewJwtHandler(audience []string, jwtParser JwtParser, validityPeriodMinutes uint, GetId (func() (string, error))) *JwtHandler {
	return &JwtHandler{
		Issuer:                issuer,
		Audience:              audience,
		ValidityPeriodMinutes: validityPeriodMinutes,
		GetId:                 GetId,
		JwtParser:             jwtParser
	}
}

func CreateJwtHandler() (*JwtHandler, error) {
	var audience, audienceExist = os.LookupEnv("JWT_AUDIENCE");
	if !audienceExist {
		return nil, fmt.Error("JWT_AUDIENCE is not set on env")
	}

	var validityPeriodMinutesStr, validityPeriodMinutesExist = os.LookupEnv("JWT_VALIDITY_PERIOD_MINUTES");
	if !validityPeriodMinutesExist {
		return nil, fmt.Error("JWT_VALIDITY_PERIOD_MINUTES is not set on env")
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
	if jwtParser == nil {
		return nil, fmt.Error("jwtParser is nil")
	}

	return NewJwtHandler(
		[]string{audience,jwtParser.Myself},
		validityPeriodMinutes,
		GetId,
		*jwtParser,
	), nil
}

(jwtHandler *JwtHandler) func GenerateAccessToken(user User, issueDate time.Time) (JwtToken, error) {
	var id, err = jwtHandler.GetId()
	if err != nil {
		return JwtToken(""), err
	}

	var expireDate = issueDate.Add(jwtHandler.ValidityPeriodMinutes * time.Minute)

	var claims = CreateClaims(
		user,
		jwtHandler.Issuer,
		jwtHandler.Audience,
		expireDate,
		issueDate,
		id,
	)

	var token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["kid"] = jwtHandler.LatestSecretKeyId

	var secret = jwtHandler.SecretMap[jwtHandler.LatestSecretKeyId]
	tokenString, err := token.SignedString(secret)
	if err != nil {
        	return JwtToken(""), err
	}

	return NewJwtToken(tokenString), nil
}
