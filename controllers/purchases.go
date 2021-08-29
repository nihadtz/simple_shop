package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/nihadtz/simple_shop/models"
	"github.com/nihadtz/simple_shop/services"
)

type Purchases struct {
}

func (pc Purchases) Get(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	user := context.Get(req, "user").(models.User)
	purchaseID, _ := strconv.Atoi(params.ByName("id"))

	var purchase models.Purchase

	err := purchase.Get(purchaseID)

	if user.ID != purchase.UserID && user.Type != "admin" {
		err := errors.New("you are not allowed to view this purchase")
		services.LogError("Error", err)
		services.Renderer.Error(res, http.StatusForbidden, err.Error())
		return
	}

	if err != nil {
		services.LogError("Error reading Purchase", err)
		services.Renderer.Error(res, http.StatusInternalServerError, err.Error())
	}

	services.Renderer.Render(res, http.StatusOK, purchase)
}

func (pc Purchases) List(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	user := context.Get(req, "user").(models.User)
	var filter models.PurchaseFilter

	if user.Type != "admin" {
		filter.UserID = user.ID
	}

	purchases, err := models.ListPurchases(filter)

	if err != nil {
		services.LogError("Error reading Purchases", err)
		services.Renderer.Error(res, http.StatusInternalServerError, err.Error())
	}

	services.Renderer.Render(res, http.StatusOK, purchases)
}

func (pc Purchases) Create(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	var purchase models.Purchase
	user := context.Get(req, "user").(models.User)

	err := json.NewDecoder(req.Body).Decode(&purchase)

	if err != nil {
		services.LogError("Error Unmarshal Purchase", err)
		services.Renderer.Error(res, http.StatusInternalServerError, err.Error())
	}

	if purchase.ProductID == 0 {
		err = errors.New("ProductID is not set")
		services.LogError("Error", err)
		services.Renderer.Error(res, http.StatusBadRequest, err.Error())
		return
	}

	if purchase.Quantity == 0 {
		err = errors.New("quantity is not set")
		services.LogError("Error", err)
		services.Renderer.Error(res, http.StatusBadRequest, err.Error())
		return
	}

	var product models.Product

	err = product.Get(purchase.ProductID)

	if err != nil {
		services.LogError("Error Reading Purchased Products", err)
		services.Renderer.Error(res, http.StatusInternalServerError, err.Error())
		return
	}

	purchase.Total = product.Price * float64(purchase.Quantity)
	purchase.Date = time.Now().Unix()
	purchase.UserID = user.ID
	purchase.Status = "new"

	err = purchase.Create()

	if err != nil {
		services.LogError("Error Creating Purchase", err)
		services.Renderer.Error(res, http.StatusInternalServerError, err.Error())
		return
	}

	services.Renderer.Render(res, http.StatusOK, purchase)
}
