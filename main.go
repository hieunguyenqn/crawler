package main

import (
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
  crawl("http://www.macasaurus.com", false)
  crawl("http://www.devbootcamp.com", false)
  crawl("http://www.digitalocean.com", false)
}

func printFirstLevel(page *crawler.Page) {
  fmt.Println("Printing First Level: ", page.URL)
  for _, p := range page.Links {
    fmt.Println("Subpage: ", p.URL)
    for _, a := range p.Assets {
      fmt.Println("    ", a.URL)
    }
  }
}

func crawl(url string, printResults bool) {
  start := time.Now()
  page, job := crawler.Crawl(url)
  if printResults {
    printFirstLevel(page)
  }
  stop := time.Now()
  duration := stop.Sub(start)
  fmt.Printf("Starting from %s, crawled %d pages in %s\n", page.URL, job.PagesCrawled, duration)
}
