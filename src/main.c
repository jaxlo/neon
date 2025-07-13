#include "player/player.h"
#include "world/world.h"
#include "renderer/renderer.h"
#include "ui/ui.h"
#include <raylib.h>
#include <stdio.h>


int main(void) {
    InitWindow(GetMonitorWidth(0), GetMonitorHeight(0), "neon");
    
    // Enable fullscreen
    SetWindowState(FLAG_FULLSCREEN_MODE);

    // Initialize renderer (0.5 = half resolution)
    InitRenderer(GetScreenWidth(), GetScreenHeight(), 1.0);

    // Initialize everything
    Player player = InitPlayer();
    World world = InitWorld();

    SetTargetFPS(200);


    // Main game loop
    while (!WindowShouldClose()) {
        // Update --- --- ---
        UpdatePlayer(&player);
        // End update --- --- ---
        
        // Draw --- --- ---
        BeginRenderFrame();
        
        BeginScene3D(player.camera);
        DrawWorld(&world);
        EndScene3D();

        // New renderer
        // Custom compute-style rendering approach for the ray traced scene (Built on top of rlgl)
        // BeginDrawing();
        //     BeginShaderMode(raytracingShader);
        //         // Draw a fullscreen quad to trigger your fragment shader
        //         DrawRectangle(0, 0, screenWidth, screenHeight, WHITE);
        //     EndShaderMode();
        // EndDrawing();
        
        EndRenderFrame();
        
        // Draw UI at full resolution
        DrawUI();
        
        EndDrawing();
        // End draw --- --- ---
    }

    // De-Initialization
    UnloadWorld(&world);
    CleanupRenderer();
    CloseWindow();

    return 0;
}
