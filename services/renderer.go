package services

import (
	"fmt"

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
