package storage

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

var (
	// bigCacheClientSessionMapping singleton pattern
	bigCacheClientSessionMapping = make(map[string]*BigCacheClient)
)

// NewBigCache init new instance
func NewBigCache(config *bigcache.Config) INoSQLKeyValue {
	hasher := &hash.Client{}
	configAsJSON, err := json.Marshal(config)
	if err != nil {
		log.Fatalln("Unable to marshal BigCache configuration: ", err)
	}
	configAsString := hasher.SHA1(string(configAsJSON))

	currentBigCacheClientSession := bigCacheClientSessionMapping[configAsString]
	if currentBigCacheClientSession == nil {
		currentBigCacheClientSession = &BigCacheClient{nil}
		client, err := bigcache.NewBigCache(*config)
		if err != nil {
			log.Fatalln("Unable to connect to BigCache: ", err)
		} else {
			currentBigCacheClientSession.Client = client
			bigCacheClientSessionMapping[configAsString] = currentBigCacheClientSession
			log.Println("Connected to BigCache")
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
		log.Println("Unable to marshal value to []byte: ", err)
		return errors.New("Unable to marshal value")
	}

	return bc.Client.Set(key, b)
}

// Get return value based on the key provided
func (bc *BigCacheClient) Get(key string) (interface{}, error) {
	b, err := bc.Client.Get(key)
	if err != nil {
		log.Println("Unable to get value: ", err)
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
		log.Println("Unable to get value: ", err)
		return err
	}

	if err := bc.Client.Delete(key); err != nil {
		log.Println("Unable to delete value: ", err)
		return err
	}

	b, err := json.Marshal(value)
	if err != nil {
		log.Println("Unable to Marshal value: ", err)
		return err
	}

	return bc.Client.Set(key, b)
}

// Append new value base on the key provide, With Append() you can concatenate multiple entries under the same key in an lock optimized way.
func (bc *BigCacheClient) Append(key string, value interface{}) error {
	b, err := json.Marshal(value)
	if err != nil {
		log.Println("Unable to Marshal value: ", err)
		return errors.New("Unable to Marshal value")
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
