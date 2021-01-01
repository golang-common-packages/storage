package storage

import (
	"time"

	"github.com/labstack/echo/v4"

	"github.com/golang-common-packages/hash"
)

// INoSQLKeyValue factory pattern interface
type INoSQLKeyValue interface {
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

// newNoSQLKeyValue factory pattern
func newNoSQLKeyValue(databaseCompany int, config *Config) interface{} {

	switch databaseCompany {
	case CUSTOM:
		return newKeyValueCustom(&config.CustomKeyValue)
	case REDIS:
		return newRedis(&config.Redis)
	case BIGCACHE:
		return newBigCache(&config.BigCache)
	}

	return nil
}
