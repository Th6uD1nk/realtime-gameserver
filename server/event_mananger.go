package main

import (
  "sync"
)

type EventType string

const (
  EventMapGenerated EventType = "map_generated"
)

type Event struct {
  Type EventType
  Data interface{}
}

type EventHandler func(event Event)

type EventManager struct {
  handlers map[EventType][]EventHandler
  mu       sync.RWMutex
}

func NewEventManager() *EventManager {
  return &EventManager{
    handlers: make(map[EventType][]EventHandler),
  }
}

func (em *EventManager) Subscribe(eventType EventType, handler EventHandler) {
  em.mu.Lock()
  defer em.mu.Unlock()
  
  em.handlers[eventType] = append(em.handlers[eventType], handler)
}

func (em *EventManager) Dispatch(event Event) {
  em.mu.RLock()
  handlers := em.handlers[event.Type]
  em.mu.RUnlock()
  
  var wg sync.WaitGroup
  for _, handler := range handlers {
    wg.Add(1)
    go func(h EventHandler) {
      defer wg.Done()
      h(event)
    }(handler)
  }
  wg.Wait()
}

func (em *EventManager) DispatchAsync(event Event) {
  go em.Dispatch(event)
}

