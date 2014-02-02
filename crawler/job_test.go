package crawler

import (
  "testing"
)

func Test_Job_Requeue(t *testing.T) {
  j := new(job)
  j.Retries = make(map[*Page]int)
  page := newTestPage()
  for i := 0; i < 5; i++ {
    j.Requeue(page)
  }
  if j.ScrapeQueue.Len() > 3 {
    t.Errorf("Job should have only been requeued 3 times.")
  }
  if j.ScrapeQueue.Len() < 3 {
    t.Errorf("Job should have been requeued 3 times. Got: %d", j.ScrapeQueue.Len())
  }
}

func Test_Job_WorkersDone(t *testing.T) {
  j := new(job)
  if !j.WorkersDone() {
    t.Errorf("Expected workers to be done.")
  }
}
