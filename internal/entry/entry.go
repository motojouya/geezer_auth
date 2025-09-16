package entry

import (
	"github.com/labstack/echo/v4"
	configBehavior "github.com/motojouya/geezer_auth/internal/behavior/config"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
	"github.com/motojouya/geezer_auth/internal/shelter/essence"
	"github.com/motojouya/geezer_auth/internal/entry/transfer/common"
	pkgUser "github.com/motojouya/geezer_auth/pkg/shelter/user"
	"fmt"
)

func Hand[C any, I any, O any](createControl func() (C, error), handleControl func(C, I, *pkgUser.Authentic) (O, error)) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("Hand called01")

		var request I
		if err := c.Bind(&request); err != nil {
			fmt.Println("Hand called02, " + err.Error())
			return err
		}

		authentic, err := getAuthentic(request)
		if err != nil {
			fmt.Println("Hand called04, " + err.Error())
			return err
		}

		control, err := createControl()
		if err != nil {
			fmt.Println("Hand called05, " + err.Error())
			return err
		}

		response, err := handleControl(control, request, authentic)
		if err != nil {
			fmt.Println("Hand called06, " + err.Error())
			return err
		}

		if closable, ok := any(control).(essence.Closable); ok {
			err := closable.Close()
			if err != nil {
				fmt.Println("Hand called07, " + err.Error())
				return err
			}
		}

		fmt.Println("Hand called08")
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
		return nil, nil;
	}

	token, err := header.GetBearerToken()
	if err != nil {
		return nil, err
	}

	return jwt.Parse(token)
}
