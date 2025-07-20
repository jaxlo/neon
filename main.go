// main.go
package main

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	width  = 1920
	height = 1080
)

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

	app := NewApplication(window)
	app.Run()
}
