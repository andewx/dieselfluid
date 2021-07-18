package scene

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"testing"
)

func TestGLTFImport(t *testing.T) {

	myScene := InitDSLScene("Minimal.gltf", 1024.0, 720.0)
	err := myScene.ImportGLTF()

	if err != nil {
		t.Errorf("Failed to load GLTF file from %s\n", myScene.Filepath)
	} else {

		//----------------------Test Scene type int --------------------------
		scnIdx := myScene.GetDefaultScene()
		if scnIdx != 0 { //Read in Scene Root
			t.Errorf("Scene Index Should be 0")
		}

		//----------------------Test Asset type Meta --------------------------
		gen := myScene.Root.Asset.Generator
		ver := myScene.Root.Asset.Version
		if gen != "Khronos glTF Blender I/O v1.6.16" {
			t.Errorf("Krohnos Asset")
		}

		if ver != "2.0" {
			t.Errorf("Version should be 2.0!\n")
		}
		//----------------------Test Accessors type Accessor[] --------------------------
		for i := 0; i < len(myScene.Root.Accessors); i++ {
			myAccessor, err := myScene.GetAccessorIx(i)

			if err != nil {
				t.Errorf("Error retrieving accessor at %d", i)
			}
			//Print out accessor information
			bview := myAccessor.BufferView
			ctype := myAccessor.ComponentType
			cnt := myAccessor.Count
			tp := myAccessor.Type
			if tp == "VEC3" {
				//	max := myAccessor.Max.([3]float32)
				//	min := myAccessor.Min.([3]float32)
			}
			fmt.Printf("Accessor(%d)::BfrView(%d)\n", i, bview)
			if ctype < 0 {
				t.Errorf("Component Type should be gl.GL_FLOAT(%d)\ngl.INT(%d)\ngl.UNSIGNED_INT(%d)\n", gl.FLOAT, gl.INT, gl.UNSIGNED_INT)
			}
			if cnt < 0 {
				t.Errorf("Count should be > 0 and an int\n")
			}
		}

		//----------------------Test BufferViews []*BufferView --------------------------
		for i := 0; i < len(myScene.Root.BufferViews); i++ {
			myBuffer, err := myScene.GetBufferViewIx(i)

			if err != nil {
				t.Errorf("Error retrieving buffer view\n")
			}
			bufferRef := myBuffer.Buffer
			bufferLength := myBuffer.ByteLength
			byteOffset := myBuffer.ByteOffset

			if &bufferRef == nil {
				t.Errorf("Buffer Reference could not be retrieved")
			}
			if &bufferLength == nil {
				t.Errorf("Buffer Reference could not be retrieved\n")
			}
			if &byteOffset == nil {
				t.Errorf("Byte offset could not be retrieved\n")
			}

		}

		//----------------------Images--------------------------
		if myScene.Root.Images == nil {
			fmt.Printf("Images passed\n")
		}

		//----------------------Caeras--------------------------
		if myScene.Root.Cameras == nil {
			fmt.Printf("Cameras passed\n")
		}

		//----------------------Test Meshes--------------------------
		for i := 0; i < len(myScene.Root.Meshes); i++ {
			mesh, err := myScene.GetMeshIx(i)

			if err != nil {
				t.Errorf("Error retrieving mesh")
			}
			if i == 0 {
				if mesh.Name != "Cube" {
					t.Errorf("Cube mesh not found\n")
				}
			} //Go through primitives

			prims, err := myScene.GetMeshPrimitives(i)
			for j := 0; j < len(prims); j++ {
				indices := prims[j].Indices
				mat := prims[j].Material
				attr := prims[j].Attributes
				if i == 0 {
					if indices != 3 {
						t.Errorf("Cube mesh indices buffer not found\n")
					}
					if mat != 0 {
						t.Errorf("Cube mesh material not found\n")
					}
					if attr["POSITION"] != 0 {
						t.Errorf("Cube mesh primitive attributes map not well formed\n")
					}
				}
			}

		}

		/*---------------Test Materials -----------------*/
		for i := 0; i < len(myScene.Root.Materials); i++ {
			myMat := myScene.Root.Materials[i]
			name := myMat.Name
			pbrMat := myMat.PbrMetallicRoughness
			metalFac := pbrMat.MetallicFactor
			roughFac := pbrMat.RoughnessFactor
			baseColor := pbrMat.BaseColorFactor

			if i == 0 {
				if metalFac != 0 {
					t.Errorf("PBR Material sucks\n")
					return
				}

				if roughFac != 0.8045452833175659 {
					t.Errorf("PBR Material sucks\n")
					return
				}

				if name != "White Clay" {
					t.Errorf("PBR Material sucks\n")
					return
				}

				if baseColor[0] != 0.800000011920929 {
					t.Errorf("PBR Material sucks\n")
					return
				}
			}
		}

		fmt.Printf("Materials Passed (still need to check Color factor)\n")

		//STILL NEED TO CHECK NODES /SAMPLERS /SCENES /SKINS / TEXTURES

		nodes := myScene.GetNodes()
		for i := 0; i < len(nodes); i++ {
			node := nodes[i]
			mesh := node.Mesh
			name := node.Name
			scale := node.Scale
			Trans := node.Translation
			rot := node.Rotation

			if i == 0 {
				if mesh != 0 {
					t.Errorf("Wrong mesh in the node ")
				}
				if name != "Cube" {
					t.Errorf("Node test didn't pass for string name cube")
					return
				}
				if scale[0] != 0.41879701614379883 {
					t.Errorf("Node test didn't pass for scale val")
					return
				}
				if Trans[0] != -2.91190242767334 {
					t.Errorf("Node test didn't pass for trans val")
					return
				}
				if rot != nil {
					t.Errorf("Rot shouldn't exist\n")
					return
				}
			}
		}

		for i := 0; i < len(myScene.Root.Scenes); i++ {
			scene := myScene.Root.Scenes[i]
			nodes := scene.Nodes
			if len(nodes) != 8 {
				t.Errorf("Not enough nodes in the scene")
				return
			}

			fmt.Printf("Scenes passed\n")
		}

		for i := 0; i < len(myScene.Root.Textures); i++ {
			tex := myScene.Root.Textures[i]
			if tex != nil {
				t.Errorf("Shouldn't be any textures")
			}
		}

		for i := 0; i < len(myScene.Root.Samplers); i++ {
			sampler := myScene.Root.Samplers[i]
			if sampler != nil {
				t.Errorf("Shoulnt be an samplers")
			}
		}

		fmt.Printf("------TEST:GLTF Importation Passed---------\n")
	} // End test internal

} //End Func

func TestGLTFBuffer(t *testing.T) {

}

func TestGLTFImage(t *testing.T) {

}

func TestFormatDSLTF(t *testing.T) {

}

func TestGLTFWrite(t *testing.T) {

}

func TestDSLTFModWrite(t *testing.T) {

}

func TestDSLTFExtension(t *testing.T) {

}

func TestAccessorBufferFeatures(t *testing.T) {

}

func TestMeshFeatures(t *testing.T) {

}

func TestNodeFeatures(t *testing.T) {

}

func TestCamera(t *testing.T) {

}

func TestSkinning(t *testing.T) {

}

func TestAnimation(t *testing.T) {

}
