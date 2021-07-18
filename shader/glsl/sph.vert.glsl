#version 330
layout (location=0) in vec3 aPos;
uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;
uniform mat3 rotX;
uniform mat3 rotY;
uniform mat4 rotOriginTrans0; //Positive Translation Rotation
uniform mat4 rotOriginTrans1; //Negative Translation Rotation Origin
uniform int fMode;

void main()
{

    //Apply rotation
    vec4 rotPos = rotOriginTrans1 * vec4(aPos,1.0);
    vec3 pointRotate = rotX * rotY * vec3(rotPos.x, rotPos.y, rotPos.z);
    vec4 nPos = rotOriginTrans0 * vec4(pointRotate,1.0);
    gl_Position = projection * view * model * nPos;

}
