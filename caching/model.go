package caching

// customCacheItem private model for custom cache record
type customCacheItem struct {
	data    interface{}
	expires int64
}
