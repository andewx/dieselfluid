# dieselfluid

### Project Status

Refactor build ongoing once modules are complete the testing, code coverage, and
project documentation will be updated and moved into a docs repository

![DieselFluid Go Logo](logo.png?raw=true "DieselFluid Go")


# README
Currently Apple OSX has deprecated Modern OpenGL and only supports OpenGL 4.1 core
profiles with future deprecation on modern M1 chipsets. Currently all render routines
run throuhg OpenGL although the calls have been wrapped in an API and interface. Future
updates could theoretically support OpenGL rendering. Proper support for abstract rendering
interfaces require further decoupling and refactors which are not being supported at this time. 

Since OSX Darwin builds do not support OpenGL 4.3 (GL COMPUTE SHADER) profiles efforts have been made to split the project builds between Darwin and all other builds. Path separators however are currently
not windows compatible although a refactor in the future could add support for windows.
OSX Darwin Builds leverage the OpenCL glow build interfaces for GPU support

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
