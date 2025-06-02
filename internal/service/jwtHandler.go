package service

import (
	"time"
	"github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/motojouya/geezer_auth/pkg/core/user"
	"github.com/motojouya/geezer_auth/pkg/core/jwt"
	"github.com/motojouya/geezer_auth/internal/io"
)

type JwtHandler interface {
	Generate(user *user.User, issueDate time.Time) (*user.Authentic, text.JwtToken, error)
}

func LoadJwtHandler(local io.Local) (JwtHandler, error) {
	var jwtHandling, err = local.GetEnv[jwt.JwtHandling]()
	if err != nil {
		return err
	}

	jwtHandling.GetId = func () (string, error) {
		uuid, err := local.GenerateUUID()
		if err != nil {
			return "", err
		}
		return uuid.String(), nil
	}

	return jwtHandling, nil
}

// !old code!
// func CreateJwtHandler() (JwtHandler, error) {
// 	var audience, audienceExist = os.LookupEnv("JWT_AUDIENCE");
// 	if !audienceExist {
// 		return nil, utility.NewSystemConfigError("JWT_AUDIENCE", "JWT_AUDIENCE is not set on env")
// 	}
// 
// 	var validityPeriodMinutesStr, validityPeriodMinutesExist = os.LookupEnv("JWT_VALIDITY_PERIOD_MINUTES");
// 	if !validityPeriodMinutesExist {
// 		return nil, utility.NewSystemConfigError("JWT_VALIDITY_PERIOD_MINUTES", "JWT_VALIDITY_PERIOD_MINUTES is not set on env")
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
// 		return nil, utility.NewNilError("jwtParser", "JwtParser is nil")
// 	}
// 
// 	return NewJwtHandler(
// 		[]string{audience,jwtParser.Myself},
// 		validityPeriodMinutes,
// 		GetId,
// 		jwtParserConfig,
// 	), nil
// }

