# dieselfluid

### Project Status

Much of the computational fluid library went under a major refactor and is now in ```model``` package and has been postponed in favor the renderer. The ```render``` package has been developed and is undergoing testing and further feature additions. We need to add a few more core features and a suitable gui system to consider this repository release ready however the repository remains open for those wishing to contribute. 

If you would like to contribute please contact the author andewx and volunteer to complete a core feature listed in the issues section of the project. The motivation for this project is to bring in a framework to experiment with physics and fluids in GoLang while also bringing the amenities of employing some modern real-time rendering techniques in the process. If there are other feature requests or contributions you would like to make don't hesitate to ask the author.

![DieselFluid Go Logo](logo.png?raw=true "DieselFluid Go")


# README


# OVERVIEW

Diesel SPH develops a fluid modeling framework in Go language. Allows for scene construction with both Grid/Particle/Hybrid Model Solvers. It supports a software CPU simulation / GPU Simulations as well and is compatible with the versatile glTF Scene description format, and even extends this format with added extensions for loading fluid frames & geometry descriptors into scenes. The renderer for DieselFluid is intended to be a suitable modern renderer utilizing modern techniques for well lit physically based scenes such as Diffuse Irradiance Mapping, Spherical Harmonics, Ambient Occlusion and Light mapping. 

DieselSPH can export glTF files but 100% compatibility is not guaranteed.

DslFluid also comes with its own render packages which can be used in other go
projects as part of this open source project. See the ```render``` package
for the main rendering methods. 

![Screeshot](rendersphere.png?raw=true "Render Sphere")

## External Dependencies:
Go-Gl
Go GLFW

## Primary Requirements
- Go CPU Grid/SPH/Hybrid Solver w/ OpenGL support
- Export to Fluid Data Exchange Format application support
- Support Node Electron DslFluid application extension for WebGL support

##  For OpenGL 4.1
```bash
$ go get -u github.com/go-gl/gl/v4.1-core/gl
```

##  For GLFW 3.2
```bash
$ go get github.com/go-gl/glfw/v3.2/glfw
```

## BRANCHES

main - Latest working build
