#version 330
// This did not work how I thought it would. Saving for possible future use


// Input
in vec2 fragTexCoord;
in vec4 fragColor;

// Uniforms
uniform sampler2D texture0;
uniform vec2 resolution;
uniform float time;
uniform float ditherStrength;

// Output
out vec4 finalColor;

// 8x8 Bayer matrix for ordered dithering
const mat4 bayerMatrix4x4 = mat4(
     0.0/16.0,  8.0/16.0,  2.0/16.0, 10.0/16.0,
    12.0/16.0,  4.0/16.0, 14.0/16.0,  6.0/16.0,
     3.0/16.0, 11.0/16.0,  1.0/16.0,  9.0/16.0,
    15.0/16.0,  7.0/16.0, 13.0/16.0,  5.0/16.0
);

// 8x8 Bayer matrix (extended)
float getBayerValue(ivec2 coord) {
    // Use modulo to tile the 4x4 matrix into 8x8 pattern
    int x = coord.x % 8;
    int y = coord.y % 8;
    
    // Convert 8x8 to 4x4 pattern with offsets
    int bx = x % 4;
    int by = y % 4;
    float base = bayerMatrix4x4[bx][by];
    
    // Add offset for extended pattern
    if (x >= 4) base += 32.0/64.0;
    if (y >= 4) base += 16.0/64.0;
    
    return base;
}

// Alternative: Simple 4x4 Bayer matrix
float getBayerValue4x4(ivec2 coord) {
    return bayerMatrix4x4[coord.x % 4][coord.y % 4];
}

// Blue noise dithering (approximation)
float blueNoise(vec2 coord) {
    return fract(sin(dot(coord, vec2(12.9898, 78.233))) * 43758.5453);
}

// Quantize color to specific bit depth
vec3 quantize(vec3 color, float levels) {
    return floor(color * levels) / levels;
}

void main() {
    // Sample the low-resolution texture
    vec4 texColor = texture(texture0, fragTexCoord);
    
    // Get screen coordinates for dithering pattern
    ivec2 screenCoord = ivec2(gl_FragCoord.xy);
    
    // Choose dithering method
    float ditherValue;
    
    // Method 1: Ordered dithering (Bayer matrix)
    ditherValue = getBayerValue4x4(screenCoord);
    
    // Method 2: Blue noise (uncomment to use)
    // ditherValue = blueNoise(gl_FragCoord.xy);
    
    // Method 3: Temporal dithering (animated)
    // ditherValue = getBayerValue4x4(screenCoord + ivec2(time * 60.0));
    
    // Apply dithering
    vec3 ditheredColor = texColor.rgb;
    
    // Quantization + dithering
    float levels = 16.0; // Adjust for more/less color banding
    vec3 quantized = quantize(ditheredColor, levels);
    
    // Add dither threshold
    float threshold = (ditherValue - 0.5) * ditherStrength;
    ditheredColor += threshold;
    
    // Re-quantize after dithering
    ditheredColor = quantize(ditheredColor, levels);
    
    // Alternative: Simple additive dithering
    // ditheredColor = texColor.rgb + (ditherValue - 0.5) * ditherStrength;
    
    // Output final color
    finalColor = vec4(clamp(ditheredColor, 0.0, 1.0), texColor.a);
    
    // Debug: Show dither pattern (uncomment to visualize)
    // finalColor = vec4(vec3(ditherValue), 1.0);
}
