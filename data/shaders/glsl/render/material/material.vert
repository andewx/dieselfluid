#version 330
layout (location = 0) in vec3 pos;
layout (location = 3) in vec3 normal;
layout (location = 6) in vec2 tex_uv;

out vec3 norm;  //model_space normal
out vec3 eye;   //view space position
out vec2 uv;    //texture uv coordinates

uniform mat4 mvp;
uniform mat4 model;
uniform mat4 view;
uniform mat4 normMat;


#define USE_COOK_TORRANCE
#define COOK_CGX
#define USE_BASE_MAP

void main()
{
    mat4 modelview = view*model;
    vec4 position = mvp * view* model * vec4(pos,1.0);
    norm = normal;
    eye = (view * model* vec4(pos,1.0)).xyz;
    uv = tex_uv;
    gl_Position = position;
}
