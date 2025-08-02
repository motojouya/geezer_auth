package config

import (
	"github.com/motojouya/geezer_auth/internal/io"
	"github.com/motojouya/geezer_auth/pkg/core/jwt"
)

type JwtHandlerGetter interface {
	GetJwtHandler() (jwt.JwtHandler, error)
}

type JwtHandlerGet struct {
	env io.JwtHandleGetter
}

func NewJwtHandlerGet(env io.JwtHandleGetter) *JwtHandlerGet {
	return &JwtHandlerGet{
		env: env,
	}
}

var jwtHandle *jwt.JwtHandle

func (getter *JwtHandlerGet) GetJwtHandler() (*jwt.JwtHandle, error) {
	if jwtHandle == nil {
		var jwtHandleObj, err = getter.env.GetJwtHandle()
		if err != nil {
			return nil, err
		}

		jwtHandle = &jwtHandleObj
	}

	return jwtHandle, nil
}

// !old code!
// func CreateJwtHandler() (JwtHandler, error) {
// 	var audience, audienceExist = os.LookupEnv("JWT_AUDIENCE");
// 	if !audienceExist {
// 		return nil, essence.NewSystemConfigError("JWT_AUDIENCE", "JWT_AUDIENCE is not set on env")
// 	}
//
// 	var validityPeriodMinutesStr, validityPeriodMinutesExist = os.LookupEnv("JWT_VALIDITY_PERIOD_MINUTES");
// 	if !validityPeriodMinutesExist {
// 		return nil, essence.NewSystemConfigError("JWT_VALIDITY_PERIOD_MINUTES", "JWT_VALIDITY_PERIOD_MINUTES is not set on env")
// 	}
//
// 	var validityPeriodMinutes, err = strconv.Atoi(validityPeriodMinutesStr)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	var GetId = func () (string, error) {
// 		token, err := uuid.NewUUID()
// 		if err != nil {
// 			return "", err
// 		}
//
// 		return token.String(), nil
// 	}
//
// 	var jwtParser, err = CreateJwtIssuerParser()
// 	if err != nil {
// 		return nil, err
// 	}
// 	if jwtParserConfig, ok := jwtParser.(jwtParserConfig); !ok {
// 		return nil, essence.NewNilError("jwtParser", "JwtParser is nil")
// 	}
//
// 	return NewJwtHandler(
// 		[]string{audience,jwtParser.Myself},
// 		validityPeriodMinutes,
// 		GetId,
// 		jwtParserConfig,
// 	), nil
// }
