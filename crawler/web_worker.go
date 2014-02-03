package crawler

import (
  "sync/atomic"
)

/////////////////////////////
// webWorker
/////////////////////////////

type webWorker struct {
  id   int
  busy safeBool
  job  *Job
  stop chan int
}

func newWebWorker(id int, j *Job) *webWorker {
  w := new(webWorker)
  w.id = id
  w.job = j
  w.stop = make(chan int)
  return w
}

func (w *webWorker) Work() {
  for {
    select {
    case <-w.stop:
      return
    case page := <-w.job.Queue:
      w.busy.True()
      success := w.Crawl(page)
      if success {
        atomic.AddInt64(&w.job.PagesCrawled, 1)
      }
      w.busy.False()
      go func() { w.job.done <- 1 }()
    }
  }
}
