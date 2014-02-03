package crawler

import (
  "testing"
  "time"
)

func Test_Job_Requeue(t *testing.T) {
  page := newTestPage()
  j := NewJob(page)

  for i := 0; i < 5; i++ {
    j.Requeue(page)
  }

  timer := time.NewTimer(1 * time.Second)
  received := []*Page{}
  func() {
    for {
      select {
      case p := <-j.Queue:
        received = append(received, p)
        if len(received) == 3 {
          return
        }
      case <-timer.C:
        t.Errorf("Received timed out.")
        return
      }
    }
  }()

  if len(received) > 3 {
    t.Errorf("Job should have only been requeued 3 times.")
  }

  if len(received) < 3 {
    t.Errorf("Job should have been requeued 3 times. Got: %d", len(received))
  }
}

func Test_Job_WorkersDone(t *testing.T) {
  j := new(Job)
  if !j.WorkersDone() {
    t.Errorf("Expected workers to be done.")
  }
}
