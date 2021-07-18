#version 330
out vec4 fragColor;

void main()
{
  float x = gl_PointCoord.x * 2.0 - 1.0;
  float y = gl_PointCoord.y * 2.0 - 1.0;

  float z = sqrt(1.0 - (pow(x, 2.0) + pow(y, 2.0)));

  vec3 position = vec3(x, y, z);

  float mag = dot(position.xy, position.xy);
  vec3 normal_col = normalize(position) * 0.5 + 0.5;
  vec3 black_col = vec3(0.1,0.1,0.1);
  fragColor = vec4(normal_col, 0.5);
  if(mag > 1.0){
      fragColor = vec4(black_col, 0.1);
  }



}
