// application.go
package main

import (
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
)

type Application struct {
	window       *glfw.Window
	renderer     *Renderer
	raytracer    *Raytracer
	textRenderer *TextRenderer

	// FPS tracking
	frameCount int
	lastTime   time.Time
	fps        float64
}

func NewApplication(window *glfw.Window) *Application {
	app := &Application{
		window:     window,
		lastTime:   time.Now(),
		frameCount: 0,
		fps:        0.0,
	}

	app.renderer = NewRenderer()
	app.raytracer = NewRaytracer()
	app.textRenderer = NewTextRenderer()

	return app
}

func (a *Application) Run() {
	for !a.window.ShouldClose() {
		a.updateFPS()
		a.update()
		a.render()
		a.window.SwapBuffers()
		glfw.PollEvents()
	}
}

func (a *Application) updateFPS() {
	a.frameCount++
	currentTime := time.Now()
	elapsed := currentTime.Sub(a.lastTime)

	if elapsed >= time.Second {
		a.fps = float64(a.frameCount) / elapsed.Seconds()
		a.frameCount = 0
		a.lastTime = currentTime
	}
}

func (a *Application) update() {
	// Update camera and other game logic here
}

func (a *Application) render() {
	// Execute raytracing compute shader
	a.raytracer.Render()

	// Render the result to screen
	a.renderer.RenderFullscreenQuad(a.raytracer.GetOutputTexture())

	// Render FPS counter
	a.textRenderer.RenderFPS(a.fps)
}
