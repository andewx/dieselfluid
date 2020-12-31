# dieselfluid
Diesel Smoothed Particle Hydrodynamics Library. Using Go Language, Go-GL / GLFW frameworks / Gonum

Not Curently GIT Hosted

Need to initialize all packages and modules still

Seeks to be SPH GO LANG Implementation framework of reference


# README

This version of Diesel SPH is still in development and not yet considered releasable for application development. However you may clone and work on your own repository. Currently we are testing just the SPH methods and rendering.
Next order of business is moving fluid logic to the GPU side of the house and dealing with GPU shaders primarily

# OVERVIEW

Diesel SPH develops a langrangian fluid simulation environment with the GO Language environment. Go language was chosen for the ability to quickly build and test frameworks and applications, it's superior runtime performance, and easy multi-threading support. Note that all context calls to the GLFW/GL frameworks requires a single locked thread once the windowed application starts.

# DEPENDENCIES:

Uses <b>Go 1.8+</b>

Currently this version uses <b>GL V4.1 / and GLFW wV3.2.</b> In the future we may look at build constraints in Go Lang in order to allow different versions of go-gl to be invoked.

##  For OpenGL 4.1
'''bash
$ go get -u github.com/go-gl/gl/v4.1-core/gl
'''

##  For GLFW 3.2
'''bash
$ go get github.com/go-gl/glfw/v3.2/glfw
'''
###Language Specification

Colon Parse SPH Language Specification

'''bash
N=100,000
D=20
Method:DFSPH
Particles:N
GridResolution:D
SceneOrigin:{0,0,0}
CameraOrigin:{0,1,1}
CameraVec:{0,0,1}
FOV:{1,1,1,1,1}
EOS:{1,1}
Sampler:Voxel
Render:GL
Notifications:ON
NotifyPeriod:1000MS
Log:true
LogByteMax:1MB
Export:Mitsuba
GeomExport:OBJ
Material:Particle
BSDF:None
ExportFinalState:False 'For Network Nodes'
VolumeMesh:MARCHING_INTERP
VolumeBSDF:WATER|SMOKE
Error:0.1
BenchmarkLog:True
'''

##  TODO

1. Turn DieselSPH in extensible production Code
2. Implement WCSPH - DFSPH Methods
2. Serialize and compress animation Keyframes
3. Implement Multiple Frame Buffers in GL for Fragment Discard
4. Export Scene for Lux Render / PBRT / Cycles Ray Tracing

## BRANCHES

Master - Latest working build
Working - Development Build