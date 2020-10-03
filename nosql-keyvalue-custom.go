package storage

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/golang-common-packages/hash"
	"github.com/golang-common-packages/linear"
)

// KeyValueCustomClient manage all custom caching actions
type KeyValueCustomClient struct {
	client *linear.Client
	close  chan struct{}
}

var (
	// keyValueCustomClientSessionMapping singleton pattern
	keyValueCustomClientSessionMapping = make(map[string]*KeyValueCustomClient)
)

// NewKeyValueCustom init new instance
func NewKeyValueCustom(config *CustomKeyValue) INoSQLKeyValue {
	hasher := &hash.Client{}
	configAsJSON, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	configAsString := hasher.SHA1(string(configAsJSON))

	currentCustomClientSession := keyValueCustomClientSessionMapping[configAsString]
	if currentCustomClientSession == nil {
		currentCustomClientSession = &KeyValueCustomClient{linear.New(config.MemorySize, config.CleaningEnable), make(chan struct{})}
		keyValueCustomClientSessionMapping[configAsString] = currentCustomClientSession
		log.Println("Custom caching is ready")

		// Check record expiration time and remove
		go func() {
			ticker := time.NewTicker(config.CleaningInterval)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					items := currentCustomClientSession.client.GetItems()
					items.Range(func(key, value interface{}) bool {
						item := value.(customKeyValueItem)

						if item.expires < time.Now().UnixNano() {
							k, _ := key.(string)
							currentCustomClientSession.client.Get(k)
						}

						return true
					})

				case <-currentCustomClientSession.close:
					return
				}
			}
		}()
	}

	return currentCustomClientSession
}

// Middleware for echo framework
func (cl *KeyValueCustomClient) Middleware(hash hash.IHash) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get(echo.HeaderAuthorization)
			key := hash.SHA512(token)

			if val, err := cl.Get(key); err != nil {
				log.Printf("Can not get accesstoken from custom caching in echo middleware: %s", err.Error())
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			} else if val == "" {
				return c.NoContent(http.StatusUnauthorized)
			}

			return next(c)
		}
	}
}

// Get return value based on the key provided
func (cl *KeyValueCustomClient) Get(key string) (interface{}, error) {
	if key == "" {
		return nil, errors.New("key must not empty")
	}

	obj, err := cl.client.Read(key)
	if err != nil {
		return nil, err
	}

	item, ok := obj.(customKeyValueItem)
	if !ok {
		return nil, errors.New("can not map object to customKeyValueItem model")
	}

	if item.expires < time.Now().UnixNano() {
		return nil, nil
	}

	return item.data, nil
}

// GetMany return value based on the list of keys provided
func (cl *KeyValueCustomClient) GetMany(keys []string) (map[string]interface{}, []string, error) {
	if len(keys) == 0 {
		return nil, nil, errors.New("keys must not empty")
	}

	var itemFound map[string]interface{}
	var itemNotFound []string

	for _, key := range keys {
		obj, err := cl.client.Read(key)
		if obj == nil && err == nil {
			itemNotFound = append(itemNotFound, key)
		}

		item, ok := obj.(customKeyValueItem)
		if !ok {
			return nil, nil, errors.New("can not map object to customKeyValueItem model")
		}

		itemFound[key] = item.data
	}

	return itemFound, itemNotFound, nil
}

// Set new record set key and value
func (cl *KeyValueCustomClient) Set(key string, value interface{}, expire time.Duration) error {
	if key == "" || value == nil {
		return errors.New("key and value must not empty")
	}

	if expire == 0 {
		expire = 24 * time.Hour
	}

	if err := cl.client.Push(key, customKeyValueItem{
		data:    value,
		expires: time.Now().Add(expire).UnixNano(),
	}); err != nil {
		return err
	}

	return nil
}

// Update new value over the key provided
func (cl *KeyValueCustomClient) Update(key string, value interface{}, expire time.Duration) error {
	if key == "" || value == nil {
		return errors.New("key and value must not empty")
	}

	_, err := cl.client.Get(key)
	if err != nil {
		return err
	}

	if expire == 0 {
		expire = 24 * time.Hour
	}

	if err := cl.client.Push(key, customKeyValueItem{
		data:    value,
		expires: time.Now().Add(expire).UnixNano(),
	}); err != nil {
		return err
	}

	return nil
}

// Delete deletes the key and its value from the memory
func (cl *KeyValueCustomClient) Delete(key string) error {
	if key == "" {
		return errors.New("key must not empty")
	}

	if _, err := cl.client.Get(key); err != nil {
		return err
	}

	return nil
}

// Range over linear data structure
func (cl *KeyValueCustomClient) Range(f func(key, value interface{}) bool) {
	fn := func(key, value interface{}) bool {
		item := value.(customKeyValueItem)

		if item.expires > 0 && item.expires < time.Now().UnixNano() {
			return true
		}

		return f(key, item.data)
	}

	cl.client.Range(fn)
}

// GetNumberOfRecords return number of records
func (cl *KeyValueCustomClient) GetNumberOfRecords() int {
	return cl.client.GetNumberOfKeys()
}

// GetCapacity method return redis database size
func (cl *KeyValueCustomClient) GetCapacity() (interface{}, error) {
	return cl.client.GetLinearCurrentSize(), nil
}

// Close the service and frees up resources.
func (cl *KeyValueCustomClient) Close() error {
	cl.close <- struct{}{}

	return nil
}
