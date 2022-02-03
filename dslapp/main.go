package dslapp

import ( //Potential Package conflict with "app"
	"dslfluid.com/dsl/render"
	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	"log"
	"runtime"
)

func main() {

	// Set logger
	l := log.New(log.Writer(), log.Prefix(), log.Flags())

	// Create astilectron
	a, err := astilectron.New(l, astilectron.Options{
		AppName:           "Diesel Fluid",
		BaseDirectoryPath: "/Users/andewx/go/src/github.com/andewx/dieselfluid/dslapp",
	})
	if err != nil {
		//r.Errorf("main: creating astilectron failed: %w", err))
	}
	defer a.Close()

	// Handle signals
	a.HandleSignals()

	// Start
	if err = a.Start(); err != nil {
		//	l.Fatal(fm//r.Errorf("main: starting astilectron failed: %w", err))
	}

	// New window
	var w *astilectron.Window
	if w, err = a.NewWindow("example/index.html", &astilectron.WindowOptions{
		Center: astikit.BoolPtr(true),
		Height: astikit.IntPtr(700),
		Width:  astikit.IntPtr(700),
	}); err != nil {
		//l.Fatal(fm//r.Errorf("main: new window failed: %w", err))
	}

	// Create windows
	if err = w.Create(); err != nil {
		//	l.Fatal(fm//r.Errorf("main: creating window failed: %w", err))
	}

	// Blocking pattern
	a.Wait()

	runtime.LockOSThread()
	Sys, _ := render.Init("MaterialSphere.gltf")

	if err := Sys.Init(1024, 720, "andewx/diselfluid/render"); err != nil {
		//r.Error(err)
	}

	if err := Sys.CompileLink(); err != nil {
		//r.Error(err)
	}

	if err := Sys.Meshes(); err != nil {
		//r.Error(err)
	}

	if err := Sys.Run(); err != nil {
		//r.Error(err)
	}
	//Initialize and Hook all GO/Server Callbacks
	//Initialize GLFW window
	//Initialize Renderer Settings
	//Initialize Fluid Settings
	//Initialize Fluid Structures
	//Initialize GL Texture Data
	//Initialize GL Shader Data
	//Initialize GL Render Contexts & Buffers
	//GL Render

}
