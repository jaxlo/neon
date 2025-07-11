#ifndef SHADERS_H
#define SHADERS_H

#include <raylib.h>

// Shader management
void InitShaders(void);
void UnloadShaders(void);

// Shader uniforms
extern Shader basicShader;
extern int timeUniform;
extern int resolutionUniform;

#endif
