package model

import "dslfluid.com/dsl/model/sph"
import "dslfluid.com/dsl/model/field"
import V "dslfluid.com/dsl/math/math32"
import G "dslfluid.com/dsl/geom/mesh"
import "dslfluid.com/dsl/sphmethod"

//SPHFluid - Highest Level Fluid Model Abstraction
type SPHFluid struct {
	Mesh   G.Mesh
	Fluid  sph.SPHCore
	Method sphmethod.SPHMethod
}
