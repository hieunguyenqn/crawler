package crawler

import "sync"

/////////////////////////////
// safeBool
/////////////////////////////

type safeBool struct {
  b    bool
  lock sync.RWMutex
}

func (s *safeBool) True() {
  s.lock.Lock()
  defer s.lock.Unlock()
  s.b = true
}

func (s *safeBool) False() {
  s.lock.Lock()
  defer s.lock.Unlock()
  s.b = false
}

func (s safeBool) Value() bool {
  s.lock.RLock()
  defer s.lock.RUnlock()
  return s.b
}
