package crawler

import (
  "fmt"
  "math/rand"
  "time"
)

func init() {
  rand.Seed(time.Now().UnixNano())
}

func newTestPage() *Page {
  p, e := NewPageFromString(fmt.Sprintf("http://www.example.com/%d", rand.Uint32()))
  if e != nil {
    panic(e)
  }
  return p
}

func newTestAsset() *Asset {
  a, e := NewAssetFromString(fmt.Sprintf("http://www.example.com/%d", rand.Uint32()))
  if e != nil {
    panic(e)
  }
  return a

}

func newTestNestedPages() *Page {
  p := newTestPage()
  for i := 0; i < 3; i++ {
    subpage := newTestPage()
    asset := newTestAsset()
    p.Links = append(p.Links, subpage)
    p.Assets = append(p.Assets, asset)
    for x := 0; x < 3; x++ {
      secondSubpage := newTestPage()
      secondAsset := newTestAsset()
      subpage.Links = append(subpage.Links, secondSubpage)
      subpage.Assets = append(subpage.Assets, secondAsset)
    }
  }
  return p
}
