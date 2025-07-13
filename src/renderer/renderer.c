#include "renderer.h"
#include <stdio.h>
#include <stdlib.h>

// Global renderer instance
Renderer g_renderer = {0};

void InitRenderer(int screenWidth, int screenHeight, float scaleFactor) {
    if (g_renderer.initialized) {
        printf("Warning: Renderer already initialized!\n");
        return;
    }

    // Setup configuration
    g_renderer.config.screenWidth = screenWidth;
    g_renderer.config.screenHeight = screenHeight;
    g_renderer.config.scaleFactor = scaleFactor;
    g_renderer.config.renderWidth = (int)(screenWidth * scaleFactor);
    g_renderer.config.renderHeight = (int)(screenHeight * scaleFactor);
    
    // Create render texture at scaled resolution
    g_renderer.renderTexture = LoadRenderTexture(g_renderer.config.renderWidth, g_renderer.config.renderHeight);
    
    // Compute shaders
    //g_renderer.raytraceShader = LoadComputeShader("raytracer.comp");

    // regular shaders
    g_renderer.lightingShader = LoadShader(0, "shaders/lighting.fs"); // Shader not implemented yet
    g_renderer.shadowShader = LoadShader(0, "shaders/shadows.fs"); // Shader not implemented yet
    g_renderer.ditherShader = LoadShader(0, "assets/shaders/dither.fs"); // Shader not implemented yet
    
    // Set shader uniforms
    float resolution[2] = {(float)g_renderer.config.screenWidth, (float)g_renderer.config.screenHeight};
    SetShaderValue(g_renderer.ditherShader, GetShaderLocation(g_renderer.ditherShader, "resolution"), resolution, SHADER_UNIFORM_VEC2);

    g_renderer.initialized = true;

    printf("Renderer initialized: %dx%d -> %dx%d (scale: %.2f)\n", 
        g_renderer.config.renderWidth,
        g_renderer.config.renderHeight,
        g_renderer.config.screenWidth,
        g_renderer.config.screenHeight,
        g_renderer.config.scaleFactor
    );
}

void CleanupRenderer(void) {
    if (!g_renderer.initialized) return;
    
    UnloadRenderTexture(g_renderer.renderTexture);
    UnloadShader(g_renderer.raytraceShader);
    UnloadShader(g_renderer.lightingShader);
    UnloadShader(g_renderer.shadowShader);
    UnloadShader(g_renderer.ditherShader);

    g_renderer.initialized = false;
    printf("Renderer cleanup complete\n");
}

void BeginRenderFrame(void) {
    if (!g_renderer.initialized) {
        printf("Error: Renderer not initialized!\n");
        return;
    }
    
    // Begin rendering to low-resolution texture
    BeginTextureMode(g_renderer.renderTexture);
    ClearBackground(DARKBLUE);
}

void EndRenderFrame(void) {
    if (!g_renderer.initialized) return;
    
    // End rendering to texture
    EndTextureMode();
    
    // Begin final screen rendering
    BeginDrawing();
    ClearBackground(BLACK);
    
    // Apply dithering shader and upscale to full resolution
    BeginShaderMode(g_renderer.ditherShader);
    
    // Draw the render texture to screen with upscaling
    Rectangle sourceRec = {0, 0, (float)g_renderer.config.renderWidth, -(float)g_renderer.config.renderHeight}; // Flip Y
    Rectangle destRec = {0, 0, (float)g_renderer.config.screenWidth, (float)g_renderer.config.screenHeight};
    
    DrawTexturePro(g_renderer.renderTexture.texture, sourceRec, destRec, (Vector2){0, 0}, 0.0f, WHITE);
    
    EndShaderMode();
}

void BeginScene3D(Camera3D camera) {
    if (!g_renderer.initialized) return;
    
    BeginMode3D(camera);
    
    // Apply lighting shader if needed
    // BeginShaderMode(g_renderer.lightingShader);
}

void EndScene3D(void) {
    if (!g_renderer.initialized) return;
    
    // EndShaderMode(); // End lighting shader if used
    EndMode3D();
}

bool IsRendererInitialized(void) {
    return g_renderer.initialized;
}

Vector2 GetRenderResolution(void) {
    return (Vector2){(float)g_renderer.config.renderWidth, (float)g_renderer.config.renderHeight};
}

Vector2 GetScreenResolution(void) {
    return (Vector2){(float)g_renderer.config.screenWidth, (float)g_renderer.config.screenHeight};
}

float GetRenderScale(void) {
    return g_renderer.config.scaleFactor;
}
