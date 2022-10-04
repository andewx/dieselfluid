/*
 Shader - OpenCL Compute Shader Executes the advection and pressure error solver
Note that the excess particles are the appended particle values which are not directly
Invoked by the compute shader
*/
#ifndef pcisph_h
#include "pcisph.h"
#endif


//Density Compute Shader
kernel void compute_density(global particle *fluid,
                            global intdata *sizes,
                            global floatdata *data,
                            global int *table,
                            global float3 *vecs){
  int x = get_global_id(0);
  fluid[x].density = interp_density(fluid[x].position,data,fluid,table,sizes,vecs);
  evaluate_pressure(x,fluid);
}
