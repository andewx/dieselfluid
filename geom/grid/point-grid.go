package grid

import (
	"fmt"

	V "github.com/andewx/dieselfluid/math/vector"
)

//Defines [I, J, K] element structured unit grid position

type Grid struct {
	Div        V.Vec //Number Divisions Per (Cast int)
	origin     V.Vec //Origin XYZ
	scale      V.Vec //Scaling Vector XYZ
	min_bounds V.Vec //Min Vec
	step       V.Vec //Step Vec
}

//Builds Grid - Scale grid around the origin for a unit cube then translate
func BuildGrid(sclVec V.Vec, transOrigin V.Vec, dimVec V.Vec) (Grid, error) {

	if dimVec[0] == 0.0 || dimVec[1] == 0.0 || dimVec[2] == 0.0 {
		return Grid{}, fmt.Errorf("Dimension passed to BuildGrid() with value of zero")
	}

	mGrid := Grid{dimVec, transOrigin, sclVec, []float32{-1 * sclVec[0], -1 * sclVec[1], -1 * sclVec[2]},
		[]float32{2 * sclVec[0] / dimVec[0], 2 * sclVec[1] / dimVec[1], 2 * sclVec[2] / dimVec[2]}}
	mGrid.min_bounds = V.Add(mGrid.min_bounds, transOrigin)
	return mGrid, nil
}

//BuildKernGrid - Builds a Grid Based on the Kernel Spacing - Returns Grid and Grid Cubic Dimensionality
func BuildKernGrid(trans V.Vec, dimVec V.Vec, kern float32) (Grid, error) {

	mGrid, err := BuildGrid([]float32{1, 1, 1}, trans, dimVec)
	if err != nil {
		return mGrid, err
	}
	inv := []float32{1 / dimVec[0], 1 / dimVec[1], 1 / dimVec[2]}
	mGrid.step = mGrid.min_bounds.Scale(-2.0).Mul(inv)
	return mGrid, nil
}

func (g Grid) Volume() float32 {
	return 2 * g.scale[0] * 2 * g.scale[1] * 2 * g.scale[2]
}

func ijk2Vec(i int, j int, k int) V.Vec {
	return V.Vec{float32(i), float32(j), float32(k)}
}

//Map 3D Index - maps a 3D I,J,K position to 1D flattened array
func (g Grid) Index(i int, j int, k int) int {
	i_w := int(g.Div[0])
	j_w := int(g.Div[1])
	return (k + i_w*(i*j_w+j))
}

//Assigns Position Based on [i][j][k] Grid Position element
func (g Grid) GridPosition(i int, j int, k int) V.Vec {
	p := ijk2Vec(i, j, k)
	return V.Add(g.min_bounds, V.ScaleVar(g.step, p))
}

//Updates the Internal Grid Component
func (g Grid) UpdateGrid(sclVar V.Vec, TransOrigin V.Vec, DimVec V.Vec) Grid {
	//Scale XYZ Components - Regarding Previous Origin
	g.min_bounds = V.Sub(g.min_bounds, g.origin)
	g.min_bounds = V.ScaleVar(g.min_bounds, sclVar)
	MinDouble := V.ScaleVar(g.min_bounds, V.Vec{-2.0, -2.0, -2.0})
	LenSides := V.Sub(g.min_bounds, MinDouble)
	InvDimVec := V.Vec{1 / DimVec[0], 1 / DimVec[1], 1 / DimVec[2]}
	g.step = V.ScaleVar(LenSides, InvDimVec)
	g.min_bounds = V.Add(g.min_bounds, V.Add(g.origin, TransOrigin))
	g.origin = V.Add(g.origin, TransOrigin)
	return g
}
