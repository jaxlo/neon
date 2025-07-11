# neon
Learning C and graphics programming with Raylib

## Development
Setup Raylib: https://github.com/raysan5/raylib/wiki/Working-on-GNU-Linux


## Notes
Use 16 bit RGB? https://rgbcolorpicker.com/565
OR use 8 bit colors?



## Sound
Have a base layer of blue noise?
How can I add layers on top to make the city feel alive?

## Graphics pipeline

1. Quarter-Resolution Raytracing 
Ray tracing via compute shaders 
Each ray traces one "super-pixel" representing a 2x2 block
For 1920x1080, there would be 480x270 rays
The result is then upscaled with bilinear filtering 

2. RGB565 Conversion - Reduces color depth with Bayer dithering
Atmospheric Fog - Adds neon-tinted fog with animated density
Bilinear Upscaling - Smoothly scales quarter-res to full resolution
Posterization - Quantizes colors to 8 levels per channel with dithering
Bloom Effect - Extracts bright areas and additively blends them

Single Compute Shader:
- Raytrace at quarter-res
- Apply fog, posterization, RGB565 conversion inline
- Write to quarter-res buffer
- Use texture filtering hardware for upscaling during final draw

### Non ray traced pipeline
1. Start with half resolution render
2. Lighting/shadow shaders (There will be more info when I implement this part)
3. Dither and upscale last

### Pipeline notes
- If there are too many artifacts, consider upscaling, then dithering (Small cost hit)
- Sadly, it looks like GPUs are optimized for 32 bit color spaces like R8G8B8A8. So RGB565(16 bit) would only slow things down.
