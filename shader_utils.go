// shader_utils.go
package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/go-gl/gl/v4.3-core/gl"
)

// ShaderType represents different types of shaders
type ShaderType int

const (
	VertexShader ShaderType = iota
	FragmentShader
	ComputeShader
	GeometryShader
	TessControlShader
	TessEvaluationShader
)

// ShaderConfig holds configuration for creating shader programs
type ShaderConfig struct {
	Vertex         string // Path to vertex shader file
	Fragment       string // Path to fragment shader file
	Geometry       string // Path to geometry shader file (optional)
	TessControl    string // Path to tessellation control shader file (optional)
	TessEvaluation string // Path to tessellation evaluation shader file (optional)
}

// getGLShaderType converts our ShaderType to OpenGL constants
func getGLShaderType(shaderType ShaderType) uint32 {
	switch shaderType {
	case VertexShader:
		return gl.VERTEX_SHADER
	case FragmentShader:
		return gl.FRAGMENT_SHADER
	case ComputeShader:
		return gl.COMPUTE_SHADER
	case GeometryShader:
		return gl.GEOMETRY_SHADER
	case TessControlShader:
		return gl.TESS_CONTROL_SHADER
	case TessEvaluationShader:
		return gl.TESS_EVALUATION_SHADER
	default:
		panic(fmt.Errorf("unknown shader type: %d", shaderType))
	}
}

// loadShaderFile loads a shader source from a file
func loadShaderFile(filepath string) (string, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to load shader file %s: %v", filepath, err)
	}
	return string(data) + "\x00", nil
}

// createShaderProgram creates a basic vertex + fragment shader program
func createShaderProgram(vertexSource, fragmentSource string) uint32 {
	vertexShader := compileShader(vertexSource, gl.VERTEX_SHADER)
	fragmentShader := compileShader(fragmentSource, gl.FRAGMENT_SHADER)

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	checkProgramLink(program)

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program
}

// createShaderProgramFromFiles creates a shader program from file paths
func createShaderProgramFromFiles(vertexPath, fragmentPath string) (uint32, error) {
	vertexSource, err := loadShaderFile(vertexPath)
	if err != nil {
		return 0, err
	}

	fragmentSource, err := loadShaderFile(fragmentPath)
	if err != nil {
		return 0, err
	}

	return createShaderProgram(vertexSource, fragmentSource), nil
}

// createComplexShaderProgram creates a shader program with multiple shader stages
func createComplexShaderProgram(config ShaderConfig) (uint32, error) {
	program := gl.CreateProgram()
	var shaders []uint32

	// Load and compile vertex shader (required)
	if config.Vertex != "" {
		source, err := loadShaderFile(config.Vertex)
		if err != nil {
			return 0, err
		}
		shader := compileShader(source, gl.VERTEX_SHADER)
		gl.AttachShader(program, shader)
		shaders = append(shaders, shader)
	}

	// Load and compile fragment shader (required)
	if config.Fragment != "" {
		source, err := loadShaderFile(config.Fragment)
		if err != nil {
			return 0, err
		}
		shader := compileShader(source, gl.FRAGMENT_SHADER)
		gl.AttachShader(program, shader)
		shaders = append(shaders, shader)
	}

	// Load and compile geometry shader (optional)
	if config.Geometry != "" {
		source, err := loadShaderFile(config.Geometry)
		if err != nil {
			return 0, err
		}
		shader := compileShader(source, gl.GEOMETRY_SHADER)
		gl.AttachShader(program, shader)
		shaders = append(shaders, shader)
	}

	// Load and compile tessellation control shader (optional)
	if config.TessControl != "" {
		source, err := loadShaderFile(config.TessControl)
		if err != nil {
			return 0, err
		}
		shader := compileShader(source, gl.TESS_CONTROL_SHADER)
		gl.AttachShader(program, shader)
		shaders = append(shaders, shader)
	}

	// Load and compile tessellation evaluation shader (optional)
	if config.TessEvaluation != "" {
		source, err := loadShaderFile(config.TessEvaluation)
		if err != nil {
			return 0, err
		}
		shader := compileShader(source, gl.TESS_EVALUATION_SHADER)
		gl.AttachShader(program, shader)
		shaders = append(shaders, shader)
	}

	gl.LinkProgram(program)
	checkProgramLink(program)

	// Clean up individual shaders
	for _, shader := range shaders {
		gl.DeleteShader(shader)
	}

	return program, nil
}

// createComputeShader creates a compute shader program
func createComputeShader(source string) uint32 {
	shader := compileShader(source, gl.COMPUTE_SHADER)

	program := gl.CreateProgram()
	gl.AttachShader(program, shader)
	gl.LinkProgram(program)

	checkProgramLink(program)

	gl.DeleteShader(shader)
	return program
}

// createComputeShaderFromFile creates a compute shader program from a file
func createComputeShaderFromFile(filepath string) (uint32, error) {
	source, err := loadShaderFile(filepath)
	if err != nil {
		return 0, err
	}
	return createComputeShader(source), nil
}

// compileShader compiles a shader from source code
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

		panic(fmt.Errorf("failed to compile %s shader: %v", getShaderTypeName(shaderType), log))
	}

	return shader
}

// checkProgramLink checks if a shader program linked successfully
func checkProgramLink(program uint32) {
	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		panic(fmt.Errorf("failed to link program: %v", log))
	}
}

// getShaderTypeName returns a human-readable name for a shader type
func getShaderTypeName(shaderType uint32) string {
	switch shaderType {
	case gl.VERTEX_SHADER:
		return "vertex"
	case gl.FRAGMENT_SHADER:
		return "fragment"
	case gl.COMPUTE_SHADER:
		return "compute"
	case gl.GEOMETRY_SHADER:
		return "geometry"
	case gl.TESS_CONTROL_SHADER:
		return "tessellation control"
	case gl.TESS_EVALUATION_SHADER:
		return "tessellation evaluation"
	default:
		return "unknown"
	}
}

// loadShadersFromDirectory loads all shaders from a directory based on naming convention
// Expected naming: vertex.glsl, fragment.glsl, geometry.glsl, etc.
func loadShadersFromDirectory(dir string) (ShaderConfig, error) {
	config := ShaderConfig{}

	// Check for vertex shader
	vertexPath := filepath.Join(dir, "vertex.glsl")
	if fileExists(vertexPath) {
		config.Vertex = vertexPath
	}

	// Check for fragment shader
	fragmentPath := filepath.Join(dir, "fragment.glsl")
	if fileExists(fragmentPath) {
		config.Fragment = fragmentPath
	}

	// Check for geometry shader
	geometryPath := filepath.Join(dir, "geometry.glsl")
	if fileExists(geometryPath) {
		config.Geometry = geometryPath
	}

	// Check for tessellation shaders
	tessControlPath := filepath.Join(dir, "tess_control.glsl")
	if fileExists(tessControlPath) {
		config.TessControl = tessControlPath
	}

	tessEvalPath := filepath.Join(dir, "tess_evaluation.glsl")
	if fileExists(tessEvalPath) {
		config.TessEvaluation = tessEvalPath
	}

	return config, nil
}

// fileExists checks if a file exists
func fileExists(filepath string) bool {
	_, err := ioutil.ReadFile(filepath)
	return err == nil
}
