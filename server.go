package main

import (
	"time"

	"github.com/alecthomas/kingpin"
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

	runas   = kingpin.Flag("runas", "Define deployment type and which configuration profile to read (dev, test, prod...)").Short('r').Default("dev").String()
	migrate = kingpin.Flag("migrate", "Migrate database to level defined in configuration").Short('m').Bool()
)

func init() {
	kingpin.Parse()
}

func main() {

	services.NewConfigurer(*runas, "conf/conf.yaml")
	services.NewLogger()
	services.NewRenderer()

	if *migrate {
		services.MigrateDB()
	}

	services.NewAccess(*runas)
	services.SetRBAC("/conf/rbac.conf", "/conf/rbac.policy")

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

	server.Run(":" + services.Configuration.Port)
	graceful.Run(":"+services.Configuration.Port, 10*time.Second, server)
}
