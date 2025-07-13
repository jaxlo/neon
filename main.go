package main

import (
	"fmt"
	"log"
	"math"
	"runtime"
	"strings"
	"time"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	width  = 1920
	height = 1080
)

var (
	// Vertex shader for fullscreen quad
	vertexShader = `#version 430 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec2 aTexCoord;

out vec2 TexCoord;

void main() {
    gl_Position = vec4(aPos, 1.0);
    TexCoord = aTexCoord;
}` + "\x00"

	// Fragment shader for text rendering
	textFragmentShader = `#version 430 core
in vec2 TexCoord;
out vec4 FragColor;

uniform sampler2D textTexture;
uniform vec3 textColor;

void main() {
    float alpha = texture(textTexture, TexCoord).r;
    FragColor = vec4(textColor, alpha);
}` + "\x00"

	// Vertex shader for text rendering
	textVertexShader = `#version 430 core
layout (location = 0) in vec4 vertex; // <vec2 pos, vec2 tex>
out vec2 TexCoord;

uniform mat4 projection;

void main() {
    gl_Position = projection * vec4(vertex.xy, 0.0, 1.0);
    TexCoord = vertex.zw;
}` + "\x00"
	fragmentShader = `#version 430 core
in vec2 TexCoord;
out vec4 FragColor;

uniform sampler2D screenTexture;

void main() {
    FragColor = texture(screenTexture, TexCoord);
}` + "\x00"

	// Compute shader for raytracing
	computeShader = `#version 430 core
layout(local_size_x = 16, local_size_y = 16) in;
layout(rgba32f, binding = 0) uniform image2D imgOutput;

uniform float uTime;
uniform vec3 uCameraPos;
uniform vec3 uCameraDir;
uniform vec3 uCameraUp;
uniform vec3 uCameraRight;

struct Sphere {
    vec3 center;
    float radius;
    vec3 color;
    float reflectivity;
};

struct Ray {
    vec3 origin;
    vec3 direction;
};

struct HitInfo {
    bool hit;
    float distance;
    vec3 point;
    vec3 normal;
    vec3 color;
    float reflectivity;
};

// Scene spheres
Sphere spheres[5] = Sphere[](
    Sphere(vec3(0, 0, -5), 1.0, vec3(1.0, 0.2, 0.2), 0.3),
    Sphere(vec3(-3, 0, -7), 1.5, vec3(0.2, 1.0, 0.2), 0.8),
    Sphere(vec3(3, 0, -6), 1.2, vec3(0.2, 0.2, 1.0), 0.6),
    Sphere(vec3(0, -101, -5), 100.0, vec3(0.8, 0.8, 0.8), 0.1), // Ground
    Sphere(vec3(-1, 2, -4), 0.8, vec3(1.0, 1.0, 0.2), 0.9)
);

HitInfo intersectSphere(Ray ray, Sphere sphere) {
    HitInfo hit;
    hit.hit = false;
    
    vec3 oc = ray.origin - sphere.center;
    float a = dot(ray.direction, ray.direction);
    float b = 2.0 * dot(oc, ray.direction);
    float c = dot(oc, oc) - sphere.radius * sphere.radius;
    
    float discriminant = b * b - 4.0 * a * c;
    
    if (discriminant < 0) return hit;
    
    float t = (-b - sqrt(discriminant)) / (2.0 * a);
    if (t < 0.001) {
        t = (-b + sqrt(discriminant)) / (2.0 * a);
        if (t < 0.001) return hit;
    }
    
    hit.hit = true;
    hit.distance = t;
    hit.point = ray.origin + t * ray.direction;
    hit.normal = normalize(hit.point - sphere.center);
    hit.color = sphere.color;
    hit.reflectivity = sphere.reflectivity;
    
    return hit;
}

HitInfo intersectScene(Ray ray) {
    HitInfo closestHit;
    closestHit.hit = false;
    closestHit.distance = 999999.0;
    
    for (int i = 0; i < 5; i++) {
        HitInfo hit = intersectSphere(ray, spheres[i]);
        if (hit.hit && hit.distance < closestHit.distance) {
            closestHit = hit;
        }
    }
    
    return closestHit;
}

vec3 reflect(vec3 direction, vec3 normal) {
    return direction - 2.0 * dot(direction, normal) * normal;
}

vec3 traceRay(Ray ray) {
    vec3 color = vec3(0.0);
    vec3 rayColor = vec3(1.0);
    
    for (int bounce = 0; bounce < 3; bounce++) {
        HitInfo hit = intersectScene(ray);
        
        if (!hit.hit) {
            // Sky color
            vec3 skyColor = mix(vec3(0.5, 0.7, 1.0), vec3(1.0, 1.0, 1.0), 
                               smoothstep(-0.5, 0.5, ray.direction.y));
            color += rayColor * skyColor;
            break;
        }
        
        // Simple lighting
        vec3 lightDir = normalize(vec3(0.5, 1.0, 0.3));
        float diffuse = max(dot(hit.normal, lightDir), 0.0);
        vec3 ambient = vec3(0.1);
        
        vec3 surfaceColor = hit.color * (ambient + diffuse * 0.7);
        
        // Add this bounce's contribution
        color += rayColor * surfaceColor * (1.0 - hit.reflectivity);
        
        // Prepare for next bounce
        rayColor *= hit.reflectivity;
        
        if (length(rayColor) < 0.01) break; // Early termination
        
        ray.origin = hit.point + hit.normal * 0.001;
        ray.direction = reflect(ray.direction, hit.normal);
    }
    
    return color;
}

void main() {
    ivec2 pixelCoords = ivec2(gl_GlobalInvocationID.xy);
    ivec2 dims = imageSize(imgOutput);
    
    if (pixelCoords.x >= dims.x || pixelCoords.y >= dims.y) return;
    
    // Convert pixel coordinates to normalized device coordinates
    vec2 uv = (vec2(pixelCoords) + 0.5) / vec2(dims);
    uv = uv * 2.0 - 1.0;
    uv.x *= float(dims.x) / float(dims.y);
    
    // Create ray
    Ray ray;
    ray.origin = uCameraPos;
    ray.direction = normalize(uCameraDir + uv.x * uCameraRight + uv.y * uCameraUp);
    
    vec3 color = traceRay(ray);
    
    // Tone mapping and gamma correction
    color = color / (color + vec3(1.0));
    color = pow(color, vec3(1.0/2.2));
    
    imageStore(imgOutput, pixelCoords, vec4(color, 1.0));
}` + "\x00"
)

type Character struct {
	textureID uint32
	size      mgl32.Vec2
	bearing   mgl32.Vec2
	advance   int32
}

type TextRenderer struct {
	program    uint32
	VAO        uint32
	VBO        uint32
	characters map[rune]Character
	projection mgl32.Mat4
}

type Renderer struct {
	window         *glfw.Window
	shaderProgram  uint32
	computeProgram uint32
	texture        uint32
	VAO            uint32
	VBO            uint32
	EBO            uint32
	textRenderer   *TextRenderer

	// FPS tracking
	frameCount int
	lastTime   time.Time
	fps        float64
}

func main() {
	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	// Get primary monitor for fullscreen
	monitor := glfw.GetPrimaryMonitor()
	mode := monitor.GetVideoMode()

	window, err := glfw.CreateWindow(mode.Width, mode.Height, "GPU Raytracer", monitor, nil)
	if err != nil {
		log.Fatalln("failed to create window:", err)
	}

	window.MakeContextCurrent()
	glfw.SwapInterval(0) // Disable V-Sync to see actual performance

	if err := gl.Init(); err != nil {
		log.Fatalln("failed to initialize OpenGL:", err)
	}

	renderer := &Renderer{
		window:     window,
		lastTime:   time.Now(),
		frameCount: 0,
		fps:        0.0,
	}
	renderer.init()

	// Main loop
	for !window.ShouldClose() {
		renderer.updateFPS()
		renderer.render()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func (r *Renderer) init() {
	// Create shader programs
	r.shaderProgram = createShaderProgram(vertexShader, fragmentShader)
	r.computeProgram = createComputeShader(computeShader)

	// Initialize text renderer
	r.textRenderer = &TextRenderer{}
	r.textRenderer.init()

	// Create texture for compute shader output
	gl.GenTextures(1, &r.texture)
	gl.BindTexture(gl.TEXTURE_2D, r.texture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA32F, width, height, 0, gl.RGBA, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	// Bind texture to image unit 0
	gl.BindImageTexture(0, r.texture, 0, false, 0, gl.WRITE_ONLY, gl.RGBA32F)

	// Create fullscreen quad
	vertices := []float32{
		// positions   // texture coords
		-1.0, -1.0, 0.0, 0.0, 0.0,
		1.0, -1.0, 0.0, 1.0, 0.0,
		1.0, 1.0, 0.0, 1.0, 1.0,
		-1.0, 1.0, 0.0, 0.0, 1.0,
	}

	indices := []uint32{
		0, 1, 2,
		2, 3, 0,
	}

	gl.GenVertexArrays(1, &r.VAO)
	gl.GenBuffers(1, &r.VBO)
	gl.GenBuffers(1, &r.EBO)

	gl.BindVertexArray(r.VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, r.VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, r.EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	// Position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// Texture coordinate attribute
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(0)
}

func (r *Renderer) updateFPS() {
	r.frameCount++
	currentTime := time.Now()
	elapsed := currentTime.Sub(r.lastTime)

	if elapsed >= time.Second {
		r.fps = float64(r.frameCount) / elapsed.Seconds()
		r.frameCount = 0
		r.lastTime = currentTime
	}
}

func (r *Renderer) render() {
	time := float32(glfw.GetTime())

	// Camera setup with simple rotation
	cameraPos := mgl32.Vec3{
		float32(5.0 * math.Cos(float64(time*0.2))),
		2.0,
		float32(5.0 * math.Sin(float64(time*0.2))),
	}
	cameraTarget := mgl32.Vec3{0, 0, -5}
	cameraUp := mgl32.Vec3{0, 1, 0}

	cameraDir := cameraTarget.Sub(cameraPos).Normalize()
	cameraRight := cameraDir.Cross(cameraUp).Normalize()
	cameraUp = cameraRight.Cross(cameraDir).Normalize()

	// Dispatch compute shader
	gl.UseProgram(r.computeProgram)
	gl.Uniform1f(gl.GetUniformLocation(r.computeProgram, gl.Str("uTime\x00")), time)
	gl.Uniform3f(gl.GetUniformLocation(r.computeProgram, gl.Str("uCameraPos\x00")),
		cameraPos.X(), cameraPos.Y(), cameraPos.Z())
	gl.Uniform3f(gl.GetUniformLocation(r.computeProgram, gl.Str("uCameraDir\x00")),
		cameraDir.X(), cameraDir.Y(), cameraDir.Z())
	gl.Uniform3f(gl.GetUniformLocation(r.computeProgram, gl.Str("uCameraUp\x00")),
		cameraUp.X(), cameraUp.Y(), cameraUp.Z())
	gl.Uniform3f(gl.GetUniformLocation(r.computeProgram, gl.Str("uCameraRight\x00")),
		cameraRight.X(), cameraRight.Y(), cameraRight.Z())

	gl.DispatchCompute(width/16, height/16, 1)
	gl.MemoryBarrier(gl.SHADER_IMAGE_ACCESS_BARRIER_BIT)

	// Render fullscreen quad
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.UseProgram(r.shaderProgram)
	gl.BindTexture(gl.TEXTURE_2D, r.texture)
	gl.BindVertexArray(r.VAO)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)

	// Render FPS counter
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	fpsText := fmt.Sprintf("FPS: %.1f", r.fps)
	r.textRenderer.RenderText(fpsText, 10, float32(height-30), 1.0, mgl32.Vec3{1, 1, 1})
	gl.Disable(gl.BLEND)
}

func createShaderProgram(vertexSource, fragmentSource string) uint32 {
	vertexShader := compileShader(vertexSource, gl.VERTEX_SHADER)
	fragmentShader := compileShader(fragmentSource, gl.FRAGMENT_SHADER)

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))
		panic(fmt.Errorf("failed to link program: %v", log))
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program
}

func createComputeShader(source string) uint32 {
	shader := compileShader(source, gl.COMPUTE_SHADER)

	program := gl.CreateProgram()
	gl.AttachShader(program, shader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))
		panic(fmt.Errorf("failed to link compute program: %v", log))
	}

	gl.DeleteShader(shader)
	return program
}

func (tr *TextRenderer) init() {
	// Create shader program for text rendering
	tr.program = createShaderProgram(textVertexShader, textFragmentShader)

	// Set up projection matrix for screen coordinates
	tr.projection = mgl32.Ortho(0, width, 0, height, -1, 1)

	// Create VAO and VBO for text rendering
	gl.GenVertexArrays(1, &tr.VAO)
	gl.GenBuffers(1, &tr.VBO)
	gl.BindVertexArray(tr.VAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, tr.VBO)
	gl.BufferData(gl.ARRAY_BUFFER, 6*4*4, nil, gl.DYNAMIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 4*4, nil)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	// Create simple bitmap font characters
	tr.characters = make(map[rune]Character)
	tr.createBitmapFont()
}

func (tr *TextRenderer) createBitmapFont() {
	// Create a simple 8x8 bitmap font for digits and letters
	fontData := map[rune][]uint8{
		'0': {0x3C, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66, 0x3C},
		'1': {0x18, 0x38, 0x18, 0x18, 0x18, 0x18, 0x18, 0x7E},
		'2': {0x3C, 0x66, 0x06, 0x0C, 0x18, 0x30, 0x60, 0x7E},
		'3': {0x3C, 0x66, 0x06, 0x1C, 0x06, 0x06, 0x66, 0x3C},
		'4': {0x0C, 0x1C, 0x2C, 0x4C, 0x7E, 0x0C, 0x0C, 0x0C},
		'5': {0x7E, 0x60, 0x60, 0x7C, 0x06, 0x06, 0x66, 0x3C},
		'6': {0x3C, 0x66, 0x60, 0x7C, 0x66, 0x66, 0x66, 0x3C},
		'7': {0x7E, 0x06, 0x0C, 0x18, 0x30, 0x30, 0x30, 0x30},
		'8': {0x3C, 0x66, 0x66, 0x3C, 0x66, 0x66, 0x66, 0x3C},
		'9': {0x3C, 0x66, 0x66, 0x66, 0x3E, 0x06, 0x66, 0x3C},
		'F': {0x7E, 0x60, 0x60, 0x7C, 0x60, 0x60, 0x60, 0x60},
		'P': {0x7C, 0x66, 0x66, 0x7C, 0x60, 0x60, 0x60, 0x60},
		'S': {0x3C, 0x66, 0x60, 0x3C, 0x06, 0x06, 0x66, 0x3C},
		':': {0x00, 0x18, 0x18, 0x00, 0x00, 0x18, 0x18, 0x00},
		' ': {0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		'.': {0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x18, 0x18},
	}

	for char, data := range fontData {
		// Convert 8x8 bitmap to 8x8 texture
		pixels := make([]uint8, 64)
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				bit := (data[y] >> (7 - x)) & 1
				pixels[y*8+x] = bit * 255
			}
		}

		var texture uint32
		gl.GenTextures(1, &texture)
		gl.BindTexture(gl.TEXTURE_2D, texture)
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RED, 8, 8, 0, gl.RED, gl.UNSIGNED_BYTE, gl.Ptr(pixels))
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

		tr.characters[char] = Character{
			textureID: texture,
			size:      mgl32.Vec2{8, 8},
			bearing:   mgl32.Vec2{0, 8},
			advance:   10,
		}
	}
}

func (tr *TextRenderer) RenderText(text string, x, y, scale float32, color mgl32.Vec3) {
	gl.UseProgram(tr.program)
	gl.UniformMatrix4fv(gl.GetUniformLocation(tr.program, gl.Str("projection\x00")), 1, false, &tr.projection[0])
	gl.Uniform3f(gl.GetUniformLocation(tr.program, gl.Str("textColor\x00")), color.X(), color.Y(), color.Z())
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindVertexArray(tr.VAO)

	for _, char := range text {
		ch, exists := tr.characters[char]
		if !exists {
			continue
		}

		xpos := x + ch.bearing.X()*scale
		ypos := y - (ch.size.Y()-ch.bearing.Y())*scale

		w := ch.size.X() * scale
		h := ch.size.Y() * scale

		vertices := []float32{
			xpos, ypos + h, 0.0, 0.0,
			xpos, ypos, 0.0, 1.0,
			xpos + w, ypos, 1.0, 1.0,

			xpos, ypos + h, 0.0, 0.0,
			xpos + w, ypos, 1.0, 1.0,
			xpos + w, ypos + h, 1.0, 0.0,
		}

		gl.BindTexture(gl.TEXTURE_2D, ch.textureID)
		gl.BindBuffer(gl.ARRAY_BUFFER, tr.VBO)
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(vertices)*4, gl.Ptr(vertices))
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.DrawArrays(gl.TRIANGLES, 0, 6)

		x += float32(ch.advance) * scale
	}

	gl.BindVertexArray(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func compileShader(source string, shaderType uint32) uint32 {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		panic(fmt.Errorf("failed to compile shader: %v", log))
	}

	return shader
}
