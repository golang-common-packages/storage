package caching

import (
	"time"

	"github.com/labstack/echo/v4"

	"github.com/golang-common-packages/database/model"
	"github.com/golang-common-packages/hash"
)

// ICaching interface for database caching factory pattern
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
	CUSTOM = iota
	BIGCACHE
	REDIS
)

// New function for database caching factory pattern
func New(databaseCompany int, config *model.Config) interface{} {
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
