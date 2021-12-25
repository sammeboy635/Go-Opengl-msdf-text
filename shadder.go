package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type DrawData struct {
	program uint32
	VAO     uint32
	VAOSize int
}

func New_Create_DrawData() DrawData {
	var drawData DrawData
	drawData.program = Create_Program("shader/frag.shadder", "shader/vert.shadder")
	drawData.VAO = Create_Dynamic_Vao(1024)
	drawData.VAOSize = 1024

	return drawData
}

func Create_Dynamic_Vao(_size int) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, _size, nil, gl.DYNAMIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}

func Create_Program(_fragmentShaderSource string, _vertexShaderSource string) uint32 {

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

	//Delete shaders
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return prog
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
