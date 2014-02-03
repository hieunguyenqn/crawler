package main

import (
  "encoding/json"
  "fmt"
  "github.com/macb/crawler/crawler"
  "os"
  "runtime/pprof"
  "time"
)

var PPROF = os.Getenv("PPROF")

func main() {
  if PPROF != "" {
    f, _ := os.Create(PPROF)
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()

  }
  crawl("http://www.macasaurus.com", true)
  crawl("http://www.devbootcamp.com", false)
  crawl("http://www.digitalocean.com", false)
}

func save(page *crawler.Page) {
  filename := "/tmp/" + page.URL.Host + ".json"
  by, e := json.Marshal(page)
  if e != nil {
    panic(e)
  }
  f, e := os.Create(filename)
  if e != nil {
    panic(e)
  }
  defer f.Close()

  _, e = f.Write(by)
  if e != nil {
    panic(e)
  }
  fmt.Println("Your results are at: " + filename)
}

func crawl(url string, saveResults bool) {
  start := time.Now()
  page, job := crawler.Crawl(url)
  stop := time.Now()

  if saveResults {
    save(page)
  }
  duration := stop.Sub(start)
  fmt.Printf("Starting from %s, crawled %d pages in %s\n", page.URL, job.PagesCrawled, duration)
}
