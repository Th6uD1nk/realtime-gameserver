package main

import (
  "log"
  "github.com/hajimehoshi/ebiten/v2"
)

func main() {
  // Start UDP client
  client, err := NewUDPClient("127.0.0.1:8888")
  if err != nil {
    log.Fatalf("Cannot create UDP client: %v", err)
  }
  client.StartReceiving()
  client.StartSending()

  // Start Ebiten GUI
  ebiten.SetWindowSize(800, 600)
  ebiten.SetWindowTitle("test")
  game := &Game{}
  if err := ebiten.RunGame(game); err != nil {
    log.Fatalf("Ebiten run failed: %v", err)
  }
}
