package scene //dslfluid.com/dsl/render

import "io/ioutil"
import "log"
import "dslfluid.com/dsl/gltf"
import "fmt"

/* dslfluid scenes are handled and managed by their GLTF Schema */

//DSelector represents active referenced indices in the GLTF root
//these indices are the top level elements available to the glTF array

type DSLScene struct {
	Root     *gltf.GlTF
	Filepath string
	Buffers  [][]byte
	BaseURI  string
}

//InitDSLScene Creates empty DSLScene struct
func InitDSLScene(base string, filepath string, w float32, h float32) DSLScene {
	return DSLScene{&gltf.GlTF{}, filepath, nil, base}
}

/*
ImportGLTF () reads scene graph of the GLTF node object and sets up a few properties
for the scene such as the camera and Buffers
*/
func (scene *DSLScene) ImportGLTF() error {

	content, err := ioutil.ReadFile(scene.BaseURI + scene.Filepath)
	if err != nil {
		fmt.Printf("Unable to load GLTF File\n")
		log.Fatal(err)
	}
	scene.Root.UnmarshalJSON(content)
	scene.Buffers = make([][]byte, len(scene.Root.Buffers))

	//----------------------Retrieve URIS --------------------------
	buffers := scene.Root.Buffers
	for i := 0; i < len(buffers); i++ {
		uri := buffers[i].Uri
		bLength := buffers[i].ByteLength
		bErr := scene.LoadURIBuffer(uri, i, bLength) //Need specify if URI is absolute or relative or deconstruct filepath
		if bErr != nil {
			fmt.Printf("Unable to load Buffer URI\n")
			return fmt.Errorf("Unable to load Buffer URI\nError Message %s", bErr.Error())
		} else {
			fmt.Printf("Loaded %d kB from URI \n", bLength/1024)
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

func (scene *DSLScene) LoadURIBuffer(uri string, bufferIndex int, bufferLength int) error {

	content, err := ioutil.ReadFile(scene.BaseURI + uri)
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
func (scene *DSLScene) ExportGLTF() error {
	content, err := scene.Root.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(scene.BaseURI+scene.Filepath, content, 0666)

	//Write Buffers / Images
	return nil
}

//-----------------------------Scene Graph API--------------------------------//
// Mostly shortcuts for getting GLTF data from the nodes and their links

//GetScene () gets default scene node
func (scene *DSLScene) GetDefaultScene() int {
	return scene.Root.Scene
}

//GetScenes get the scene array
func (scene *DSLScene) GetScenes() []*gltf.Scene {
	return scene.Root.Scenes
}

//GetSceneIx gets the scene object at the specified index
func (scene *DSLScene) GetSceneIx(index int) (*gltf.Scene, error) {
	ls := len(scene.GetScenes())
	if index < 0 || index > ls {
		return &gltf.Scene{}, fmt.Errorf("Invalid scene")
	}

	return scene.GetScenes()[index], nil
}

func (scene *DSLScene) GetAccessors() []*gltf.Accessor {
	return scene.Root.Accessors
}

//GetSceneIx gets the scene object at the specified index
func (scene *DSLScene) GetAccessorIx(index int) (*gltf.Accessor, error) {
	ls := len(scene.GetAccessors())
	if index < 0 || index > ls {
		return &gltf.Accessor{}, fmt.Errorf("Invalid accessor index")
	}
	return scene.GetAccessors()[index], nil
}

//GetAccessorBufferView returns and accessor and its associated buffer view as a tuple + the error state
//Error states return empty objects
func (scene *DSLScene) GetAccessorBufferView(accessor_index int) (*gltf.Accessor, *gltf.BufferView, error) {
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
func (scene *DSLScene) GetBuffers() []*gltf.Buffer {
	return scene.Root.Buffers
}

//GetBufferIx gets the buffer descr at the specified index - note this isn't the actual
//scene buffer data storage component just the description with the URI/Bytelengths
func (scene *DSLScene) GetBufferIx(index int) (*gltf.Buffer, error) {
	ls := len(scene.GetBuffers())
	if index < 0 || index > ls {
		return &gltf.Buffer{}, fmt.Errorf("Invalid buffer index %d", index)
	}
	return scene.GetBuffers()[index], nil
}

//GetBufferViews gets the list of buffer views
func (scene *DSLScene) GetBufferViews() []*gltf.BufferView {
	return scene.Root.BufferViews
}

//GetBufferViewIx gets the buffer descr at the specified index - note this isn't the actual
//scene buffer data storage component just the description with the URI/Bytelengths
func (scene *DSLScene) GetBufferViewIx(index int) (*gltf.BufferView, error) {
	ls := len(scene.GetBufferViews())
	if index < 0 || index > ls {
		return &gltf.BufferView{}, fmt.Errorf("Invalid BufferView index %d", index)
	}
	return scene.GetBufferViews()[index], nil
}

//GetBufferDataIx gets the Buffer Views data reference as a slice pointer
func (scene *DSLScene) GetBufferDataIx(buffer_view_index int) ([]byte, error) {
	bV, err := scene.GetBufferViewIx(buffer_view_index)

	if err != nil { //Returns empty pair on error
		return nil, err
	}

	return scene.Buffers[bV.Buffer], nil
}

//GetMeshes returns list of buffers associated with scene
func (scene *DSLScene) GetMeshes() []*gltf.Mesh {
	return scene.Root.Meshes
}

//GetBufferIx gets the buffer descr at the specified index - note this isn't the actual
//scene buffer data storage component just the description with the URI/Bytelengths
func (scene *DSLScene) GetMeshIx(index int) (*gltf.Mesh, error) {
	ls := len(scene.GetMeshes())
	if index < 0 || index > ls {
		return &gltf.Mesh{}, fmt.Errorf("Invalid mesh index %d", index)
	}
	return scene.GetMeshes()[index], nil
}

//

//GetBufferIx gets the buffer descr at the specified index - note this isn't the actual
//scene buffer data storage component just the description with the URI/Bytelengths
func (scene *DSLScene) GetMeshPrimitives(index int) ([]*gltf.MeshPrimitive, error) {
	ls := len(scene.GetMeshes())
	if index < 0 || index > ls {
		return nil, fmt.Errorf("Invalid mesh index %d", index)
	}
	prims := scene.GetMeshes()[index].Primitives
	return prims, nil
}

func (scene *DSLScene) GetNodes() []*gltf.Node {
	return scene.Root.Nodes
}

func (scene *DSLScene) GetNodeIx(index int) (*gltf.Node, error) {
	ls := len(scene.GetNodes())
	if index < 0 || index > ls {
		return &gltf.Node{}, fmt.Errorf("Invalid node index %d", index)
	}
	return scene.GetNodes()[index], nil
}

func (scene *DSLScene) GetImages() []*gltf.Image {
	return scene.Root.Images
}

func (scene *DSLScene) GetImageIx(index int) (*gltf.Image, error) {
	ls := len(scene.GetImages())
	if index < 0 || index > ls {
		return &gltf.Image{}, fmt.Errorf("Invalid node index %d", index)
	}
	return scene.GetImages()[index], nil
}

func (scene *DSLScene) GetMaterials() []*gltf.Material {
	return scene.Root.Materials
}

func (scene *DSLScene) GetMaterialIx(index int) (*gltf.Material, error) {
	ls := len(scene.GetMaterials())
	if index < 0 || index > ls {
		return &gltf.Material{}, fmt.Errorf("Invalid node index %d", index)
	}
	return scene.GetMaterials()[index], nil
}

func (scene *DSLScene) GetSamplers() []*gltf.Sampler {
	return scene.Root.Samplers
}

func (scene *DSLScene) GetSamplerIx(index int) (*gltf.Sampler, error) {
	ls := len(scene.GetSamplers())
	if index < 0 || index > ls {
		return &gltf.Sampler{}, fmt.Errorf("Invalid node index %d", index)
	}
	return scene.GetSamplers()[index], nil
}

func (scene *DSLScene) GetNodeChildren(index int) ([]int, error) {
	node, err := scene.GetNodeIx(index)

	if err != nil {
		return nil, err
	}
	return node.Children, nil
}
