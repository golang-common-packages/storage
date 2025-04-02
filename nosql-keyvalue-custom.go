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
	client *linear.Linear
	close  chan struct{}
}

var (
	// keyValueCustomClientSessionMapping singleton pattern
	keyValueCustomClientSessionMapping = make(map[string]*KeyValueCustomClient)
)

// newKeyValueCustom init new instance
func newKeyValueCustom(config *CustomKeyValue) INoSQLKeyValue {
	hasher := &hash.Client{}
	configAsJSON, err := json.Marshal(config)
	if err != nil {
		log.Fatalln("Unable to marshal service configuration: ", err)
	}
	configAsString := hasher.SHA1(string(configAsJSON))

	currentCustomClientSession := keyValueCustomClientSessionMapping[configAsString]
	if currentCustomClientSession == nil {
		currentCustomClientSession = &KeyValueCustomClient{linear.New(config.MemorySize, config.CleaningEnable), make(chan struct{})}
		keyValueCustomClientSessionMapping[configAsString] = currentCustomClientSession
		log.Println("Key-value custom is ready")

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
				log.Println("Unable to get value: ", err)
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			} else if val == "" {
				return c.NoContent(http.StatusUnauthorized)
			}

			return next(c)
		}
	}
}

// Get retrieves a value from the cache based on the key provided
// Returns nil, nil if the key exists but the value is expired
func (cl *KeyValueCustomClient) Get(key string) (interface{}, error) {
	if cl.client == nil {
		return nil, errors.New("client is not initialized")
	}
	
	if key == "" {
		return nil, errors.New("key cannot be empty")
	}

	obj, err := cl.client.Read(key)
	if err != nil {
		return nil, err
	}
	
	if obj == nil {
		return nil, errors.New("key not found: " + key)
	}

	item, ok := obj.(customKeyValueItem)
	if !ok {
		log.Printf("Warning: Unable to map object to customKeyValueItem model for key: %s", key)
		return nil, errors.New("invalid cache item format")
	}

	// Check if item is expired
	if item.expires < time.Now().UnixNano() {
		// Automatically remove expired items
		cl.client.Get(key)
		return nil, nil
	}

	return item.data, nil
}

// GetMany returns values based on the list of keys provided
// Returns a map of found items, a slice of keys not found, and any error encountered
func (cl *KeyValueCustomClient) GetMany(keys []string) (map[string]interface{}, []string, error) {
	if len(keys) == 0 {
		return nil, nil, errors.New("keys cannot be empty")
	}

	if cl.client == nil {
		return nil, nil, errors.New("client is not initialized")
	}

	itemFound := make(map[string]interface{})
	var itemNotFound []string

	for _, key := range keys {
		if key == "" {
			continue // Skip empty keys
		}
		
		obj, err := cl.client.Read(key)
		if obj == nil || err != nil {
			itemNotFound = append(itemNotFound, key)
			continue
		}

		item, ok := obj.(customKeyValueItem)
		if !ok {
			log.Printf("Warning: Unable to map object to customKeyValueItem model for key: %s", key)
			itemNotFound = append(itemNotFound, key)
			continue
		}

		// Check if item is expired
		if item.expires < time.Now().UnixNano() {
			itemNotFound = append(itemNotFound, key)
			continue
		}

		itemFound[key] = item.data
	}

	return itemFound, itemNotFound, nil
}

// Set creates a new record with the specified key, value, and expiration
func (cl *KeyValueCustomClient) Set(key string, value interface{}, expire time.Duration) error {
	if cl.client == nil {
		return errors.New("client is not initialized")
	}
	
	if key == "" {
		return errors.New("key cannot be empty")
	}
	
	if value == nil {
		return errors.New("value cannot be nil")
	}

	// Set default expiration if not provided
	if expire <= 0 {
		expire = 24 * time.Hour
	}

	expirationTime := time.Now().Add(expire).UnixNano()
	
	item := customKeyValueItem{
		data:    value,
		expires: expirationTime,
	}
	
	if err := cl.client.Push(key, item); err != nil {
		log.Printf("Unable to push data for key %s: %v", key, err)
		return err
	}

	return nil
}

// Update modifies an existing key with a new value and expiration
// Returns an error if the key doesn't exist
func (cl *KeyValueCustomClient) Update(key string, value interface{}, expire time.Duration) error {
	if cl.client == nil {
		return errors.New("client is not initialized")
	}
	
	if key == "" {
		return errors.New("key cannot be empty")
	}
	
	if value == nil {
		return errors.New("value cannot be nil")
	}

	// Check if key exists
	_, err := cl.client.Get(key)
	if err != nil {
		log.Printf("Key %s not found for update: %v", key, err)
		return errors.New("key not found: " + key)
	}

	// Set default expiration if not provided
	if expire <= 0 {
		expire = 24 * time.Hour
	}

	expirationTime := time.Now().Add(expire).UnixNano()
	
	item := customKeyValueItem{
		data:    value,
		expires: expirationTime,
	}
	
	if err := cl.client.Push(key, item); err != nil {
		log.Printf("Unable to update data for key %s: %v", key, err)
		return err
	}

	return nil
}

// Delete removes a key from the cache
// Returns an error if the key doesn't exist
func (cl *KeyValueCustomClient) Delete(key string) error {
	if cl.client == nil {
		return errors.New("client is not initialized")
	}
	
	if key == "" {
		return errors.New("key cannot be empty")
	}

	// The Get method in linear.Linear removes the item if found
	_, err := cl.client.Get(key)
	if err != nil {
		log.Printf("Key %s not found for deletion: %v", key, err)
		return errors.New("key not found: " + key)
	}

	return nil
}

// Range iterates over all non-expired items in the cache
// The provided function is called for each key-value pair
func (cl *KeyValueCustomClient) Range(f func(key, value interface{}) bool) {
	if cl.client == nil || f == nil {
		return
	}

	fn := func(key, value interface{}) bool {
		// Skip if value is nil
		if value == nil {
			return true
		}
		
		// Try to convert to customKeyValueItem
		item, ok := value.(customKeyValueItem)
		if !ok {
			log.Printf("Warning: Unable to map object to customKeyValueItem model for key: %v", key)
			return true
		}

		// Skip expired items
		if item.expires > 0 && item.expires < time.Now().UnixNano() {
			// Optionally remove expired items during iteration
			if k, ok := key.(string); ok {
				cl.client.Get(k)
			}
			return true
		}

		// Call the user-provided function with the actual data
		return f(key, item.data)
	}

	cl.client.Range(fn)
}

// GetNumberOfRecords returns the total number of records in the cache
// Note: This includes expired records that haven't been cleaned up yet
func (cl *KeyValueCustomClient) GetNumberOfRecords() int {
	if cl.client == nil {
		return 0
	}
	return cl.client.GetNumberOfKeys()
}

// GetCapacity returns the current size of the cache in bytes
func (cl *KeyValueCustomClient) GetCapacity() (interface{}, error) {
	if cl.client == nil {
		return 0, errors.New("client is not initialized")
	}
	return cl.client.GetLinearCurrentSize(), nil
}

// Close stops the background cleaning process and frees up resources
func (cl *KeyValueCustomClient) Close() error {
	if cl.client == nil {
		return errors.New("client is not initialized")
	}
	
	// Signal the cleaning goroutine to stop
	select {
	case cl.close <- struct{}{}:
		// Successfully sent close signal
	default:
		// Channel is already closed or full
		return errors.New("close channel is unavailable")
	}

	return nil
}
