package render

import "github.com/EngoEngine/ecs" //GoLang Entity Component System
import "github.com/andewx/dieselfluid/render/scene"
import "github.com/andewx/dieselfluid/render/glr"
import "github.com/andewx/dieselfluid/render/defs"
import "github.com/go-gl/gl/v4.1-core/gl"
import "github.com/go-gl/glfw/v3.3/glfw"

import "github.com/andewx/dieselfluid/math/mgl"
import "fmt"
import "time"
import "runtime"
import "math"

const (
	MIN_PARAM_SIZE  = 10
	BYTES_PER_FLOAT = 4
	COORD_4         = 4
	COORD_3         = 3
	COORD_2         = 2
	INITIALIZED     = 1
)

//RenderSystem provides ECS Render Mechanism
type RenderSystem struct {
	MyRenderer         defs.OGLRenderer
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
}

//Gather() - Package customized initialization routine function. Gather is an
//initialization like function moniker used package by package to initialize
//a main package construct for API usage. Here Gather returns the main package
//type which is the RenderSystem and initializes with a GLTF scene file name
//located in the "dslfluid.com/resources/" directory
func Init(scn string) (RenderSystem, error) {
	graph, err := scene.InitScene("../data/", scn)
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
	MainRenderer.Light = defs.Light{[]float32{50, 50, 50, 1}, []float32{2.0, 2.0, 2.0}}
	MainRenderer.Time = time.Now()
	return MainRenderer, err
}

func (r *RenderSystem) Init(width int, height int, name string) error {
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

	//GPU Memory Arena Declaration
	if err := r.MyRenderer.Layout(len(meshes), len(bufs), len(imgs)); err != nil {
		return err
	}
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
			path := "../resources/" + uri
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
	path := "../resources/fluidmap.png"
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
		newMaterial := defs.NewMaterial(name, defs.RGB{1.0, 1.0, 1.0})
		newPBR := defs.NewPBRMaterial(defs.RGB{1.0, 1.0, 1.0}, 0.5, 0.5)

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
	meshComp.TransformComponent.Model = mgl.Mat4(1.0)
	meshComp.TransformComponent.Position = mgl.Vec3()

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

func (r *RenderSystem) CompileLink() error {
	r.Shaders.VertexShaderPaths = make(map[string]string, MIN_PARAM_SIZE)
	r.Shaders.FragmentShaderPaths = make(map[string]string, MIN_PARAM_SIZE)
	r.Shaders.FragShaderID = make(map[string]uint32, MIN_PARAM_SIZE)
	r.Shaders.VertShaderID = make(map[string]uint32, MIN_PARAM_SIZE)
	r.Shaders.ProgramID = make(map[string]uint32, MIN_PARAM_SIZE)
	r.Shaders.ShaderUniforms = make(map[string]int32, MIN_PARAM_SIZE)
	r.Shaders.ProgramLinks = make(map[uint32]uint32, MIN_PARAM_SIZE)

	r.Shaders.VertexShaderPaths["default"] = "../shader/glsl/material.vert"
	r.Shaders.FragmentShaderPaths["default"] = "../shader/glsl/material.frag"

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
	for key, _ := range r.Shaders.ShaderUniforms {
		id := r.Shaders.ProgramID["default"]
		uniform_id, _ := r.MyRenderer.GetUniformLocation(id, key)
		r.Shaders.ShaderUniforms[key] = uniform_id
	}

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

//Run() RenderSystem Runtime
func (r *RenderSystem) Run() error {

	if r.MyRenderer == nil {
		return fmt.Errorf("Run()- MyRenderer not available\n")
	}

	if err := r.MyRenderer.Status(); err != nil {
		return err
	}

	for !r.MyRenderer.ShouldClose() {
		r.MyRenderer.Draw(r.MeshEntities, &r.Shaders, r.Materials, r.Textures.TexIDMap["sky"])
		r.MyRenderer.SwapBuffers()
		glfw.PollEvents()
		r.Update(time.Now().Sub(r.Time).Seconds())
	}

	runtime.UnlockOSThread()

	return nil

}
