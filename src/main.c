#include "raylib.h"
#include "raymath.h"
#include <math.h>
#include <stdlib.h>

// Simple sphere structure
typedef struct {
    Vector3 center;
    float radius;
    Color color;
} Sphere;

// Simple ray structure
typedef struct {
    Vector3 origin;
    Vector3 direction;
} Ray3D;

// Scene data
Sphere spheres[3];
Vector3 lightPos = {2.0f, 4.0f, 2.0f};

// Ray-sphere intersection
float RaySphereIntersect(Ray3D ray, Sphere sphere) {
    Vector3 oc = Vector3Subtract(ray.origin, sphere.center);
    float a = Vector3DotProduct(ray.direction, ray.direction);
    float b = 2.0f * Vector3DotProduct(oc, ray.direction);
    float c = Vector3DotProduct(oc, oc) - sphere.radius * sphere.radius;
    float discriminant = b * b - 4 * a * c;
    
    if (discriminant < 0) return -1.0f;
    
    float t1 = (-b - sqrtf(discriminant)) / (2.0f * a);
    float t2 = (-b + sqrtf(discriminant)) / (2.0f * a);
    
    return (t1 > 0.001f) ? t1 : ((t2 > 0.001f) ? t2 : -1.0f);
}

// Simple lighting calculation
Color CalculateLighting(Vector3 hitPoint, Vector3 normal, Color surfaceColor) {
    Vector3 lightDir = Vector3Normalize(Vector3Subtract(lightPos, hitPoint));
    float lightIntensity = fmaxf(0.0f, Vector3DotProduct(normal, lightDir));

    // Ambient + diffuse
    float ambient = 0.2f;
    float diffuse = lightIntensity * 0.8f;
    float totalLight = ambient + diffuse;

    return (Color){
        (unsigned char)(surfaceColor.r * totalLight),
        (unsigned char)(surfaceColor.g * totalLight),
        (unsigned char)(surfaceColor.b * totalLight),
        255
    };
}

// Cast a ray and return color
Color TraceRay(Ray3D ray) {
    float closestT = 1000.0f;
    int hitSphere = -1;
    
    // Check intersection with all spheres
    for (int i = 0; i < 3; i++) {
        float t = RaySphereIntersect(ray, spheres[i]);
        if (t > 0 && t < closestT) {
            closestT = t;
            hitSphere = i;
        }
    }
    
    if (hitSphere == -1) {
        // Sky gradient
        float t = 0.5f * (ray.direction.y + 1.0f);
        return (Color){
            (unsigned char)(255 * (1.0f - t) + 135 * t),
            (unsigned char)(255 * (1.0f - t) + 206 * t),
            (unsigned char)(255 * (1.0f - t) + 235 * t),
            255
        };
    }
    
    // Calculate hit point and normal
    Vector3 hitPoint = Vector3Add(ray.origin, Vector3Scale(ray.direction, closestT));
    Vector3 normal = Vector3Normalize(Vector3Subtract(hitPoint, spheres[hitSphere].center));
    
    return CalculateLighting(hitPoint, normal, spheres[hitSphere].color);
}

// Convert screen coordinates to ray
Ray3D ScreenToRay(int x, int y, int screenWidth, int screenHeight, float time) {
    // Camera setup
    Vector3 cameraPos = {
        3.0f * cosf(time * 0.5f),
        2.0f,
        3.0f * sinf(time * 0.5f)
    };
    Vector3 cameraTarget = {0.0f, 0.0f, 0.0f};
    Vector3 cameraUp = {0.0f, 1.0f, 0.0f};

    // Convert to normalized device coordinates
    float u = (float)x / screenWidth;
    float v = (float)y / screenHeight;

    // Convert to world coordinates
    float aspect = (float)screenWidth / screenHeight;
    float fov = 45.0f * PI / 180.0f;

    Vector3 w = Vector3Normalize(Vector3Subtract(cameraPos, cameraTarget));
    Vector3 u_axis = Vector3Normalize(Vector3CrossProduct(cameraUp, w));
    Vector3 v_axis = Vector3CrossProduct(w, u_axis);

    Vector3 horizontal = Vector3Scale(u_axis, 2.0f * aspect * tanf(fov / 2.0f));
    Vector3 vertical = Vector3Scale(v_axis, 2.0f * tanf(fov / 2.0f));
    Vector3 lower_left = Vector3Subtract(Vector3Subtract(Vector3Subtract(cameraPos, Vector3Scale(horizontal, 0.5f)), Vector3Scale(vertical, 0.5f)), w);

    Vector3 rayDir = Vector3Normalize(Vector3Subtract(Vector3Add(Vector3Add(lower_left, Vector3Scale(horizontal, u)), Vector3Scale(vertical, v)), cameraPos));

    return (Ray3D){cameraPos, rayDir};
}

// Alternate rays cast per frame. This results in 50% of the work
bool ShouldCastRay(int x, int y, int frame) {
    return ((x + y + frame) % 2) == 0;
}

//------------------------------------------------------------------------------------
// Program main entry point
//------------------------------------------------------------------------------------
int main(void)
{
    // Initialization
    //--------------------------------------------------------------------------------------
    const int screenWidth = 800;
    const int screenHeight = 450;

    InitWindow(screenWidth, screenHeight, "Dithered Raytracer");

    SetTargetFPS(60);
    
    // Initialize scene
    spheres[0] = (Sphere){{0.0f, 0.0f, 0.0f}, 1.0f, PURPLE};
    spheres[1] = (Sphere){{-2.0f, 0.0f, -1.0f}, 0.5f, GREEN};
    spheres[2] = (Sphere){{2.0f, 0.0f, -1.0f}, 0.5f, BLUE};
    
    // Create render texture for our raytraced image
    RenderTexture2D target = LoadRenderTexture(screenWidth, screenHeight);
    
    // Frame buffers for temporal dithering
    Color* currentFrame = (Color*)malloc(screenWidth * screenHeight * sizeof(Color));
    Color* previousFrame = (Color*)malloc(screenWidth * screenHeight * sizeof(Color));
    
    // Initialize with black
    for (int i = 0; i < screenWidth * screenHeight; i++) {
        currentFrame[i] = BLACK;
        previousFrame[i] = BLACK;
    }
    
    int frameCount = 0;
    //--------------------------------------------------------------------------------------

    // Main game loop
    while (!WindowShouldClose())
    {
        // Update
        //----------------------------------------------------------------------------------
        frameCount++;
        float time = GetTime();
        
        // Animate light position
        lightPos.x = 2.0f * cosf(time);
        lightPos.z = 2.0f * sinf(time);
        //----------------------------------------------------------------------------------

        // Raytracing
        //----------------------------------------------------------------------------------
        BeginTextureMode(target);
        
        // Clear with black
        ClearBackground(BLACK);
        
        // Raytrace with dithering
        for (int y = 0; y < screenHeight; y++) {
            for (int x = 0; x < screenWidth; x++) {
                Color pixelColor;
                
                if (ShouldCastRay(x, y, frameCount)) {
                    // Cast ray for this pixel
                    Ray3D ray = ScreenToRay(x, y, screenWidth, screenHeight, time);
                    pixelColor = TraceRay(ray);
                    currentFrame[y * screenWidth + x] = pixelColor;
                } else {
                    // Use previous frame data (temporal dithering)
                    pixelColor = previousFrame[y * screenWidth + x];
                    currentFrame[y * screenWidth + x] = pixelColor;
                }
                
                // Draw the pixel
                DrawPixel(x, y, pixelColor);
            }
        }
        
        EndTextureMode();
        
        // Swap frame buffers
        Color* temp = previousFrame;
        previousFrame = currentFrame;
        currentFrame = temp;
        //----------------------------------------------------------------------------------

        // Draw
        //----------------------------------------------------------------------------------
        BeginDrawing();

            ClearBackground(RAYWHITE);
            
            // Draw our raytraced image
            DrawTexture(target.texture, 0, 0, WHITE);
            
            // Draw info
            DrawText(TextFormat("FPS: %d", GetFPS()), 10, 10, 20, YELLOW);
            DrawText("50% rays per frame (dithered)", 10, 35, 20, YELLOW);

        EndDrawing();
        //----------------------------------------------------------------------------------
    }

    // De-Initialization
    //--------------------------------------------------------------------------------------
    UnloadRenderTexture(target);
    free(currentFrame);
    free(previousFrame);
    CloseWindow();
    //--------------------------------------------------------------------------------------

    return 0;
}
