package test

//Main Testing Package Scene Construction and Initialization

import "github.com/andewx/dieselfluid/model/sph"
import "github.com/andewx/dieselfluid/model/field"
import "github.com/andewx/dieselfluid/sphmethod/wcsph"
import V "github.com/andewx/dieselfluid/math/mgl"
import G "github.com/andewx/dieselfluid/geom/mesh"

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
