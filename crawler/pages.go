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
  pages.safeMap = newSafeMap()
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

func (p *Page) FlattenGraph() map[string]interface{} {
  var f func(p *Page)
  output := make(map[string]interface{})
  visited := make(map[*Page]bool)

  f = func(p *Page) {
    visited[p] = true
    output[p.Path] = p.Flatten()
    for _, l := range p.Links {
      if !visited[l] {
        f(l)
      }
    }
  }
  f(p)
  return output
}

func (p Page) Flatten() map[string]interface{} {
  output := make(map[string]interface{})
  output["Links"] = []string{}
  for _, p := range p.Links {
    output["Links"] = append(output["Links"].([]string), p.Path)
  }
  output["Assets"] = []string{}
  for _, a := range p.Assets {
    output["Assets"] = append(output["Assets"].([]string), a.Path)
  }
  return output
}
