package crawler

import (
  "net/url"
  "sync"
)

///////////////////////////////
// Pages
///////////////////////////////
type Pages struct {
  safeMap
}

func (p *Pages) Add(key string, page *Page) {
  p.safeMap.Add(key, page)
}

func (p Pages) Value(key *url.URL) *Page {
  ret := p.safeMap.Value(key.String())
  if ret == nil {
    return nil
  }
  return ret.(*Page)
}

func (p *Pages) NewPage(u *url.URL) (*Page, bool) {
  if page := p.Value(u); page != nil {
    return page, false
  }
  p.lock.Lock()
  defer p.lock.Unlock()
  page := NewPage(u)
  p.safeMap.data[u.String()] = page
  return page, true
}

///////////////////////////////
// Assets
///////////////////////////////
type Assets struct {
  safeMap
}

func (a *Assets) Add(key string, asset *Asset) {
  a.safeMap.Add(key, asset)
}

func (a Assets) Value(key string) (asset *Asset) {
  var ok bool
  if asset, ok = a.safeMap.Value(key).(*Asset); ok {
    return asset
  }
  return nil
}

func (a *Assets) New(url *url.URL) *Asset {
  if asset := a.Value(url.String()); asset != nil {
    return asset
  }
  a.lock.Lock()
  defer a.lock.Unlock()

  asset := new(Asset)
  asset.URL = url
  a.safeMap.data[url.String()] = asset
  return asset
}

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
