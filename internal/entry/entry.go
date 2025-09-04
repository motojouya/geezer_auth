package entry

import (
	"github.com/labstack/echo"
)

func createHandle[C any, I any, O any](createControl func () C, handleControl func (C, I, *pkgUser.Authentic) (O, error)) echo.HandlerFunc {
	return func(c echo.Context) error {

		var control = createControl()

		var request I
		if err := c.Bind(&request); err != nil {
			return err
		}

		// TODO *pkgUser.Authenticの解決

		var response, err = handleControl(control, request, c.Get("authentic").(*pkgUser.Authentic))
		if err != nil {
			return err
		}

		return c.Json(200, response)
	}
}
