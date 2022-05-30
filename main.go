package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)


var	win_width int = 1280
var	win_height int = 900


var mapping map[byte]Glyph
var textRendered bool
var textLength int
var cubeRendered bool

var (
	triangle = []float32{
		-0.5, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, 0.5, 0,
		0.5, -0.5, 0,
	}
)
var game Game

type Game struct {
	win *glfw.Window
	quadRender QuadRender
	textRender TextRender
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


func (g *Game) Main_Loop() {
	gl.ClearColor(0.5, 0.5, 0.5, 1)
	for !g.win.ShouldClose() {
		time.Sleep(100 * time.Millisecond)

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)	
		start := time.Now()
	
		g.quadRender.Draw_Quad()
		g.textRender.Draw_Text()
		//g.cubeRender.Draw_Cube()
		elapsed := time.Since(start)
		fmt.Println(elapsed)
		g.win.SwapBuffers()
		glfw.PollEvents()
	}
}

func main() {
	runtime.LockOSThread() // Graphic communication needs to be off one CPU

	game.Create_Window()
	game.quadRender.Init()
	game.textRender.Init()

	defer glfw.Terminate()

	textRendered = false
	cubeRendered = false

	game.Main_Loop()
	PrintMemUsage()
}