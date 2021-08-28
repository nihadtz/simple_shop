package controllers

import (
	"net/http"
	"os"

	"github.com/casbin/casbin"
	"github.com/julienschmidt/httprouter"
	"github.com/nihadtz/simple_shop/services"
)

type Provider struct {
	rules *casbin.Enforcer
}

func (m *Provider) SetRBAC(model, policy string) error {
	path, err := os.Getwd()

	if err != nil {
		return err
	}

	m.rules, err = casbin.NewEnforcerSafe(path+model, path+policy, false)

	return err
}

func (m Provider) RenderSomething(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	services.Renderer.Render(res, http.StatusOK, map[string]interface{}{"api": "Hello world"})
}
