package crawler

import (
  "sync/atomic"
  "time"
)

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
      if w.job.ScrapeQueue.Len() > 0 {
        w.busy.True()
        if page := w.job.ScrapeQueue.Pop(); page != nil {
          success := w.Crawl(page)
          if success {
            atomic.AddInt64(&w.job.PagesScraped, 1)
          }
        }
        w.busy.False()
      }
    }
    time.Sleep(5 * time.Millisecond)
  }
}
