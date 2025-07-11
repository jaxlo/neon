#include "world.h"
#include <raylib.h>
#include <stdlib.h>
#include <time.h>

World InitWorld(void) {
    World world = {0};
    
    // Initialize random seed
    srand(time(NULL));
    
    // Generate buildings
    GenerateBuildings(&world);
    
    return world;
}

void GenerateBuildings(World *world) {
    world->buildingCount = 0;
    
    // Create a simple city grid
    for (int x = -5; x <= 5; x += 2) {
        for (int z = -5; z <= 5; z += 2) {
            if (world->buildingCount >= MAX_BUILDINGS) break;
            if (x == 0 && z == 0) continue; // Skip center for player spawn
            
            Building building = {0};
            building.position = (Vector3){x * 4.0f, 0.0f, z * 4.0f};
            
            // Random building dimensions
            building.size.x = 2.0f + (float)(rand() % 3);
            building.size.y = 3.0f + (float)(rand() % 8);
            building.size.z = 2.0f + (float)(rand() % 3);
            
            // Random building colors
            Color colors[] = {
                LIGHTGRAY, MAROON, RED, ORANGE, YELLOW, 
                GREEN, LIME, DARKGREEN, SKYBLUE, BLUE,
                DARKBLUE, PURPLE, VIOLET, DARKPURPLE, BEIGE, BROWN
            };
            building.color = colors[rand() % (sizeof(colors) / sizeof(Color))];
            
            world->buildings[world->buildingCount] = building;
            world->buildingCount++;

        }
    }
}

void DrawWorld(World *world) {
    // Draw ground plane
    DrawPlane((Vector3){0, 0, 0}, (Vector2){100, 100}, DARKGRAY);
    
    // Draw buildings
    for (int i = 0; i < world->buildingCount; i++) {
        Building *building = &world->buildings[i];
        
        Vector3 buildingPos = {
            building->position.x,
            building->position.y + building->size.y / 2.0f,
            building->position.z
        };
        
        DrawCube(buildingPos, building->size.x, building->size.y, building->size.z, building->color);
        DrawCubeWires(buildingPos, building->size.x, building->size.y, building->size.z, BLACK);
    }
    
    // Draw some additional environmental elements
    DrawGrid(20, 2.0f);
}

void UnloadWorld(World *world) {
    // Nothing to unload for now since we're using primitive shapes
}
