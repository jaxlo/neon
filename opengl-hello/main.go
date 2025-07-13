package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	windowWidth   = 800
	windowHeight  = 600
	workGroupSize = 16
)

// Vertex data for a full-screen quad
var quadVertices = []float32{
	-1.0, -1.0, 0.0, 0.0,
	1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, 1.0, 1.0,
	-1.0, 1.0, 0.0, 1.0,
}

var quadIndices = []uint32{
	0, 1, 2,
	2, 3, 0,
}

// Compute shader for raytracing
const computeShaderSource = `#version 430

layout(local_size_x = 16, local_size_y = 16) in;
layout(rgba32f, binding = 0) uniform image2D img_output;

uniform float u_time;
uniform vec2 u_resolution;
uniform vec3 u_camera_pos;
uniform vec3 u_camera_dir;

struct Ray {
    vec3 origin;
    vec3 direction;
};

struct Sphere {
    vec3 center;
    float radius;
    vec3 color;
    float reflectivity;
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
const int NUM_SPHERES = 4;
Sphere spheres[NUM_SPHERES] = Sphere[](
    Sphere(vec3(0.0, 0.0, -5.0), 1.0, vec3(1.0, 0.2, 0.2), 0.3),
    Sphere(vec3(-2.0, 0.0, -6.0), 0.8, vec3(0.2, 1.0, 0.2), 0.5),
    Sphere(vec3(2.0, 0.0, -4.0), 0.6, vec3(0.2, 0.2, 1.0), 0.2),
    Sphere(vec3(0.0, -1001.0, -5.0), 1000.0, vec3(0.8, 0.8, 0.8), 0.0)
);

// Light position
vec3 lightPos = vec3(2.0, 2.0, -2.0);

HitInfo intersectSphere(Ray ray, Sphere sphere) {
    HitInfo hit;
    hit.hit = false;
    
    vec3 oc = ray.origin - sphere.center;
    float a = dot(ray.direction, ray.direction);
    float b = 2.0 * dot(oc, ray.direction);
    float c = dot(oc, oc) - sphere.radius * sphere.radius;
    
    float discriminant = b * b - 4.0 * a * c;
    
    if (discriminant < 0.0) {
        return hit;
    }
    
    float t = (-b - sqrt(discriminant)) / (2.0 * a);
    
    if (t > 0.001) {
        hit.hit = true;
        hit.distance = t;
        hit.point = ray.origin + t * ray.direction;
        hit.normal = normalize(hit.point - sphere.center);
        hit.color = sphere.color;
        hit.reflectivity = sphere.reflectivity;
    }
    
    return hit;
}

HitInfo intersectScene(Ray ray) {
    HitInfo closestHit;
    closestHit.hit = false;
    closestHit.distance = 10000.0;
    
    for (int i = 0; i < NUM_SPHERES; i++) {
        HitInfo hit = intersectSphere(ray, spheres[i]);
        if (hit.hit && hit.distance < closestHit.distance) {
            closestHit = hit;
        }
    }
    
    return closestHit;
}

vec3 trace(Ray ray, int maxDepth) {
    vec3 color = vec3(0.0);
    vec3 attenuation = vec3(1.0);
    
    for (int depth = 0; depth < maxDepth; depth++) {
        HitInfo hit = intersectScene(ray);
        
        if (!hit.hit) {
            // Sky color
            vec3 skyColor = mix(vec3(0.5, 0.7, 1.0), vec3(1.0, 1.0, 1.0), 
                              0.5 * (ray.direction.y + 1.0));
            color += attenuation * skyColor;
            break;
        }
        
        // Lighting calculation
        vec3 lightDir = normalize(lightPos - hit.point);
        float diffuse = max(dot(hit.normal, lightDir), 0.0);
        
        // Shadow ray
        Ray shadowRay;
        shadowRay.origin = hit.point + hit.normal * 0.001;
        shadowRay.direction = lightDir;
        
        HitInfo shadowHit = intersectScene(shadowRay);
        if (shadowHit.hit && shadowHit.distance < length(lightPos - hit.point)) {
            diffuse *= 0.2; // In shadow
        }
        
        vec3 ambient = vec3(0.1);
        vec3 lighting = ambient + diffuse * vec3(1.0);
        
        color += attenuation * hit.color * lighting * (1.0 - hit.reflectivity);
        
        if (hit.reflectivity > 0.0) {
            attenuation *= hit.reflectivity;
            ray.origin = hit.point + hit.normal * 0.001;
            ray.direction = reflect(ray.direction, hit.normal);
        } else {
            break;
        }
    }
    
    return color;
}

void main() {
    ivec2 pixel_coords = ivec2(gl_GlobalInvocationID.xy);
    
    if (pixel_coords.x >= int(u_resolution.x) || pixel_coords.y >= int(u_resolution.y)) {
        return;
    }
    
    // Convert pixel coordinates to NDC
    vec2 uv = (vec2(pixel_coords) + 0.5) / u_resolution;
    uv = uv * 2.0 - 1.0;
    uv.x *= u_resolution.x / u_resolution.y;
    
    // Animate camera
    vec3 cameraPos = u_camera_pos + vec3(sin(u_time * 0.5) * 0.5, 0.0, 0.0);
    
    // Create ray
    Ray ray;
    ray.origin = cameraPos;
    ray.direction = normalize(vec3(uv.x, uv.y, -1.0));
    
    // Trace the ray
    vec3 color = trace(ray, 5);
    
    // Gamma correction
    color = pow(color, vec3(1.0 / 2.2));
    
    imageStore(img_output, pixel_coords, vec4(color, 1.0));
}
` + "\x00"

// Vertex shader for rendering the texture
const vertexShaderSource = `#version 430 core

layout (location = 0) in vec2 aPos;
layout (location = 1) in vec2 aTexCoord;

out vec2 TexCoord;

void main() {
    gl_Position = vec4(aPos, 0.0, 1.0);
    TexCoord = aTexCoord;
}
` + "\x00"

// Fragment shader for rendering the texture
const fragmentShaderSource = `#version 430 core

in vec2 TexCoord;
out vec4 FragColor;

uniform sampler2D ourTexture;

void main() {
    FragColor = texture(ourTexture, TexCoord);
}
` + "\x00"

type Raytracer struct {
	window            *glfw.Window
	computeProgram    uint32
	renderProgram     uint32
	texture           uint32
	vao, vbo, ebo     uint32
	timeUniform       int32
	resolutionUniform int32
	cameraPosUniform  int32
	cameraDirUniform  int32
}

func init() {
	runtime.LockOSThread()
}

func (r *Raytracer) initWindow() error {
	if err := glfw.Init(); err != nil {
		return fmt.Errorf("failed to initialize GLFW: %v", err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.Resizable, glfw.False)

	var err error
	r.window, err = glfw.CreateWindow(windowWidth, windowHeight, "GPU Raytracer", nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create window: %v", err)
	}

	r.window.MakeContextCurrent()
	glfw.SwapInterval(1)

	if err := gl.Init(); err != nil {
		return fmt.Errorf("failed to initialize OpenGL: %v", err)
	}

	return nil
}

func (r *Raytracer) compileShader(shaderType uint32, source string) (uint32, error) {
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
		return 0, fmt.Errorf("failed to compile shader: %v", log)
	}

	return shader, nil
}

func (r *Raytracer) createProgram(vertexShader, fragmentShader uint32) (uint32, error) {
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
		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	return program, nil
}

func (r *Raytracer) initShaders() error {
	// Compile compute shader
	computeShader, err := r.compileShader(gl.COMPUTE_SHADER, computeShaderSource)
	if err != nil {
		return err
	}
	defer gl.DeleteShader(computeShader)

	r.computeProgram = gl.CreateProgram()
	gl.AttachShader(r.computeProgram, computeShader)
	gl.LinkProgram(r.computeProgram)

	var status int32
	gl.GetProgramiv(r.computeProgram, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(r.computeProgram, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(r.computeProgram, logLength, nil, gl.Str(log))
		return fmt.Errorf("failed to link compute program: %v", log)
	}

	// Get compute shader uniforms
	r.timeUniform = gl.GetUniformLocation(r.computeProgram, gl.Str("u_time\x00"))
	r.resolutionUniform = gl.GetUniformLocation(r.computeProgram, gl.Str("u_resolution\x00"))
	r.cameraPosUniform = gl.GetUniformLocation(r.computeProgram, gl.Str("u_camera_pos\x00"))
	r.cameraDirUniform = gl.GetUniformLocation(r.computeProgram, gl.Str("u_camera_dir\x00"))

	// Compile vertex and fragment shaders for rendering
	vertexShader, err := r.compileShader(gl.VERTEX_SHADER, vertexShaderSource)
	if err != nil {
		return err
	}
	defer gl.DeleteShader(vertexShader)

	fragmentShader, err := r.compileShader(gl.FRAGMENT_SHADER, fragmentShaderSource)
	if err != nil {
		return err
	}
	defer gl.DeleteShader(fragmentShader)

	r.renderProgram, err = r.createProgram(vertexShader, fragmentShader)
	if err != nil {
		return err
	}

	return nil
}

func (r *Raytracer) initBuffers() {
	// Create VAO
	gl.GenVertexArrays(1, &r.vao)
	gl.BindVertexArray(r.vao)

	// Create VBO
	gl.GenBuffers(1, &r.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, r.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(quadVertices)*4, gl.Ptr(quadVertices), gl.STATIC_DRAW)

	// Create EBO
	gl.GenBuffers(1, &r.ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, r.ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(quadIndices)*4, gl.Ptr(quadIndices), gl.STATIC_DRAW)

	// Position attribute
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 4*4, nil)
	gl.EnableVertexAttribArray(0)

	// Texture coordinate attribute
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(2*4))
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(0)
}

func (r *Raytracer) initTexture() {
	gl.GenTextures(1, &r.texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, r.texture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA32F, windowWidth, windowHeight, 0, gl.RGBA, gl.FLOAT, nil)
	gl.BindImageTexture(0, r.texture, 0, false, 0, gl.WRITE_ONLY, gl.RGBA32F)
}

func (r *Raytracer) Init() error {
	if err := r.initWindow(); err != nil {
		return err
	}

	if err := r.initShaders(); err != nil {
		return err
	}

	r.initBuffers()
	r.initTexture()

	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	return nil
}

func (r *Raytracer) Render() {
	for !r.window.ShouldClose() {
		glfw.PollEvents()

		// Run compute shader
		gl.UseProgram(r.computeProgram)

		time := float32(glfw.GetTime())
		gl.Uniform1f(r.timeUniform, time)
		gl.Uniform2f(r.resolutionUniform, windowWidth, windowHeight)
		gl.Uniform3f(r.cameraPosUniform, 0.0, 0.0, 0.0)
		gl.Uniform3f(r.cameraDirUniform, 0.0, 0.0, -1.0)

		gl.DispatchCompute(
			(windowWidth+workGroupSize-1)/workGroupSize,
			(windowHeight+workGroupSize-1)/workGroupSize,
			1)

		gl.MemoryBarrier(gl.SHADER_IMAGE_ACCESS_BARRIER_BIT)

		// Render the texture to screen
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.UseProgram(r.renderProgram)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, r.texture)
		gl.Uniform1i(gl.GetUniformLocation(r.renderProgram, gl.Str("ourTexture\x00")), 0)

		gl.BindVertexArray(r.vao)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
		gl.BindVertexArray(0)

		r.window.SwapBuffers()
	}
}

func (r *Raytracer) Cleanup() {
	gl.DeleteProgram(r.computeProgram)
	gl.DeleteProgram(r.renderProgram)
	gl.DeleteTextures(1, &r.texture)
	gl.DeleteBuffers(1, &r.vbo)
	gl.DeleteBuffers(1, &r.ebo)
	gl.DeleteVertexArrays(1, &r.vao)

	r.window.Destroy()
	glfw.Terminate()
}

func main() {
	raytracer := &Raytracer{}

	if err := raytracer.Init(); err != nil {
		log.Fatalf("Failed to initialize raytracer: %v", err)
	}
	defer raytracer.Cleanup()

	fmt.Println("GPU Raytracer initialized successfully!")
	fmt.Println("Press ESC to exit")

	raytracer.window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if key == glfw.KeyEscape && action == glfw.Press {
			w.SetShouldClose(true)
		}
	})

	raytracer.Render()
}
