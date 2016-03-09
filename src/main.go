package main

import (
    "fmt"
    gl "github.com/chsc/gogl/gl33"
    "github.com/veandco/go-sdl2/sdl"
    "lmath"
    "runtime"
)

const (
    winTitle           = "mib's go-playground"
    winWidth           = 640
    winHeight          = 480
    vertexShaderSource = `
        attribute vec3 pos;
        attribute vec3 vertexColor;
        uniform mat4 viewProjMat;
        uniform mat4 worldMat;
        varying vec3 fragmentColor;
        void main() { 
            gl_Position = viewProjMat * (worldMat * vec4(pos, 1.0));
            fragmentColor = vertexColor;
        }`
    fragmentShaderSource = `
        varying vec3 fragmentColor;
        void main() {
            gl_FragColor = vec4(fragmentColor, 1.0);
        }`
)

var IViewProj gl.Int
var IWorld gl.Int

var matViewProj lmath.Mat44
var matWorld lmath.Mat44
var rot lmath.Quat

func main() {
    var window *sdl.Window
    var context sdl.GLContext
    var event sdl.Event
    var running bool
    var err error
    runtime.LockOSThread()
    if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
        panic(err)
    }
    defer sdl.Quit()

    // set antialiasing
    sdl.GL_SetAttribute(sdl.GL_MULTISAMPLEBUFFERS, 1);
    sdl.GL_SetAttribute(sdl.GL_MULTISAMPLESAMPLES, 4);

    // create window and ogl context
    window, err = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED,
        sdl.WINDOWPOS_UNDEFINED,
        winWidth, winHeight, sdl.WINDOW_OPENGL)
    if err != nil {
        panic(err)
    }
    defer window.Destroy()
    context, err = sdl.GL_CreateContext(window)
    if err != nil {
        panic(err)
    }
    defer sdl.GL_DeleteContext(context)

    // init the renderstate
    gl.Init()
    gl.Viewport(0, 0, gl.Sizei(winWidth), gl.Sizei(winHeight))
    gl.ClearColor(0.0, 0.1, 0.0, 1.0)
    gl.Enable(gl.DEPTH_TEST)
    gl.DepthFunc(gl.LESS)
    gl.Enable(gl.BLEND)
    gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

    // vertexbuffer setup
    var vertexbuffer gl.Uint
    gl.GenBuffers(1, &vertexbuffer)
    gl.BindBuffer(gl.ARRAY_BUFFER, vertexbuffer)
    gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(len(triangle_vertices)*4), gl.Pointer(&triangle_vertices[0]), gl.STATIC_DRAW)

    var colourbuffer gl.Uint
    gl.GenBuffers(1, &colourbuffer)
    gl.BindBuffer(gl.ARRAY_BUFFER, colourbuffer)
    gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(len(triangle_colours)*4), gl.Pointer(&triangle_colours[0]), gl.STATIC_DRAW)

    program := createprogram()

    // vertex array format binding
    var VertexArrayID gl.Uint
    gl.GenVertexArrays(1, &VertexArrayID)
    gl.BindVertexArray(VertexArrayID)
    gl.EnableVertexAttribArray(0)
    gl.BindBuffer(gl.ARRAY_BUFFER, vertexbuffer)
    gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
    gl.EnableVertexAttribArray(1)
    gl.BindBuffer(gl.ARRAY_BUFFER, colourbuffer)
    gl.VertexAttribPointer(1, 3, gl.FLOAT, gl.FALSE, 0, nil)

    // get uniform locations and bind program
    IViewProj = gl.GetUniformLocation(program, gl.GLString("viewProjMat"))
    fmt.Printf("ViewProj Link: %v\n", IViewProj+1)

    IWorld = gl.GetUniformLocation(program, gl.GLString("worldMat"))
    fmt.Printf("World Link: %v\n", IWorld+1)

    gl.UseProgram(program)

    // setup the transformations
    matViewProj = lmath.Mat44PerspLH(60, winWidth/winHeight, 0.1, 1000.0);
    matWorld = lmath.Mat44{1,0,0,0, 0,1,0,0, 0,0,1,0, 0,0,0,1}
    rot = lmath.Quat{0,0,0,1}

    running = true
    for running {
        for event = sdl.PollEvent(); event != nil; event =
            sdl.PollEvent() {
            switch t := event.(type) {
            case *sdl.QuitEvent:
                running = false
            case *sdl.MouseMotionEvent:                
                fmt.Printf("[%dms]MouseMotion\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\n", t.Timestamp, t.Which, t.X, t.Y, t.XRel, t.YRel)
                
                // rotate based on mousemove
                rot.RotateInY(float32(t.YRel) * 0.01)
                rot.RotateInX(float32(t.XRel) * 0.01)
            }
        }
        drawgl()
        sdl.GL_SwapWindow(window)
    }
}

func drawgl() {
    // compose the transform-matrix
    rot.RotateInY(1 * lmath.HALF_DEG2RAD)
    matWorld.Recompose(
        rot,
        lmath.Vec3{1,1,1},
        lmath.Vec3{0,0,-7},
    )

    // update uniform data
    gl.UniformMatrix4fv(IViewProj, 1, gl.FALSE, (*gl.Float)(&matViewProj[0]));
    gl.UniformMatrix4fv(IWorld, 1, gl.FALSE, (*gl.Float)(&matWorld[0]));

    // clear and draw
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
    gl.DrawArrays(gl.TRIANGLES, gl.Int(0), gl.Sizei(len(triangle_vertices)*4))
}

func compileShader(shader_id gl.Uint) {
    var status gl.Int
    gl.CompileShader(shader_id)
    gl.GetShaderiv(shader_id, gl.COMPILE_STATUS, &status)
    fmt.Printf("Compiled Vertex Shader: %v\n", status)
    if status == gl.FALSE {
        var info_log_length gl.Int
        gl.GetShaderiv(shader_id, gl.INFO_LOG_LENGTH, &info_log_length)
        if info_log_length == 0 {
            panic("Program linking failed but OpenGL has no log about it.")
        } else {
            info_log_gl := gl.GLStringAlloc(gl.Sizei(info_log_length))
            defer gl.GLStringFree(info_log_gl)
            gl.GetShaderInfoLog(shader_id, gl.Sizei(info_log_length), nil, info_log_gl)
            info_log := gl.GoString(info_log_gl)
            panic(info_log)
        }
    }
}

func createprogram() gl.Uint {
    vs := gl.CreateShader(gl.VERTEX_SHADER)
    vs_source := gl.GLString(vertexShaderSource)
    gl.ShaderSource(vs, 1, &vs_source, nil)    
    compileShader(vs)

    fs := gl.CreateShader(gl.FRAGMENT_SHADER)
    fs_source := gl.GLString(fragmentShaderSource)
    gl.ShaderSource(fs, 1, &fs_source, nil)
    compileShader(fs)

    program := gl.CreateProgram()
    gl.AttachShader(program, vs)
    gl.AttachShader(program, fs)
    gl.LinkProgram(program)
    var linkstatus gl.Int
    gl.GetProgramiv(program, gl.LINK_STATUS, &linkstatus)

    fmt.Printf("Program Link: %v\n", linkstatus)
    if linkstatus == gl.FALSE {
        var info_log_length gl.Int
        gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &info_log_length)
        if info_log_length == 0 {
            panic("Program linking failed but OpenGL has no log about it.")
        } else {
            info_log_gl := gl.GLStringAlloc(gl.Sizei(info_log_length))
            defer gl.GLStringFree(info_log_gl)
            gl.GetProgramInfoLog(program, gl.Sizei(info_log_length), nil, info_log_gl)
            info_log := gl.GoString(info_log_gl)
            panic(info_log)
        }
    }
    return program
}

var triangle_vertices = []gl.Float{
    -1.0, -1.0, -1.0,
    -1.0, -1.0, 1.0,
    -1.0, 1.0, 1.0,
    1.0, 1.0, -1.0,
    -1.0, -1.0, -1.0,
    -1.0, 1.0, -1.0,
    1.0, -1.0, 1.0,
    -1.0, -1.0, -1.0,
    1.0, -1.0, -1.0,
    1.0, 1.0, -1.0,
    1.0, -1.0, -1.0,
    -1.0, -1.0, -1.0,
    -1.0, -1.0, -1.0,
    -1.0, 1.0, 1.0,
    -1.0, 1.0, -1.0,
    1.0, -1.0, 1.0,
    -1.0, -1.0, 1.0,
    -1.0, -1.0, -1.0,
    -1.0, 1.0, 1.0,
    -1.0, -1.0, 1.0,
    1.0, -1.0, 1.0,
    1.0, 1.0, 1.0,
    1.0, -1.0, -1.0,
    1.0, 1.0, -1.0,
    1.0, -1.0, -1.0,
    1.0, 1.0, 1.0,
    1.0, -1.0, 1.0,
    1.0, 1.0, 1.0,
    1.0, 1.0, -1.0,
    -1.0, 1.0, -1.0,
    1.0, 1.0, 1.0,
    -1.0, 1.0, -1.0,
    -1.0, 1.0, 1.0,
    1.0, 1.0, 1.0,
    -1.0, 1.0, 1.0,
    1.0, -1.0, 1.0}

var triangle_colours = []gl.Float{
    0.583, 0.771, 0.014,
    0.609, 0.115, 0.436,
    0.327, 0.483, 0.844,
    0.822, 0.569, 0.201,
    0.435, 0.602, 0.223,
    0.310, 0.747, 0.185,
    0.597, 0.770, 0.761,
    0.559, 0.436, 0.730,
    0.359, 0.583, 0.152,
    0.483, 0.596, 0.789,
    0.559, 0.861, 0.639,
    0.195, 0.548, 0.859,
    0.014, 0.184, 0.576,
    0.771, 0.328, 0.970,
    0.406, 0.615, 0.116,
    0.676, 0.977, 0.133,
    0.971, 0.572, 0.833,
    0.140, 0.616, 0.489,
    0.997, 0.513, 0.064,
    0.945, 0.719, 0.592,
    0.543, 0.021, 0.978,
    0.279, 0.317, 0.505,
    0.167, 0.620, 0.077,
    0.347, 0.857, 0.137,
    0.055, 0.953, 0.042,
    0.714, 0.505, 0.345,
    0.783, 0.290, 0.734,
    0.722, 0.645, 0.174,
    0.302, 0.455, 0.848,
    0.225, 0.587, 0.040,
    0.517, 0.713, 0.338,
    0.053, 0.959, 0.120,
    0.393, 0.621, 0.362,
    0.673, 0.211, 0.457,
    0.820, 0.883, 0.371,
    0.982, 0.099, 0.879}