package crawler

import (
  "net/url"
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
// Page
///////////////////////////////
type Page struct {
  *url.URL
  Links  []*Page
  Assets []*Asset
}

func NewPage(u *url.URL) *Page {
  return &Page{URL: u}
}
