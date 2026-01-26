package core

import (
  "math"
  "github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
  Position       mgl32.Vec3
  Target         mgl32.Vec3
  Up             mgl32.Vec3
  Yaw            float32 // h
  Pitch          float32
  // tmp
  targetPosition      mgl32.Vec3
  isMoving            bool
  moveProgress        float32
  interpolationSpeed  float32
}

func NewCamera(position, target mgl32.Vec3) *Camera {
  return &Camera{
    Position:           position,
    Target:             target,
    Up:                 mgl32.Vec3{0, 1, 0},
    Yaw:                -90.0,
    Pitch:              0.0,
    // tmp
    targetPosition:     position,
    isMoving:           false,
    moveProgress:       0.0,
    interpolationSpeed: 0.05,
  }
}

func (c *Camera) Update() {
  if c.isMoving {
    c.moveProgress += c.interpolationSpeed
    
    if c.moveProgress >= 1.0 {
      c.Position = c.targetPosition
      c.isMoving = false
      c.moveProgress = 1.0
    } else {
      c.Position = c.Position.Add(
        c.targetPosition.Sub(c.Position).Mul(c.interpolationSpeed),
      )
    }
  }
}

func (c *Camera) SetInterpolationSpeed(speed float32) {
  c.interpolationSpeed = speed
}

func (c *Camera) GetViewMatrix() mgl32.Mat4 {
  return mgl32.LookAtV(c.Position, c.Position.Add(c.Target), c.Up)
}

func (c *Camera) GetProjectionMatrix(fov, aspect, near, far float32) mgl32.Mat4 {
  return mgl32.Perspective(mgl32.DegToRad(fov), aspect, near, far)
}

func (c *Camera) GetTransformMatrix(fov, aspect, near, far float32) mgl32.Mat4 {
  view := c.GetViewMatrix()
  projection := c.GetProjectionMatrix(fov, aspect, near, far)
  return projection.Mul4(view)
}

func (c *Camera) updateCameraVectors() {
  yawRad := mgl32.DegToRad(c.Yaw)
  pitchRad := mgl32.DegToRad(c.Pitch)
  
  front := mgl32.Vec3{
    float32(math.Cos(float64(yawRad)) * math.Cos(float64(pitchRad))),
    float32(math.Sin(float64(pitchRad))),
    float32(math.Sin(float64(yawRad)) * math.Cos(float64(pitchRad))),
  }
  c.Target = front.Normalize()
  
  right := c.Target.Cross(mgl32.Vec3{0, 1, 0}).Normalize()
  c.Up = right.Cross(c.Target).Normalize()
}

func (c *Camera) RotateLeft(angle float32) {
  c.Yaw -= angle
  c.updateCameraVectors()
}

func (c *Camera) RotateRight(angle float32) {
  c.Yaw += angle
  c.updateCameraVectors()
}

func (c *Camera) RotateDown(angle float32) {
  c.Pitch -= angle
  if c.Pitch < -89.0 {
    c.Pitch = -89.0
  }
  c.updateCameraVectors()
}

func (c *Camera) RotateUp(angle float32) {
  c.Pitch += angle
  if c.Pitch > 89.0 {
    c.Pitch = 89.0
  }
  c.updateCameraVectors()
}

func (c *Camera) SetPosition(newPosition mgl32.Vec3) {
  c.targetPosition = newPosition
  c.isMoving = true
  c.moveProgress = 0.0
}


func (c *Camera) FollowPosition(targetPos mgl32.Vec3, maxDistance float32) {
  direction := targetPos.Sub(c.Position)
  distance := direction.Len()
  if distance <= maxDistance {
    return
  }
  normalizedDirection := direction.Normalize()
  newPosition := targetPos.Sub(normalizedDirection.Mul(maxDistance))
  c.SetPosition(newPosition)
}
