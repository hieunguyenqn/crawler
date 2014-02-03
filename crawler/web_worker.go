package crawler

import (
  "sync/atomic"
  "time"
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
  ticker := time.NewTicker(50 * time.Millisecond)
  for {
    select {
    case <-w.stop:
      return
    case <-ticker.C:
      if w.job.Queue.Len() > 0 {
        w.busy.True()
        if page := w.job.Queue.Pop(); page != nil {
          success := w.Crawl(page)
          if success {
            atomic.AddInt64(&w.job.PagesCrawled, 1)
          }
        }
        w.busy.False()
      }
    }
  }
}
