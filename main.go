package main

import (
  "fmt"
  "github.com/PuerkitoBio/goquery"
  "log"
  "net/url"
  "strings"
  "sync/atomic"
  "time"
)

const MAX_WEB_WORKERS int64 = 5

var WEB_WORKERS int64 = 0

var pages = new(Pages)
var assets = new(Assets)

func init() {
  pages.data = make(map[string]interface{})
  assets.data = make(map[string]interface{})
}

func main() {
  Start("http://www.macasaurus.com")
  Start("http://www.digitalocean.com")
}

func Start(u string) {
  done := make(chan bool, 1)
  page := new(Page)
  parsedUrl, _ := url.Parse(u)
  page.URL = &URL{parsedUrl}

  go page.Crawl(done)
  <-done

  for _, p := range page.Links {
    fmt.Printf("Page: %s\n", p.URL)
    for _, a := range p.Assets {
      fmt.Printf("    Asset: %s\n", a.URL)
    }
  }
}

func (p *Page) Crawl(done chan bool) {
  for atomic.LoadInt64(&WEB_WORKERS) >= MAX_WEB_WORKERS {
    time.Sleep(5 * time.Millisecond)
  }
  if p.Visited.Value() {
    finishCrawl(done, fmt.Sprintf("Already Visted: %s", p.URL))
    return
  }

  atomic.AddInt64(&WEB_WORKERS, 1)
  p.Visited.Visit()
  fmt.Printf("Starting Crawl: %s\n", p.URL)
  subpageDone := make(chan bool, 1)
  var subpageCount int64

  doc, e := goquery.NewDocument(p.URL.String())
  if e != nil {
    finishCrawl(done, fmt.Sprintf("Error: %s", e.Error()))
    p.Visited.Unvisit()
    p.Crawl(done)
    return
  }

  doc.Find("a").Each(func(i int, s *goquery.Selection) {
    u, ok := s.Attr("href")
    if !ok {
      return
    }

    if strings.Contains(u, "mailto") {
      return
    }

    parsedUrl, _ := url.Parse(u)
    if parsedUrl.Scheme == "" {
      var err error
      if parsedUrl, err = p.ParseRelative(u); err != nil {
        return
      }
    }

    // Skip any subpages of different domains
    if parsedUrl.Host != p.URL.Host() {
      return
    }

    subpage := pages.NewPage(parsedUrl)
    atomic.AddInt64(&subpageCount, 1)

    // Go gettem' tiger
    go subpage.Crawl(subpageDone)

    p.Links = append(p.Links, subpage)
  })

  assetTags := map[string]string{
    "script": "src",
    "link":   "href",
    "img":    "src",
  }

  for assetTag, urlTag := range assetTags {
    doc.Find(assetTag).Each(func(i int, s *goquery.Selection) {
      u, ok := s.Attr(urlTag)
      if !ok {
        return
      }

      parsedUrl, _ := url.Parse(u)
      if parsedUrl.Scheme == "" {
        parsedUrl, _ = p.ParseRelative(u)
      }
      asset := assets.New(parsedUrl)
      p.Assets = append(p.Assets, asset)
    })
  }

  atomic.AddInt64(&WEB_WORKERS, -1)

  var i int64
  for i = 0; i < subpageCount; i++ {
    <-subpageDone
  }
  fmt.Printf("Finished Crawl: %s\n", p.URL.String())
  done <- true
}

func finishCrawl(done chan bool, message string) {
  if message != "" {
    log.Println(message)
  }
  done <- true
  atomic.AddInt64(&WEB_WORKERS, -1)
}
