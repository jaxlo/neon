#ifndef WORLD_H
#define WORLD_H

#include <raylib.h>

#define MAX_BUILDINGS 20

typedef struct {
    Vector3 position;
    Vector3 size;
    Color color;
} Building;

typedef struct {
    Building buildings[MAX_BUILDINGS];
    int buildingCount;
    Model groundModel;
    Texture2D groundTexture;
} World;

// Function declarations
World InitWorld(void);
void DrawWorld(World *world);
void UnloadWorld(World *world);
void GenerateBuildings(World *world);

#endif
