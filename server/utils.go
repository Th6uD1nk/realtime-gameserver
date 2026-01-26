package main

import (
  "math"
  "math/rand"
  "net"
)

func randomSpawn(id string, userType UserType, conn *net.UDPConn,
  minX, maxX, minY, maxY, minZ, maxZ float32) *User {
  
  var location Vector3
  
  location.x = float32(math.Round(float64(minX + rand.Float32()*(maxX-minX))))
  location.y = float32(math.Round(float64(minY + rand.Float32()*(maxY-minY))))
  location.z = float32(math.Round(float64(minZ + rand.Float32()*(maxZ-minZ))))
  
  orientation := rand.Float32() * 360.0
  
  user := NewUser(id, userType, conn)
  user.location = location
  user.orientation = orientation
  
  return user
}

