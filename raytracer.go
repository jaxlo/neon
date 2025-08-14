// raytracer.go
package main

import (
	"io/ioutil"
	"log"
	"math"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Raytracer struct {
	computeProgram uint32
	outputTexture  uint32
}

func NewRaytracer() *Raytracer {
	r := &Raytracer{}
	r.init()
	return r
}

func (r *Raytracer) init() {
	// Load compute shader from file
	computeShaderSource, err := ioutil.ReadFile("raytracer.comp")
	if err != nil {
		log.Fatalf("Failed to load compute shader: %v", err)
	}

	// Add null terminator for OpenGL
	shaderSource := string(computeShaderSource) + "\x00"
	r.computeProgram = createComputeShader(shaderSource)

	// Create output texture
	gl.GenTextures(1, &r.outputTexture)
	gl.BindTexture(gl.TEXTURE_2D, r.outputTexture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA32F, width, height, 0, gl.RGBA, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	// Bind texture to image unit 0
	gl.BindImageTexture(0, r.outputTexture, 0, false, 0, gl.WRITE_ONLY, gl.RGBA32F)
}

func (r *Raytracer) Render() {
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

	// Set uniforms and dispatch compute shader
	gl.UseProgram(r.computeProgram)
	gl.Uniform1f(gl.GetUniformLocation(r.computeProgram, gl.Str("uTime\x00")), time)
	gl.Uniform3f(gl.GetUniformLocation(r.computeProgram, gl.Str("uCameraPos\x00")), cameraPos.X(), cameraPos.Y(), cameraPos.Z())
	gl.Uniform3f(gl.GetUniformLocation(r.computeProgram, gl.Str("uCameraDir\x00")), cameraDir.X(), cameraDir.Y(), cameraDir.Z())
	gl.Uniform3f(gl.GetUniformLocation(r.computeProgram, gl.Str("uCameraUp\x00")), cameraUp.X(), cameraUp.Y(), cameraUp.Z())
	gl.Uniform3f(gl.GetUniformLocation(r.computeProgram, gl.Str("uCameraRight\x00")), cameraRight.X(), cameraRight.Y(), cameraRight.Z())

	gl.DispatchCompute(width/16, height/16, 1)
	gl.MemoryBarrier(gl.SHADER_IMAGE_ACCESS_BARRIER_BIT)
}

func (r *Raytracer) GetOutputTexture() uint32 {
	return r.outputTexture
}
