package scene //dslfluid.com/dsl/render

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/andewx/dieselfluid/common"
	"github.com/andewx/dieselfluid/gltf"
)

/* dslfluid scenes are handled and managed by their GLTF Schema */

//DSelector represents active referenced indices in the GLTF root
//these indices are the top level elements available to the glTF array

type Scene struct {
	Root     *gltf.GlTF
	Filepath string
	Buffers  [][]byte
	BaseURI  string
}

//InitScene Creates empty Scene struct
func InitScene(filepath string) (Scene, error) {
	scn := Scene{}
	scn.Root = &gltf.GlTF{}
	scn.Filepath = filepath
	scn.BaseURI = strings.Clone(common.Cd(filepath) + "/")
	err := scn.ImportGLTF()
	return scn, err
}

/*
ImportGLTF () reads scene graph of the GLTF node object and sets up a few properties
for the scene such as the camera and Buffers
*/
func (scene *Scene) ImportGLTF() error {

	content, err := ioutil.ReadFile(scene.Filepath)
	if err != nil {
		fmt.Printf("Unable to load GLTF File\n")
		log.Fatal(err)
	}
	jsonErr := scene.Root.UnmarshalJSON(content)

	if jsonErr != nil {
		fmt.Println(jsonErr)
		return jsonErr
	}
	scene.Buffers = make([][]byte, len(scene.Root.Buffers))

	//----------------------Retrieve URIS --------------------------
	buffers := scene.Root.Buffers
	for i := 0; i < len(buffers); i++ {
		uri := buffers[i].Uri
		bLength := buffers[i].ByteLength
		bErr := scene.LoadURIBuffer(scene.BaseURI+uri, i, bLength) //Need specify if URI is absolute or relative or deconstruct filepath
		if bErr != nil {
			fmt.Printf("Unable to load Buffer URI\n")
			return fmt.Errorf("Unable to load Buffer URI\nError Message %s", bErr.Error())
		}
	}

	//-------------------Load In Camera Node Property------------------------//
	nodes := scene.GetNodes()
	for i := 0; i < len(nodes); i++ {
		thisNode := nodes[i]
		if thisNode.Name == "camera" {
			//	scene.Cam.Node = thisNode
		}
	}

	return nil
}

func (scene *Scene) LoadURIBuffer(uri string, bufferIndex int, bufferLength int) error {
	content, err := ioutil.ReadFile(uri)
	if err != nil {
		fmt.Printf("Buffer URI Unavailable\n")
		return fmt.Errorf("Buffer unavailable, check URI\n")
	}

	if scene.Buffers != nil {
		scene.Buffers[bufferIndex] = make([]byte, bufferLength)
		scene.Buffers[bufferIndex] = content
	}

	if bufferLength != len(content) {
		return fmt.Errorf("Buffer size mismatch\n")
	}

	return nil
}

/*ExportGLTF () marshals scene graph of the GLTF node object*/
func (scene *Scene) ExportGLTF() error {
	content, err := scene.Root.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(scene.BaseURI+"/"+scene.Filepath, content, 0666)

	//Write Buffers / Images
	return nil
}

func (scene *Scene) Info() error {
	content, err := scene.Root.MarshalJSON()
	if err != nil {
		return err
	}
	fmt.Printf("\n%s\n", content)
	return err
}

//-----------------------------Scene Graph API--------------------------------//
// Mostly shortcuts for getting GLTF data from the nodes and their links

//GetScene () gets default scene node
func (scene *Scene) GetDefaultScene() int {
	return scene.Root.Scene
}

//GetScenes get the scene array
func (scene *Scene) GetScenes() []*gltf.Scene {
	return scene.Root.Scenes
}

//GetScenes get the scene array
func (scene *Scene) GetTextures() []*gltf.Texture {
	return scene.Root.Textures
}

//GetSceneIx gets the scene object at the specified index
func (scene *Scene) GetSceneIx(index int) (*gltf.Scene, error) {
	ls := len(scene.GetScenes())
	if index < 0 || index >= ls {
		return &gltf.Scene{}, fmt.Errorf("Invalid scene")
	}

	return scene.GetScenes()[index], nil
}

func (scene *Scene) GetAccessors() []*gltf.Accessor {
	return scene.Root.Accessors
}

//GetSceneIx gets the scene object at the specified index
func (scene *Scene) GetAccessorIx(index int) (*gltf.Accessor, error) {
	ls := len(scene.GetAccessors())
	if index < 0 || index >= ls {
		return &gltf.Accessor{}, fmt.Errorf("Invalid accessor index")
	}
	return scene.GetAccessors()[index], nil
}

//GetAccessorBufferView returns and accessor and its associated buffer view as a tuple + the error state
//Error states return empty objects
func (scene *Scene) GetAccessorBufferView(accessor_index int) (*gltf.Accessor, *gltf.BufferView, error) {
	acc, err := scene.GetAccessorIx(accessor_index)

	if err != nil { //Returns empty pair on error
		return &gltf.Accessor{}, &gltf.BufferView{}, err
	}

	bufV, mErr := scene.GetBufferViewIx(acc.BufferView)

	if err != nil { //Returns empty pair on error
		return &gltf.Accessor{}, &gltf.BufferView{}, mErr
	}

	return acc, bufV, nil

}

//GetBuffers returns list of buffers associated with scene
func (scene *Scene) GetBuffers() []*gltf.Buffer {
	return scene.Root.Buffers
}

//GetBufferIx gets the buffer descr at the specified index - note this isn't the actual
//scene buffer data storage component just the description with the URI/Bytelengths
func (scene *Scene) GetBufferIx(index int) (*gltf.Buffer, error) {
	ls := len(scene.GetBuffers())
	if index < 0 || index >= ls {
		return &gltf.Buffer{}, fmt.Errorf("Invalid buffer index %d", index)
	}
	return scene.GetBuffers()[index], nil
}

//GetBufferViews gets the list of buffer views
func (scene *Scene) GetBufferViews() []*gltf.BufferView {
	return scene.Root.BufferViews
}

//GetBufferViewIx gets the buffer descr at the specified index - note this isn't the actual
//scene buffer data storage component just the description with the URI/Bytelengths
func (scene *Scene) GetBufferViewIx(index int) (*gltf.BufferView, error) {
	ls := len(scene.GetBufferViews())
	if index < 0 || index >= ls {
		return &gltf.BufferView{}, fmt.Errorf("Invalid BufferView index %d", index)
	}
	return scene.GetBufferViews()[index], nil
}

//GetBufferDataIx gets the Buffer Views data reference as a slice pointer
func (scene *Scene) GetBufferDataIx(buffer_view_index int) ([]byte, error) {
	bV, err := scene.GetBufferViewIx(buffer_view_index)

	if err != nil { //Returns empty pair on error
		return nil, err
	}

	return scene.Buffers[bV.Buffer], nil
}

//GetMeshes returns list of buffers associated with scene
func (scene *Scene) GetMeshes() []*gltf.Mesh {
	return scene.Root.Meshes
}

//GetBufferIx gets the buffer descr at the specified index - note this isn't the actual
//scene buffer data storage component just the description with the URI/Bytelengths
func (scene *Scene) GetMeshIx(index int) (*gltf.Mesh, error) {
	ls := len(scene.GetMeshes())
	if index < 0 || index >= ls {
		return &gltf.Mesh{}, fmt.Errorf("Invalid mesh index %d", index)
	}
	return scene.GetMeshes()[index], nil
}

//

//GetBufferIx gets the buffer descr at the specified index - note this isn't the actual
//scene buffer data storage component just the description with the URI/Bytelengths
func (scene *Scene) GetMeshPrimitives(index int) ([]*gltf.MeshPrimitive, error) {
	ls := len(scene.GetMeshes())
	if index < 0 || index >= ls {
		return nil, fmt.Errorf("Invalid mesh index %d", index)
	}
	prims := scene.GetMeshes()[index].Primitives
	return prims, nil
}

func (scene *Scene) GetNodes() []*gltf.Node {
	return scene.Root.Nodes
}

func (scene *Scene) GetNodeIx(index int) (*gltf.Node, error) {
	ls := len(scene.GetNodes())
	if index < 0 || index >= ls {
		return &gltf.Node{}, fmt.Errorf("Invalid node index %d", index)
	}
	return scene.GetNodes()[index], nil
}

func (scene *Scene) GetImages() []*gltf.Image {
	return scene.Root.Images
}

func (scene *Scene) GetImageIx(index int) (*gltf.Image, error) {
	ls := len(scene.GetImages())
	if index < 0 || index >= ls {
		return &gltf.Image{}, fmt.Errorf("Invalid node index %d", index)
	}
	return scene.GetImages()[index], nil
}

func (scene *Scene) GetMaterials() []*gltf.Material {
	return scene.Root.Materials
}

func (scene *Scene) GetMaterialIx(index int) (*gltf.Material, error) {
	ls := len(scene.GetMaterials())
	if index < 0 || index >= ls {
		return &gltf.Material{}, fmt.Errorf("Invalid node index %d", index)
	}
	return scene.GetMaterials()[index], nil
}

func (scene *Scene) GetSamplers() []*gltf.Sampler {
	return scene.Root.Samplers
}

func (scene *Scene) GetSamplerIx(index int) (*gltf.Sampler, error) {
	ls := len(scene.GetSamplers())
	if index < 0 || index >= ls {
		return &gltf.Sampler{}, fmt.Errorf("Invalid node index %d", index)
	}
	return scene.GetSamplers()[index], nil
}

func (scene *Scene) GetNodeChildren(index int) ([]int, error) {
	node, err := scene.GetNodeIx(index)

	if err != nil {
		return nil, err
	}
	return node.Children, nil
}
