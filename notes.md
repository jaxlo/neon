# Development notes


## ToDO Next
Make the object system extensible 


## Organization

Put the engine in a go module? With it's own go modules? lol
Or only make one layer of modules for modularity

scene.go - for scene management and object definitions
objects.go - for renderable object definitions
camera.go - for camera-related functionality
materials.go - for material definitions
lights.go - for lighting setup



## OpenGL notes
OpenGL 4.3 is used because it is new enough to support compute shaders and old enough to support hardware going back to around 2012.
It's also well supported by https://github.com/go-gl and has lots of documentation/examples online.

File extensions:
.glsl - Generic "OpenGL Shading Language" file
.fs - Fragment Shader
.vs - Vertex Shader
.comp - Compute Shader

V-sync - Have this as a setting somewhere?
glfw.SwapInterval(1) - V-Sync enabled, capped at monitor refresh rate
glfw.SwapInterval(0) - V-Sync disabled, uncapped framerate


## Raytracing
Implemented via a compute shader.
Only have 3 bounces and upscale from quarter resolution?


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
Instead of chrome, violence and sexualization, I want the feeling of digital human connection.

Like Them Midnight's music. Their Wave song is a good example of this.

Synthwave is too 80s lol
