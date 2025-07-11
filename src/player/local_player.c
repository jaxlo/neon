#include "player.h"
#include <raylib.h>

Player InitPlayer(void) {
    Player player = {0};
    
    // Initialize player position
    player.position = (Vector3){0.0f, 2.0f, 0.0f};
    player.velocity = (Vector3){0.0f, 0.0f, 0.0f};
    
    // Initialize camera
    player.camera.position = player.position;
    player.camera.target = (Vector3){0.0f, 2.0f, 1.0f};
    player.camera.up = (Vector3){0.0f, 1.0f, 0.0f};
    player.camera.fovy = 60.0f;
    player.camera.projection = CAMERA_PERSPECTIVE;
    
    // Initialize movement parameters
    player.speed = 5.0f;
    player.jumpSpeed = 8.0f;
    player.gravity = 20.0f;
    player.isGrounded = false;
    
    // Initialize mouse look
    player.mouseX = 0.0f;
    player.mouseY = 0.0f;
    player.pitch = 0.0f;
    player.yaw = 0.0f;
    
    DisableCursor();
    
    return player;
}

void UpdatePlayer(Player *player) {
    HandleInput(player);
    
    // Apply gravity
    if (!player->isGrounded) {
        player->velocity.y -= player->gravity * GetFrameTime();
    }
    
    // Update position
    player->position.x += player->velocity.x * GetFrameTime();
    player->position.z += player->velocity.z * GetFrameTime();
    player->position.y += player->velocity.y * GetFrameTime();
    
    // Simple ground collision (y = 0)
    if (player->position.y <= 2.0f) {
        player->position.y = 2.0f;
        player->velocity.y = 0.0f;
        player->isGrounded = true;
    } else {
        player->isGrounded = false;
    }
    
    // Update camera after position changes
    UpdatePlayerCamera(player);
}
