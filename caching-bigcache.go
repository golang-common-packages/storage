package database

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/allegro/bigcache/v2"
	"github.com/labstack/echo/v4"

	"github.com/golang-common-packages/hash"
)

// BigCacheClient manage all BigCache actions
type BigCacheClient struct {
	Client *bigcache.BigCache
}

// bigCacheClientSessionMapping singleton pattern
var bigCacheClientSessionMapping map[string]*BigCacheClient

// NewBigCache init new instance
func NewBigCache(config *Config) ICaching {
	hasher := &hash.Client{}
	configAsJSON, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	configAsString := hasher.SHA1(string(configAsJSON))

	currentBigCacheClientSession := bigCacheClientSessionMapping[configAsString]
	if currentBigCacheClientSession == nil {
		currentBigCacheClientSession = &BigCacheClient{nil}
		client, err := bigcache.NewBigCache(config.BigCache)
		if err != nil {
			panic(err)
		} else {
			currentBigCacheClientSession.Client = client
			bigCacheClientSessionMapping[configAsString] = currentBigCacheClientSession
			log.Println("Connected to BigCache Server")
		}
	}

	return currentBigCacheClientSession
}

// Middleware for echo framework
func (bc *BigCacheClient) Middleware(hash hash.IHash) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get(echo.HeaderAuthorization)
			key := hash.SHA512(token)

			val, err := bc.Get(key)
			if err != nil {
				if err.Error() == "Entry not found" {
					return c.NoContent(http.StatusUnauthorized)
				}

				log.Printf("Can not get accesstoken from BigCache in echo middleware: %s", err.Error())
				return c.NoContent(http.StatusInternalServerError)
			} else if val == "" {
				return c.NoContent(http.StatusUnauthorized)
			}

			return next(c)
		}
	}
}

// Set new record set key and value
func (bc *BigCacheClient) Set(key string, value interface{}, expire time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return errors.New("can not marshal value to []byte")
	}

	return bc.Client.Set(key, b)
}

// Get return value based on the key provided
func (bc *BigCacheClient) Get(key string) (interface{}, error) {
	b, err := bc.Client.Get(key)
	if err != nil {
		return nil, err
	}

	var value interface{}
	json.Unmarshal(b, value)

	return value, nil
}

// Update new value over the key provided
func (bc *BigCacheClient) Update(key string, value interface{}, expire time.Duration) error {
	_, err := bc.Client.Get(key)
	if err != nil {
		return err
	}

	if err := bc.Client.Delete(key); err != nil {
		return nil
	}

	if b, err := json.Marshal(value); err != nil {
		return bc.Client.Set(key, b)
	}

	return nil
}

// Append new value base on the key provide, With Append() you can concatenate multiple entries under the same key in an lock optimized way.
func (bc *BigCacheClient) Append(key string, value interface{}) error {
	b, err := json.Marshal(value)
	if err != nil {
		return errors.New("can not marshal value to []byte")
	}

	return bc.Client.Append(key, b)
}

// Delete function will delete value based on the key provided
func (bc *BigCacheClient) Delete(key string) error {
	return bc.Client.Delete(key)
}

// GetNumberOfRecords return number of records
func (bc *BigCacheClient) GetNumberOfRecords() int {
	return bc.Client.Len()
}

// GetCapacity method return redis database size
func (bc *BigCacheClient) GetCapacity() (interface{}, error) {
	return bc.Client.Capacity(), nil
}

// Close function will close BigCache connection
func (bc *BigCacheClient) Close() error {
	return bc.Client.Close()
}
