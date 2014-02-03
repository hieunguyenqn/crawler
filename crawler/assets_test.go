package crawler

import (
  "net/url"
  "testing"
)

func Test_NewAssets(t *testing.T) {
  a := NewAssets()
  if a.safeMap.data == nil {
    t.Errorf("SafeMap was not initialized properly")
  }
}

func Test_NewAsset(t *testing.T) {
  rawUrl := "http://www.example.com"
  u, e := url.Parse(rawUrl)
  if e != nil {
    t.Errorf("Url parse error: %s", e.Error())
  }
  a := NewAsset(u)
  if a.Path != rawUrl {
    t.Errorf("%s did not match %s", a.Path, rawUrl)
  }
}

func Test_Assets_Value(t *testing.T) {
  assets := NewAssets()
  asset, e := NewAssetFromString("http://www.example.com")
  if e != nil {
    t.Error(e)
  }
  assets.data[asset.String()] = asset
  ret := assets.Value(asset.String())
  if ret != asset {
    t.Errorf("%v did not match %v", ret, asset)
  }
}
