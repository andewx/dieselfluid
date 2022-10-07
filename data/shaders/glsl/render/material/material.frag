#version 330

layout(location = 0) out vec4 fragColor;


in vec3 norm;  //model_space normal
in vec3 eye;   //view space position
in vec2 uv;    //texture uv coordinates

uniform vec4  baseColor;
uniform float metallness;
uniform float roughness;
uniform float fresnel_rim;

uniform vec3 lightColor;
uniform vec3 lightPos;


uniform mat4 mvp;
uniform mat4 model;
uniform mat4 view;

uniform sampler2D colorTex;
uniform sampler2D normTex;
uniform samplerCube cube;



#define PI 3.1415926


// simple blinn specular calculation with normalization
vec3 blinn_specular(in float NdH, in vec3 specular, in float roughness)
{
    float k = 1.999 / (roughness * roughness);

    return min(1.0, 3.0 * 0.0398 * k) * pow(NdH, min(10000.0, k)) * specular;
}

// phong (lambertian) diffuse term
float phong_diffuse()
{
    return (1.0 / PI);
}

// compute fresnel specular factor for given base specular and product
// product could be NdV or VdH depending on used technique
float fresnel_factor(in float f0, in float ndv)
{
    return f0 + (1-f0)*pow(1-ndv,5.0);
}

//---------------------------------------------
// PBR Metal-Roughness Distribution Microfacet Model Stubs
//---------------------------------------------

//Blinn roughness distribution
float D_blinn(in float roughness, in float NdH)
{
    float m = roughness * roughness;
    float m2 = m * m;
    float n = 2.0 / m2 - 2.0;
    return (n + 2.0) / (2.0 * PI) * pow(NdH, n);
}

//Returns normalized distribution parameter for NdH calc
float D_beckmann(in float roughness, in float NdH)
{
    float m = roughness * roughness;
    float m2 = m * m;
    float NdH2 = NdH * NdH;
    return exp((NdH2 - 1.0) / (m2 * NdH2)) / (PI * m2 * NdH2 * NdH2);
}

//Returns normalized distribution for NdH
float D_CGX(in float roughness, in float NdH)
{
    float m = roughness * roughness;
    float m2 = m * m;
    float d = (NdH * NdH) * (m2-1)  + 1.0;
    return m2 / (PI * d * d);
}



//Returns schlick G() scaling parameter
float G_schlick(in float roughness, in float NdV, in float NdL)
{
    float k = roughness * roughness * 0.5;
    float V = NdV * (1.0 - k) + k;
    float L = NdL * (1.0 - k) + k;
    return 0.25 / (V * L);
}


//Geometric Attenutation term
float G_term(in float HdN, in float VdN, in float VdH, in float LdN, in float rough){
  float m3 = rough*rough;
  float a = 2*HdN*VdN*m3;
  float b = 2*HdN*LdN*m3;
  float c = min(1,a/VdH);
  return min(c,b/VdH);
}


// simple phong specular calculation with normalization
vec3 phong_specular(in vec3 V, in vec3 L, in vec3 N, in vec3 specular, in float roughness)
{
    vec3 R = reflect(-L, N);
    float spec = max(0.0, dot(V, R));

    float k = 1.999 / (roughness * roughness);

    return min(1.0, 3.0 * 0.0398 * k) * pow(spec, min(10000.0, k)) * specular;
}


// cook-torrance specular calculation
float cooktorrance_specular(in float HdN, in float VdN, in float VdH, in float LdN, in float NdH, in float NdV, in float NdL, in float roughness, in float F)
{
  float D = D_CGX(roughness, NdH);
  float G = G_term(HdN,VdN,VdH,LdN, roughness);
  float P = 1/mix(1.0 - roughness*0.9,1.0,NdV);
  return ((D*F*G*P)/(3.1415926*VdN*NdL));
}

//Calculate Tangent-Bitangent Normal Matrix aligned with the Model UV texture
//coordinates from shader code. Ultimately TBN Matrixes should be handled by
//application software host so the shader routines don't have to sample each pixel
//so many times
mat3 tbn_matrix( vec3 N, vec3 P, vec2 UV )
{
    // Calcuate Point Curve Deltas and UV direction deltas
    vec3 deltaPos1 =    dFdx(P);
    vec3 deltaPos2 =    dFdy(P);
    vec2 deltaUV1  =    dFdx(UV);
    vec2 deltaUV2  =    dFdy(UV);

    // Create TBN 3x3 Matrix
    vec3 X = cross(deltaPos2, N );
    vec3 Y = cross( N, deltaPos1 );
    vec3 T = X * deltaUV1.x + X * deltaUV2.x;
    vec3 B = Y * deltaUV1.y + X * deltaUV2.y;
    mat3 tbn = mat3(T,B,N);

    //Inverse is same as transpose
    return transpose(tbn);
}

mat3 cotangent_frame( vec3 N, vec3 p, vec2 uv )
{
    // get edge vectors of the pixel triangle
    vec3 dp1 = dFdx( p );
    vec3 dp2 = dFdy( p );
    vec2 duv1 = dFdx( uv );
    vec2 duv2 = dFdy( uv );

    // solve the linear system
    vec3 dp2perp = cross( dp2, N );
    vec3 dp1perp = cross( N, dp1 );
    vec3 T = dp2perp * duv1.x + dp1perp * duv2.x;
    vec3 B = dp2perp * duv1.y + dp1perp * duv2.y;

    // construct a scale-invariant frame
    float invmax = inversesqrt( max( dot(T,T), dot(B,B) ) );
    return mat3( T * invmax, B * invmax, N );
}

//normal mapping function returns perturbed normal according to TBN matrix value
vec3 normal_map( vec3 N, vec3 V, vec2 uv)
{
    vec3 map = texture(normTex,uv).rgb*2.0-1.0;
    mat3 tbn = cotangent_frame(N,-V,uv);
    return normalize(tbn * map);
}


void main() {
    // Attenuation of Point Light (This gets normalized to the position at each)
    vec3 l_pos0 = (view * vec4(lightPos,1.0)).xyz;
    float mL = 1/length(l_pos0);
    float mLE = length(l_pos0-eye);
    float A = 1/(mLE*mL);
    // L, V, H vectors
    vec3 L = normalize(l_pos0 - eye);
    vec3 V = normalize(-eye);
    vec3 N = normalize(norm);
    vec3 H = normalize(L + V);
    vec3 R = reflect(-V,N);

    float metallic = clamp(0.01, 1.0, metallness);
    float rough = clamp(0.01, 1.0, roughness);

    N = normal_map(N,V,uv);
    vec3 base = texture(colorTex,uv).xyz;

    // roughness
    #ifdef USE_ROUGHNESS_MAP
         rough = texture(roughTex, uv).y * roughness;
    #else
         rough = roughness;
    #endif

    //Get specular highlight interpolation color
    vec3 sp_color = mix(base*(1/rough), vec3(1.0), metallic);

    float NdL = max(0.0, dot(N, L));
    float NdV = max(0.001, dot(N, V));
    float NdH = max(0.001, dot(N, H));
    float LdV = max(0.001, dot(L, V));
    float LdN = max(dot(L,N), 0.0);
    float HdN = max(0.001, dot(H, N));
    float HdV = max(0.001, dot(H, V));
    float VdH = max(dot(V,H),0.000001);
    float VdN = max(dot(V,N), 0.0);


    /*--------------------------------------------------*/
    //    Specular
    /*--------------------------------------------------*/

    //Compute Specular Power
    vec3 light_color = lightColor * A;
    sp_color = light_color * sp_color;
    float F0 = (1.0002 - 0.273)/(1.0002+0.273); //fresnel gold
    F0 *= F0;
    float fresnel = fresnel_factor(F0, VdN);
    float power = cooktorrance_specular(HdN,VdN,VdH, LdN, NdH, NdV,NdL,rough,fresnel);
    power = max(0.00,power);
    vec3 specular = vec3(power)*sp_color;
    specular *= vec3(NdL);

    /*--------------------------------------------------*/
    //    Diffuse + Ambient
    /*--------------------------------------------------*/
    //Compute phong shading diffuse terms for single light source


    vec3 lambert = light_color *phong_diffuse() * NdL;
    vec3 refl_light = vec3(0);
    vec3 diff_light = light_color * 0.1;

    refl_light += specular;
    diff_light += lambert;

    #ifdef USE_DIFFUSE_IRRADIANCE_MAP
        vec2 brdf = texture(diffuseTex, vec2(rough, 1.0 - NdV)).xy;
        vec3 spec = min(vec3(0.99), fresnel_factor(sp_color, NdV) * brdf.x + brdf.y);
        refl_light += spec * envspec;
        diff_light += envdiff * (1.0 / PI);
    #endif

    vec3 cube_color = texture(cube,R).rgb;
    float corr = clamp(0.0,1.0,(metallic-rough)/metallic);
    vec3 res = (diff_light * base)+ refl_light + cube_color*corr*(1-rough);
    fragColor = vec4(res,1.0);
}
