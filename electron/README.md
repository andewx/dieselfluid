## Electron GUI

#OVERVIEW

This is a node electron GPU WebGL application for GPU fluid computation. This layer
acts as a client which initializes off of the GoLang data structures through a JSON
exchange with the DslFluid application server.

The use case for such a hierarchy is for easy integration of pre-computed data Structures
injected into a WebGL instance for GPU based computational purposes.

### Electron Events
The Node Electron Application instance will need to ingest the intialization conditions
and is expected to pass buffers to the GPU side and load the correct GLSL shader pipeline

#### FDXF - Fluid Data Exchange Format

```
"fluid" : {
  type: 'sph3d',
  "sph" : {
    Positions: [...],
    Length: 1024
    Densities: [...],
    Velocities: [...],
    Pressures: [...],
    Exp: 7.1,
    Kern: 0.3,
    Rad: 0.1,
    TargetDens: 1.0
    SimTime: 123124142141, //ms
    TStep: 0.0000004,
    NegScale: 0.000000001
  }
}
```
This is the construction for a grid solver:

```
"fluid" : {
  "type": "grid3d",
  "grid":{
    type: "static",
    centroids: [1.0,1.0,1.0...],
    dims: [1.0,1.0,1.0],
    densities: [...],
    pressures: [...],
    velocities: [...],
    mass: [...],
  }
  TargetDens: 1.0
  SimTime: 0,
  TStep: 0.0000004,
  NegScale: 0.000000001
}
```

## Scene data
```
"scene" : {
  Camera: [...],
  Origin: [...],
  FOV: 45,
  RotMatrix: [...],
  PerMatrix: [...],
  VecScale: [...]
}
```

## Model data
JSON 3D Exchange Format
```
"model":
  {
     "meshes":
     [
        {
            "meta":
            {
               "vertNum": 8,
               "faceNum": 12
            },
            "name": "Box1",
            "isCollider": "true",
            "drawCollider": "true",
            "node": 1,
            "verts": [-50,0,25,50,0,25,-50,0,-25, ...],
            "vertElement":
            {
               "vertIndices": [0,2,3, ...],
               "normals": [0,-1,0,0,-1,0, ...],
               "uvs": [[1,0,1,1, ...]]
            },
            "face":
            {
               "vertElementIndices": [0,1,2, ...],
               "groups":
               [
                 {
                    "start": 0,
                    "count": 36,
                    "materialIndex": 0
                 }
               ]
            },
        }
     ]
  }
  ```


#### Go Server HTTP Response Format
The Go Sim Server will be responsible for passing the RAW data objects back to the Electron Client from the
specified URLs
