package crawler

import (
  "net/url"
)

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
