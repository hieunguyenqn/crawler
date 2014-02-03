package crawler

import (
  "testing"
)

func Test_safeMap_NewsafeMap(t *testing.T) {
  s := newSafeMap()
  if s.data == nil {
    t.Errorf("Map should be initialized.")
  }
}
