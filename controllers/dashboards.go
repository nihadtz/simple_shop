package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/nihadtz/simple_shop/models"
	"github.com/nihadtz/simple_shop/services"
)

type Dashboards struct {
}

func (d Dashboards) ListAdminDashboard(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	dashboardProducts, err := models.ListAdminDashboard()

	if err != nil {
		services.LogError("Error reading DashboardItems", err)
		services.Renderer.Error(res, http.StatusInternalServerError, err.Error())
	}

	services.Renderer.Render(res, http.StatusOK, dashboardProducts)
}

func (d Dashboards) ListPublicDashboard(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	products, err := models.ListProducts()

	if err != nil {
		services.LogError("Error reading Products", err)
		services.Renderer.Error(res, http.StatusInternalServerError, err.Error())
	}

	popularProducts, err := models.PopularProducts()

	if err != nil {
		services.LogError("Error reading Popular Products", err)
		services.Renderer.Error(res, http.StatusInternalServerError, err.Error())
	}

	services.Renderer.Render(res, http.StatusOK, map[string]interface{}{"products": products, "popular_products": popularProducts})
}
