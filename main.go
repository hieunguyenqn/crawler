package main

import (
  "fmt"
  "github.com/macb/crawler/crawler"
)

func main() {
  crawler.Scrape("http://www.macasaurus.com")
  crawler.Scrape("http://www.digitalocean.com")
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
