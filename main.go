package main

import (
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	width  = 500
	height = 500
)

var (
	triangle = []float32{
		-0.5, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, 0.5, 0,
		0.5, -0.5, 0,
	}
)

type Game struct {
	win      *glfw.Window
	drawData DrawData
}

func main() {

	var game Game
	runtime.LockOSThread()
	game.win = Create_Window()
	defer glfw.Terminate()

	game.drawData = New_Create_DrawData()

	for !game.win.ShouldClose() {
		draw(&game)
	}
}

func draw(game *Game) {

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(game.drawData.program)

	cube := New_Cube(0.0, 0.0, 0.0)
	cube = append(cube, New_Cube(-1.0, -1.0, 0)...)
	//gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, game.drawData.VAO)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(cube)*4, gl.Ptr(cube))
	//gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle)/3))
	//gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 8)
	gl.DrawArraysInstanced(gl.TRIANGLE_STRIP, 0, 6, 8)

	glfw.PollEvents()
	game.win.SwapBuffers()
}

// initGlfw initializes glfw and returns a Window to use.

// makeVao initializes and returns a vertex array from the points provided.
