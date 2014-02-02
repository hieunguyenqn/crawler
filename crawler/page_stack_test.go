package crawler

import (
  "testing"
)

func Test_pageStack_Push(t *testing.T) {
  s := newTestpageStack()
  page := newTestPage()
  s.Push(page)
  if len(s.data) != 1 {
    t.Errorf("Expected 1 Page, got %d", len(s.data))
  }
}

func Test_pageStack_Pop(t *testing.T) {
  s := newTestpageStack()
  page := newTestPage()
  s.Push(page)
  poppedPage := s.Pop()
  if page != poppedPage {
    t.Errorf("Popped (%v) != Pushed (%v)", poppedPage, page)
  }
}

func Test_pageStack_Len(t *testing.T) {
  s := newTestpageStack()
  page := newTestPage()
  s.Push(page)
  if page != poppedPage {
    t.Errorf("Popped (%v) != Pushed (%v)", poppedPage, page)
  }
}

func newTestpageStack() *pageStack {
  return new(pageStack)
}
