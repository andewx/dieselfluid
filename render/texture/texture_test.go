/*
Golang package implementing quaternion math
Purpose is to provide quaternion support under the MIT license as existing
Go quaternion packages are under more restrictive or unspecified licenses.
This project is licensed under the terms of the MIT license.
*/

package texture

import (
	"fmt"
	"testing"
)

func TestLoad(t *testing.T) {

	myLib := NewTexLibrary()
	err := myLib.Load("/../../resources/logo.png", 0)
	if err != nil {
		fmt.Printf("Image load fail\n%s\n", err.Error())
		t.Fail()
	}
	myLib.HasDevice = false

	myLib.CommitTexLibGL()
	myLib.RemoveTexLibGL()

}
