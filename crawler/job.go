package crawler

import (
  "sync"
)

const MAX_WEB_WORKERS int = 10

/////////////////////////////
// Job
/////////////////////////////

type Job struct {
  StartPage    *Page
  Queue        chan *Page
  WebWorkers   []*webWorker
  retryLock    sync.Mutex
  Retries      map[*Page]int
  Pages        *Pages
  Assets       *Assets
  PagesCrawled int64
  done         chan int
}

func NewJob(page *Page) *Job {
  j := new(Job)
  j.StartPage = page

  // Initialize Pages and Assets store.
  j.Pages = NewPages()
  j.Assets = NewAssets()

  // Retry count map. Must be available to all goroutines.
  j.Retries = make(map[*Page]int)

  // Initialize to be scraped Queue, worker done channel, and tr
  j.Queue = make(chan *Page, 100)
  j.done = make(chan int, MAX_WEB_WORKERS)

  // Seed worker queue with first page.
  j.Queue <- page

  // Initialize MAX_WEB_WORKERS number of workers.
  for i := 0; i < MAX_WEB_WORKERS; i++ {
    w := newWebWorker(i, j)
    j.WebWorkers = append(j.WebWorkers, w)
  }

  return j
}

func (j *Job) Start() {
  // Start workers and loop forever checking if done.
  j.startWorkers()
  for {
    select {
    // When a worker reports they finished a job, check to see if job is done.
    case <-j.done:
      if j.Done() {
        j.Stop()
        return
      }
    }
  }
}

func (j *Job) Stop() {
  j.stopWorkers()
}

func (j *Job) Done() bool {
  // TODO Not entirely sure this is thread safe since we're using a buffered channel.
  if len(j.Queue) == 0 && j.WorkersDone() {
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

  // TODO Better solution.
  // This will use 4KB of memory for every request that can't be put into the
  // queue immediately. Probably messes with garbage collection as well since
  // page can't be cleaned up.

  // Send in a goroutine to avoid any chance of being blocked.
  go func() { j.Queue <- page }()
}

func (j *Job) startWorkers() {
  for _, w := range j.WebWorkers {
    go w.Work()
  }
}

func (j *Job) stopWorkers() {
  for _, w := range j.WebWorkers {
    w.stop <- 1
  }
}
