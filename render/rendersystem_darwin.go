package render

import (
	"fmt"
	"math" //GoLang Entity Component System
	"runtime"
	"time"
	"unsafe"

	"github.com/EngoEngine/ecs"
	"github.com/andewx/dieselfluid/common"
	"github.com/andewx/dieselfluid/geom/mesh"
	"github.com/andewx/dieselfluid/math/matrix"
	"github.com/andewx/dieselfluid/math/vector"
	"github.com/andewx/dieselfluid/render/defs"
	"github.com/andewx/dieselfluid/render/glr"
	"github.com/andewx/dieselfluid/render/scene"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	MIN_PARAM_SIZE  = 20
	BYTES_PER_FLOAT = 4
	COORD_4         = 4
	COORD_3         = 3
	COORD_2         = 2
	INITIALIZED     = 1
)

//RenderSystem provides ECS Render Mechanism
type RenderSystem struct {
	MyRenderer         *glr.GLRenderer
	MeshEntities       []*defs.MeshEntity
	Shaders            defs.ShadersMap
	Graph              *scene.Scene
	Textures           defs.TexturesMap
	Materials          map[string]*defs.Material
	MaterialIndexNames []*string
	CubeMap            uint32
	Light              defs.Light
	Time               time.Time
	Angle              float32
	num_vao            int
	num_vbo            int
	num_tex            int
}

func byteSliceToFloat32Slice(src []byte) []float32 {
	if len(src) == 0 {
		return nil
	}

	l := len(src) / 4
	ptr := unsafe.Pointer(&src[0])
	// It is important to keep in mind that the Go garbage collector
	// will not interact with this data, and that if src if freed,
	// the behavior of any Go code using the slice is nondeterministic.
	// Reference: https://github.com/golang/go/wiki/cgo#turning-c-arrays-into-go-slices
	return (*[1 << 26]float32)((*[1 << 26]float32)(ptr))[:l:l]
}

//Gather() - Package customized initialization routine function. Gather is an
//initialization like function moniker used package by package to initialize
//a main package construct for API usage. Here Gather returns the main package
//type which is the RenderSystem and initializes with a GLTF scene file name
//located in the "dslfluid.com/resources/" directory
func Init(scn string) (RenderSystem, error) {
	graph, err := scene.InitScene(scn)
	mRenderer := glr.Renderer()
	var MainRenderer RenderSystem

	if err != nil {
		return MainRenderer, err
	}
	shaderDict := defs.ShadersMap{}
	MainRenderer = RenderSystem{}
	MainRenderer.MyRenderer = mRenderer
	MainRenderer.Graph = &graph
	MainRenderer.Shaders = shaderDict
	MainRenderer.Light = defs.Light{Pos: []float32{50, 50, 50, 1}, Color: []float32{2.0, 2.0, 2.0}}
	MainRenderer.Time = time.Now()
	return MainRenderer, err
}

func (r *RenderSystem) Init(width int, height int, name string, particle_system bool) error {
	r.MyRenderer.Setup(width, height, name)
	bufs := r.Graph.GetBuffers()
	meshes := r.Graph.GetMeshes()
	imgs := r.Graph.GetImages()
	mats := r.Graph.GetMaterials()

	r.Materials = make(map[string]*defs.Material)
	r.MaterialIndexNames = make([]*string, 0, 10)
	r.Textures.TexIDMap = make(map[string]uint32)
	r.Textures.TexUnitMap = make(map[string]uint32)
	r.Textures.TexID = make([]uint32, 0, 10)
	r.Textures.TexUnit = make([]int32, 0, 10)

	p := int(0)
	if particle_system {
		p = 1
	}

	//GPU Memory Arena Declaration
	if err := r.MyRenderer.Layout(len(meshes)+p, len(bufs)+p, len(imgs)); err != nil {
		return err
	}

	r.num_tex = len(imgs)
	r.num_vbo = len(bufs) + p
	r.num_vao = len(meshes) + p
	//Buffer Binding Point
	for i := range bufs {
		bufRef := bufs[i]
		r.MyRenderer.BufferArrayData(i, bufRef.ByteLength, 0, r.Graph.Buffers[i])
	}

	//Get Textures
	for i := range imgs {
		img, _ := r.Graph.GetImageIx(i)
		if img.Uri != "" {
			uri := img.Uri
			path := r.Graph.BaseURI + uri
			fmt.Printf("Loaded Image URI: %s\n", uri)
			if texID, err := r.MyRenderer.LoadTexture(path, i); err != nil {
				return err
			} else {
				r.Textures.TexID = append(r.Textures.TexID, texID)    //Image Index Store Texture
				r.Textures.Names = append(r.Textures.Names, img.Name) //Image Index Store Texture
			}
		} else {
			return fmt.Errorf("Encoded images not yet supported in development\nImage: %v", img)
		}
	}

	//Default cubemap - we can make more procedural soon
	path := common.ProjectRelativePath("data/textures/fluidmap.png")
	if texID, err := r.MyRenderer.LoadEnvironment(path, 0); err != nil {
		return err
	} else {
		r.Textures.TexID = append(r.Textures.TexID, texID)
		r.Textures.TexIDMap["sky"] = texID
		r.Textures.TexUnitMap["sky"] = 2
	}

	//Configure Materials - GLTF wants everything to be Metallic Roughness but some parameters make sense as general material params :)
	for i := range mats {
		mat, _ := r.Graph.GetMaterialIx(i)
		name := mat.Name
		newMaterial := defs.NewMaterial(name, defs.RGB{R: 1.0, G: 1.0, B: 1.0})
		newPBR := defs.NewPBRMaterial(defs.RGB{R: 1.0, G: 1.0, B: 1.0}, 0.5, 0.5)

		baseIndex := mat.PbrMetallicRoughness.BaseColorTexture.Index
		coord := mat.PbrMetallicRoughness.BaseColorTexture.TexCoord
		glindex := r.Textures.TexID[baseIndex]
		r.Textures.TexIDMap["colorTex"] = glindex
		r.Textures.TexUnitMap["colorTex"] = 0
		baseTex := defs.Texture{baseIndex, glindex, 0, coord, "colorTex", INITIALIZED}
		newMaterial.ColorTexture = &baseTex

		normIndex := mat.NormalTexture.Index
		coord = mat.NormalTexture.TexCoord
		glindex = r.Textures.TexID[normIndex]
		r.Textures.TexIDMap["normTex"] = glindex
		r.Textures.TexUnitMap["normTex"] = 1
		normTex := defs.Texture{normIndex, glindex, 0, coord, "normTex", INITIALIZED}
		newMaterial.NormalTexture = &normTex

		newPBR.Metallic = float32(mat.PbrMetallicRoughness.MetallicFactor)
		newPBR.Roughness = float32(mat.PbrMetallicRoughness.RoughnessFactor)
		newMaterial.MetallicRoughMaterial = newPBR

		r.Materials[name] = newMaterial
	}

	return nil
}

/*
RegisterMesh - Parses GLTF scene file and fills MeshComponent entity with appropriate
GPU intialization parameters in memory
*/
func (r *RenderSystem) RegisterMesh(meshIndex int, primitiveIndex int, scn *scene.Scene) (*defs.MeshComponent, error) {

	meshComp := new(defs.MeshComponent)
	mesh, _ := scn.GetMeshIx(meshIndex)
	primitive := mesh.Primitives[primitiveIndex]

	if posAccessorIdx, ok := primitive.Attributes["POSITION"]; ok {

		PosAccessor, PosBufferView, AccessError := scn.GetAccessorBufferView(posAccessorIdx)

		if AccessError != nil {
			return nil, AccessError
		}

		if PosAccessor.ComponentType == gl.FLOAT {
			meshComp.VertexComponent = new(defs.VertexComponent)
			meshComp.VertexComponent.Vertices = scn.Buffers[PosBufferView.Buffer]
			meshComp.VertexComponent.VertexByteLength = PosBufferView.ByteLength
			meshComp.VertexComponent.VertexByteOffset = PosBufferView.ByteOffset
			meshComp.VertexComponent.VerticesVBO = PosBufferView.Buffer

		} else {
			return nil, fmt.Errorf("RegisterMesh() - GLTF Format Error: Non gl.FLOAT type not supported\n")
		}
	}
	if normAccessorIdx, ok := primitive.Attributes["NORMAL"]; ok {

		normAccessor, normBufferView, AccessError := scn.GetAccessorBufferView(normAccessorIdx)

		if AccessError != nil {
			return nil, AccessError
		}

		if normAccessor.ComponentType == gl.FLOAT {
			meshComp.NormalComponent = new(defs.NormalComponent)
			meshComp.NormalComponent.Normals = scn.Buffers[normBufferView.Buffer]
			meshComp.NormalComponent.NormalsByteLength = normBufferView.ByteLength
			meshComp.NormalComponent.NormalsByteOffset = normBufferView.ByteOffset
			meshComp.NormalComponent.NormalsVBO = normBufferView.Buffer
		} else {
			return nil, fmt.Errorf("RegisterMesh() GLTF Format Error: Normal Non gl.FLOAT type not supported\n")
		}
	}
	if texAccessorIdx, ok := primitive.Attributes["TEXCOORD_0"]; ok {

		texAccessor, texBufferView, AccessError := scn.GetAccessorBufferView(texAccessorIdx)

		if AccessError != nil {
			return nil, AccessError
		}

		if texAccessor.ComponentType == gl.FLOAT {
			meshComp.TexComponent = new(defs.TexComponent)
			meshComp.TexComponent.TexCoords = scn.Buffers[texBufferView.Buffer]
			meshComp.TexComponent.TexCoordsByteLength = texBufferView.ByteLength
			meshComp.TexComponent.TexCoordsByteOffset = texBufferView.ByteOffset
			meshComp.TexComponent.TexVBO = texBufferView.Buffer
		} else {
			return nil, fmt.Errorf("RegisterMesh() - GLTF Format Error: TexCoordinate Non gl.FLOAT type not supported\n")
		}
	}

	//Generate Mesh Location
	meshComp.TransformComponent = new(defs.TransformComponent)
	meshComp.TransformComponent.Model = matrix.Mat4(1.0)
	meshComp.TransformComponent.Position = vector.Vec3()

	if primitive.Indices == 0 {
		return nil, fmt.Errorf("RegisterMesh() - GLTF Format Error: No primitive indices reference. Indices index is 0. GLTF file may not use a indices accessor reference of zero\n")
	}

	indicesAccessorIdx := primitive.Indices
	_, IdxBufferView, _ := scn.GetAccessorBufferView(indicesAccessorIdx)
	meshComp.IndiceComponent = new(defs.IndiceComponent)
	meshComp.IndiceComponent.Indices = scn.Buffers[IdxBufferView.Buffer]
	meshComp.IndiceComponent.IndexByteLength = IdxBufferView.ByteLength
	meshComp.IndiceComponent.IndexByteOffset = IdxBufferView.ByteOffset
	meshComp.IndiceComponent.IndicesVBO = meshIndex //VAO

	materialId := primitive.Material

	meshComp.MaterialComponent = new(defs.MaterialComponent)
	matname, _ := r.Graph.GetMaterialIx(materialId)
	meshComp.MaterialComponent.Name = matname.Name

	return meshComp, nil
}

/*
Meshes() error - Processes current MeshEntity list and configures mesh VAO pointer
state for OpenGL usage. Also performs Mesh Preprocessing if neccessary
*/
func (r *RenderSystem) Meshes() error {

	for index := range r.Graph.GetMeshes() {
		if mesh, err := r.Graph.GetMeshIx(index); err != nil {
			return err
		} else {
			for j := range mesh.Primitives {
				if prim, err := r.RegisterMesh(index, j, r.Graph); err != nil {
					return err
				} else {
					r.Add(ecs.NewBasic(), *prim, index)
				}
			}
		}
	}

	//Set GL GPU VAO State Pointers for each Mesh VAO
	for index := range r.MeshEntities {

		mesh := r.MeshEntities[index]
		mesh.VAO = index
		r.MyRenderer.BindVertexArray(mesh.VAO)

		if mesh.Mesh.VertexComponent != nil {
			vc := mesh.Mesh.VertexComponent
			r.MyRenderer.BindArrayBuffer(vc.VerticesVBO)
			r.MyRenderer.VertexArrayAttr(0, 3, gl.FLOAT, vc.VertexByteOffset)
		}
		if mesh.Mesh.NormalComponent != nil {
			nc := mesh.Mesh.NormalComponent
			r.MyRenderer.BindArrayBuffer(nc.NormalsVBO)
			r.MyRenderer.VertexArrayAttr(3, 3, gl.FLOAT, nc.NormalsByteOffset)
		}
		if mesh.Mesh.TexComponent != nil {
			tc := mesh.Mesh.TexComponent
			r.MyRenderer.BindArrayBuffer(tc.TexVBO)
			r.MyRenderer.VertexArrayAttr(6, 2, gl.FLOAT, tc.TexCoordsByteOffset)
		}
		ic := mesh.Mesh.IndiceComponent
		r.MyRenderer.BufferIndexData(ic.IndicesVBO, ic.IndexByteLength, ic.IndexByteOffset, ic.Indices)
		r.MyRenderer.BindVertexArray(-1)
	}

	return nil
}

//Registers unique particle system (positions only)
func (r *RenderSystem) RegisterParticleSystem(positions []float32, layout_location int) (uint32, error) {

	//Registers particle system positions to a vertex array object

	r.MyRenderer.BindVertexArray(r.num_vao)
	r.MyRenderer.BindArrayBuffer(r.num_vbo)
	r.MyRenderer.VertexArrayAttr(layout_location, 3, gl.FLOAT, 0)
	r.MyRenderer.BufferArrayFloat(r.num_vbo, len(positions)*4, 0, positions)
	r.MyRenderer.BindVertexArray(0)
	return r.MyRenderer.GetVBO(r.num_vbo), nil
}

//Get opengl buffer location id
func (r *RenderSystem) GetParticleBufferId() uint32 {
	return r.MyRenderer.GetVBO(r.num_vbo)
}

func (r *RenderSystem) UpdateParticleSystem(positions []float32) {
	r.MyRenderer.BindVertexArray(r.num_vao)
	r.MyRenderer.BindArrayBuffer(r.num_vbo)
	r.MyRenderer.BufferSubArrayFloat(r.num_vbo, len(positions)*4, 0, positions)
	r.MyRenderer.BindVertexArray(-1)
}

func (r *RenderSystem) CompileLink() error {
	r.Shaders.VertexShaderPaths = make(map[string]string, MIN_PARAM_SIZE)
	r.Shaders.FragmentShaderPaths = make(map[string]string, MIN_PARAM_SIZE)
	r.Shaders.FragShaderID = make(map[string]uint32, MIN_PARAM_SIZE)
	r.Shaders.VertShaderID = make(map[string]uint32, MIN_PARAM_SIZE)
	r.Shaders.ProgramID = make(map[string]uint32, MIN_PARAM_SIZE)
	r.Shaders.ShaderUniforms = make(map[string]int32, MIN_PARAM_SIZE)
	r.Shaders.ProgramLinks = make(map[uint32]uint32, MIN_PARAM_SIZE)

	r.Shaders.VertexShaderPaths["default"] = common.ProjectRelativePath("data/shaders/glsl/render/material/material.vert")
	r.Shaders.FragmentShaderPaths["default"] = common.ProjectRelativePath("data/shaders/glsl/render/material/material.frag")
	r.Shaders.VertexShaderPaths["particle"] = common.ProjectRelativePath("data/shaders/glsl/render/particle_fluid/particle_fluid.vert")
	r.Shaders.FragmentShaderPaths["particle"] = common.ProjectRelativePath("data/shaders/glsl/render/particle_fluid/particle_fluid.frag")

	//----------Uniforms ------------------------------//
	r.Shaders.ShaderUniforms["mvp"] = 0
	r.Shaders.ShaderUniforms["view"] = 0
	r.Shaders.ShaderUniforms["model"] = 0
	r.Shaders.ShaderUniforms["cube"] = 0
	r.Shaders.ShaderUniforms["colorTex"] = 0
	r.Shaders.ShaderUniforms["normTex"] = 0
	r.Shaders.ShaderUniforms["lightPos"] = 0
	r.Shaders.ShaderUniforms["lightColor"] = 0
	r.Shaders.ShaderUniforms["baseColor"] = 0
	r.Shaders.ShaderUniforms["fresnel_rim"] = 0
	r.Shaders.ShaderUniforms["metallness"] = 0
	r.Shaders.ShaderUniforms["roughness"] = 0
	r.Shaders.ShaderUniforms["normMat"] = 0
	r.Shaders.ShaderUniforms["_mvp"] = 0
	r.Shaders.ShaderUniforms["_view"] = 0
	r.Shaders.ShaderUniforms["_model"] = 0

	//----------Compile Shader Targets-----------------//
	for key, value := range r.Shaders.VertexShaderPaths {
		if vtxID, err := r.MyRenderer.AddShader(value, gl.VERTEX_SHADER); err == nil {
			r.Shaders.VertShaderID[key] = vtxID
		} else {
			return err
		}
	}

	for key, value := range r.Shaders.FragmentShaderPaths {
		if frgID, err := r.MyRenderer.AddShader(value, gl.FRAGMENT_SHADER); err == nil {
			r.Shaders.FragShaderID[key] = frgID
		} else {
			return err
		}
	}

	//-------------Links programs by vertex/frag shaders sharing same key, enters program id into same key entry-----//
	for key, value := range r.Shaders.FragShaderID {
		var err error
		fID := value
		vID := r.Shaders.VertShaderID[key]
		r.Shaders.ProgramLinks[vID] = fID
		r.Shaders.ProgramID[key], err = r.MyRenderer.LinkShaders(vID, fID)
		if err != nil {
			return err
		}

	}

	//-------------Assign default shader uniform locations-------------------//
	for key := range r.Shaders.ShaderUniforms {
		if key != "_mvp" && key != "_view" && key != "_model" {
			id := r.Shaders.ProgramID["default"]
			uniform_id, _ := r.MyRenderer.GetUniformLocation(id, key)
			r.Shaders.ShaderUniforms[key] = uniform_id
		} else {
			id := r.Shaders.ProgramID["particle"]
			uniform_id, _ := r.MyRenderer.GetUniformLocation(id, key)
			r.Shaders.ShaderUniforms[key] = uniform_id
		}
	}

	r.MyRenderer.ShaderLog(r.Shaders.ProgramID["default"], "Default glsl")
	r.MyRenderer.ShaderLog(r.Shaders.ProgramID["particle"], "particle glsl")

	gl.UseProgram(r.Shaders.ProgramID["default"])

	baseColor := []float32{1.0, 1.0, 1.0, 1.0}
	gl.Uniform4fv(r.Shaders.ShaderUniforms["baseColor"], 1, &baseColor[0])
	gl.Uniform1f(r.Shaders.ShaderUniforms["fresnel_rim"], 0.05)
	gl.Uniform3fv(r.Shaders.ShaderUniforms["lightColor"], 1, &r.Light.Color[0])
	gl.Uniform3fv(r.Shaders.ShaderUniforms["lightPos"], 1, &r.Light.Pos[0])

	//Set Material Location Uniforms -- Material shader only accepts 1 Materials
	for _, material := range r.Materials {

		if material.ColorTexture.Status == INITIALIZED {
			material.ColorTexture.GLID = r.Textures.TexIDMap["colorTex"]
			material.ColorTexture.Key = "colorTex"
			material.ColorTexture.GLTEXUNIT = 0
			gl.Uniform1i(r.Shaders.ShaderUniforms["colorTex"], 0)
			r.MyRenderer.SetActiveTexture(r.Textures.TexIDMap["colorTex"], 0)
		}
		if material.NormalTexture.Status == INITIALIZED {
			material.ColorTexture.GLID = r.Textures.TexIDMap["normTex"]
			material.ColorTexture.Key = "normTex"
			material.ColorTexture.GLTEXUNIT = 1
			gl.Uniform1i(r.Shaders.ShaderUniforms["normTex"], 1)
			r.MyRenderer.SetActiveTexture(r.Textures.TexIDMap["normTex"], 1)
		}

	}

	//Set skybox
	gl.Uniform1i(r.Shaders.ShaderUniforms["cube"], 2)
	r.MyRenderer.SetActiveTexture(r.Textures.TexIDMap["sky"], 2)

	return nil
}

//Update() - Scene time step update with frame render rendering components associated
func (r *RenderSystem) Update(elapsed float64) {
	r.MyRenderer.Update(0.01)
	a := float32(math.Cos(float64(r.Angle)))
	b := float32(math.Sin(float64(r.Angle)))
	x := a * r.Light.Pos[0]
	y := b * r.Light.Pos[1]
	z := r.Light.Pos[2]
	nLight := []float32{x, y, z}
	gl.Uniform3fv(r.Shaders.ShaderUniforms["lightPos"], 1, &nLight[0])
	r.Angle += 0.01
	if r.Angle > 2*3.141529 {
		r.Angle = 0
	}
}

func (r *RenderSystem) Remove(basic ecs.BasicEntity) {
	var delete int = -1
	for index, entity := range r.MeshEntities {
		if entity.Basic.ID() == basic.ID() {
			meshComponent := entity.Mesh
			if &meshComponent.VertexComponent != nil {
				vert := meshComponent.VertexComponent
				r.MyRenderer.InvalidateBuffer(vert.VerticesVBO, vert.VertexByteLength)
			}
			if &meshComponent.NormalComponent != nil {
				norm := meshComponent.NormalComponent
				r.MyRenderer.InvalidateBuffer(norm.NormalsVBO, norm.NormalsByteLength)
			}
			if &meshComponent.ColorComponent != nil {
				color := meshComponent.ColorComponent
				r.MyRenderer.InvalidateBuffer(color.ColorsVBO, color.ColorsByteLength)
			}
			if &meshComponent.TexComponent != nil {
				tex := meshComponent.TexComponent
				r.MyRenderer.InvalidateBuffer(tex.TexVBO, tex.TexCoordsByteLength)
			}
			if &meshComponent.IndiceComponent != nil {
				ind := meshComponent.IndiceComponent
				r.MyRenderer.InvalidateBuffer(ind.IndicesVBO, ind.IndexByteLength)
			}
			delete = index
			break
		}
	}
	if delete >= 0 {
		r.MeshEntities = append(r.MeshEntities[:delete], r.MeshEntities[delete+1:]...)
	}
}

//Add() - Adds MeshEntity to RenderSystem
func (r *RenderSystem) Add(e ecs.BasicEntity, mesh defs.MeshComponent, i int) []*defs.MeshEntity {
	entity := defs.MeshEntity{e, mesh, i}
	ret := append(r.MeshEntities, &entity)
	r.MeshEntities = ret
	return ret
}

func (r *RenderSystem) GetColliderMeshes() []*mesh.Mesh {

	colliderMeshes := make([]*mesh.Mesh, len(r.MeshEntities))
	for i, meshy := range r.MeshEntities {
		vertices := byteSliceToFloat32Slice(meshy.Mesh.Vertices)
		tris := make([]vector.Vec, len(vertices)/3)
		for i := 0; i < len(tris); i++ {
			x := i * 3

			tris[i] = vector.Vec{vertices[x], vertices[x+1], vertices[x+2]}
		}
		m_mesh := mesh.InitMesh(tris, []float32{0, 0, 0})
		colliderMeshes[i] = &m_mesh
	}
	return colliderMeshes
}

//Run() RenderSystem Runtime
func (r *RenderSystem) Run(message chan string) error {

	if r.MyRenderer == nil {
		return fmt.Errorf("Run()- MyRenderer not available\n")
	}

	if err := r.MyRenderer.Status(); err != nil {
		return err
	}

	frame_timer := 0.0
	for !r.MyRenderer.ShouldClose() {
		r.MyRenderer.Draw(r.MeshEntities, &r.Shaders, r.Materials, r.Textures.TexIDMap["sky"])
		r.MyRenderer.SwapBuffers()
		glfw.PollEvents()
		seconds := time.Now().Sub(r.Time).Seconds()
		r.Update(seconds)

		frame_timer += seconds
		if frame_timer < 0.1 {
			select { //Non blocking send
			case message <- "ACK":
			default:
			}
		} else {
			frame_timer = 0.0
			select { //Non blocking send
			case message <- "REFRESH_PARTICLES":
			default:
			}
		}
	}

	runtime.UnlockOSThread()

	return nil

}
