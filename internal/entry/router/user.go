package router

import (
	"github.com/labstack/echo/v4"
	userCtrl "github.com/motojouya/geezer_auth/internal/control/user"
	"github.com/motojouya/geezer_auth/internal/entry"
)

func UserRoute(e *echo.Group) {
	e.POST("/register", entry.Hand(userCtrl.CreateRegisterControl, userCtrl.RegisterExecute))
	e.POST("/verify_email", entry.Hand(userCtrl.CreateVerifyEmailControl, userCtrl.EmailVerifyExecute))
	e.GET("/self", entry.Hand(userCtrl.CreateGetUserControl, userCtrl.GetUserExecute))
	e.POST("/change", entry.Hand(userCtrl.CreateChangeNameControl, userCtrl.ChangeNameExecute))
	e.POST("/change_password", entry.Hand(userCtrl.CreateChangePasswordControl, userCtrl.ChangePasswordExecute))
	e.POST("/change_email", entry.Hand(userCtrl.CreateChangeEmailControl, userCtrl.ChangeEmailExecute))
}
