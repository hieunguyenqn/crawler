package main

import (
  "fmt"
  "github.com/macb/crawler/crawler"
  "time"
)

func main() {
  scrape("http://www.macasaurus.com", false)
  scrape("http://www.devbootcamp.com", false)
  scrape("http://www.spirent.com", false)
  scrape("http://www.digitalocean.com", false)
  scrape("http://www.apple.com", false)
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

func scrape(url string, printResults bool) {
  start := time.Now()
  page, job := crawler.Scrape(url)
  if printResults {
    printFirstLevel(page)
  }
  stop := time.Now()
  duration := stop.Sub(start)
  fmt.Printf("Starting from %s, scraped %d pages in %s\n", page.URL, job.PagesScraped, duration)
}
