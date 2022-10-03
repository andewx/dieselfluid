//go:build darwin
// +build darwin

package shader

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"io/ioutil"
	"strings"
	"unsafe"
)

const (
	FRAG_SHADER    = gl.FRAGMENT_SHADER
	VERT_SHADER    = gl.VERTEX_SHADER
	GEOM_SHADER    = gl.GEOMETRY_SHADER
	COMPUTE_SHADER = gl.COMPUTE_SHADER
)

type Shader struct {
	filename      string
	path          string
	name          string
	contents      string
	gpu_shader_id uint32
	message       string
	log           string
	shader_type   uint32
	compiled      bool
}

type Program struct {
	name           string
	gpu_program_id uint32
	links          map[string]*Shader
	message        string
	log            string
	linked         bool
	uniforms       map[string]*uint32
}

//-----------------------------------
//          Shader
//-----------------------------------
//New shader attempts to create and compile a new shader when instantiated
func NewShader(filename string, path string, name string, shader_type uint32) (*Shader, error) {

	var err error
	var bytes []byte

	sh := &Shader{filename, path, name, "", 0, "", "", shader_type, false}
	full_path := path + "/" + filename
	bytes, err = ioutil.ReadFile(full_path)
	if err != nil {
		sh.message = "Shader Invalid Path Specified: " + full_path
		sh.log = sh.log + "\nNewShader() - Invalid Path (see message)"
		return sh, err
	}
	sh.log += "File Found " + filename
	sh.contents = string(bytes) + "\x00"

	if shader_type != FRAG_SHADER && shader_type != VERT_SHADER && shader_type != GEOM_SHADER && shader_type != gl.COMPUTE_SHADER {
		sh.log += "Not A Valid Shader Type" + string(shader_type) + "\n"
		return sh, fmt.Errorf("Not a valid GL Shader Type %d", shader_type)
	}

	if shader_type == gl.COMPUTE_SHADER {
		sh.log += "Compute Shader Specified. Please ensure you are running a working version of GL 4.3 or higher\n"
		version := gl.GoStr(gl.GetString(gl.VERSION))
		fmt.Printf("OpenGL Version: %s\n", version)
	}

	err = sh.Compile()
	if err != nil {
		fmt.Printf("Error compiling %s path %s type %d", sh.name, full_path, sh.shader_type)
		fmt.Printf(err.Error())
		return sh, err
	}
	sh.compiled = true
	sh.message = "Compiled"
	sh.log += "Compiled"
	sh.contents = ""
	return sh, nil
}

/*
 Invokes GL Shader Compilation return 0 and error on compilation error
*/
func (m *Shader) Compile() error {
	var status int32

	shader := gl.CreateShader(m.shader_type)
	csources, free := gl.Strs(m.contents)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)
	check_error("Compiled Shader[" + m.name + "]")
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		m.message = "Compilation Error See Log\n"
		m.log += "\n" + log
		return fmt.Errorf("Failed to Compile %v: %v", m.contents, m.log)
	}
	m.gpu_shader_id = shader
	return nil
}

func (m *Shader) GetMessage() string {
	return m.message
}

func (m *Shader) GetLog() string {
	return m.log
}

func (m *Shader) Id() uint32 {
	return m.gpu_shader_id
}

func (m *Shader) IsCompiled() bool {
	return m.compiled
}

func (m *Shader) PrintMessage() {
	fmt.Printf("Message: %s\nLog: %s\n", m.message, m.log)
}

//-----------------------------------
//          Program
//-----------------------------------

func NewProgram(name string) *Program {
	pg := &Program{name, 0, make(map[string]*Shader), "", "", false, make(map[string]*uint32, 20)}
	pg.gpu_program_id = gl.CreateProgram()
	return pg
}

/*Adds Shader*/
func (m *Program) AddShader(sh *Shader) error {
	if sh.IsCompiled() {
		m.links[sh.name] = sh
		return nil
	}
	return fmt.Errorf("Cannot add a non-compiled shader to a program\n")
}

/*Deletes Shader*/
func (m *Program) DeleteShader(name string) {
	delete(m.links, name)
}

/*Adds Uniform*/
func (m *Program) AddUniform(name string) {
	*m.uniforms[name] = 0
}

func (m *Program) Address(name string) *uint32 {
	return m.uniforms[name]
}

/*Links programs with all attached valid compiled shaders, assumes that all shaders and uniforms
expected have been set.*/
func (m *Program) Link() error {
	var err error
	for key, val := range m.links {
		gl.AttachShader(m.gpu_program_id, val.gpu_shader_id)
		check_error("Attach shader[" + key + "]")
	}
	gl.LinkProgram(m.gpu_program_id)
	check_error("Link Program[" + m.name + "]")
	if m.gpu_program_id == gl.INVALID_VALUE || m.gpu_program_id == gl.INVALID_OPERATION {
		err = fmt.Errorf("Invalid Linking")
		return err
	}
	m.linked = true

	for key, _ := range m.uniforms {
		loc, uerr := m.GetUniformLocation(key)
		if uerr == nil {
			*m.uniforms[key] = loc
		}
	}
	return nil
}

func (m *Program) Use() {
	gl.UseProgram(m.gpu_program_id)
}

func (m *Program) Halt() {
	gl.UseProgram(0)
}

func (m *Program) IsLinked() bool {
	return m.linked
}

func (m *Program) Location(name string) uint32 {
	return *m.uniforms[name]
}

func (m *Program) GetUniformLocation(name string) (uint32, error) {
	loc := uint32(gl.GetUniformLocation(m.gpu_program_id, gl.Str(name+"\x00")))
	if loc == gl.INVALID_VALUE || loc == gl.INVALID_OPERATION {
		err := fmt.Errorf("Uniform Location %s Not Found\n", name)
		return loc, err
	}
	return loc, nil
}

func (m *Program) BindShaderStorage(name string, layout uint32, length int, data unsafe.Pointer) {
	gl.GenBuffers(1, m.Address(name))
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, m.Location(name))
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, length, data, gl.STREAM_DRAW)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, layout, m.Location(name))
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, 0)
}

func (m *Program) GetMaxWorkGroup() [3]int32 {
	work_group_size := [3]int32{}
	gl.GetIntegeri_v(gl.MAX_COMPUTE_WORK_GROUP_SIZE, 0, &work_group_size[0])
	gl.GetIntegeri_v(gl.MAX_COMPUTE_WORK_GROUP_SIZE, 1, &work_group_size[1])
	gl.GetIntegeri_v(gl.MAX_COMPUTE_WORK_GROUP_SIZE, 2, &work_group_size[2])
	return work_group_size
}

func (m *Program) DispatchCompute(x int32, y int32, z int32) {
	gl.DispatchCompute(uint32(x), uint32(y), uint32(z))
}

func (m *Program) ShaderLog() {
	var logLength = int32(2048)
	log := strings.Repeat("\x00", int(logLength+1))
	gl.GetProgramInfoLog(m.gpu_program_id, logLength, nil, gl.Str(log))
	fmt.Printf("%s", log)
	active := int32(0)
	gl.GetProgramiv(m.gpu_program_id, gl.ACTIVE_UNIFORMS, &active)
	fmt.Printf("SHADER UNIFORMS[%d]\n", active)
}

func check_error(op string) int {
	error := gl.GetError()
	if error == gl.NO_ERROR {
		return 0
	}
	fmt.Printf(op+"GL Error %d: ", error)
	return int(error)
}
