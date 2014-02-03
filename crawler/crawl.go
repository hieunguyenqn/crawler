package crawler

import (
  "github.com/PuerkitoBio/goquery"
  "strings"
)

func Crawl(u string) (*Page, *Job) {
  page, err := NewPageFromString(u)
  if err != nil {
    panic(err)
  }
  j := NewJob(page)
  j.Start()
  return page, j
}

func (w *webWorker) Crawl(p *Page) bool {
  doc, e := goquery.NewDocument(p.URL.String())
  if e != nil {
    // Only retries 3 times.
    w.job.Requeue(p)
    return false
  }

  doc.Find("a").Each(func(i int, s *goquery.Selection) {
    u, ok := s.Attr("href")
    if !ok {
      return
    }
    // Skip mails and on-page navigation
    if strings.Contains(u, "mailto") || len(u) > 0 && u[0] == '#' {
      return
    }

    parsedUrl, _ := p.Parse(u)
    if parsedUrl.Host != p.URL.Host {
      return
    }

    subpage, newPage := w.job.Pages.NewPage(parsedUrl)

    // Scrape new pages that share the same host.
    if newPage && subpage.Host == p.Host {

      // TODO Better solution.
      // This will use 4KB of memory for every request that can't be put into the
      // queue immediately. Probably messes with garbage collection as well since
      // page can't be cleaned up. Something like an infinite queue.

      // Send in a goroutine to avoid any chance of being blocked.
      go func() { w.job.Queue <- subpage }()
    }

    // Add all subpages as links, even if they weren't crawled.
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
