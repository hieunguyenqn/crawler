package crawler

import (
  "fmt"
  "github.com/PuerkitoBio/goquery"
  "net/url"
  "strings"
  "sync/atomic"
  "time"
)

func main() {
  Scrape("http://www.macasaurus.com")
  Scrape("http://www.digitalocean.com")
}

func Scrape(u string) *Page {
  page := new(Page)
  parsedUrl, _ := url.Parse(u)
  page.URL = parsedUrl
  j := newJob(page)
  j.Start()
  time.Sleep(1 * time.Second)
  func() {
    for {
      fmt.Println("Pages Scraped: ", atomic.LoadInt64(&j.PagesScraped))
      if j.Done() {
        j.Stop()
        return
      }
      time.Sleep(10 * time.Millisecond)
    }
  }()

  return page
}

func (w *webWorker) Crawl(p *Page) bool {
  doc, e := goquery.NewDocument(p.URL.String())
  if e != nil {
    // TODO Inspect error, don't blindly push.
    fmt.Println("Error: ", e)
    // Only retries 3 times.
    w.job.Requeue(p)
    return false
  }

  doc.Find("a").Each(func(i int, s *goquery.Selection) {
    u, ok := s.Attr("href")
    if !ok {
      return
    }
    if strings.Contains(u, "mailto") {
      return
    }

    parsedUrl, _ := p.Parse(u)
    if parsedUrl.Host != p.URL.Host {
      return
    }
    // Skip any subpages of different domains
    subpage, newPage := w.job.Pages.NewPage(parsedUrl)
    // Go gettem' tiger
    if newPage {
      w.job.ScrapeQueue.Push(subpage)
    }
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

      parsedUrl, _ := p.URL.Parse(u)
      asset := w.job.Assets.New(parsedUrl)
      p.Assets = append(p.Assets, asset)
    })
  }
  return true
}
