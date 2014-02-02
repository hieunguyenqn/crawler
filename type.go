package main

import (
  "fmt"
  "net/url"
  "sync"
  "sync/atomic"
  "time"
)

///////////////////////////////
// Pages
///////////////////////////////
type Pages struct {
  safeMap
}

func (p *Pages) Add(key string, page *Page) {
  p.safeMap.Add(key, page)
}

func (p Pages) Value(key *url.URL) *Page {
  ret := p.safeMap.Value(key.String())
  if ret == nil {
    return nil
  }
  return ret.(*Page)
}

func (p *Pages) NewPage(url *url.URL) (*Page, bool) {
  if page := p.Value(url); page != nil {
    return page, false
  }
  p.lock.Lock()
  defer p.lock.Unlock()
  page := new(Page)
  page.URL = &URL{url}
  p.safeMap.data[url.String()] = page
  return page, true
}

///////////////////////////////
// Assets
///////////////////////////////
type Assets struct {
  safeMap
}

func (a *Assets) Add(key string, asset *Asset) {
  a.safeMap.Add(key, asset)
}

func (a Assets) Value(key string) (asset *Asset) {
  var ok bool
  if asset, ok = a.safeMap.Value(key).(*Asset); ok {
    return asset
  }
  return nil
}

func (a *Assets) New(url *url.URL) *Asset {
  if asset := a.Value(url.String()); asset != nil {
    return asset
  }
  a.lock.Lock()
  defer a.lock.Unlock()

  asset := new(Asset)
  asset.URL = &URL{url}
  a.safeMap.data[url.String()] = asset
  return asset
}

///////////////////////////////
// Page
///////////////////////////////
type Page struct {
  *URL
  Links  []*Page
  Assets []*Asset
}

///////////////////////////////
// Asset
///////////////////////////////
type Asset struct {
  *URL
}

///////////////////////////////
// URL
///////////////////////////////
type URL struct {
  URL *url.URL
}

func (u URL) ParseRelative(path string) (*url.URL, error) {
  switch {
  case path == "":
    return u.URL, nil
  case path[0] == '#':
    return nil, fmt.Errorf("ID Path: %s", path)
  case len(path) > 2 && path[0:2] == "//":
    return url.Parse(u.Scheme() + ":" + path)
  case path == "/":
    return nil, fmt.Errorf("Root Path: %s", path)
  default:
    return url.Parse(u.subpage(path))
  }
}

func (u URL) subpage(path string) string {
  return fmt.Sprintf("%s://%s%s", u.URL.Scheme, u.URL.Host, path)
}

func (u URL) Host() string {
  return u.URL.Host
}

func (u URL) Scheme() string {
  return u.URL.Scheme
}

func (u URL) String() string {
  return u.URL.String()
}

/////////////////////////////
// safeBool
/////////////////////////////

type safeBool struct {
  b    bool
  lock sync.RWMutex
}

func (s *safeBool) True() {
  s.lock.Lock()
  defer s.lock.Unlock()
  s.b = true
}

func (s *safeBool) False() {
  s.lock.Lock()
  defer s.lock.Unlock()
  s.b = false
}

func (s safeBool) Value() bool {
  s.lock.RLock()
  defer s.lock.RUnlock()
  return s.b
}

/////////////////////////////
// safeMap
/////////////////////////////

type safeMap struct {
  data map[string]interface{}
  lock sync.RWMutex
}

func (s *safeMap) Add(key string, value interface{}) {
  s.lock.Lock()
  defer s.lock.Unlock()
  s.data[key] = value
}

func (s *safeMap) Value(key string) interface{} {
  s.lock.RLock()
  defer s.lock.RUnlock()
  return s.data[key]
}

/////////////////////////////
// pageStack
/////////////////////////////

type pageStack struct {
  data []*Page
  lock sync.Mutex
}

func (s *pageStack) Len() int {
  s.lock.Lock()
  defer s.lock.Unlock()
  return len(s.data)
}

func (s *pageStack) Push(page *Page) {
  s.lock.Lock()
  defer s.lock.Unlock()
  s.data = append(s.data, page)
}

func (s *pageStack) Pop() (page *Page) {
  s.lock.Lock()
  defer s.lock.Unlock()
  switch len(s.data) {
  case 0:
    return nil
  case 1:
    page = s.data[0]
    s.data = make([]*Page, 0)
  default:
    page, s.data = s.data[0], s.data[1:len(s.data)]
  }
  return page
}

/////////////////////////////
// job
/////////////////////////////

type job struct {
  ScrapeQueue  pageStack
  WebWorkers   []*webWorker
  Pages        Pages
  Assets       Assets
  PagesScraped int64
}

func newJob(page *Page) *job {
  j := new(job)
  j.Pages.safeMap.data = make(map[string]interface{})
  j.Assets.safeMap.data = make(map[string]interface{})
  for i := 0; i < MAX_WEB_WORKERS; i++ {
    w := new(webWorker)
    w.job = j
    w.stop = make(chan int)
    j.WebWorkers = append(j.WebWorkers, w)
  }
  j.ScrapeQueue.Push(page)
  return j
}

func (j *job) Start() {
  j.startWorkers()
}

func (j *job) Stop() {
  j.stopWorkers()
}

func (j *job) Done() bool {
  if j.ScrapeQueue.Len() == 0 && j.WorkersDone() {
    return true
  }
  return false
}

func (j *job) WorkersDone() bool {
  for _, w := range j.WebWorkers {
    if w.busy.Value() {
      return false
    }
  }
  return true
}

func (j *job) startWorkers() {
  for _, w := range j.WebWorkers {
    go w.Scrape()
  }
}

func (j *job) stopWorkers() {
  for _, w := range j.WebWorkers {
    w.stop <- 1
  }
}

/////////////////////////////
// webWorker
/////////////////////////////

type webWorker struct {
  busy safeBool
  job  *job
  stop chan int
}

func (w *webWorker) Scrape() {
  for {
    select {
    case <-w.stop:
      return
    default:
      if page := w.job.ScrapeQueue.Pop(); page != nil {
        w.busy.True()
        success := w.Crawl(page)
        if success {
          atomic.AddInt64(&w.job.PagesScraped, 1)
        }
        w.busy.False()
      }
    }
    time.Sleep(1 * time.Millisecond)
  }
}
