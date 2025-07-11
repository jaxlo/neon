#include "shaders.h"
#include <raylib.h>

// Global shader variables
Shader basicShader = {0};
int timeUniform = 0;
int resolutionUniform = 0;

void InitShaders(void) {
    // For now, we'll just initialize without loading custom shaders
    // In a more complex project, you would load vertex and fragment shaders here
    
    // Example of how you might load a custom shader:
    // basicShader = LoadShader("shaders/basic.vs", "shaders/basic.fs");
    // timeUniform = GetShaderLocation(basicShader, "time");
    // resolutionUniform = GetShaderLocation(basicShader, "resolution");
    
    // For this basic example, we'll use default rendering
}

void UnloadShaders(void) {
    // Unload any custom shaders if they were loaded
    // if (basicShader.id != 0) {
    //     UnloadShader(basicShader);
    // }
}
