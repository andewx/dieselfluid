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
  float density;
}

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
      sum += fluid_data.mass * f(x, fluid_data.h0);
    }
  }
  return sum;
}

float particle_density(int x){
  float sum = f(0);
  for(int i = 0; i < 150; i++){
    int j = get_sample(fluid.particles[x].position,i);
    if (j != 0){
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



//Density Compute Shader
void main(){
  uint x = gl_LocalInvocationIndex;
  fluid.particles[x].density = interp_density(fluid.particles[x].position);
  evaluate_pressure(int(x));
}