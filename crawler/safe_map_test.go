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

func Test_safeMap_Add(t *testing.T) {
  type Temp struct {
    Name string
  }
  temp := &Temp{Name: "test"}
  s := newSafeMap()

  s.Add(temp.Name, temp)
  ret := s.data[temp.Name]
  if ret != temp {
    t.Errorf("%v did not match %v", ret, temp)
  }
}
