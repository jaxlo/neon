#include "player/player.h"
#include "world/world.h"
#include "renderer/renderer.h"
#include <raylib.h>
#include <stdio.h>


int main(void) {
    // Make a window half the size of 1920x1080 to start
    InitWindow(960, 540, "neon");

    // Get monitor size for fullscreen
    int monitorWidth = GetMonitorWidth(0);
    int monitorHeight = GetMonitorHeight(0);
    
    // Render at half resolution
    int renderWidth = monitorWidth / 2;
    int renderHeight = monitorHeight / 2;
    printf("DEV: Render resolution: %dx%d\n", renderWidth, renderHeight);
    
    // Enable fullscreen
    SetWindowState(FLAG_FULLSCREEN_MODE);

    // Initialize renderer (0.5 = half resolution)
    InitRenderer(GetScreenWidth(), GetScreenHeight(), 0.5f);

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
        
        EndRenderFrame();
        
        // Draw UI at full resolution
        DrawFPS(10, 10);
        
        EndDrawing();
        // End draw --- --- ---
    }

    // De-Initialization
    UnloadWorld(&world);
    CleanupRenderer();
    CloseWindow();

    return 0;
}
