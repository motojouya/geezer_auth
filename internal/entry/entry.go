package entry

import (
	"errors"
	"github.com/labstack/echo/v4"
	configBehavior "github.com/motojouya/geezer_auth/internal/behavior/config"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
	"github.com/motojouya/geezer_auth/internal/shelter/essence"
	pkgUser "github.com/motojouya/geezer_auth/pkg/shelter/user"
	"strings"
)

var tokenPrefix = "Bearer "

type RequestHeader struct {
	token string `header:"Authorization"`
}

func (r *RequestHeader) GetToken() (string, error) {
	if !strings.HasPrefix(r.token, tokenPrefix) {
		return "", errors.New("invalid token") // TODO 独自エラーにする
	}
	return strings.TrimPrefix(r.token, tokenPrefix), nil
}

func Hand[C any, I any, O any](createControl func() (C, error), handleControl func(C, I, *pkgUser.Authentic) (O, error)) echo.HandlerFunc {
	return func(c echo.Context) error {

		var request I
		if err := c.Bind(&request); err != nil {
			return err
		}

		var header RequestHeader
		if err := c.Bind(&header); err != nil {
			return err
		}

		authentic, err := getAuthentic(header)
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

		var con any = control
		if closable, ok := con.(essence.Closable); ok {
			err := closable.Close()
			if err != nil {
				return err
			}
		}

		return c.JSON(200, response)
	}
}

func getAuthentic(header RequestHeader) (*pkgUser.Authentic, error) {
	var env = localPkg.CreateEnvironment()
	jwt, err := configBehavior.NewJwtHandlerGet(env).GetJwtHandler()
	if err != nil {
		return err
	}

	token, err := header.GetToken()
	if err != nil {
		return err
	}

	return jwt.Parse(token)
}
