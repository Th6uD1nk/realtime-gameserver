package core

import (
  "encoding/json"
  "fmt"
  "net"
  "time"
  "strings"
)

type UDPClient struct {
  Conn       *net.UDPConn
  WorldState *WorldState
  LocalUserID   string
}

func NewUDPClient(addr string, worldState *WorldState) (*UDPClient, error) {
  
  serverAddr, err := net.ResolveUDPAddr("udp", addr)
  if err != nil {
    return nil, err
  }
  
  conn, err := net.DialUDP("udp", nil, serverAddr)
  if err != nil {
    return nil, err
  }

  return &UDPClient{
    Conn:       conn,
    WorldState: worldState,
  }, nil
}

func (client *UDPClient) GetLocalUser() *User {
  if client.LocalUserID == "" {
    return nil
  }
  
  return client.WorldState.GetUser(client.LocalUserID)
}

type ServerMessage struct {
  Type    string        `json:"type"`
  Users   []UserUpdate  `json:"users,omitempty"`
  UserID  string        `json:"user_id,omitempty"`
}

type UserUpdate struct {
  ID          string      `json:"id"`
  UserType    string      `json:"user_type"`
  Location    [3]float32  `json:"location"`
  Orientation float32     `json:"orientation"`
  IsActive    bool        `json:"is_active"`
}

func (client *UDPClient) displayUserCount() {
  client.WorldState.mu.RLock()
  defer client.WorldState.mu.RUnlock()

  count := len(client.WorldState.Users)
  fmt.Printf("+ User count: %d users\n", count)
}

func (client *UDPClient) StartReceiving() {
  
  go func() {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()
    for range ticker.C {
      client.displayUserCount()
    }
  }()

  go func() {
    buffer := make([]byte, 4096)
    for {
      client.Conn.SetReadDeadline(time.Now().Add(5 * time.Second))

      n, err := client.Conn.Read(buffer)
      if err != nil {
        if ne, ok := err.(net.Error); ok && ne.Timeout() {
          fmt.Println("UDP Read timeout: can't reach server")
          continue
        }
        
        errMessage := err.Error()
        formatted := strings.ReplaceAll(errMessage, ": ", "\n\t")
        fmt.Println("Receive error:\n\t" + formatted)
        continue
      }
      
      var msg ServerMessage
      if err := json.Unmarshal(buffer[:n], &msg); err != nil {
        fmt.Printf("Parse error: %v\n", err)
        continue
      }
      
      if msg.Type == "connection_confirm" {
        client.LocalUserID = msg.UserID
        fmt.Printf("+ Connected as user: %s\n", client.LocalUserID)
      }

      if msg.Type == "world_update" {
        receivedIDs := map[string]bool{}
        for _, userUpdate := range msg.Users {
          user := &User{
            ID:           userUpdate.ID,
            UserType:     UserType(userUpdate.UserType),
            Location:     Vec3{
                            X: float64(userUpdate.Location[0]),
                            Y: float64(userUpdate.Location[1]),
                            Z: float64(userUpdate.Location[2]),
                          },
            Orientation:  userUpdate.Orientation,
            IsActive:     userUpdate.IsActive,
            LastUpdate:   time.Now(),
            Color:        GetColorForUserType(UserType(userUpdate.UserType)),
          }
          client.WorldState.UpdateUser(user)
          receivedIDs[user.ID] = true
        }

        for id := range client.WorldState.Users {
          if !receivedIDs[id] {
            delete(client.WorldState.Users, id)
            fmt.Printf("- User removed: %s\n", id)
          }
        }
      }
      
    }
  }()
}

func (client *UDPClient) StartSending() {
  go func() {
    counter := 1
    ticker := time.NewTicker(2 * time.Second)
    defer ticker.Stop()
    for range ticker.C {
      message := fmt.Sprintf("Message %d", counter)
      _, err := client.Conn.Write([]byte(message))
      if err != nil {
        fmt.Printf("Send error: %v\n", err)
        continue
      }
      fmt.Printf("- Sent: %s\n", message)
      counter++
    }
  }()
}

