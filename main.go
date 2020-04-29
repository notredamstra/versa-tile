package main

import (
	"3d-tile-editor/opengl"
	"3d-tile-editor/window"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"log"
	"runtime"
)

var cubeVertices = []float32{
	-1.0, -1.0, -1.0, 0.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,

	// Top
	-1.0, 1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, 1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, 1.0, 1.0, 1.0,

	// Front
	-1.0, -1.0, 1.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, 1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,

	// Back
	-1.0, -1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, -1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	-1.0, 1.0, -1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 1.0,

	// Left
	-1.0, -1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0, -1.0, 0.0, 0.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,

	// Right
	1.0, -1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, -1.0, 0.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	1.0, 1.0, -1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 0.0, 1.0,
}

func main(){
	// Make sure main() runs on main thread (glfw requirement)
	runtime.LockOSThread()

	// Initialize new window with settings
	window := window.NewWindow(1280, 720, "3D Tile Editor")

	// make sure window doesn't close while the loop runs
	defer glfw.Terminate()

	// Initialize OpenGL bindings (must come after setting Window context)
	if err := gl.Init(); err != nil {
		panic(err)
	}

	// Start the actual application loop
	err := programLoop(window)
	if err != nil {
		log.Fatalln(err)
	}
}

func programLoop(window *window.Window) error {
	vertShader, err := opengl.NewShaderFromFile("assets/shaders/basic.vert.glsl", gl.VERTEX_SHADER)
	if err != nil {
		return err
	}

	fragShader, err := opengl.NewShaderFromFile("assets/shaders/basic.frag.glsl", gl.FRAGMENT_SHADER)
	if err != nil {
		return err
	}

	shaders := []*opengl.Shader{vertShader, fragShader}

	// Create new OpenGL program with given shaders
	program, err := opengl.NewProgram(shaders)
	if err != nil {
		return err
	}

	defer program.Delete()

	model := opengl.NewModel("assets/models/DiscoCharacter.obj")

	// Ensure that triangles that are "behind" others do not show in front
	gl.Enable(gl.DEPTH_TEST)

	camera := opengl.NewCamera(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, -90, 0, window.InputManager())

	for !window.ShouldClose() {

		window.StartFrame()

		// update camera position and direction
		camera.Update(window.SinceLastFrame())

		gl.ClearColor(0, 0, 0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)  // depth buffer needed for DEPTH_TEST

		model.Draw(program, camera, window)
	}

	return nil
}