# dieselfluid

### Project Status

Much of the computational fluid library went under a major refactor and is now in ```model``` package and has been postponed in favor the renderer. The ```render``` package has been developed and is undergoing testing and further feature additions. We need to add a few more core features and a suitable gui system to consider this repository release ready however the repository remains open for those wishing to contribute.

If you would like to contribute please contact the author andewx and volunteer to complete a core feature listed in the issues section of the project. The motivation for this project is to bring in a framework to experiment with physics and fluids in GoLang.

![DieselFluid Go Logo](logo.png?raw=true "DieselFluid Go")


# README


# OVERVIEW

Diesel SPH develops a fluid modeling framework in Go language with environment
render capabilities. The API functionality is intended to be exposed to Electron
API viewports and therefore high level API requests respond to TCP hanlders.
![Screeshot](rendersphere.png?raw=true "Render Sphere")

## Primary Requirements
- Asset Management / Render Controls / Task Organizing / Event Arch
- Advanced Lighting / Shadows both realtime and offline render capabilities
- Areal/Volumetric Real Time Lighting / Ambient Occlusions / Light Pre-Processing
- Advanced Volumetric Rendering
- GLTF Advanced 3D Scenes & Assets
- SPH/Fluid Grid Interface Solvers
- SPH/Fluid Solver GPU Computation Exports
- SPH/Fluid Implicit Collision Technqiues / SDF Processing
- SPH/Fluid Surfaces
- Fluid Animation Export
- Modular Scene Assembly

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
