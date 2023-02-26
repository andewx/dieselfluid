//PCISPHs Shader - GLSL Compute Shader Executes the advection and pressure error solver
//Note that the excess particles are the appended particle values which are not directly
//Invoked by the compute shader
#ifndef pcisph_h
#include "pcisph.h"
#endif

//Compute Gradient Level Forces + Projected Intermediate Velocities and Positions
kernel void predict_correct(global float3 *positions, global float3 *velocities, global float3 *forces, global float *densities, global float *pressures,
                            global intdata* sizes,
                            global floatdata* data,
                            global int* table,
                            global float3* vecs,
                            global temp_particle* temp){
  int x = get_global_id(0);
  float ts = cfl(data->maxVel);
  particles m_particles = {&positions[x], &velocities[x],&forces[x], &densities[x], &pressures[x]};
  pressure_solve(x,temp,data, &m_particles,table,sizes,vecs);
  velocities[x] = velocities[x] + ts* forces[x]/data->mass;
  positions[x] = positions[x] + ts*velocities[x];
  pressures[x] = 0.0;
  forces[x]  = 0.0;
  //CFL Condtions
  if (length(velocities[x]) > data->maxVel){
    data->maxVel = length(velocities[x]);
  }
}
