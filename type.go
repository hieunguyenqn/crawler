package main

import (
  "fmt"
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

func (p *Pages) NewPage(url *url.URL) *Page {
  if page := p.Value(url); page != nil {
    return page
  }
  p.lock.Lock()
  defer p.lock.Unlock()
  page := new(Page)
  page.URL = &URL{url}
  p.safeMap.data[url.String()] = page
  return page
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
  asset.URL = &URL{url}
  a.safeMap.data[url.String()] = asset
  return asset
}

///////////////////////////////
// Page
///////////////////////////////
type Page struct {
  *URL
  Links   []*Page
  Assets  []*Asset
  Visited safeBool
}

///////////////////////////////
// Asset
///////////////////////////////
type Asset struct {
  *URL
}

///////////////////////////////
// URL
///////////////////////////////
type URL struct {
  URL *url.URL
}

func (u URL) ParseRelative(path string) (*url.URL, error) {
  switch {
  case path == "":
    return u.URL, nil
  case path[0] == '#':
    return nil, fmt.Errorf("ID Path: %s", path)
  case len(path) > 2 && path[0:2] == "//":
    return url.Parse(u.Scheme() + ":" + path)
  case path == "/":
    return nil, fmt.Errorf("Root Path: %s", path)
  default:
    return url.Parse(u.subpage(path))
  }
}

func (u URL) subpage(path string) string {
  return fmt.Sprintf("%s://%s%s", u.URL.Scheme, u.URL.Host, path)
}

func (u URL) Host() string {
  return u.URL.Host
}

func (u URL) Scheme() string {
  return u.URL.Scheme
}

func (u URL) String() string {
  return u.URL.String()
}

/////////////////////////////
// safeBool
/////////////////////////////

type safeBool struct {
  b    bool
  lock sync.RWMutex
}

func (s *safeBool) Visit() {
  s.lock.Lock()
  defer s.lock.Unlock()
  s.b = true
}

func (s *safeBool) Unvisit() {
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
