package main

import (
  "encoding/binary"
  "fmt"
  "math/rand"
  "os"
  "path/filepath"
  "sync"
  "time"
)

type MapData struct {
  Width  int
  Height int
  MaxVal int
  Data   [][]float32
}

type MapGenerator struct {
  eventManager *EventManager
  currentMap   *MapData
  mu           sync.RWMutex
  lastFilename string
}

func NewMapGenerator(eventManager *EventManager) *MapGenerator {
  return &MapGenerator{
    eventManager: eventManager,
  }
}

func (mg *MapGenerator) Generate(width, height, maxVal int) *MapData {
  rand.Seed(time.Now().UnixNano())
  
  data := make([][]float32, height)
  for y := 0; y < height; y++ {
    data[y] = make([]float32, width)
    for x := 0; x < width; x++ {
      val := float32(rand.Intn(maxVal + 1))
      data[y][x] = val
    }
  }
  
  return &MapData{
    Width:  width,
    Height: height,
    MaxVal: maxVal,
    Data:   data,
  }
}

func (mg *MapGenerator) SaveToFile(mapData *MapData, filename string) error {
  dir := filepath.Dir(filename)
  if err := os.MkdirAll(dir, 0755); err != nil {
    return fmt.Errorf("error while creating folder: %v", err)
  }
  
  file, err := os.Create(filename)
  if err != nil {
    return fmt.Errorf("error while creating file: %v", err)
  }
  defer file.Close()
  
  if err := binary.Write(file, binary.LittleEndian, int32(mapData.Width)); err != nil {
    return err
  }
  if err := binary.Write(file, binary.LittleEndian, int32(mapData.Height)); err != nil {
    return err
  }
  if err := binary.Write(file, binary.LittleEndian, int32(mapData.MaxVal)); err != nil {
    return err
  }
  
  for y := 0; y < mapData.Height; y++ {
    for x := 0; x < mapData.Width; x++ {
      if err := binary.Write(file, binary.LittleEndian, mapData.Data[y][x]); err != nil {
        return err
      }
    }
  }
  
  return nil
}


func (mg *MapGenerator) GenerateAndSave(width, height, maxVal int, filename string) error {
  fmt.Printf("Generates map %dx%d (max: %d)...\n", width, height, maxVal)
  
  mapData := mg.Generate(width, height, maxVal)
  
  if err := mg.SaveToFile(mapData, filename); err != nil {
    return err
  }
  
  mg.mu.Lock()
  mg.currentMap = mapData
  mg.lastFilename = filename
  mg.mu.Unlock()
  
  fmt.Printf("Map saved at %s\n", filename)
  
  event := Event{
    Type: EventMapGenerated,
    Data: map[string]interface{}{
      "filename": filename,
      "width":    width,
      "height":   height,
      "maxVal":   maxVal,
    },
  }
  
  mg.eventManager.DispatchAsync(event)
  
  return nil
}

func LoadMapFromFile(filename string) (*MapData, error) {
  file, err := os.Open(filename)
  if err != nil {
    return nil, fmt.Errorf("error opening files: %v", err)
  }
  defer file.Close()
  
  var width, height, maxVal int32
  
  if err := binary.Read(file, binary.LittleEndian, &width); err != nil {
    return nil, err
  }
  if err := binary.Read(file, binary.LittleEndian, &height); err != nil {
    return nil, err
  }
  if err := binary.Read(file, binary.LittleEndian, &maxVal); err != nil {
    return nil, err
  }
  
  data := make([][]float32, height)
  for y := int32(0); y < height; y++ {
    data[y] = make([]float32, width)
    for x := int32(0); x < width; x++ {
      if err := binary.Read(file, binary.LittleEndian, &data[y][x]); err != nil {
        return nil, err
      }
    }
  }
  
  return &MapData{
    Width:  int(width),
    Height: int(height),
    MaxVal: int(maxVal),
    Data:   data,
  }, nil
}

// todo review
func (mg *MapGenerator) GetMapData() (*MapData, error) {
  mg.mu.RLock()
  if mg.currentMap != nil {
    defer mg.mu.RUnlock()
    return mg.currentMap, nil
  }
  lastFile := mg.lastFilename
  mg.mu.RUnlock()
  
  if lastFile == "" {
    return nil, fmt.Errorf("no map found in memory")
  }
  
  mapData, err := LoadMapFromFile(lastFile)
  if err != nil {
    return nil, err
  }
  
  mg.mu.Lock()
  mg.currentMap = mapData
  mg.mu.Unlock()
  
  return mapData, nil
}
