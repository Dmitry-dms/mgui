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
uniform int textureType;

float median(float r, float g, float b) {
    return max(min(r, g), min(max(r, g), b));
}

float screenPxRange() {
    return 4.5;
}

void main()
{
    if (fTexId > 0) {
        if (textureType == 1) {
            vec3 msd = texture(Texture, fTexCoords).rgb;
            float sd = median(msd.r, msd.g, msd.b);
            float screenPxDistance = screenPxRange()*(sd - 0.5);
            float opacity = clamp(screenPxDistance + 0.5, 0.0, 1.0);
            color = fColor * opacity;
        } else {
            vec4 tex = texture(Texture, fTexCoords);
            color = fColor * tex;
        }
    } else {
        color = fColor;
    }

}
