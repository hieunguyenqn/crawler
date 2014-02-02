package crawler

import (
  "sync"
)

/////////////////////////////
// pageStack
/////////////////////////////

type pageStack struct {
  data []*Page
  lock sync.RWMutex
}

func (s *pageStack) Len() int {
  s.lock.RLock()
  defer s.lock.RUnlock()
  return len(s.data)
}

func (s *pageStack) Push(page *Page) {
  s.lock.Lock()
  defer s.lock.Unlock()
  s.data = append(s.data, page)
}

func (s *pageStack) Pop() (page *Page) {
  s.lock.Lock()
  defer s.lock.Unlock()
  switch len(s.data) {
  case 0:
    return nil
  case 1:
    page = s.data[0]
    s.data = make([]*Page, 0)
  default:
    page, s.data = s.data[0], s.data[1:len(s.data)]
  }
  return page
}
