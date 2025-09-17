package entry

import (
	"github.com/labstack/echo/v4"
	configBehavior "github.com/motojouya/geezer_auth/internal/behavior/config"
	"github.com/motojouya/geezer_auth/internal/entry/transfer/common"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
	"github.com/motojouya/geezer_auth/internal/shelter/essence"
	pkgUser "github.com/motojouya/geezer_auth/pkg/shelter/user"
)

func Hand[C any, I any, O any](createControl func() (C, error), handleControl func(C, I, *pkgUser.Authentic) (O, error)) echo.HandlerFunc {
	return func(c echo.Context) error {

		var request I
		if err := c.Bind(&request); err != nil {
			return err
		}

		authentic, err := getAuthentic(request)
		if err != nil {
			return err
		}

		control, err := createControl()
		if err != nil {
			return err
		}

		response, err := handleControl(control, request, authentic)
		if err != nil {
			return err
		}

		if closable, ok := any(control).(essence.Closable); ok {
			err := closable.Close()
			if err != nil {
				return err
			}
		}

		return c.JSON(200, response)
	}
}

func getAuthentic[I any](request I) (*pkgUser.Authentic, error) {
	var env = localPkg.CreateEnvironment()
	jwt, err := configBehavior.NewJwtHandlerGet(env).GetJwtHandler()
	if err != nil {
		return nil, err
	}

	header, ok := any(request).(common.BearerTokenGetter)
	if !ok {
		return nil, nil
	}

	token, err := header.GetBearerToken()
	if err != nil {
		return nil, err
	}

	return jwt.Parse(token)
}
