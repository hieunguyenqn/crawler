package main

import (
  "fmt"
  "github.com/macb/crawler/crawler"
  "time"
)

func main() {
  crawl("http://www.macasaurus.com", false)
  crawl("http://www.devbootcamp.com", false)
  crawl("http://www.spirent.com", false)
  crawl("http://www.digitalocean.com", false)
  crawl("http://www.apple.com", false)
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
