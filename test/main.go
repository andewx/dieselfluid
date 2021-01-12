package test

//Main Testing Package Scene Construction and Initialization

import "dslfluid.com/dsl/model/sph"
import "dslfluid.com/dsl/model/field"
import "dslfluid.com/dsl/sphmethod/wcsph"
import V "dslfluid.com/dsl/math/math32"
import G "dslfluid.com/dsl/geom/mesh"

type MainTest struct {
	Mesh  G.Mesh
	Fluid sph.SPHCore
	WCSPH wcsph.WCSPH
	Field field.SPHField
}

func Init() {
	appMain := MainTest{}
}

func Run() {

}
