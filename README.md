# neon
Learning C and graphics programming with Raylib

## Development
Setup Raylib: https://github.com/raysan5/raylib/wiki/Working-on-GNU-Linux


## Notes
Use 16 bit RGB? https://rgbcolorpicker.com/565
Is transparency needed?



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

### Dithering
