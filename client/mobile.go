//go:build mobile

package main

import (
  "os"
  "log"
  "fmt"
  "strings"
  "golang.org/x/mobile/app"
  "golang.org/x/mobile/event/lifecycle"
  "golang.org/x/mobile/event/paint"
  "golang.org/x/mobile/event/size"
  "golang.org/x/mobile/gl"
  "rtgs-client/rgl"
  "rtgs-client/core"
)

func loadConfig() (string, error) {
  path := "/sdcard/Download/rtgs-config.txt"
  data, err := os.ReadFile(path)
  if err != nil {
    fmt.Println("Can't read configuration file");
    return "", err
  }
  return strings.TrimSpace(string(data)), nil
}

func main() {
  
  worldState := core.NewWorldState()

  addr, err := loadConfig()
  if err != nil {
    log.Fatalf("Cannot create UDP client: %v", err)
  }
  
  client, err := core.NewUDPClient(addr, worldState)
  if err != nil {
    log.Fatalf("Cannot create UDP client: %v", err)
  }
  defer client.Conn.Close()

  client.StartReceiving()
  client.StartSending()

  var game *core.Game

  app.Main(func(a app.App) {
    var glctx gl.Context
    var width, height int

    for e := range a.Events() {
      switch e := a.Filter(e).(type) {

      case lifecycle.Event:
        if e.Crosses(lifecycle.StageVisible) == lifecycle.CrossOn {
          glctx, _ = e.DrawContext.(gl.Context)
          
          if game == nil && glctx != nil {
            
            shaders := &core.Shaders{
              Vertex: `
                attribute vec3 aPosition;
                uniform mat4 uMVP;
                void main() {
                  gl_Position = uMVP * vec4(aPosition, 1.0);
                }
              `,
              Fragment: `
                precision mediump float;
                uniform vec4 uColor;
                void main() {
                  gl_FragColor = uColor;
                }
              `,
            }
            rgl.Init(glctx);
            game = core.NewGame(worldState, shaders)
          }

          a.Send(paint.Event{})
        }
      
      case size.Event:
        width = int(e.WidthPx)
        height = int(e.HeightPx)

      case paint.Event:
        if glctx == nil || e.External || game == nil {
          continue
        }
        
        game.Draw(width, height)

        a.Publish()
        a.Send(paint.Event{})
      }
    }
  })
}
