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

func (s *safeMap) Add(key string, value interface{}) {
  s.lock.Lock()
  defer s.lock.Unlock()
  s.data[key] = value
}

func (s *safeMap) Value(key string) interface{} {
  s.lock.RLock()
  defer s.lock.RUnlock()
  return s.data[key]
}
