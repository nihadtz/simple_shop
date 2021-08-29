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
	provider   controllers.Provider
	users      controllers.Users
	products   controllers.Products
	purchases  controllers.Purchases
	payments   controllers.Payments
	dashboards controllers.Dashboards
)

func main() {
	services.NewLogger()
	services.NewRenderer()
	services.NewAccess()

	provider.SetRBAC("/conf/rbac.conf", "/conf/rbac.policy")

	var PORT = "3000"

	server := negroni.Classic()
	mux := httprouter.New()

	recovery := negroni.NewRecovery()
	recovery.PrintStack = false

	server.Use(recovery)
	server.Use(negroni.HandlerFunc(provider.JWT))

	server.UseHandler(mux)

	mux.POST("/register", users.Create)
	mux.GET("/user/:id", users.Get)

	mux.POST("/login", users.Login)
	mux.GET("/logout", users.Logout)

	mux.POST("/product", products.Create)
	mux.PUT("/product", products.Update)
	mux.GET("/product/:id", products.Get)
	mux.GET("/products", products.List)

	mux.POST("/purchase", purchases.Create)
	mux.GET("/purchase/:id", purchases.Get)
	mux.GET("/purchases", purchases.List)

	mux.POST("/payment/stripe/:id", payments.ViaStripe)

	mux.GET("/dashboard/admin", dashboards.ListAdminDashboard)
	mux.GET("/dashboard/public", dashboards.ListPublicDashboard)

	server.Run(":" + PORT)
	graceful.Run(":"+PORT, 10*time.Second, server)
}
