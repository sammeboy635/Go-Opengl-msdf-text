package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	width  = 500
	height = 500
)

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

type Game struct {
	win        *glfw.Window
	cubeRender CubeRender
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

func main() {

	var game Game
	runtime.LockOSThread()

	game.Create_Window()
	game.quadRender.Init()
	game.textRender.Init()
	game.cubeRender.Init()

	defer glfw.Terminate()

	textRendered = false
	cubeRendered = false

	game.Main_Loop()
	PrintMemUsage()
}
func (g *Game) Main_Loop() {
	gl.ClearColor(0.5, 1, 1.0, 1)
	for !g.win.ShouldClose() {
		time.Sleep(100 * time.Millisecond)
		//PrintMemUsage()
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		//g.textRender.Draw_Text()
		//g.quadRender.Draw_Quad()
		g.cubeRender.Draw_Cube()

		g.win.SwapBuffers()
		glfw.PollEvents()
	}
}

func (g *Game) Init() {

}
