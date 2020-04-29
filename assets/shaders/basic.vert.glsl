#version 410 core

layout (location = 0) in vec3 position;

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

void main()
{
    gl_Position = projection * camera * model * vec4(position, 1.0);
}