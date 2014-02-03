package main

import (
  "fmt"
  "github.com/macb/crawler/crawler"
  "time"
)

func main() {
  scrape("http://www.macasaurus.com", false)
  scrape("http://www.digitalocean.com", false)
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
  page := crawler.Scrape(url)
  if printResults {
    printFirstLevel(page)
  }
  stop := time.Now()
  duration := stop.Sub(start)
  fmt.Printf("Scrape took: %s\n", duration)
}
