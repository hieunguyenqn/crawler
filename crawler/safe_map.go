package crawler

import (
  "sync"
)

/////////////////////////////
// safeMap
/////////////////////////////

type safeMap struct {
  data map[string]interface{}
  lock sync.RWMutex
}

func newSafeMap() *safeMap {
  s := new(safeMap)
  s.data = make(map[string]interface{})
  return s
}

func (s *safeMap) Add(key string, value interface{}) (interface{}, bool) {
  s.lock.Lock()
  defer s.lock.Unlock()
  if val := s.data[key]; val != nil {
    return val, false
  }

  s.data[key] = value
  return value, true
}

func (s *safeMap) Value(key string) interface{} {
  s.lock.RLock()
  defer s.lock.RUnlock()
  return s.data[key]
}
