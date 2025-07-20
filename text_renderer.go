// text_renderer.go
package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	// Vertex shader for text rendering
	textVertexShader = `#version 430 core
layout (location = 0) in vec4 vertex; // <vec2 pos, vec2 tex>
out vec2 TexCoord;

uniform mat4 projection;

void main() {
    gl_Position = projection * vec4(vertex.xy, 0.0, 1.0);
    TexCoord = vertex.zw;
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

func NewTextRenderer() *TextRenderer {
	tr := &TextRenderer{}
	tr.init()
	return tr
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

func (tr *TextRenderer) RenderFPS(fps float64) {
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	fpsText := fmt.Sprintf("FPS: %.1f", fps)
	tr.RenderText(fpsText, 10, float32(height-30), 1.0, mgl32.Vec3{1, 1, 1})
	gl.Disable(gl.BLEND)
}
