package storage

import (
	"time"

	"github.com/labstack/echo/v4"

	"github.com/golang-common-packages/hash"
)

// ICaching caching factory pattern interface
type ICaching interface {
	Middleware(hash hash.IHash) echo.MiddlewareFunc
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, expire time.Duration) error
	Update(key string, value interface{}, expire time.Duration) error
	Delete(key string) error
	GetNumberOfRecords() int
	GetCapacity() (interface{}, error)
	Close() error
}

const (
	// CUSTOM caching on local memory
	CUSTOM = iota
	// BIGCACHE database
	BIGCACHE
	// REDIS database
	REDIS
)

// NewCaching factory pattern
func NewCaching(databaseCompany int, config *Config) interface{} {
	switch databaseCompany {
	case CUSTOM:
		return NewCustom(config)
	case REDIS:
		return NewRedis(config)
	case BIGCACHE:
		return NewBigCache(config)
	}

	return nil
}
