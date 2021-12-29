package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type DrawData struct {
	program uint32
	VAO     uint32
	VAOSize int
	VBO     uint32
	image   uint32
}

func (d *DrawData) Create_Program(_fragmentShaderSource string, _vertexShaderSource string) {

	//Todo have a debug const for printing errors
	//version := gl.GoStr(gl.GetString(gl.VERSION))
	//log.Println("OpenGL version", version)

	vertexShaderData := Read_File(_vertexShaderSource)
	fragmentShaderData := Read_File(_fragmentShaderSource)

	vertexShader, err := Compile_Shader(vertexShaderData, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := Compile_Shader(fragmentShaderData, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	//Create Program and attach shadders
	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)

	//Error Handling for linking the program
	var status int32
	gl.GetProgramiv(prog, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(prog, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(prog, logLength, nil, gl.Str(log))

		fmt.Println("failed to link program:", log)
	}

	//Delete the shaders after they have been binded
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	d.program = prog
}

func Compile_Shader(source string, shaderType uint32) (uint32, error) {
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

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

//File needs to be shaders/filename.ext
func Read_File(_location string) string {
	data, err := ioutil.ReadFile(_location)
	if err != nil {
		panic("File Error: ")
	}
	return (string(data) + "\x00")
}

// loadImage opens an image file, upload it the the GPU and returns the texture id
//Loading image from file
func (d *DrawData) Load_Image(_file string, _texId uint32) {
	imgFile, err := os.Open(_file)
	if err != nil {
		println("Problem opening image:")
		panic(err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		println("Problem decoding the image:")
		panic(err)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("incorret stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	d.image = loadTexture(rgba, _texId)
}
func loadTexture(_rgba *image.RGBA, _texId uint32) uint32 {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(_texId)

	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(_rgba.Rect.Size().X),
		int32(_rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(_rgba.Pix))

	return texture
}
