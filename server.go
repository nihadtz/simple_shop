package main

import (
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/nihadtz/simple_shop/controllers"
	"github.com/nihadtz/simple_shop/services"
	"github.com/tylerb/graceful"
	"github.com/urfave/negroni"
)

var (
	provider controllers.Provider
	users    controllers.Users
)

func main() {
	services.NewLogger()
	services.NewRenderer()
	services.NewAccess()

	provider.SetRBAC("conf/rbac.conf", "conf/rbac.policy")

	var PORT = "3000"

	server := negroni.Classic()
	mux := httprouter.New()

	recovery := negroni.NewRecovery()
	recovery.PrintStack = false
	server.Use(recovery)

	server.UseHandler(mux)

	mux.POST("/register", users.Create)
	mux.GET("/user/:id", users.Get)

	mux.GET("/", provider.RenderSomething)

	server.Run(":" + PORT)
	graceful.Run(":"+PORT, 10*time.Second, server)
}
