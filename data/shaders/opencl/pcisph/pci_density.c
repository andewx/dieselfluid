/*
 Shader - OpenCL Compute Shader Executes the advection and pressure error solver
Note that the excess particles are the appended particle values which are not directly
Invoked by the compute shader
*/
#ifndef pcisph_h
#include "pcisph.h"
#endif


//Density Compute Shader
kernel void compute_density(global float3 *positions, global float3 *velocities, global float3 *forces, global float *densities, global float *pressures,
                            global intdata *sizes,
                            global floatdata *data,
                            global int *table,
                            global float3 *vecs){
  int x = get_global_id(0);
  particles m_particles = {positions, velocities,forces, densities, pressures};
  densities[x] = interp_density(positions[x],data, &m_particles,table,sizes,vecs);
  evaluate_pressure(x, &m_particles);
}
