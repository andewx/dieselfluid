//PCISPHs Shader - GLSL Compute Shader Executes the advection and pressure error solver
//Note that the excess particles are the appended particle values which are not directly
//Invoked by the compute shader
#ifndef pcisph_h
#include "pcisph.h"
#endif

//Compute Gradient Level Forces + Projected Intermediate Velocities and Positions
kernel void predict_correct(global particle* fluid,
                            global intdata* sizes,
                            global floatdata* data,
                            global int* table,
                            global float3* vecs,
                            global temp_particle* temp){
  int x = get_global_id(0);
  float ts = cfl(data->maxVel);
  viscosity_force(x, data,fluid,table,sizes,vecs);
  external_force(x,data ,fluid);
  pressure_solve(x,temp,data,fluid,table,sizes,vecs);
  fluid[x].velocity = fluid[x].velocity + ts* fluid[x].force/data->mass;
  fluid[x].position = fluid[x].position + ts*fluid[x].velocity;
  //CFL Condtions
  if (length(fluid[x].velocity) > data->maxVel){
    data->maxVel = length(fluid[x].velocity);
  }
}
