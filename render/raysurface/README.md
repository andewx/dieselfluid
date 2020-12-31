#RAYSURFACE Module
###dieselfluid/render/raysurface

##Overview:

Implements Ray-Tracing Surface Fluid construction the CPU which allows for surface reconstruction from the resolution of the image space. With GPU implemented algorithms this means that we could construct fluid surface construction in realtime.

Based on "Graphics Symposium on Parallel Graphics and Visualization (2018)" - Direct Raytracing of Particle Based Fluid Surfaces Using Anisoptric kernels (Biedert, Sohns, Schroder, Et Al)

This approach is extended to the FLIP/PIC Hybrid Solver Resconstruction Method which allows a straighforward implementation of Interpolated Marching Cube Construction which is combined with the Ray Traced Surface.

Here Smoothing Kernels Define the anisoptirc iso-surface value that may be associated with a particle.

We allow for geometry construction, but it is better to let the offline renderer (MITSUBA/PBRT) have access to the ISOSURFACE scalar field and particle sets, and extend the Ray components to handle these fields within the underlying system. This means that the fluid context for the frame can be exported and tagged into the render environment and handled from there.
