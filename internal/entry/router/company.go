package router

import (
	"github.com/labstack/echo/v4"
	companyCtrl "github.com/motojouya/geezer_auth/internal/control/company"
	companyUserCtrl "github.com/motojouya/geezer_auth/internal/control/companyUser"
	"github.com/motojouya/geezer_auth/internal/entry"
)

func CompanyRoute(e *echo.Group) {
	e.POST("/create", entry.Hand(companyCtrl.CreateCreateControl, companyCtrl.CreateExecute))
	e.GET("/:identifier", entry.Hand(companyCtrl.CreateGetCompanyControl, companyCtrl.GetCompanyExecute))
	e.GET("/:identifier/user", entry.Hand(companyCtrl.CreateGetUserControl, companyCtrl.GetUserExecute))
	e.POST("/:identifier/invite", entry.Hand(companyUserCtrl.CreateInviteControl, companyUserCtrl.InviteExecute))
	e.POST("/:identifier/accept", entry.Hand(companyUserCtrl.CreateAcceptControl, companyUserCtrl.AcceptExecute))
	e.POST("/:identifier/assign", entry.Hand(companyUserCtrl.CreateAssignControl, companyUserCtrl.AssignExecute))
}
