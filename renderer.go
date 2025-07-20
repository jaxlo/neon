// renderer.go
package main

import (
	"github.com/go-gl/gl/v4.3-core/gl"
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

	// Fragment shader for displaying the raytraced image
	fragmentShader = `#version 430 core
in vec2 TexCoord;
out vec4 FragColor;

uniform sampler2D screenTexture;

void main() {
    FragColor = texture(screenTexture, TexCoord);
}` + "\x00"
)

type Renderer struct {
	shaderProgram uint32
	VAO           uint32
	VBO           uint32
	EBO           uint32
}

func NewRenderer() *Renderer {
	r := &Renderer{}
	r.init()
	return r
}

func (r *Renderer) init() {
	// Create shader program
	r.shaderProgram = createShaderProgram(vertexShader, fragmentShader)

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

func (r *Renderer) RenderFullscreenQuad(texture uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.UseProgram(r.shaderProgram)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.BindVertexArray(r.VAO)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)
}
