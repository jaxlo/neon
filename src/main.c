#include "player/player.h"
#include "world/world.h"
#include "shaders/shaders.h"
#include <raylib.h>


int main(void)
{
    // Initialization
    const int screenWidth = 640;
    const int screenHeight = 480;

    InitWindow(screenWidth, screenHeight, "neon");

    // jump to fullscreen for now
    ToggleFullscreen();

    // Initialize everything
    Player player = InitPlayer();
    World world = InitWorld();
    InitShaders();

    SetTargetFPS(60);

    // Main game loop
    while (!WindowShouldClose())
    {
        // Update
        UpdatePlayer(&player);
        
        // Draw
        BeginDrawing();
        
            ClearBackground(SKYBLUE);
            
            BeginMode3D(player.camera);
            
                DrawWorld(&world);
                
            EndMode3D();
            
            // Draw UI
            DrawFPS(10, 10);
            
        EndDrawing();
    }

    // De-Initialization
    UnloadWorld(&world);
    UnloadShaders();
    CloseWindow();

    return 0;
}
