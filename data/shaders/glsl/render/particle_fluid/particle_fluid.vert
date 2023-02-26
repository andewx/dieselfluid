#version 330
layout (location=13) in vec3 aPos;
uniform mat4 _mvp;
uniform mat4 _view;
uniform mat4 _model;


void main()
{
    gl_Position = _mvp * _view * _model * vec4(aPos, 1.0);

}