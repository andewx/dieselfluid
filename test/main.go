package test

//Main Testing Package Scene Construction and Initialization

import (
	"github.com/andewx/dieselfluid/model/field"
	"github.com/andewx/dieselfluid/model/sph"
	"github.com/andewx/dieselfluid/solver/wcsph"
)

type MainTest struct {
	Fluid sph.SPH
	WCSPH wcsph.WCSPH
	Field field.SPHField
}

func Init() *MainTest {
	appMain := MainTest{}
	return &appMain
}
