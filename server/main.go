package main

import (
  "fmt"
)

func main() {
  fmt.Println("=== RTGS Server ===")
  
  server, err := NewServer(8888)
  if err != nil {
    fmt.Printf("Error on create: %v\n", err)
    return
  }
  
  server.Start()
}

