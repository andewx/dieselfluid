# dieselfluid

### This Repository is still in development

![DieselFluid Go Logo](logo.png?raw=true "DieselFluid Go")


# README

This version of Diesel SPH is still in development and is not yet build ready.

# OVERVIEW

Diesel SPH develops a fluid modeling framework in Go language. Allows for scene construction with both Grid/Particle/Hybrid Model Solvers. It supports a software CPU simulation / GPU Simulations as well and is compatible with the versatile glTF
Scene description format, and even extends this format with added extensions for
loading fluid frames & geometry descriptors into scenes.

DslFluid also comes with its own render packages which can be used in other go
projects as part of this open source project. See the ```render/glr``` package
for the main rendering methods. 

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
