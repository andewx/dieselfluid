#version 430 core


//IISPH Shader - GLSL Compute Shader Executes the advection and pressure error solver
//Note that the excess particles are the appended particle values which are not directly
//Invoked by the compute shader

#define PI 3.1415926
#define MAX_FLOAT 99999.0
#define LOCAL_GROUP_SIZE 4

layout(local_size_x = LOCAL_GROUP_SIZE, local_size_y = LOCAL_GROUP_SIZE, local_size_z = LOCAL_GROUP_SIZE) in;

struct Particle{
  vec3 position;
  vec3 velocity;
  vec3 force;
  float pressure;
  float density;
};

struct TempParticle{
  vec3 vel;
  vec3 pos;
  float pressure;
}

layout(std430, binding = 0) buffer FluidBuffer
{
  Particle particles[];
}fluid;



layout(std430, binding = 1) buffer FluidBuffer
{
  Particle particles[];
}fluid;



layout(std430, binding = 2) buffer IntData
{
  int particles;
  int boundarys;
  int buckets;
  int bucket_size;
}sizes;


layout(std430, binding = 3) buffer LocalityHash{
  int data[];
}locality;

layout(std430, binding = 4) buffer Hyperplane{
  vec3 vec[8];
}vectors;

layout(std430, binding = 5) buffer FloatData{
  float time_step;
  float mass;
  float delta;
  float maxVel;
  float h0;
}fluid_data;

layout(std430, binding = 6) buffer PCISPH{
  TempParticle data[];
}_temp;




//Kernel Functions
float f(float x, float h){
  if (x >= h){
    return 0.0;
  }
  float q = (1-(x*x/h*h));
  return ((315.0)/(64*PI*h*h*h))*q*q;
}

float f(float x){
  float h = fluid_data.h0;
  if (x >= h){
    return 0.0;
  }
  float q = (1-(x*x/h*h));
  return ((315.0)/(64*PI*h*h*h))*q*q;
}

float o1d(float x, float h){
  if (x >= h){
    return 0.0;
  }
  float q = (1-(x/h));
  return ((-45.0)/(PI*h*h*h*h))*q*q;
}

float o1d(float x){
  float h = fluid_data.h0;
  if (x >= h){
    return 0.0;
  }
  float q = (1-(x/h));
  return ((-45.0)/(PI*h*h*h*h))*q*q;
}

float o2d(float x, float h){
  if (x >= h){
    return 0.0;
  }
  float q = (1-(x/h));
  return ((90.0)/(PI*h*h*h*h*h))*q;
}

float o2d(float x){
  float h = fluid_data.h0;
  if (x >= h){
    return 0.0;
  }
  float q = (1-(x/h));
  return ((90.0)/(PI*h*h*h*h*h))*q;
}

vec3 f_grad(float x, float h,vec3 dir){
  if (x >= h){
    return vec3(0.0,0.0,0.0);
  }
  return -o1d(x,h)*dir;
}

vec3 f_grad(float x, vec3 dir){
  float h = fluid_data.h0;
  if (x >= h){
    return vec3(0.0,0.0,0.0);
  }
  return -o1d(x,h)*dir;
}


//Hashing Function given hyperplane
int sgn(float x){
  if (x <= 0){
    return 0;
  }else{
    return 1;
  }
}

//Simplified matching CFL constraint since GPU shader is evaluating the CFL condition
float cfl(){
  if(fluid_data.maxVel != 0.0){
  	return (1.0 / fluid_data.maxVel);
  }
  return 0.5;
}


//Hashing Algorithm of LSH
int lsh_h(vec3 v){
  int hash_val = 0;
  for(int i = 0; i < 8; i++){
    hash_val << 1;
    hash_val += sgn(dot(vectors.vec[i], v));
  }
  return hash_val;
}

int get_sample(vec3 position, int offset){
  int m_hash = lsh_h(position);
  return locality.data[m_hash*sizes.buckets + offset];
}

float interp_density(vec3 position){
  float sum = 0.0;
  for(int i = 0; i < 150; i++){
    int j = get_sample(position,i);
    if (j != 0){
      float x = length(position - fluid.particles[j].position);
      sum += fluid_data.mass * f(x);
    }
  }
  return sum;
}

float particle_density(int x){
  float sum = f(0);
  for(int i = 0; i < 150; i++){
    int j = get_sample(fluid.particles[x].position,i);
    if (j != 0 && j != x){
      float x = length(fluid.particles[x].position - fluid.particles[j].position);
      sum += fluid_data.mass * f(x);
    }
  }
  return sum;
}

float tait_eos(float x){
  float g = 7.16;
  float w = 2.15;
  float d0 = 1000.0;
  float p0 = 1013.25;
  float y = (w/g)*(pow(x/d0,g)-1.0);
  if (y <= 0.0) {
    return 0.0;
  }
  return y;
}
void evaluate_pressure(int particle){
      fluid.particles[particle].pressure = tait_eos(fluid.particles[particle].density);
}


vec3 pressure_force(int x){
  vec3 force = vec3(0.0,0.0,0.0);
  float di = fluid.particles[x].density;
  float di2 = di*di;
  float m = -(fluid_data.mass*fluid_data.mass);
  float pi = fluid.particles[x].pressure;
  float pi2 = pi * pi;

  for (int i = 0; i < 150; i++){
    int j = get_sample(fluid.particles[x].position, i);
    if (j != 0 && j != x){
      float dj = fluid.particles[j].density;
      float dj2 = dj*dj;
      float pj = fluid.particles[j].pressure;
      vec3 dir = fluid.particles[x].position - fluid.particles[j].position;
      force += m*(pi/di2+pj/dj2)*f_grad(length(dir), 1.0, dir);
    }
  }
  fluid.particles[x].force += force;
  return force;
}

vec3 viscosity_force(int x){
  vec3 force = vec3(0.0,0.0,0.0);
  float m = (fluid_data.mass);

  for (int i = 0; i < 150; i++){
    int j = get_sample(fluid.particles[x].position,i);
    if( j != 0 && j != x){
      float dj = 1/fluid.particles[j].density;
      vec3 dir = fluid.particles[x].position - fluid.particles[j].position;
      vec3 vel_diff = (fluid.particles[j].velocity -  fluid.particles[x].velocity)*(dj);
      force += vel_diff*m*o2d(length(dir));
    }
  }
  fluid.particles[x].force += force;
  return force;
}

void external_force(int x){
  fluid.particles[x].force += vec3(0.0, -fluid_data.mass*-9.81,0.0);
}

float clamp_greater(float x, float a){
  if (x >= a){
    return  x;
  }else{
    return a;
  }
}



void pressure_solve(int x){
  float d0 = 1000.0;
  float w = 0.5;
  int iters = 0;
  int max_iters = 5;
  float error = fluid.particles[x].pressure-d0;
  float delta = 0.1;

  vec3 predict_vel = fluid.particles[x].velocity + cfl() * fluid.particles[x].force/fluid_data.mass;
  vec3 predict_pos = fluid.particles[x].position + cfl() * fluid.particles[x].velocity;
  //Predict Velocity and Positions
  _velocity.data[x] = predict_vel;
  _position.data[x] = predict_pos;

  while(error > 0.2 && iters <= max_iters){

    //Compute Pressure from the density error
    float density = particle_density(x);
    float error = density - d0;
    float pressure = fluid_data.delta * error;
    fluid.particles[x].pressure += pressure;
    fluid.particles[x].density = density;

    //Update predicted positions and velocity
    pressure_force(x);
    _velocity.data[x] = _velocity.data[x] + fluid_data.time_step * fluid.particles[x].force/fluid_data.mass;
    _position.data[x] = _position.data[x] + fluid_data.time_step * fluid.particles[x].velocity;
    iters++;
  }
}


//Compute Gradient Level Forces + Projected Intermediate Velocities and Positions
void main(){
  int x = int(gl_LocalInvocationIndex);
  float ts = cfl();
  viscosity_force(x);
  external_force(x);
  pressure_solve(x);
  fluid.particles[x].velocity = fluid.particles[x].velocity + ts* fluid.particles[x].force/fluid_data.mass;
  fluid.particles[x].position = fluid.particles[x].position + ts*fluid.particles[x].velocity;
  //CFL Condtions
  if (length(fluid.particles[x].velocity) > fluid_data.maxVel){
    fluid_data.maxVel = length(fluid.particles[x].velocity);
  }
}
