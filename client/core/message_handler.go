package core

import (
  "encoding/json"
  "fmt"
  "time"
)

type MessageHandler struct {
  client *UDPClient
}

func NewMessageHandler(client *UDPClient) *MessageHandler {
  return &MessageHandler{
    client: client,
  }
}

func (h *MessageHandler) HandleMessage(buffer []byte, n int) {
  var msg ServerMessage
  if err := json.Unmarshal(buffer[:n], &msg); err != nil {
    fmt.Printf("Parse error: %v\n", err)
    return
  }
  
  switch msg.Type {
  case "connection_confirm":
    h.handleConnectionConfirm(&msg)
  case "world_update":
    h.handleWorldUpdate(&msg)
  default:
    fmt.Printf("Unknown message type: %s\n", msg.Type)
  }
}

func (h *MessageHandler) handleConnectionConfirm(msg *ServerMessage) {
  h.client.LocalUserID = msg.UserID
  fmt.Printf("+ Connected as user: %s\n", h.client.LocalUserID)
}

func (h *MessageHandler) handleWorldUpdate(msg *ServerMessage) {
  receivedIDs := map[string]bool{}
  
  for _, userUpdate := range msg.Users {
    user := &User{
      ID:      userUpdate.ID,
      UserType:  UserType(userUpdate.UserType),
      Location: Vec3{
        X: float64(userUpdate.Location[0]),
        Y: float64(userUpdate.Location[1]),
        Z: float64(userUpdate.Location[2]),
      },
      Orientation: userUpdate.Orientation,
      IsActive:  userUpdate.IsActive,
      LastUpdate:  time.Now(),
      Color:     GetColorForUserType(UserType(userUpdate.UserType)),
    }
    h.client.WorldState.UpdateUser(user)
    receivedIDs[user.ID] = true
  }
  
  for id := range h.client.WorldState.Users {
    if !receivedIDs[id] {
      delete(h.client.WorldState.Users, id)
      fmt.Printf("- User removed: %s\n", id)
    }
  }
}
