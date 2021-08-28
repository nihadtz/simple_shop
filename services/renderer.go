package services

import (
	"fmt"
	"net/http"

	"github.com/unrolled/render"
)

var Renderer *RendererCtrl

type RendererCtrl struct {
	r *render.Render
}

func NewRenderer() {
	fmt.Println("Initializing Renderer")

	rend := new(RendererCtrl)

	rend.r = render.New(render.Options{
		IndentJSON: true,
	})

	Renderer = rend
}

func (r *RendererCtrl) Render(res http.ResponseWriter, status int, v interface{}) {
	res.Header().Set("Access-Control-Allow-Origin", "*")

	r.r.JSON(res, status, v)
}
