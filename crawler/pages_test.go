package crawler

import (
  "encoding/json"
  "testing"
)

func Test_Page_Flatten(t *testing.T) {
  p := newTestNestedPages()
  flat := p.Flatten()
  if p.Links[0].Path != flat["Links"].([]string)[0] {
    t.Errorf("Link paths did not match.")
  }
}

func Test_Page_FlattenGraph(t *testing.T) {
  p := newTestNestedPages()
  flat := p.FlattenGraph()
  if _, e := json.Marshal(flat); e != nil {
    t.Errorf("It should be JSON marshal-able.")
  }
}
