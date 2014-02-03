package crawler

import (
  "sync"
  "time"
)

const MAX_WEB_WORKERS int = 30

/////////////////////////////
// Job
/////////////////////////////

type Job struct {
  StartPage    *Page
  ScrapeQueue  pageStack
  WebWorkers   []*webWorker
  retryLock    sync.Mutex
  Retries      map[*Page]int
  Pages        *Pages
  Assets       *Assets
  PagesScraped int64
}

func NewJob(page *Page) *Job {
  j := new(Job)
  j.StartPage = page
  j.Pages = NewPages()
  j.Assets = NewAssets()
  j.Retries = make(map[*Page]int)
  for i := 0; i < MAX_WEB_WORKERS; i++ {
    w := newWebWorker(i, j)
    j.WebWorkers = append(j.WebWorkers, w)
  }
  j.ScrapeQueue.Push(page)
  return j
}

func (j *Job) Start() {
  j.startWorkers()
  ticker := time.NewTicker(50 * time.Millisecond)
  func() {
    for {
      select {
      case <-ticker.C:
        if j.Done() {
          j.Stop()
          return
        }
      }
    }
  }()
}

func (j *Job) Stop() {
  j.stopWorkers()
}

func (j *Job) Done() bool {
  if j.ScrapeQueue.Len() == 0 && j.WorkersDone() {
    return true
  }
  return false
}

func (j *Job) WorkersDone() bool {
  for _, w := range j.WebWorkers {
    if w.busy.Value() {
      return false
    }
  }
  return true
}

func (j *Job) Requeue(page *Page) {
  j.retryLock.Lock()
  defer j.retryLock.Unlock()
  val := j.Retries[page]
  if val > 2 {
    return
  }
  j.Retries[page] = val + 1
  j.ScrapeQueue.Push(page)
}

func (j *Job) startWorkers() {
  for _, w := range j.WebWorkers {
    go w.Scrape()
  }
}

func (j *Job) stopWorkers() {
  for _, w := range j.WebWorkers {
    w.stop <- 1
  }
}
