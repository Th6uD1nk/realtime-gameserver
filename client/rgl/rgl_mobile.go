//go:build mobile
package rgl

import (
  "unsafe"
  gl "golang.org/x/mobile/gl"
)

var glctx gl.Context

func Init(ctx gl.Context) error {
  glctx = ctx
  return nil
}

var (
  shaderRegistry  = make(map[uint32]gl.Shader)
  programRegistry = make(map[uint32]gl.Program)
)

const (
  FALSE                  = gl.FALSE
  TRUE                   = gl.TRUE
  INFO_LOG_LENGTH        = gl.INFO_LOG_LENGTH
  COMPILE_STATUS         = gl.COMPILE_STATUS
  LINK_STATUS            = gl.LINK_STATUS
  VERTEX_SHADER          = gl.VERTEX_SHADER
  FRAGMENT_SHADER        = gl.FRAGMENT_SHADER
  ARRAY_BUFFER           = gl.ARRAY_BUFFER
  STATIC_DRAW            = gl.STATIC_DRAW
  TRIANGLES              = gl.TRIANGLES
  LINES                  = gl.LINES
  FLOAT                  = gl.FLOAT
  COLOR_BUFFER_BIT       = gl.COLOR_BUFFER_BIT
  DEPTH_BUFFER_BIT       = gl.DEPTH_BUFFER_BIT
  DEPTH_TEST             = gl.DEPTH_TEST
  BLEND                  = gl.BLEND
  SRC_ALPHA              = gl.SRC_ALPHA
  ONE_MINUS_SRC_ALPHA    = gl.ONE_MINUS_SRC_ALPHA
)

func Str(s string) string {
  return s
}

func GetUniformLocation(program uint32, name string) int32 {
  uniform := glctx.GetUniformLocation(programRegistry[program], name)
  return int32(uniform.Value)
}

func GetAttribLocation(program uint32, name string) int32 {
  attrib := glctx.GetAttribLocation(programRegistry[program], name)
  return int32(attrib.Value)
}

func CreateProgram() uint32 {
  program := glctx.CreateProgram()
  if program.Value == 0 {
    return 0
  }
  programRegistry[program.Value] = program
  return uint32(program.Value)
}

func AttachShader(program uint32, shader uint32) {
  glctx.AttachShader(
    programRegistry[program],
    shaderRegistry[shader],
  )
}

func LinkProgram(program uint32) {
  glctx.LinkProgram(programRegistry[program])
}

func CreateShader(shaderType uint32) uint32 {
  shader := glctx.CreateShader(gl.Enum(shaderType))
  if shader.Value == 0 {
    return 0
  }
  shaderRegistry[shader.Value] = shader
  return shader.Value
}

func Strs(sources ...string) (**uint8, func()) {
  lastShaderSource = sources[0]
  return nil, func() {}
}

var lastShaderSource string

func ShaderSource(shader uint32, count int32, source **uint8, length *int32) {
  glctx.ShaderSource(shaderRegistry[shader], lastShaderSource)
}

func CompileShader(shader uint32) {
  glctx.CompileShader(shaderRegistry[shader])
}

func DeleteShader(shader uint32) {
  glctx.DeleteShader(shaderRegistry[shader])
}

func GetProgramiv(program uint32, pname uint32, params *int32) {
  *params = int32(glctx.GetProgrami(programRegistry[program], gl.Enum(pname)))
}

func GetProgramInfoLog(program uint32, maxLength int32, length *int32, infoLog *byte) {
  logStr := glctx.GetProgramInfoLog(programRegistry[program])
  n := len(logStr)
  if maxLength <= 0 {
    if length != nil {
      *length = 0
    }
    return
  }
  if int32(n) > maxLength-1 {
    n = int(maxLength - 1)
  }
  copy((*[1 << 30]byte)(unsafe.Pointer(infoLog))[:n:n], logStr[:n])
  if length != nil {
    *length = int32(n)
  }
}

func GetShaderiv(shader uint32, pname uint32, params *int32) {
  *params = int32(glctx.GetShaderi(shaderRegistry[shader], gl.Enum(pname)))
}

func GetShaderInfoLog(shader uint32, maxLength int32, length *int32, infoLog *byte) {
  logStr := glctx.GetShaderInfoLog(shaderRegistry[shader])
  n := len(logStr)
  if int32(n) > maxLength-1 {
    n = int(maxLength - 1)
  }
  copy((*[1 << 30]byte)(unsafe.Pointer(infoLog))[:n:n], logStr[:n])
  if length != nil {
    *length = int32(n)
  }
}

func UseProgram(program uint32) {
  glctx.UseProgram(programRegistry[program])
}

func DeleteProgram(program uint32) {
  glctx.DeleteProgram(programRegistry[program])
}

func UniformMatrix4fv(location int32, count int32, transpose bool, value *float32) {
  data := unsafe.Slice(value, 16*count)
  glctx.UniformMatrix4fv(gl.Uniform{Value: location}, data)
}

func Uniform4f(location int32, v0, v1, v2, v3 float32) {
  glctx.Uniform4f(gl.Uniform{Value: location}, v0, v1, v2, v3)
}

func GenBuffers(n int32, buffers *uint32) {
  bufs := unsafe.Slice(buffers, n)
  for i := range bufs {
    buf := glctx.CreateBuffer()
    bufs[i] = uint32(buf.Value)
  }
}

func BindBuffer(target uint32, buffer uint32) {
  glctx.BindBuffer(gl.Enum(target), gl.Buffer{Value: buffer})
}

func BufferData(target uint32, size int, data unsafe.Pointer, usage uint32) {

  var byteData []byte
  if data != nil {
    byteData = unsafe.Slice((*byte)(data), size)
  } else {
    byteData = make([]byte, size)
  }
  glctx.BufferData(gl.Enum(target), byteData, gl.Enum(usage))
}

func Ptr(data interface{}) unsafe.Pointer {

  switch v := data.(type) {
  case []float32:
    if len(v) > 0 {
      return unsafe.Pointer(&v[0])
    }
  case []uint16:
    if len(v) > 0 {
      return unsafe.Pointer(&v[0])
    }
  case []byte:
    if len(v) > 0 {
      return unsafe.Pointer(&v[0])
    }
  }
  return nil
}

func EnableVertexAttribArray(index uint32) {
  glctx.EnableVertexAttribArray(gl.Attrib{Value: uint(index)})
}

func VertexAttribPointer(index uint32, size int32, xtype uint32, normalized bool, stride int32, pointer unsafe.Pointer) {
  offset := int(uintptr(pointer))
  glctx.VertexAttribPointer(gl.Attrib{Value: uint(index)}, int(size), gl.Enum(xtype), normalized, int(stride), offset)
}

func DrawArrays(mode uint32, first int32, count int32) {
  glctx.DrawArrays(gl.Enum(mode), int(first), int(count))
}

func DisableVertexAttribArray(index uint32) {
  glctx.DisableVertexAttribArray(gl.Attrib{Value: uint(index)})
}

func DeleteBuffers(n int32, buffers *uint32) {
  bufs := unsafe.Slice(buffers, n)
  for _, buf := range bufs {
    glctx.DeleteBuffer(gl.Buffer{Value: buf})
  }
}

func Viewport(x, y, width, height int32) {
  glctx.Viewport(int(x), int(y), int(width), int(height))
}

func ClearColor(r, g, b, a float32) {
  glctx.ClearColor(r, g, b, a)
}

func Clear(mask gl.Enum) {
  glctx.Clear(mask)
}

func Enable(cap gl.Enum) {
  glctx.Enable(cap)
}

func Disable(cap gl.Enum) {
  glctx.Disable(cap)
}

func BlendFunc(sfactor, dfactor gl.Enum) {
  glctx.BlendFunc(sfactor, dfactor)
}
