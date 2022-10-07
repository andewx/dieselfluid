#version 330
layout (location=8) in vec3 aPos;
uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;
uniform int mode;
out vec4 fluidColor; // specify a color output to the fragment shader


void main()
{
    gl_Position = projection * view * model * vec4(aPos, 1.0);

    fluidColor = vec4(0.3,0.3,0.3, 0.2); // Blue Points

    if(mode == 0 ){
    fluidColor = vec4(0.3,0.4,0.8, 1.0); // Blue Points
    }

}