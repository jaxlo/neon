#ifndef PLAYER_H
#define PLAYER_H

#include <raylib.h>

typedef struct {
    Vector3 position;
    Vector3 velocity;
    Camera3D camera;
    float speed;
    float jumpSpeed;
    float gravity;
    bool isGrounded;
    float mouseX;
    float mouseY;
    float pitch;
    float yaw;
} Player;

// Function declarations
Player InitPlayer(void);
void UpdatePlayer(Player *player);
void HandleInput(Player *player);
void UpdatePlayerCamera(Player *player);

#endif
