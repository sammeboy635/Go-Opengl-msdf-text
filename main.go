package main

import (
	"runtime"

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
	win       *glfw.Window
	drawBlock DrawData
	drawText  DrawData
}

func main() {

	var game Game
	runtime.LockOSThread()
	game.win = Create_Window()
	defer glfw.Terminate()

	//game.drawBlock = New_Create_DrawData_Block()
	game.drawText = New_Create_DrawData_Text()

	for !game.win.ShouldClose() {

		Draw_Text(&game)
		//Draw_Cube(&game)
	}
}
