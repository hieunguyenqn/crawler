package crawler_test

import (
  "github.com/macb/crawler/crawler"
  "net/http"
  "net/http/httptest"
  "testing"
)

func Test_Scrape(t *testing.T) {
  ts := setupSinglePage()
  defer ts.Close()
  page, _ := crawler.Crawl(ts.URL)
  if len(page.Links) > 0 {
    t.Errorf("Page should not have links.")
  }
}

func Test_Scrape5(t *testing.T) {
  ts := setupFiveSubpages()
  defer ts.Close()

  page, _ := crawler.Crawl(ts.URL)
  if len(page.Links) != 5 {
    t.Errorf("Page should have 5 subpages.")
  }
  for _, p := range page.Links {
    if len(p.Links) > 0 {
      t.Errorf("Subpage should have no links.")
    }
  }
}

func Benchmark_Scrape5(b *testing.B) {

  b.StopTimer()
  ts := setupFiveSubpages()
  defer ts.Close()
  b.StartTimer()

  for i := 0; i < b.N; i++ {
    crawler.Crawl(ts.URL)
  }
}

func Test_ScrapeCircular(t *testing.T) {
  ts := setupCircularSubpages()
  defer ts.Close()

  page, _ := crawler.Crawl(ts.URL)
  if len(page.Links) != 1 {
    t.Errorf("Page should have 1 subpage.")
  }
  for _, p := range page.Links {
    if len(p.Links) != 1 {
      t.Errorf("Subpage should have 1 link.")
    }
  }
}

func setupSinglePage() *httptest.Server {
  return httptest.NewServer(http.FileServer(http.Dir("single_page/")))
}

func setupFiveSubpages() *httptest.Server {
  return httptest.NewServer(http.FileServer(http.Dir("five_subpages/")))
}

func setupCircularSubpages() *httptest.Server {
  return httptest.NewServer(http.FileServer(http.Dir("circular_subpages/")))
}
