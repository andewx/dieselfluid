package dslapp

import (
	"github.com/maxence-charriere/go-app/v8/pkg/app" //Potential Package conflict with "app"
	"log"
	"net/http"
)

type UIComponent struct {
	app.Compo
}

func (panel *UIComponent) Render() app.UI {
	return app.H1().Text("Go-App World")
}

func main() {

	//Opens and closes server connection
	server := &http.Server{Addr: ":8000", Handler: &app.Handler{
		Name:        "DslFluid",
		Description: "DslFluid GUI Toolkit Window",
	}}
	app.Route("/", &UIComponent{})

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println("HTTP Error: ", err)
		}
	}()

	server.Close()

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
