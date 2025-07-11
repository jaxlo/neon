#include "player.h"
#include <raylib.h>
#include <raymath.h>

void UpdatePlayerCamera(Player *player) {
    // Calculate camera target based on pitch and yaw
    Vector3 target = {
        player->position.x + cosf(player->yaw),
        player->position.y + sinf(player->pitch),
        player->position.z + sinf(player->yaw)
    };
    
    player->camera.target = target;
    player->camera.position = player->position;
}
