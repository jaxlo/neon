// shader_utils.go
package main

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.3-core/gl"
)

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
