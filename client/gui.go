package main

import (
  "image/color"
  "math"
  "github.com/hajimehoshi/ebiten/v2"
  "github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
  angle float64
}

type Vec3 struct {
  X, Y, Z float64
}

func project(v Vec3, screenW, screenH float64) (float32, float32) {
  angleX := math.Pi / 6
  angleY := math.Pi / 4

  x := v.X*math.Cos(angleY) - v.Z*math.Sin(angleY)
  z := v.X*math.Sin(angleY) + v.Z*math.Cos(angleY)
  y := v.Y

  y2 := y*math.Cos(angleX) - z*math.Sin(angleX)

  scale := 50.0
  px := x*scale + screenW/2
  py := -y2*scale + screenH/2

  return float32(px), float32(py)
}

func (g *Game) Update() error { return nil }

func (g *Game) Draw(screen *ebiten.Image) {
  screen.Fill(color.RGBA{30, 30, 40, 255})

  w, h := screen.Size()
  screenW, screenH := float64(w), float64(h)

  gridSize := 5
  for i := -gridSize; i <= gridSize; i++ {
    x1, y1 := project(Vec3{float64(i), 0, float64(-gridSize)}, screenW, screenH)
    x2, y2 := project(Vec3{float64(i), 0, float64(gridSize)}, screenW, screenH)
    vector.StrokeLine(screen, x1, y1, x2, y2, 1, color.RGBA{60, 60, 80, 255}, false)

    z1, w1 := project(Vec3{float64(-gridSize), 0, float64(i)}, screenW, screenH)
    z2, w2 := project(Vec3{float64(gridSize), 0, float64(i)}, screenW, screenH)
    vector.StrokeLine(screen, z1, w1, z2, w2, 1, color.RGBA{60, 60, 80, 255}, false)
  }

  cubeVertices := []Vec3{
    {-0.5, 0, -0.5}, {0.5, 0, -0.5}, {0.5, 0, 0.5}, {-0.5, 0, 0.5},
    {-0.5, 1, -0.5}, {0.5, 1, -0.5}, {0.5, 1, 0.5}, {-0.5, 1, 0.5},
  }

  var projected [][2]float32
  for _, v := range cubeVertices {
    px, py := project(v, screenW, screenH)
    projected = append(projected, [2]float32{px, py})
  }

  edges := [][2]int{
    {0, 1}, {1, 2}, {2, 3}, {3, 0},
    {4, 5}, {5, 6}, {6, 7}, {7, 4},
    {0, 4}, {1, 5}, {2, 6}, {3, 7},
  }

  for _, edge := range edges {
    p1 := projected[edge[0]]
    p2 := projected[edge[1]]
    vector.StrokeLine(screen, p1[0], p1[1], p2[0], p2[1],
      2, color.RGBA{0, 255, 100, 255}, false)
  }
}

func (g *Game) Layout(w, h int) (int, int) { return 800, 600 }
