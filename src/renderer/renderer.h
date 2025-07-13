#ifndef RENDERER_H
#define RENDERER_H

#include <raylib.h>

// Renderer configuration
typedef struct {
    int screenWidth;
    int screenHeight;
    int renderWidth;
    int renderHeight;
    float scaleFactor;
} RendererConfig;

// Renderer state
typedef struct {
    RendererConfig config;
    RenderTexture2D renderTexture;
    Shader raytraceShader;
    Shader lightingShader;
    Shader shadowShader;
    Shader ditherShader;
    bool initialized;
} Renderer;

// Global renderer instance
extern Renderer g_renderer;

// Function declarations
void InitRenderer(int screenWidth, int screenHeight, float scaleFactor);
void CleanupRenderer(void);
void BeginRenderFrame(void);
void EndRenderFrame(void);
void BeginScene3D(Camera3D camera);
void EndScene3D(void);
bool IsRendererInitialized(void);

// Utility functions
Vector2 GetRenderResolution(void);
Vector2 GetScreenResolution(void);
float GetRenderScale(void);

#endif // RENDERER_H
