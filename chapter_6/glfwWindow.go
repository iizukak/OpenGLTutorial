package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"fmt"
)

type glfwWindow struct {
	width int
	height int
	title string
	window *glfw.Window
}

func (w glfwWindow) init () {
	// make an application window
	window, err := glfw.CreateWindow(w.width, w.height, w.title, nil, nil)
	w.window = window
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// init gl
	if err := gl.Init(); err != nil {
		panic(err)
	}

	fmt.Println("OpenGL version:\t", gl.GoStr(gl.GetString(gl.VERSION)))
	fmt.Println("GLSL version:\t", gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION)))
	fmt.Println("GLFW version:\t", glfw.GetVersionString())

	glfw.SwapInterval(1)
}