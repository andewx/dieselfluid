#version 330
out vec4 fragColor;
in vec4 fluidColor; // the input variable from the vertex shader (same name and same type)

void main()
{
    fragColor = fluidColor;
}