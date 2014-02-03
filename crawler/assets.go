package crawler

import (
  "net/url"
)

///////////////////////////////
// Assets
///////////////////////////////
type Assets struct {
  *safeMap
}

func NewAssets() *Assets {
  assets := new(Assets)
  assets.safeMap = newSafeMap()
  return assets
}

func (a Assets) Value(key string) (asset *Asset) {
  var ok bool
  if asset, ok = a.safeMap.Value(key).(*Asset); ok {
    return asset
  }
  return nil
}

func (a *Assets) New(u *url.URL) *Asset {
  asset := NewAsset(u)
  ret, newEntry := a.Add(asset.String(), asset)
  if !newEntry {
    asset = ret.(*Asset)
  }
  return asset
}

///////////////////////////////
// Asset
///////////////////////////////
type Asset struct {
  Path     string
  *url.URL `json:"-"` // Don't want to encode the URL.
}

func NewAsset(u *url.URL) *Asset {
  return &Asset{URL: u, Path: u.String()}
}
