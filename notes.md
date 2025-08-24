# Development notes
Previous attempts
- Raylib. But it does not support true Ray tracing
- Go-gl with compute shaders for Ray tracing. Too complicated for the absolute minimum with little documentation.

This latest iteration uses Vulkan and modern C++.
So this has lots of examples, docs and tooling because this is the industry standard way to build an engine/renderer.
And it's possible to use RT cores! (RT cores will be required to run the game)



## Organization





## Startup notes
Settings to have in the TUI:
- Refresh Rate
- Resolution
- Borderless Windowed or Fullscreen
- Monitor selection
- V-Sync toggle

Get this info from SDL:
Monitor list, resolutions (if possible with borderless), refresh rates



## Sound
Have a base layer of blue noise?
How can I add layers on top to make the city feel alive?
Spatial audio would be so cool




## Dithering
Dithering was created as a way to provide the perception of more color with reduced color spaces.
I am using dithering (or dither-like algorithms) for style and possibly for performance gains.

The main problem with "dot effects" is that they look bad if they change when the camera moves. How can I keep them mostly the same on objects like the Obra Dinn?

Resources:
Different ditherpunk effects: https://surma.dev/things/ditherpunk/
Obra Dinn tech explanation: https://forums.tigsource.com/index.php?topic=40832.800

Instead of dithering, what if I checkerboard ray trace? (I would need to do something for motion)

## Colors
Sadly, it looks like GPUs are optimized for 32 bit color spaces like R8G8B8A8. So RGB565(16 bit) would only slow things down.
This is worth testing again now that compute shaders are being used. raytracer.comp is currently using RGBA32F. (32 bits per channel!)

R8G8B8 +
1 bit for transparency (opaque/transparent flag)
3 bits for reflectivity (8 levels from non-reflective to mirror-like)
4 bits for roughness (16 levels of surface roughness)

BUT RGB565 might be worth it if there is more data to store


## Not cyberpunk. Neonpunk?
Cyberpunk is neat but I'm drawn to neon asthenic of it more than anything else.
Instead of chrome, violence and sexualization, I want the feeling of optomism and human connection.

Like Them Midnight's music. Their Wave song is a good example of this.