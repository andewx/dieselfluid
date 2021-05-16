# dieselfluid

### This Repository is still in development

![DieselFluid Go Logo](logo.png?raw=true "DieselFluid Go")


# README

This version of Diesel SPH is still in development and is not yet build ready.

# OVERVIEW

Diesel SPH develops a fluid modeling framework in Go language. Allows for scene construction with both Grid/Particle/Hybrid Model Solvers. It supports a software CPU simulations, but also supports a HTTP initialization fluid exchange format so that GPU applications can access initialized fluid data Structures
in order to run GPU based fluid computation.

## External Dependencies:
Go-Gl
Go GLFW
GoNum
PBRT

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
