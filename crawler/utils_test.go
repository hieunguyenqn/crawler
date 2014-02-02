package crawler

import (
  "net/url"
)

func newTestPage() *Page {
  u, _ := url.Parse("http://www.example.com")
  return NewPage(u)
}
