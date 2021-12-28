package main

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	width  = 500
	height = 500
)

var mapping map[byte]Glyph
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

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func main() {
	mapping = Text_Json_Parsing("custom-msdf/custom-msdf.json")
	Text_Json_Parsing("custom-msdf/custom-msdf.json")
	var game Game
	runtime.LockOSThread()
	game.win = Create_Window()
	defer glfw.Terminate()

	game.drawBlock = New_Create_DrawData_Block()
	game.drawText = New_Create_DrawData_Text()
	PrintMemUsage()
	gl.ClearColor(0.5, 1, 1.0, 1)
	for !game.win.ShouldClose() {

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		Draw_Text(&game)
		Draw_Cube(&game)
		glfw.PollEvents()
		game.win.SwapBuffers()
	}
	PrintMemUsage()
}
