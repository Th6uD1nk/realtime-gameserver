package core

import (
  "rtgs-client/rgl"
  "github.com/go-gl/mathgl/mgl32"
)

type Game struct {
  renderer   *Renderer
  udpClient  *UDPClient
}

func NewGame(udpClient *UDPClient, shaders *Shaders) *Game {
  return &Game{
    renderer:  NewRenderer(shaders),
    udpClient: udpClient,
  }
}

func (g *Game) Draw(width, height int) {
  
  worldState := g.udpClient.WorldState;
  
  rgl.Viewport(0, 0, int32(width), int32(height))
  rgl.ClearColor(0.118, 0.118, 0.157, 1.0)
  rgl.Clear(rgl.COLOR_BUFFER_BIT | rgl.DEPTH_BUFFER_BIT)
  rgl.Enable(rgl.DEPTH_TEST)
  
  aspect := float32(width) / float32(height)
  
  localUser := g.udpClient.GetLocalUser()
  if localUser != nil {
    location := mgl32.Vec3{
      float32(localUser.Location.X),
      float32(localUser.Location.Y),
      float32(localUser.Location.Z),
    }
    g.renderer.CameraFollowLocation(location);
  }
  
  g.renderer.UpdateCamera()
  mvp := g.renderer.GetMVP(aspect)

  gridVerts := g.renderer.GetGridVertices(10)
  g.renderer.DrawVertices(gridVerts, [4]float32{0.235, 0.235, 0.314, 1.0}, rgl.LINES, mvp)

  for _, user := range worldState.GetUsers() {
    if !user.IsActive {
      continue
    }

    rgl.Enable(rgl.BLEND)
    rgl.BlendFunc(rgl.SRC_ALPHA, rgl.ONE_MINUS_SRC_ALPHA)

    pos := Vec3{
      X: user.Location.X + 0.5,
      Y: user.Location.Y + 0.5,
      Z: user.Location.Z + 0.5,
    }
    cubeVerts := g.renderer.GetCubeVertices(pos)
    color := GetColorForUserType(user.UserType)
    g.renderer.DrawVertices(cubeVerts, [4]float32{
      float32(color[0]) / 255.0,
      float32(color[1]) / 255.0,
      float32(color[2]) / 255.0,
      0.5,
    }, rgl.TRIANGLES, mvp)

    edges := g.renderer.GetCubeEdgesFromVertices(pos)
    g.renderer.DrawVertices(edges, [4]float32{
      float32(color[0]) / 255.0,
      float32(color[1]) / 255.0,
      float32(color[2]) / 255.0,
      1.0,
    }, rgl.LINES, mvp)

    rgl.Disable(rgl.BLEND)
  }
}
