package main

import (
  graphiz "code.google.com/p/gographviz"
  "code.google.com/p/gographviz/ast"
  "encoding/json"
  "fmt"
  "github.com/macb/crawler/crawler"
  "io/ioutil"
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
  crawl("http://www.devbootcamp.com", true)
  crawl("http://www.digitalocean.com", true)
}

func save(page *crawler.Page) {
  filename := "/tmp/" + page.URL.Host + ".json"
  by, e := json.Marshal(page.FlattenGraph())
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

  duration := stop.Sub(start)
  fmt.Printf("Starting from %s, crawled %d pages in %s\n", page.URL, job.PagesCrawled, duration)

  if saveResults {
    save(page)
    a := graph(page)
    filename := "/tmp/" + page.Path + ".dot"
    ioutil.WriteFile(filename, []byte(a.String()), 0755)
    fmt.Printf("Your .dot file is at: %s\n", filename)
  }
}

func graph(page *crawler.Page) *ast.Graph {
  var f func(page *crawler.Page)
  visited := make(map[*crawler.Page]bool)
  g := graphiz.NewGraph()
  graphName := `"` + page.Path + `"`
  g.SetName(graphName)
  g.SetDir(true)
  g.AddNode(graphName, quote(page.Path), nil)

  f = func(page *crawler.Page) {
    visited[page] = true
    g.AddNode(graphName, quote(page.Path), nil)
    for _, l := range page.Links {
      if !visited[l] {
        g.AddNode(graphName, quote(l.Path), nil)
        f(l)
      }
      g.AddEdge(quote(page.Path), "", quote(l.Path), "", true, nil)
    }
  }
  f(page)
  return g.WriteAst()
}

func quote(s string) string {
  return `"` + s + `"`
}
