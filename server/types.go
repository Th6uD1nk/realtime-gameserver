package main

import (
  "net"
  "time"
  "sync"
)

type Vector3 struct {
  x float32
  y float32
  z float32
}

type WorldUpdate struct {
  Type  string     `json:"type"`
  Users []UserData `json:"users"`
}

type UserData struct {
  ID          string     `json:"id"`
  UserType    string     `json:"user_type"`
  Location    [3]float32 `json:"location"`
  Orientation float32    `json:"orientation"`
  IsActive    bool       `json:"is_active"`
}

type ConnectionConfirm struct {
  Type        string `json:"type"`
  LocalUserID string `json:"user_id"`
}

type Client struct {
  addr     *net.UDPAddr
  lastSeen time.Time
  user     *User
}

type Server struct {
  conn    *net.UDPConn
  clients map[string]*Client
  mu      sync.RWMutex
}
