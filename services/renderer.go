package services

import (
	"fmt"
	"net/http"
	"strconv"

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

func (rend *RendererCtrl) Render(res http.ResponseWriter, status int, v interface{}) {
	res.Header().Set("Access-Control-Allow-Origin", "*")

	rend.r.JSON(res, status, v)
}

func (rend *RendererCtrl) Error(res http.ResponseWriter, status int, err string) {

	if err != "" {
		rend.Render(res, status, map[string]string{"status": strconv.Itoa(status), "error": err})
	} else {
		rend.Render(res, status, nil)
	}
}
