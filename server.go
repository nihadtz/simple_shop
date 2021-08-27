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
)

func main() {
	services.NewLogger()
	services.NewRenderer()

	provider.SetRBAC("conf/rbac.conf", "conf/rbac.policy")

	const PORT = "3000"
	const IP = "localhost"

	server := negroni.Classic()
	mux := httprouter.New()

	recovery := negroni.NewRecovery()
	recovery.PrintStack = false
	server.Use(recovery)

	server.UseHandler(mux)

	server.Run(IP + ":" + PORT)
	graceful.Run(IP+":"+PORT, 10*time.Second, server)
}
