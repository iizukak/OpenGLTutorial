package main

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"strings"
	"os"
	"bufio"
)

const windowWidth = 600
const windowHeight = 480

var aspect = float32(windowWidth) / float32(windowHeight)
var scale float32 = 100.0
var size = [2]float32{float32(windowWidth), float32(windowHeight)}

var points = []float32{
	-0.5, -0.5,
	0.5, -0.5,
	0.5, 0.5,
	-0.5, 0.5,
}

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	window := initGlfw(windowWidth, windowHeight)
	window.SetSizeCallback(resize)
	defer glfw.Terminate()
	defer window.Destroy()

	// Configure the vertex and fragment shaders
	var vertexShader = readFile("./chapter_6/point.vert")
	var fragmentShader = readFile("./chapter_6/point.frag")
	program, err := newProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}
	var aspectLock = gl.GetUniformLocation(program, gl.Str("aspect\x00"))
	var sizeLock = gl.GetUniformLocation(program, gl.Str("size\x00"))
	var scaleLock = gl.GetUniformLocation(program, gl.Str("scale\x00"))

	fmt.Println("OpenGL version:\t", gl.GoStr(gl.GetString(gl.VERSION)))
	fmt.Println("GLSL version:\t", gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION)))
	fmt.Println("GLFW version:\t", glfw.GetVersionString())

	w, h := window.GetSize()
	fw, fh := window.GetFramebufferSize()
	fmt.Println("aspect ratio: ", aspect)
	fmt.Printf("width: %d, height: %d\n", w, h)
	fmt.Printf("frame buffer width: %d, frame buffer height: %d\n", fw, fh)

	vao := makeVao(points)
	defer gl.DeleteBuffers(1, &vao)

	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	gl.Viewport(0, 0, int32(fw), int32(fh))

	for !window.ShouldClose() {
		gl.Uniform1f(aspectLock, aspect)

		fw, fh := window.GetFramebufferSize()
		gl.Uniform2f(sizeLock, float32(fw), float32(fh))
		gl.Uniform1f(scaleLock, scale)

		draw(vao, window, program)
	}
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)

	gl.BindAttribLocation(program, 0, gl.Str("position\x00"))
	gl.BindFragDataLocation(program, 0, gl.Str("fragment\x00"))

	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
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

func readFile(path string) string {
	// ファイルを読み出し用にオープン
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var out string
	for scanner.Scan() {
		out += scanner.Text()
		out += "\n"
	}

	out += "\x00"

	return out
}

func initGlfw(width, height int) *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// make an application window
	window, err := glfw.CreateWindow(width, height, "Hello", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	return window
}

func makeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, nil)

	return vao
}

func draw(vao uint32, window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.LINE_LOOP, 0,4 )

	glfw.PollEvents()
	window.SwapBuffers()
}

func resize(w *glfw.Window, width, height int) {
	// Retina Display must use FrameBufferSize
	fw, fh := w.GetFramebufferSize()
	gl.Viewport(0, 0, int32(fw), int32(fh))
	size[0] = float32(fw)
	size[1] = float32(fh)

	// fmt.Printf("resize called. width: %d, height: %d, aspect: %f\n", int32(width), int32(height), aspect)
	aspect = float32(width) / float32(height)
}