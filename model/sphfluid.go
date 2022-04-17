package model

import (
	"github.com/andewx/dieselfluid/model/sph"

	G "github.com/andewx/dieselfluid/geom/mesh"
	"github.com/andewx/dieselfluid/sphmethod"
)

//SPHFluid - Highest Level Fluid Model Abstraction
type SPHFluid struct {
	Mesh   G.Mesh
	Fluid  sph.SPHCore
	Method sphmethod.SPHMethod
}
