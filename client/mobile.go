//go:build mobile

package main

import (
  "log"
  "golang.org/x/mobile/app"
  "golang.org/x/mobile/event/lifecycle"
  "golang.org/x/mobile/event/paint"
  "golang.org/x/mobile/event/size"
  "golang.org/x/mobile/gl"
  "rtgs-client/rgl"
)

func main() {
  worldState := NewWorldState()

  /*
  client, err := NewUDPClient("127.0.0.1:8888", worldState)
  if err != nil {
    log.Fatalf("Cannot create UDP client: %v", err)
  }
  defer client.Conn.Close()

  client.StartReceiving()
  client.StartSending()
  */

  var game *Game

  app.Main(func(a app.App) {
    var glctx gl.Context
    var width, height int

    for e := range a.Events() {
      switch e := a.Filter(e).(type) {

      case lifecycle.Event:
        if e.Crosses(lifecycle.StageVisible) == lifecycle.CrossOn {
          glctx, _ = e.DrawContext.(gl.Context)
          
          if game == nil && glctx != nil {
            rgl.Init(glctx);
            game = NewGame(worldState)
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
