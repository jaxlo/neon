# neon
A raytraced game built with Go and OpenGL 4.3


## Development environment

To run the code:
``go install``
``go run main.go``




--- --- --- Design Notes --- --- ---

## OpenGL notes
OpenGL 4.3 is used because it is new enough to support compute shaders and old enough to support hardware going back to around 2012.
It's also well supported by https://github.com/go-gl and has lots of documentation/examples online.

File extensions:
.glsl - Generic "OpenGL Shading Language" file
.fs - Fragment Shader
.vs - Vertex Shader
.cs or .comp - Compute Shader


## Raytracing
Implemented via a compute shader.
Only have 3 bounces and upscale from quarter resolution?


## Sound
Have a base layer of blue noise?
How can I add layers on top to make the city feel alive?





## Dithering
Dithering was created as a way to provide the perception of more color with reduced color spaces.
I am using dithering (or dither-like algorithms) for style and possibly for performance gains.

Should I upscale with predictable random noise? This could be better

Resources:
Different ditherpunk effects: https://surma.dev/things/ditherpunk/
Obra Dinn tech explanation: https://forums.tigsource.com/index.php?topic=40832.800


## Colors
Sadly, it looks like GPUs are optimized for 32 bit color spaces like R8G8B8A8. So RGB565(16 bit) would only slow things down.


## Not cyberpunk. Neonpunk?
Cyberpunk is neat but I'm drawn to neon asthenic of it more than anything else.
Instead of chrome, violence and sexualization, I want the feeling of digital human connection.

Like Them Midnight's music. Their Wave song is a good example of this.
