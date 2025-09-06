package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/motojouya/geezer_auth/internal/entry/router"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
	configBehavior "github.com/motojouya/geezer_auth/internal/behavior/config"
)

func main() {
	var env = localPkg.CreateEnvironment()
	serverConf, err := configBehavior.NewServerGet(env).GetServer()
	if err != nil {
		fmt.Println("failed to get server config:", err)
	}

	e := echo.New()
	router.Route(e)
	e.Start(serverConf.GetEchoPort())

	fmt.Println("start server!")
}
