package model

import "github.com/andewx/dieselfluid/model/sph"
import "github.com/andewx/dieselfluid/model/field"
import V "github.com/andewx/dieselfluid/math/mgl"
import G "github.com/andewx/dieselfluid/geom/mesh"
import "github.com/andewx/dieselfluid/sphmethod"

//SPHFluid - Highest Level Fluid Model Abstraction
type SPHFluid struct {
	Mesh   G.Mesh
	Fluid  sph.SPHCore
	Method sphmethod.SPHMethod
}
