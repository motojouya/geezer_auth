package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	authCtrl "github.com/motojouya/geezer_auth/internal/control/auth"
	roleCtrl "github.com/motojouya/geezer_auth/internal/control/role"
	"github.com/motojouya/geezer_auth/internal/entry"
)

func Route(e *echo.Echo) {
	// e := echo.New()
	// e.Start(":1323")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("2M")) // FIXME テキトー
	// prefixed := e.Group("/auth") TODO サーバの配置によってはprefixなpathが必要になるかも

	UserRoute(e.Group("/user"))

	CompanyRoute(e.Group("/company"))

	e.GET("/role", entry.Hand(roleCtrl.CreateGetRoleControl, roleCtrl.GetRoleExecute))
	e.POST("/login", entry.Hand(authCtrl.CreateLoginControl, authCtrl.LoginExecute))
	e.POST("/refresh", entry.Hand(authCtrl.CreateRefreshAccessTokenControl, authCtrl.RefreshAccessTokenExecute))
}
