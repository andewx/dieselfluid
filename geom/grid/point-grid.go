package grid

import V "github.com/andewx/dieselfluid/math/vector"

//Defines [I, J, K] element structured unit grid position

type Grid struct {
	Origin   V.Vec //Origin XYZ
	ScaleXYZ V.Vec //Scaling Vector XYZ
	DimXYZ   V.Vec //Number Divisions Per (Cast int)
	MinXYZ   V.Vec //Min Vec
	StepXYZ  V.Vec //Step Vec
}

//Builds Grid - Scale grid around the origin for a unit cube then translate
func BuildGrid(sclVec V.Vec, TransOrigin V.Vec, DimVec V.Vec) Grid {
	MyGrid := Grid{}
	MyGrid.Origin = TransOrigin
	MyGrid.MinXYZ = V.Vec{float32(-1.0), float32(-1.0), float32(-1.0)}
	MyGrid.MinXYZ = V.ScaleVar(MyGrid.MinXYZ, sclVec)
	MinDouble := V.ScaleVar(MyGrid.MinXYZ, V.Vec{-2.0, -2.0, -2.0})
	LenSides := V.Sub(MyGrid.MinXYZ, MinDouble)
	InvDimVec := V.Vec{1 / DimVec[0], 1 / DimVec[1], 1 / DimVec[2]}
	MyGrid.StepXYZ = V.ScaleVar(LenSides, InvDimVec)
	MyGrid.MinXYZ = V.Add(MyGrid.MinXYZ, TransOrigin)
	return MyGrid
}

//BuildKernGrid - Builds a Grid Based on the Kernel Spacing - Returns Grid and Grid Cubic Dimensionality
func BuildKernGrid(sclVar V.Vec, TransOrigin V.Vec, kern float32) (Grid, int) {
	MyGrid := Grid{}
	MyGrid.Origin = TransOrigin
	MyGrid.MinXYZ = V.Vec{float32(-1.0), float32(-1.0), float32(-1.0)}
	MyGrid.MinXYZ = V.ScaleVar(MyGrid.MinXYZ, sclVar)
	MinDouble := V.ScaleVar(MyGrid.MinXYZ, V.Vec{-2.0, -2.0, -2.0})
	LenSides := V.Sub(MyGrid.MinXYZ, MinDouble)
	DimVec := V.Vec{LenSides[0] / kern, LenSides[1] / kern, LenSides[2] / kern}
	Dimensionality := int(LenSides[0] / kern)
	InvDimVec := V.Vec{1 / DimVec[0], 1 / DimVec[1], 1 / DimVec[2]}
	MyGrid.StepXYZ = V.ScaleVar(LenSides, InvDimVec)
	MyGrid.MinXYZ = V.Add(MyGrid.MinXYZ, TransOrigin)
	return MyGrid, Dimensionality
}

func ijk2Vec(i int, j int, k int) V.Vec {
	return V.Vec{float32(i), float32(j), float32(k)}
}

//Assigns Position Based on [i][j][k] Grid Position element
func (g Grid) GridPosition(i int, j int, k int) V.Vec {
	p := ijk2Vec(i, j, k)
	return V.Add(g.MinXYZ, V.ScaleVar(g.StepXYZ, p))
}

//Updates the Internal Grid Component
func (g Grid) UpdateGrid(sclVar V.Vec, TransOrigin V.Vec, DimVec V.Vec) Grid {
	//Scale XYZ Components - Regarding Previous Origin
	g.MinXYZ = V.Sub(g.MinXYZ, g.Origin)
	g.MinXYZ = V.ScaleVar(g.MinXYZ, sclVar)
	MinDouble := V.ScaleVar(g.MinXYZ, V.Vec{-2.0, -2.0, -2.0})
	LenSides := V.Sub(g.MinXYZ, MinDouble)
	InvDimVec := V.Vec{1 / DimVec[0], 1 / DimVec[1], 1 / DimVec[2]}
	g.StepXYZ = V.ScaleVar(LenSides, InvDimVec)
	g.MinXYZ = V.Add(g.MinXYZ, V.Add(g.Origin, TransOrigin))
	g.Origin = V.Add(g.Origin, TransOrigin)
	return g
}
