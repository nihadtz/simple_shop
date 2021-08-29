package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/nihadtz/simple_shop/models"
	"github.com/nihadtz/simple_shop/services"
)

type Products struct {
}

func (pc Products) Get(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	productID, _ := strconv.Atoi(params.ByName("id"))

	var product models.Product

	err := product.Get(productID)

	if err != nil {
		services.LogError("Error reading Product", err)
		services.Renderer.Error(res, http.StatusInternalServerError, err.Error())
	}

	services.Renderer.Render(res, http.StatusOK, product)
}

func (pc Products) List(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

	products, err := models.ListProducts()

	if err != nil {
		services.LogError("Error reading Products", err)
		services.Renderer.Error(res, http.StatusInternalServerError, err.Error())
	}

	services.Renderer.Render(res, http.StatusOK, products)
}

func (pc Products) Create(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	var product models.Product

	err := json.NewDecoder(req.Body).Decode(&product)

	if err != nil {
		services.LogError("Error Unmarshal Product", err)
		services.Renderer.Error(res, http.StatusInternalServerError, err.Error())
	}

	if len(product.Name) == 0 {
		err = errors.New("Name is not set")
		services.LogError("Error", err)
		services.Renderer.Error(res, http.StatusBadRequest, err.Error())
		return
	}

	if len(product.Description) == 0 {
		err = errors.New("Descriptionme is not set")
		services.LogError("Error", err)
		services.Renderer.Error(res, http.StatusBadRequest, err.Error())
		return
	}

	if len(product.Category) == 0 {
		err = errors.New("Category is not set")
		services.LogError("Error", err)
		services.Renderer.Error(res, http.StatusBadRequest, err.Error())
		return
	}

	err = product.Create()

	if err != nil {
		services.LogError("Error Creating Product", err)
		services.Renderer.Error(res, http.StatusInternalServerError, err.Error())
		return
	}

	services.Renderer.Render(res, http.StatusOK, product)
}

func (pc Products) Update(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	var product models.Product

	err := json.NewDecoder(req.Body).Decode(&product)

	if err != nil {
		services.LogError("Error Unmarshal Product", err)
		services.Renderer.Error(res, http.StatusInternalServerError, err.Error())
	}

	if len(product.Name) == 0 {
		err = errors.New("Name is not set")
		services.LogError("Error", err)
		services.Renderer.Error(res, http.StatusBadRequest, err.Error())
		return
	}

	if len(product.Description) == 0 {
		err = errors.New("Descriptionme is not set")
		services.LogError("Error", err)
		services.Renderer.Error(res, http.StatusBadRequest, err.Error())
		return
	}

	if len(product.Category) == 0 {
		err = errors.New("Category is not set")
		services.LogError("Error", err)
		services.Renderer.Error(res, http.StatusBadRequest, err.Error())
		return
	}

	err = product.Update()

	if err != nil {
		services.LogError("Error updating Product", err)
		services.Renderer.Error(res, http.StatusInternalServerError, err.Error())
		return
	}

	services.Renderer.Render(res, http.StatusOK, product)
}
