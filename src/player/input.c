#include "player.h"
#include <raylib.h>
#include <raymath.h>

void HandleInput(Player *player) {
    // Get mouse delta with sensitivity
    Vector2 mouseDelta = GetMouseDelta();
    float mouseSensitivity = 0.002f; // Reduced sensitivity for smoother movement
    
    // Update mouse look
    player->yaw -= mouseDelta.x * -mouseSensitivity;
    player->pitch += mouseDelta.y * -mouseSensitivity;
    
    // Clamp pitch to prevent over-rotation
    if (player->pitch > 1.5f) player->pitch = 1.5f;
    if (player->pitch < -1.5f) player->pitch = -1.5f;
    
    // Calculate forward and right vectors (normalized)
    Vector3 forward = {
        cosf(player->yaw),
        0.0f,  // Keep forward vector on horizontal plane
        sinf(player->yaw)
    };
    
    Vector3 right = {
        cosf(player->yaw + PI/2),
        0.0f,
        sinf(player->yaw + PI/2)
    };
    
    // Reset horizontal velocity
    player->velocity.x = 0.0f;
    player->velocity.z = 0.0f;
    
    // Movement input with normalized vectors
    Vector3 moveVector = {0.0f, 0.0f, 0.0f};
    
    if (IsKeyDown(KEY_W)) {
        moveVector.x += forward.x;
        moveVector.z += forward.z;
    }
    if (IsKeyDown(KEY_S)) {
        moveVector.x -= forward.x;
        moveVector.z -= forward.z;
    }
    if (IsKeyDown(KEY_A)) {
        moveVector.x -= right.x;
        moveVector.z -= right.z;
    }
    if (IsKeyDown(KEY_D)) {
        moveVector.x += right.x;
        moveVector.z += right.z;
    }
    
    // Normalize diagonal movement to maintain consistent speed
    float moveLength = sqrtf(moveVector.x * moveVector.x + moveVector.z * moveVector.z);
    if (moveLength > 0.0f) {
        moveVector.x = (moveVector.x / moveLength) * player->speed;
        moveVector.z = (moveVector.z / moveLength) * player->speed;
    }
    
    player->velocity.x = moveVector.x;
    player->velocity.z = moveVector.z;
    
    // Jump input
    if (IsKeyPressed(KEY_SPACE) && player->isGrounded) {
        player->velocity.y = player->jumpSpeed;
        player->isGrounded = false;
    }
    
    // Exit on ESC
    if (IsKeyPressed(KEY_ESCAPE)) {
        EnableCursor();
    }
    
    // Re-enable mouse look on click
    if (IsMouseButtonPressed(MOUSE_BUTTON_LEFT)) {
        DisableCursor();
    }
}
