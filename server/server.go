package main

import (
  "encoding/json"
  "fmt"
  "net"
  "time"
)

func NewServer(port int) (*Server, error) {
  fmt.Println("Creating UDP address...")
  addr := net.UDPAddr{
    Port: port,
    IP:   net.ParseIP("0.0.0.0"),
  }
  
  fmt.Printf("Binding to %s:%d...\n", addr.IP, addr.Port)
  conn, err := net.ListenUDP("udp", &addr)
  if err != nil {
    return nil, err
  }
  
  fmt.Println("UDP socket bound successfully")
  
  return &Server{
    conn:    conn,
    clients: make(map[string]*Client),
  }, nil
}

func (server *Server) listClients() {
  server.mu.RLock()
  defer server.mu.RUnlock()
  
  fmt.Println("\n=== Clients list ===")
  if len(server.clients) == 0 {
    fmt.Println("No connected client")
  } else {
    for key, client := range server.clients {
      fmt.Printf("- %s\n", key)
      fmt.Printf("  Type: %s\n", client.user.userType)
      fmt.Printf("  Location: (%.2f, %.2f, %.2f)\n", 
        client.user.location.x, client.user.location.y, client.user.location.z)
      fmt.Printf("  Orientation: %.2f°\n", client.user.orientation)
      fmt.Printf("  Active: %t\n", client.user.isActive)
      fmt.Printf("  Last activity: %s\n", time.Since(client.lastSeen).Round(time.Second))
    }
  }
  fmt.Println("========================\n")
}

func (server *Server) cleanInactiveClients(timeout time.Duration) {
  server.mu.Lock()
  defer server.mu.Unlock()
  
  now := time.Now()
  for key, client := range server.clients {
    if now.Sub(client.lastSeen) > timeout {
      fmt.Printf("x Client timeout: %s\n", key)
      delete(server.clients, key)
    }
  }
}

func (server *Server) broadcastWorldState() {
  server.mu.RLock()
  
  worldUpdate := WorldUpdate{
    Type:  "world_update",
    Users: make([]UserData, 0, len(server.clients)),
  }
  
  for _, client := range server.clients {
    if client.user == nil {
      fmt.Printf("WARNING: Client %s has NIL user!\n", client.addr.String())
      continue
    }
    worldUpdate.Users = append(worldUpdate.Users, UserData{
      ID:          client.user.id,
      UserType:    string(client.user.userType),
      Location:    [3]float32{client.user.location.x, client.user.location.y, client.user.location.z},
      Orientation: client.user.orientation,
      IsActive:    client.user.isActive,
    })
  }
  
  server.mu.RUnlock()
  
  data, err := json.Marshal(worldUpdate)
  if err != nil {
    fmt.Printf("JSON marshal error: %v\n", err)
    return
  }
  
  server.mu.RLock()
  for _, client := range server.clients {
    _, err := server.conn.WriteToUDP(data, client.addr)
    if err != nil {
      fmt.Printf("x Broadcast error to %s: %v\n", client.addr.String(), err)
    }
  }
  server.mu.RUnlock()
}

func (server *Server) sendConnectionConfirm(addr *net.UDPAddr, clientID string) {
  confirmMsg := ConnectionConfirm{
    Type:        "connection_confirm",
    LocalUserID: clientID,
  }
  
  confirmData, err := json.Marshal(confirmMsg)
  if err != nil {
    fmt.Printf("JSON marshal error for connection_confirm: %v\n", err)
    return
  }
  
  _, err = server.conn.WriteToUDP(confirmData, addr)
  if err != nil {
    fmt.Printf("x Failed to send connection_confirm to %s: %v\n", addr.String(), err)
  } else {
    fmt.Printf("+ Sent connection_confirm to %s (ID: %s)\n", addr.String(), clientID)
  }
}

func (server *Server) startBackgroundTasks() {
  // Clean inactive clients
  go func() {
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()
    for range ticker.C {
      server.cleanInactiveClients(10 * time.Second)
    }
  }()
  
  // Display client list periodically
  go func() {
    ticker := time.NewTicker(15 * time.Second)
    defer ticker.Stop()
    for range ticker.C {
      server.listClients()
    }
  }()
  
  // Broadcast world state
  go func() {
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()
    for range ticker.C {
      server.broadcastWorldState()
    }
  }()
}

func (server *Server) handleNewClient(addr *net.UDPAddr, clientKey string) {
  user := randomSpawn(clientKey, UserTypePlayer, server.conn, 0, 10, 0, 0, 0, 10)
  
  client := &Client{
    addr:     addr,
    lastSeen: time.Now(),
    user:     user,
  }
  
  server.clients[clientKey] = client
  
  fmt.Printf("+ new client spawned at (%.2f, %.2f, %.2f) orientation: %.2f\n",
    user.location.x, user.location.y, user.location.z, user.orientation)
  
  server.sendConnectionConfirm(addr, clientKey)
}

func (server *Server) Start() {
  defer server.conn.Close()
  
  fmt.Printf("UDP server started on port %d\n", server.conn.LocalAddr().(*net.UDPAddr).Port)
  
  server.startBackgroundTasks()
  
  buffer := make([]byte, 1024)
  
  for {
    nByte, addr, err := server.conn.ReadFromUDP(buffer)
    if err != nil {
      fmt.Printf("Read error %v\n", err)
      continue
    }
    
    clientKey := addr.String()
    
    server.mu.Lock()
    client, exists := server.clients[clientKey]
    if !exists {
      server.handleNewClient(addr, clientKey)
    } else {
      client.lastSeen = time.Now()
    }
    server.mu.Unlock()
    
    message := string(buffer[:nByte])
    fmt.Printf("+ received from %s: %s\n", addr.String(), message)
  }
}
