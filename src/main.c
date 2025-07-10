#include "raylib.h"
#include "raymath.h"
#include <math.h>
#include <stdlib.h>

// RGB565 color structure for reduced color depth
typedef struct {
    unsigned short value; // 5 bits red, 6 bits green, 5 bits blue
} RGB565;

// Neon light structure for emissive elements
typedef struct {
    Vector3 position;
    float radius;
    Color color;
    float intensity;
    bool flicker; // For animated neon signs
} NeonLight;

// Simple sphere structure for geometry
typedef struct {
    Vector3 center;
    float radius;
    Color color;
    float reflectivity; // 0.0 = matte, 1.0 = mirror
    bool isEmissive;
} Sphere;

// Ray structure
typedef struct {
    Vector3 origin;
    Vector3 direction;
} Ray3D;

// Scene data
Sphere spheres[8];
NeonLight neonLights[6];
int sphereCount = 8;
int neonLightCount = 6;

// Bayer dithering matrix (4x4)
static const int bayerMatrix[16] = {
    0, 8, 2, 10,
    12, 4, 14, 6,
    3, 11, 1, 9,
    15, 7, 13, 5
};

// Convert RGB888 to RGB565 with dithering
RGB565 ColorToRGB565(Color color, int x, int y) {
    // Get dither threshold from Bayer matrix
    int ditherThreshold = bayerMatrix[(y % 4) * 4 + (x % 4)];
    
    // Add dithering noise and clamp
    int r = color.r + ((ditherThreshold - 8) * 2);
    int g = color.g + ((ditherThreshold - 8) * 2);
    int b = color.b + ((ditherThreshold - 8) * 2);
    
    r = (r < 0) ? 0 : (r > 255) ? 255 : r;
    g = (g < 0) ? 0 : (g > 255) ? 255 : g;
    b = (b < 0) ? 0 : (b > 255) ? 255 : b;
    
    // Convert to RGB565
    unsigned short r565 = (r >> 3) & 0x1F;      // 5 bits
    unsigned short g565 = (g >> 2) & 0x3F;      // 6 bits
    unsigned short b565 = (b >> 3) & 0x1F;      // 5 bits
    
    return (RGB565){(r565 << 11) | (g565 << 5) | b565};
}

// Convert RGB565 back to RGB888 for display
Color RGB565ToColor(RGB565 rgb565) {
    unsigned short r = (rgb565.value >> 11) & 0x1F;
    unsigned short g = (rgb565.value >> 5) & 0x3F;
    unsigned short b = rgb565.value & 0x1F;
    
    // Scale back to 8-bit (with proper rounding)
    return (Color){
        (r * 255 + 15) / 31,
        (g * 255 + 31) / 63,
        (b * 255 + 15) / 31,
        255
    };
}

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

// Calculate neon lighting with atmospheric scattering
Color CalculateNeonLighting(Vector3 hitPoint, Vector3 normal, Color surfaceColor, float reflectivity, float time) {
    Vector3 totalLight = {0.0f, 0.0f, 0.0f};
    
    // Ambient city glow (purple-blue night sky)
    float ambient = 0.1f;
    totalLight.x += ambient * 0.2f; // Red
    totalLight.y += ambient * 0.1f; // Green  
    totalLight.z += ambient * 0.4f; // Blue
    
    // Add contribution from each neon light
    for (int i = 0; i < neonLightCount; i++) {
        Vector3 lightDir = Vector3Subtract(neonLights[i].position, hitPoint);
        float distance = Vector3Length(lightDir);
        lightDir = Vector3Normalize(lightDir);
        
        // Distance attenuation with atmospheric scattering
        float attenuation = 1.0f / (1.0f + 0.1f * distance + 0.01f * distance * distance);
        
        // Add atmospheric scattering effect (more visible in fog)
        float scattering = expf(-distance * 0.05f);
        attenuation *= scattering;
        
        // Neon flicker effect
        float flicker = 1.0f;
        if (neonLights[i].flicker) {
            flicker = 0.8f + 0.2f * sinf(time * 10.0f + i * 2.0f);
        }
        
        // Lambertian diffuse
        float diffuse = fmaxf(0.0f, Vector3DotProduct(normal, lightDir));
        
        // Simple reflection for shiny surfaces
        if (reflectivity > 0.0f) {
            Vector3 reflection = Vector3Subtract(Vector3Scale(normal, 2.0f * Vector3DotProduct(lightDir, normal)), lightDir);
            Vector3 viewDir = Vector3Normalize(Vector3Subtract(hitPoint, (Vector3){0, 0, 5})); // Assume camera at origin
            float spec = powf(fmaxf(0.0f, Vector3DotProduct(reflection, Vector3Negate(viewDir))), 32.0f);
            diffuse += spec * reflectivity;
        }
        
        float lightContrib = attenuation * flicker * diffuse * neonLights[i].intensity;
        
        totalLight.x += lightContrib * neonLights[i].color.r / 255.0f;
        totalLight.y += lightContrib * neonLights[i].color.g / 255.0f;
        totalLight.z += lightContrib * neonLights[i].color.b / 255.0f;
    }
    
    // Apply lighting to surface color
    return (Color){
        (unsigned char)fminf(255, surfaceColor.r * totalLight.x),
        (unsigned char)fminf(255, surfaceColor.g * totalLight.y),
        (unsigned char)fminf(255, surfaceColor.b * totalLight.z),
        255
    };
}

// Trace a ray through the scene
Color TraceRay(Ray3D ray, float time) {
    float closestT = 1000.0f;
    int hitSphere = -1;
    
    // Find closest intersection
    for (int i = 0; i < sphereCount; i++) {
        float t = RaySphereIntersect(ray, spheres[i]);
        if (t > 0 && t < closestT) {
            closestT = t;
            hitSphere = i;
        }
    }
    
    if (hitSphere == -1) {
        // Night sky with neon city glow
        float skyGlow = 0.5f + 0.3f * sinf(time * 0.1f);
        return (Color){
            (unsigned char)(20 * skyGlow),   // Dark red
            (unsigned char)(10 * skyGlow),   // Very dark green
            (unsigned char)(60 * skyGlow),   // Deep blue
            255
        };
    }
    
    // Calculate hit point and normal
    Vector3 hitPoint = Vector3Add(ray.origin, Vector3Scale(ray.direction, closestT));
    Vector3 normal = Vector3Normalize(Vector3Subtract(hitPoint, spheres[hitSphere].center));
    
    // Handle emissive objects (neon signs)
    if (spheres[hitSphere].isEmissive) {
        float glow = 0.8f + 0.2f * sinf(time * 5.0f);
        return (Color){
            (unsigned char)(spheres[hitSphere].color.r * glow),
            (unsigned char)(spheres[hitSphere].color.g * glow),
            (unsigned char)(spheres[hitSphere].color.b * glow),
            255
        };
    }
    
    return CalculateNeonLighting(hitPoint, normal, spheres[hitSphere].color, spheres[hitSphere].reflectivity, time);
}

// Convert screen coordinates to ray (for quarter resolution)
Ray3D ScreenToRay(int x, int y, int quarterWidth, int quarterHeight, float time) {
    // Animated camera position (moving through neon city)
    Vector3 cameraPos = {
        2.0f * sinf(time * 0.3f),
        1.5f + 0.5f * sinf(time * 0.2f),
        5.0f + 1.0f * cosf(time * 0.1f)
    };
    Vector3 cameraTarget = {0.0f, 0.0f, 0.0f};
    Vector3 cameraUp = {0.0f, 1.0f, 0.0f};
    
    // Convert to normalized device coordinates
    float u = (float)x / quarterWidth;
    float v = (float)y / quarterHeight;
    
    // Calculate ray direction
    float aspect = (float)quarterWidth / quarterHeight;
    float fov = 60.0f * PI / 180.0f;
    
    Vector3 w = Vector3Normalize(Vector3Subtract(cameraPos, cameraTarget));
    Vector3 u_axis = Vector3Normalize(Vector3CrossProduct(cameraUp, w));
    Vector3 v_axis = Vector3CrossProduct(w, u_axis);
    
    Vector3 horizontal = Vector3Scale(u_axis, 2.0f * aspect * tanf(fov / 2.0f));
    Vector3 vertical = Vector3Scale(v_axis, 2.0f * tanf(fov / 2.0f));
    Vector3 lower_left = Vector3Subtract(Vector3Subtract(Vector3Subtract(cameraPos, Vector3Scale(horizontal, 0.5f)), Vector3Scale(vertical, 0.5f)), w);
    
    Vector3 rayDir = Vector3Normalize(Vector3Subtract(Vector3Add(Vector3Add(lower_left, Vector3Scale(horizontal, u)), Vector3Scale(vertical, v)), cameraPos));
    
    return (Ray3D){cameraPos, rayDir};
}

// Apply dithered atmospheric fog
Color ApplyAtmosphericFog(Color originalColor, float depth, int x, int y, float time) {
    // Base fog parameters
    float fogDensity = 0.05f + 0.02f * sinf(time * 0.1f); // Animated fog density
    float fogDistance = expf(-depth * fogDensity);
    
    // Neon-tinted fog color
    Color fogColor = {
        40 + (int)(20 * sinf(time * 0.3f)),   // Purple-red fog
        20 + (int)(15 * sinf(time * 0.2f)),   // Slight green
        80 + (int)(30 * sinf(time * 0.4f)),   // Blue-purple fog
        255
    };
    
    // Apply dithering to fog blend
    int ditherValue = bayerMatrix[(y % 4) * 4 + (x % 4)];
    float ditherNoise = (ditherValue - 8) * 0.02f;
    fogDistance += ditherNoise;
    fogDistance = (fogDistance < 0.0f) ? 0.0f : (fogDistance > 1.0f) ? 1.0f : fogDistance;
    
    // Blend original color with fog
    return (Color){
        (unsigned char)(originalColor.r * fogDistance + fogColor.r * (1.0f - fogDistance)),
        (unsigned char)(originalColor.g * fogDistance + fogColor.g * (1.0f - fogDistance)),
        (unsigned char)(originalColor.b * fogDistance + fogColor.b * (1.0f - fogDistance)),
        255
    };
}

// Bilinear upscaling with dithering
Color BilinearUpscale(RGB565* quarterBuffer, int qWidth, int qHeight, int fullX, int fullY, int fullWidth, int fullHeight) {
    // Map full resolution coordinates to quarter resolution
    float qx = (float)fullX * qWidth / fullWidth;
    float qy = (float)fullY * qHeight / fullHeight;
    
    // Get integer coordinates and fractional parts
    int x0 = (int)qx;
    int y0 = (int)qy;
    int x1 = x0 + 1;
    int y1 = y0 + 1;
    
    // Clamp coordinates
    x0 = (x0 < 0) ? 0 : (x0 >= qWidth) ? qWidth - 1 : x0;
    y0 = (y0 < 0) ? 0 : (y0 >= qHeight) ? qHeight - 1 : y0;
    x1 = (x1 < 0) ? 0 : (x1 >= qWidth) ? qWidth - 1 : x1;
    y1 = (y1 < 0) ? 0 : (y1 >= qHeight) ? qHeight - 1 : y1;
    
    // Get fractional parts
    float fx = qx - x0;
    float fy = qy - y0;
    
    // Sample four neighboring pixels
    Color c00 = RGB565ToColor(quarterBuffer[y0 * qWidth + x0]);
    Color c10 = RGB565ToColor(quarterBuffer[y0 * qWidth + x1]);
    Color c01 = RGB565ToColor(quarterBuffer[y1 * qWidth + x0]);
    Color c11 = RGB565ToColor(quarterBuffer[y1 * qWidth + x1]);
    
    // Bilinear interpolation
    Color top = (Color){
        (unsigned char)(c00.r * (1.0f - fx) + c10.r * fx),
        (unsigned char)(c00.g * (1.0f - fx) + c10.g * fx),
        (unsigned char)(c00.b * (1.0f - fx) + c10.b * fx),
        255
    };
    
    Color bottom = (Color){
        (unsigned char)(c01.r * (1.0f - fx) + c11.r * fx),
        (unsigned char)(c01.g * (1.0f - fx) + c11.g * fx),
        (unsigned char)(c01.b * (1.0f - fx) + c11.b * fx),
        255
    };
    
    Color result = (Color){
        (unsigned char)(top.r * (1.0f - fy) + bottom.r * fy),
        (unsigned char)(top.g * (1.0f - fy) + bottom.g * fy),
        (unsigned char)(top.b * (1.0f - fy) + bottom.b * fy),
        255
    };
    
    return result;
}

// Apply posterization with dithering
Color PosterizeColor(Color input, int x, int y) {
    // Reduce color levels (8 levels per channel = 512 colors total)
    int levels = 8;
    float step = 255.0f / (levels - 1);
    
    // Add dithering noise
    int ditherValue = bayerMatrix[(y % 4) * 4 + (x % 4)];
    float ditherNoise = (ditherValue - 8) * step * 0.25f;
    
    // Quantize with dithering
    int r = (int)((input.r + ditherNoise) / step + 0.5f) * step;
    int g = (int)((input.g + ditherNoise) / step + 0.5f) * step;
    int b = (int)((input.b + ditherNoise) / step + 0.5f) * step;
    
    // Clamp values
    r = (r < 0) ? 0 : (r > 255) ? 255 : r;
    g = (g < 0) ? 0 : (g > 255) ? 255 : g;
    b = (b < 0) ? 0 : (b > 255) ? 255 : b;
    
    return (Color){r, g, b, 255};
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
    
    // Quarter resolution for raytracing
    const int quarterWidth = screenWidth / 4;
    const int quarterHeight = screenHeight / 4;

    InitWindow(screenWidth, screenHeight, "Neon City Quarter-Res Raytracer");
    SetTargetFPS(60);
    
    // Initialize neon city scene
    // Ground plane (wet asphalt)
    spheres[0] = (Sphere){{0.0f, -1000.0f, 0.0f}, 1000.0f, (Color){30, 30, 35, 255}, 0.3f, false};
    
    // Building facades
    spheres[1] = (Sphere){{-3.0f, 0.0f, -2.0f}, 1.0f, (Color){60, 60, 70, 255}, 0.1f, false};
    spheres[2] = (Sphere){{3.0f, 0.0f, -1.0f}, 0.8f, (Color){50, 50, 60, 255}, 0.1f, false};
    spheres[3] = (Sphere){{0.0f, 0.0f, -3.0f}, 1.2f, (Color){40, 40, 50, 255}, 0.2f, false};
    
    // Neon signs (emissive spheres)
    spheres[4] = (Sphere){{-2.0f, 2.0f, 0.0f}, 0.3f, (Color){255, 0, 100, 255}, 0.0f, true};  // Pink neon
    spheres[5] = (Sphere){{2.0f, 1.5f, 0.0f}, 0.25f, (Color){0, 255, 150, 255}, 0.0f, true};  // Cyan neon
    spheres[6] = (Sphere){{0.0f, 2.5f, -1.0f}, 0.2f, (Color){255, 255, 0, 255}, 0.0f, true};  // Yellow neon
    
    // Reflective car/puddle
    spheres[7] = (Sphere){{1.0f, -0.5f, 1.0f}, 0.5f, (Color){20, 20, 30, 255}, 0.8f, false};
    
    // Initialize neon lights (these cast light on the scene)
    neonLights[0] = (NeonLight){{-2.0f, 2.0f, 0.0f}, 2.0f, (Color){255, 0, 100, 255}, 1.5f, true};
    neonLights[1] = (NeonLight){{2.0f, 1.5f, 0.0f}, 1.8f, (Color){0, 255, 150, 255}, 1.2f, true};
    neonLights[2] = (NeonLight){{0.0f, 2.5f, -1.0f}, 1.5f, (Color){255, 255, 0, 255}, 1.0f, false};
    neonLights[3] = (NeonLight){{-4.0f, 3.0f, -2.0f}, 3.0f, (Color){255, 50, 0, 255}, 0.8f, true};
    neonLights[4] = (NeonLight){{4.0f, 2.0f, -3.0f}, 2.5f, (Color){0, 100, 255, 255}, 1.1f, false};
    neonLights[5] = (NeonLight){{0.0f, 4.0f, 2.0f}, 4.0f, (Color){200, 0, 255, 255}, 0.9f, true};
    
    // Create render targets
    RenderTexture2D quarterTarget = LoadRenderTexture(quarterWidth, quarterHeight);
    RenderTexture2D fullTarget = LoadRenderTexture(screenWidth, screenHeight);
    
    // Allocate RGB565 buffer for quarter resolution
    RGB565* quarterBuffer = (RGB565*)malloc(quarterWidth * quarterHeight * sizeof(RGB565));
    
    int frameCount = 0;
    //--------------------------------------------------------------------------------------

    // Main game loop
    while (!WindowShouldClose())
    {
        // Update
        //----------------------------------------------------------------------------------
        frameCount++;
        float time = GetTime();
        //----------------------------------------------------------------------------------

        // STEP 1: Render quarter resolution raytracing to RGB565 buffer
        //----------------------------------------------------------------------------------
        for (int y = 0; y < quarterHeight; y++) {
            for (int x = 0; x < quarterWidth; x++) {
                // Cast ray for this quarter-res pixel
                Ray3D ray = ScreenToRay(x, y, quarterWidth, quarterHeight, time);
                Color rayColor = TraceRay(ray, time);
                
                // Apply atmospheric fog
                float depth = 5.0f; // Approximate depth for fog calculation
                Color foggedColor = ApplyAtmosphericFog(rayColor, depth, x, y, time);
                
                // Convert to RGB565 with dithering
                quarterBuffer[y * quarterWidth + x] = ColorToRGB565(foggedColor, x, y);
            }
        }
        //----------------------------------------------------------------------------------

        // STEP 2: Upscale to full resolution with bilinear filtering
        //----------------------------------------------------------------------------------
        BeginTextureMode(fullTarget);
        ClearBackground(BLACK);
        
        for (int y = 0; y < screenHeight; y++) {
            for (int x = 0; x < screenWidth; x++) {
                // Bilinear upscale from quarter resolution
                Color upscaledColor = BilinearUpscale(quarterBuffer, quarterWidth, quarterHeight, x, y, screenWidth, screenHeight);
                
                // Apply posterization with dithering
                Color finalColor = PosterizeColor(upscaledColor, x, y);
                
                DrawPixel(x, y, finalColor);
            }
        }
        EndTextureMode();
        //----------------------------------------------------------------------------------

        // STEP 3: Apply bloom effect (simple additive bloom)
        //----------------------------------------------------------------------------------
        BeginTextureMode(quarterTarget);
        ClearBackground(BLACK);
        
        // Render bright areas for bloom
        for (int y = 0; y < quarterHeight; y++) {
            for (int x = 0; x < quarterWidth; x++) {
                Color originalColor = RGB565ToColor(quarterBuffer[y * quarterWidth + x]);
                
                // Extract bright areas (simple threshold)
                int brightness = (originalColor.r + originalColor.g + originalColor.b) / 3;
                if (brightness > 120) {
                    Color bloomColor = (Color){
                        (unsigned char)(originalColor.r * 0.3f),
                        (unsigned char)(originalColor.g * 0.3f),
                        (unsigned char)(originalColor.b * 0.3f),
                        255
                    };
                    DrawPixel(x, y, bloomColor);
                }
            }
        }
        EndTextureMode();
        //----------------------------------------------------------------------------------

        // STEP 4: Final composition and display
        //----------------------------------------------------------------------------------
        BeginDrawing();
        ClearBackground(BLACK);
        
        // Draw main image
        DrawTextureRec(fullTarget.texture, (Rectangle){0, 0, screenWidth, -screenHeight}, (Vector2){0, 0}, WHITE);
        
        // Add bloom on top (upscaled and blended)
        DrawTextureRec(quarterTarget.texture, (Rectangle){0, 0, quarterWidth, -quarterHeight}, (Vector2){0, 0}, (Color){255, 255, 255, 60});
        
        // Draw debug info
        DrawText(TextFormat("FPS: %d", GetFPS()), 10, 10, 20, YELLOW);
        DrawText(TextFormat("Quarter-Res: %dx%d", quarterWidth, quarterHeight), 10, 35, 20, YELLOW);
        DrawText("RGB565 + Dithering + Bloom", 10, 60, 20, YELLOW);
        DrawText("Neon City Pipeline", 10, 85, 20, MAGENTA);
        
        EndDrawing();
        //----------------------------------------------------------------------------------
    }

    // De-Initialization
    //--------------------------------------------------------------------------------------
    UnloadRenderTexture(quarterTarget);
    UnloadRenderTexture(fullTarget);
    free(quarterBuffer);
    CloseWindow();
    //--------------------------------------------------------------------------------------

    return 0;
}
