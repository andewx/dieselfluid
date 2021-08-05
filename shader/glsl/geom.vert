#version 330
layout (location = 0) in vec3 position;
layout (location = 3) in vec3 normal;
out vec4 color;
uniform mat4 mvp;
uniform mat4 model;
uniform mat4 viewMat;

//TODO :: Textures / UV Load / Normal Maps / Interpolation
//TODO :: PBR BRSDF Metal-Diffuse Shader Reps
//TODO :: HDR CUBE MAP SAMPLERS
void main()
{

    vec3 light = vec3(100,100,100);
    float lightDist = length(light);
    float falloff = 1/lightDist;
    float cosT = clamp(dot(normal,light)*falloff, 0.0,1.0 );
    vec3 ambient= vec3(0.2, 0.2, 0.2);
    vec3 diffuse = vec3(0.5, 0.5, 0.5);
    vec3 ambientDiffuse = (ambient)  + (diffuse * cosT);
    color = vec4(ambientDiffuse,1.0);
    gl_Position = mvp * viewMat* model * vec4(position,1.0);

}
