#type vertex
#version 420 core

layout (location=0) in vec3 aPos;
layout (location=1) in vec4 aColor;
layout (location=2) in vec2 aTexCoords;
layout (location=3) in float aTexId;

uniform mat4 uProjection;

out vec4 fColor;
out vec2 fTexCoords;
out float fTexId;

void main()
{
    fColor = aColor;
    fTexCoords = aTexCoords;
    fTexId = aTexId;
    gl_Position = uProjection * vec4(aPos,1.0);
}

#type fragment
#version 420 core

in vec4 fColor;
in vec2 fTexCoords;
in float fTexId;
out vec4 color;

uniform sampler2D Texture;

void main()
{
    if (fTexId > 0) {
//        color = fColor * texture(Texture,fTexCoords);
        vec4 tC = texture(Texture,fTexCoords);
        color =  fColor * tC;
    } else {
        color = fColor;
    }
    // if (fTexId > 0) {
    //    int id = int(fTexId);
    //    float c = texture(uTextures[id], fTexCoords).a;
        // color = vec4(1, 1, 1, c) * fColor;
        // color = fColor * texture(uTextures[id],fTexCoords);
        // color = vec4(fTexCoords, 0, 1);//чтобы узнать u v координаты
    // } else {
        // color = fColor;
    // }

}
