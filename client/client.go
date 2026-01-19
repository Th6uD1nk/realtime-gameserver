package main

import (
  "fmt"
  "net"
  "time"
)

type UDPClient struct {
  Conn *net.UDPConn
}

func NewUDPClient(addr string) (*UDPClient, error) {
  serverAddr, err := net.ResolveUDPAddr("udp", addr)
  if err != nil {
    return nil, err
  }
  conn, err := net.DialUDP("udp", nil, serverAddr)
  if err != nil {
    return nil, err
  }
  return &UDPClient{Conn: conn}, nil
}

func (c *UDPClient) StartReceiving() {
  go func() {
    buffer := make([]byte, 1024)
    for {
      n, err := c.Conn.Read(buffer)
      if err != nil {
        fmt.Printf("Receive error: %v\n", err)
        continue
      }
      fmt.Printf("+ ACK received: %s\n", string(buffer[:n]))
    }
  }()
}

func (c *UDPClient) StartSending() {
  go func() {
    counter := 1
    ticker := time.NewTicker(2 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
      message := fmt.Sprintf("Message %d", counter)
      _, err := c.Conn.Write([]byte(message))
      if err != nil {
        fmt.Printf("Send error: %v\n", err)
        continue
      }
      fmt.Printf("- Sent: %s\n", message)
      counter++
    }
  }()
}
