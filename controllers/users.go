package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/nihadtz/simple_shop/models"
	"github.com/nihadtz/simple_shop/services"
)

type Users struct {
}

func (ac Users) Get(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	user := context.Get(req, "user").(models.User)
	userID, _ := strconv.Atoi(params.ByName("id"))

	var userGet models.User

	err := userGet.Get(userID)

	if err != nil {
		services.LogError("Error reading User", err)
		services.Renderer.Error(res, http.StatusInternalServerError, err.Error())
	}

	if user.ID != userGet.ID && user.Type != "admin" {
		err := errors.New("you are not allowed to view this user")
		services.LogError("Error", err)
		services.Renderer.Error(res, http.StatusForbidden, err.Error())
		return
	}

	services.Renderer.Render(res, http.StatusOK, userGet)
}

func (ac Users) Create(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	var user models.User

	err := json.NewDecoder(req.Body).Decode(&user)

	if err != nil {
		services.LogError("Error Unmarshal User", err)
		services.Renderer.Error(res, http.StatusInternalServerError, err.Error())
	}

	if len(user.Name) == 0 {
		err = errors.New("name is not set")
		services.LogError("Error", err)
		services.Renderer.Error(res, http.StatusBadRequest, err.Error())
		return
	}

	regexpEmail := regexp.MustCompile("^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$")

	if !regexpEmail.MatchString(user.Email) {
		err = errors.New("wrong email format")
		services.LogError("Error ", err)
		services.Renderer.Error(res, http.StatusBadRequest, err.Error())
		return
	}

	if len(user.Password) == 0 || len(user.Password) < 10 {
		err = errors.New("password is not set or less than 10 chars")
		services.LogError("Error", err)
		services.Renderer.Error(res, http.StatusBadRequest, err.Error())
		return
	}

	user.Type = "customer"
	user.Token = ""
	user.Algorithm = ""
	user.Issued = 0

	err = user.Create()

	if err != nil {
		services.LogError("Error Creating User", err)
		services.Renderer.Error(res, http.StatusInternalServerError, err.Error())
		return
	}

	services.Renderer.Render(res, http.StatusOK, user)
}

func (ac Users) Login(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	user := context.Get(req, "user").(models.User)

	services.Renderer.Render(res, http.StatusOK,
		map[string]interface{}{
			"token": user.Token,
		},
	)
}

func (ac Users) Logout(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	user := context.Get(req, "user").(models.User)

	err := user.ClearToken()

	if err != nil {
		services.Renderer.Error(res, http.StatusInternalServerError, err.Error())
	}

	services.Renderer.Render(res, http.StatusOK,
		map[string]interface{}{
			"status": "ok",
			"error":  "",
		},
	)
}
