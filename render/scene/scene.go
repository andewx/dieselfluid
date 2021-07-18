package scene //dslfluid.com/dsl/render

import "io/ioutil"
import "log"
import "dslfluid.com/dsl/gltf"
import "dslfluid.com/dsl/render/camera"
import "fmt"
import "dslfluid.com/dsl/math/math32"

/* dslfluid scenes are handled and managed by their GLTF Schema */

//DSelector represents active referenced indices in the GLTF root
//these indices are the top level elements available to the glTF array

type DSLScene struct {
	Root      *gltf.GlTF
	Filepath  string
	Buffers   []*[]byte
	BaseURI   string
	Cam       camera.Camera
	Model     math32.Mat4
	View      math32.Mat4
	RotX      math32.Mat3
	RotY      math32.Mat3
	Rot0      math32.Mat4
	Rot1      math32.Mat4
	RotOrigin math32.Vec
}

//InitDSLScene Creates empty DSLScene struct
func InitDSLScene(filepath string, w float32, h float32) DSLScene {
	model := math32.Identity4()
	view := math32.Identity4()
	rotX := math32.Identity3()
	rotY := math32.Identity3()
	rot0 := math32.Identity4()
	rot1 := math32.Identity4()
	orig := math32.Vec{}
	return DSLScene{&gltf.GlTF{}, filepath, nil, "/Users/briananderson/go/src/github.com/andewx/dslfluid/resources/",
		camera.InitCamera(w, h), model, view, rotX, rotY, rot0, rot1, orig}
}

/*
ImportGLTF () reads scene graph of the GLTF node object and sets up a few properties
for the scene such as the camera and Buffers
*/
func (scene *DSLScene) ImportGLTF() error {

	content, err := ioutil.ReadFile(scene.BaseURI + scene.Filepath)
	if err != nil {
		log.Fatal(err)
	}
	scene.Root.UnmarshalJSON(content)
	scene.Buffers = make([]*[]byte, len(scene.Root.Buffers))

	//----------------------Retrieve URIS --------------------------
	buffers := scene.Root.Buffers
	for i := 0; i < len(buffers); i++ {
		uri := buffers[i].Uri
		bLength := buffers[i].ByteLength
		bErr := scene.LoadURIBuffer(scene.BaseURI+uri, i, bLength) //Need specify if URI is absolute or relative or deconstruct filepath
		if bErr != nil {
			return fmt.Errorf("Unable to load URI\nError Message %s", bErr.Error())
		} else {
			fmt.Printf("Loaded %d kB from URI...\n", bLength/1024)
		}
	}

	//-------------------Load In Camera Node Property------------------------//
	nodes := scene.GetNodes()
	for i := 0; i < len(nodes); i++ {
		thisNode := nodes[i]
		if thisNode.Name == "camera" {
			scene.Cam.Node = thisNode
		}
	}

	return nil
}

func (scene *DSLScene) LoadURIBuffer(uri string, bufferIndex int, bufferLength int) error {

	content, err := ioutil.ReadFile(uri)
	if err != nil {
		return fmt.Errorf("Buffer unavailable, check URI\n")
	}

	if scene.Buffers != nil {
		scene.Buffers[bufferIndex] = &content
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

	ioutil.WriteFile(scene.Filepath, content, 0666)

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
func (scene *DSLScene) GetBufferDataIx(buffer_view_index int) (*[]byte, error) {
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

func (scene *DSLScene) GetNodeChildren(index int) ([]int, error) {
	node, err := scene.GetNodeIx(index)

	if err != nil {
		return nil, err
	}
	return node.Children, nil
}
