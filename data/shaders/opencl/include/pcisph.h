
#ifndef pcisph_h
#define pcisph_h
#endif

//Function headers and structure definitions

#define PI 3.1415926
#define MAX_FLOAT 99999.0
#define LARGE_PRIME 721546
#define LOCAL_GROUP_SIZE 4

typedef struct m_particle{
  float3 position;
  float3 velocity;
  float3 force;
  float pressure;
  float density;
}particle;

typedef struct t_particle{
  float3 vel;
  float3 pos;
  float density;
}temp_particle;

typedef struct m_intdata{
  int particles;
  int boundarys;
  int buckets;
  int bucket_size;
}intdata;

typedef struct m_floatdata{
  float time_step;
  float mass;
  float delta;
  float maxVel;
  float h0;
}floatdata;



//Function prototypes
float f(float x, float h);

float o1d(float x, float h);

float o2d(float x, float h);

float3 f_grad(float x, float h,float3 dir);

int sgn(float x);

int hash(float3 v,global float3* n);

int get_sample(float3 position, int offset, global int* table,global intdata *sizes, global float3* rands);

float interp_density(float3 position,
                     global floatdata *fluid_data,
                     global particle *fluid,
                     global int* table,
                     global intdata *sizes ,
                     global float3* rands);

float particle_density(int x,
                      global floatdata *fluid_data,
                      global particle *fluid,
                      global int* table,
                      global intdata *sizes ,
                      global float3* rands);

float tait_eos(float x);

void evaluate_pressure(int x, global particle *fluid);


float3 pressure_force(int x,
                      global floatdata *fluid_data, global particle *fluid,
                      global int* table,
                      global intdata *sizes ,
                      global float3* rands);

float3 viscosity_force(int x,
                      global floatdata *fluid_data, global particle *fluid,
                      global int* table,
                      global intdata *sizes ,
                      global float3* rands);


void external_force(int x, global floatdata *fluid_data, global particle *fluid);

float clamp_greater(float x, float a);

void pressure_solve(int x,
                    global temp_particle *temp,
                    global floatdata *fluid_data, global particle *fluid,
                    global int* table,
                    global intdata *sizes ,
                    global float3* rands);

float cfl(float maxvel);
//Kernel Functions
float f(float x, float h){
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


float o2d(float x, float h){
  if (x >= h){
    return 0.0;
  }
  float q = (1-(x/h));
  return ((90.0)/(PI*h*h*h*h*h))*q;
}


float3 f_grad(float x, float h,float3 dir){
  if (x >= h){
    return float3(0.0,0.0,0.0);
  }
  return -o1d(x,h)*dir;
}

float cfl(float maxvel){
  if (maxvel > 2.0){
    return 1/maxvel;
  }
  return 0.5;
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
int hash(float3 v,global float3* n){
  int hash_val = 0;
  for(int i = 0; i < 8; i++){
    hash_val = hash_val << 1;
    hash_val += sgn(dot(n[i], v));
  }
  return hash_val;
}

int get_sample(float3 position, int offset, global int* table, global intdata *sizes, global float3* rands){
  int m_hash = hash(position, rands);
  int h = (m_hash*sizes->buckets + offset)%(sizes->buckets*sizes->bucket_size);
  return table[h];
}


float interp_density(float3 position, global floatdata *fluid_data, global particle *fluid, global int* table, global intdata *sizes , global float3* rands){
  float sum = 0.0;
  for(int i = 0; i < 150; i++){
    int j = get_sample(position,i, table, sizes,rands);
    if (j != 0){
      float dist = length(position - fluid[j].position);
      sum += fluid_data->mass * f(dist, fluid_data->h0);
    }
  }
  return sum;
}

float particle_density(int x,
                      global floatdata *fluid_data, global particle *fluid,
                      global int* table,
                      global intdata *sizes ,
                      global float3* rands){
  float sum = f(0,fluid_data->h0);
  for(int i = 0; i < 150; i++){
    int j = get_sample(fluid[x].position,i, table, sizes,rands);
    if (j != 0){
      float dist = length(fluid[x].position - fluid[j].position);
      sum += fluid_data->mass * f(dist, fluid_data->h0);
    }
  }
  return sum;
}


float tait_eos(float x){
  float g = 7.16;
  float w = 2.15;
  float d0 = 1000.0;
  float p0 = 0.0;
  float y = (w/g)*(pow(x/d0,g)-1.0) + p0;
  if (y <= 0.0) {
    return 0.0;
  }
  return y;
}

void evaluate_pressure(int x, global particle *fluid){
      fluid[x].pressure = tait_eos(fluid[x].density);
}


float3 pressure_force(int x, global floatdata *fluid_data, global particle *fluid, global int* table, global intdata *sizes, global float3* rands){
  float3 force = float3(0.0,0.0,0.0);
  float di = fluid[x].density;
  float di2 = di*di;
  float m = -(fluid_data->mass*fluid_data->mass);
  float pi = fluid[x].pressure;

  for (int i = 0; i < 150; i++){
    int j = get_sample(fluid[x].position, i,table,sizes,rands);
    if (j != 0 && j != x){
      float dj = fluid[j].density;
      float dj2 = dj*dj;
      float pj = fluid[j].pressure;
      float3 dir = fluid[x].position - fluid[j].position;
      force += m*(pi/di2+pj/dj2)*f_grad(length(dir), fluid_data->h0, dir);
    }
  }
  fluid[x].force += force;
  return force;
}

float3 viscosity_force(int x, global floatdata *fluid_data, global particle *fluid, global int* table, global intdata *sizes, global float3* rands){
  float3 force = float3(0.0,0.0,0.0);
  float m = (fluid_data->mass);

  for (int i = 0; i < 150; i++){
    int j = get_sample(fluid[x].position,i,table,sizes,rands);
    if( j != 0 && j != x){
      float dj = 1/fluid[j].density;
      float3 dir = fluid[x].position - fluid[j].position;
      float3 vel_diff = (fluid[j].velocity -  fluid[x].velocity)*(dj);
      force += vel_diff*m*o2d(length(dir),fluid_data->h0);
    }
  }
  fluid[x].force += force;
  return force;
}

void external_force(int x, global floatdata *fluid_data, global particle *fluid){
  fluid[x].force += float3(0.0, -fluid_data->mass*-9.81,0.0);
}

float clamp_greater(float x, float a){
  if (x >= a){
    return  x;
  }else{
    return a;
  }
}

void pressure_solve(int x , global temp_particle *temp, global floatdata *fluid_data, global particle *fluid, global int* table, global intdata *sizes, global float3* rands){
  float d0 = 1000.0;
  int iters = 0;
  int max_iters = 5;
  float error = fluid[x].pressure-d0;

  float3 predict_vel = fluid[x].velocity + cfl(fluid_data->maxVel) * fluid[x].force/fluid_data->mass;
  float3 predict_pos = fluid[x].position + cfl(fluid_data->maxVel) * fluid[x].velocity;
  //Predict Velocity and Positions
  temp[x].vel = predict_vel;
  temp[x].pos = predict_pos;

  while(error > 0.2 && iters <= max_iters){

    //Compute Pressure from the density error
    float density = particle_density(x, fluid_data, fluid,table,sizes,rands);
    float error = density - d0;
    float pressure = fluid_data->delta * error;
    fluid[x].pressure += pressure;
    fluid[x].density = density;

    //Update predicted positions and velocity
    pressure_force(x, fluid_data,fluid,table,sizes,rands);
    temp[x].vel = temp[x].vel + fluid_data->time_step * fluid[x].force/fluid_data->mass;
    temp[x].pos = temp[x].pos + fluid_data->time_step * fluid[x].velocity;
    iters++;
  }
}
