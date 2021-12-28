package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func Create_Window() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	//Init gl
	if err := gl.Init(); err != nil {
		panic(err)
	}
	//glfw.WindowHint(glfw.OpenGLDebugContext, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Sam Waite", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	//gl.Enable(gl.DEPTH_TEST)
	//window.SetMouseButtonCallback(callback.mouse_button())

	window.SetFramebufferSizeCallback(window_size_callback)
	window.SetKeyCallback(window_key_callback)
	window.SetMouseButtonCallback(window_mouse_button_callback)

	return window
}

func window_size_callback(win *glfw.Window, width int, height int) {
	println("hello")
}
func window_key_callback(win *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	println(key, action, mods)
}
func window_mouse_button_callback(win *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	println("Mouse click")
}
