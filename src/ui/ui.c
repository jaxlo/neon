#include <raylib.h>


void DrawUI(void) {
    // UI is drawn after EndRenderFrame() but before EndDrawing()
    // This ensures UI is rendered at full resolution
    DrawText(TextFormat("FPS: %d", GetFPS()), 10, 10, 20, YELLOW);
}

// TODO add a debug area with a blackish background later