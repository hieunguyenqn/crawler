package crawler

import (
  "net/url"
)

///////////////////////////////
// Pages
///////////////////////////////
type Pages struct {
  *safeMap
}

func NewPages() *Pages {
  pages := new(Pages)
  pages.safeMap = NewsafeMap()
  return pages
}

func (p Pages) Value(key *url.URL) *Page {
  ret := p.safeMap.Value(key.String())
  if ret == nil {
    return nil
  }
  return ret.(*Page)
}

func (p *Pages) NewPage(u *url.URL) (*Page, bool) {
  page := NewPage(u)
  ret, newEntry := p.Add(page.String(), page)
  if !newEntry {
    page = ret.(*Page)
  }
  return page, newEntry
}

///////////////////////////////
// Page
///////////////////////////////
type Page struct {
  Path     string
  *url.URL `json:"-"` // Don't want to encode the URL.
  Links    []*Page
  Assets   []*Asset
}

func NewPage(u *url.URL) *Page {
  return &Page{URL: u, Path: u.String()}
}

func NewPageFromString(u string) (*Page, error) {
  parsedUrl, e := url.Parse(u)
  if e != nil {
    return nil, e
  }

  return NewPage(parsedUrl), nil
}
