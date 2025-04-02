package storage

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"

	"github.com/golang-common-packages/hash"
)

// RedisClient manage all redis actions
type RedisClient struct {
	Client *redis.Client
}

var (
	// redisClientSessionMapping singleton pattern
	redisClientSessionMapping = make(map[string]*RedisClient)
)

// newRedis init new instance
func newRedis(config *Redis) INoSQLKeyValue {
	hasher := &hash.Client{}
	configAsJSON, err := json.Marshal(config)
	if err != nil {
		log.Fatalln("Unable to marshal Redis configuration: ", err)
	}
	configAsString := hasher.SHA1(string(configAsJSON))

	currentRedisClientSession := redisClientSessionMapping[configAsString]
	if currentRedisClientSession == nil {
		currentRedisClientSession = &RedisClient{nil}
		client, err := currentRedisClientSession.connect(config)
		if err != nil {
			log.Fatalln("Unable to connect to Redis: ", err)
		} else {
			currentRedisClientSession.Client = client
			redisClientSessionMapping[configAsString] = currentRedisClientSession
			log.Println("Connected to Redis")
		}
	}

	return currentRedisClientSession
}

func (r *RedisClient) connect(data *Redis) (client *redis.Client, err error) {
	if r.Client == nil {
		client = redis.NewClient(&redis.Options{
			Addr:       data.Host,
			Password:   data.Password,
			DB:         data.DB,
			MaxRetries: data.MaxRetries,
		})

		_, err := client.Ping().Result()
		if err != nil {
			log.Fatalln("Unable to connect to Redis: ", err)
			return nil, err
		}
	} else {
		client = r.Client
		err = nil
	}
	return
}

// Middleware for echo framework
func (r *RedisClient) Middleware(hash hash.IHash) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get(echo.HeaderAuthorization)
			key := hash.SHA512(token)

			if val, err := r.Get(key); err != nil {
				log.Println("Can not get accesstoken from redis in echo middleware: ", err)
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			} else if val == "" {
				return c.NoContent(http.StatusUnauthorized)
			}

			return next(c)
		}
	}
}

// Get retrieves a value from Redis based on the key provided
func (r *RedisClient) Get(key string) (interface{}, error) {
	if r.Client == nil {
		return nil, errors.New("redis client is not initialized")
	}
	
	if key == "" {
		return nil, errors.New("key cannot be empty")
	}
	
	result, err := r.Client.Get(key).Result()
	if err == redis.Nil {
		return "", nil // Key does not exist
	}
	return result, err
}

// Set creates a new record with the specified key, value, and expiration
func (r *RedisClient) Set(key string, value interface{}, expire time.Duration) error {
	if r.Client == nil {
		return errors.New("redis client is not initialized")
	}
	
	if key == "" {
		return errors.New("key cannot be empty")
	}
	
	return r.Client.Set(key, value, expire).Err()
}

// Update modifies an existing key with a new value and expiration
// Returns an error if the key doesn't exist
func (r *RedisClient) Update(key string, value interface{}, expire time.Duration) error {
	if r.Client == nil {
		return errors.New("redis client is not initialized")
	}
	
	if key == "" {
		return errors.New("key cannot be empty")
	}
	
	// Check if key exists
	exists, err := r.Client.Exists(key).Result()
	if err != nil {
		log.Printf("Unable to check if key exists: %v", err)
		return err
	}
	
	if exists == 0 {
		return errors.New("key does not exist: " + key)
	}

	return r.Client.Set(key, value, expire).Err()
}

// Append adds a string value to the end of an existing string key
func (r *RedisClient) Append(key string, value interface{}) error {
	if r.Client == nil {
		return errors.New("redis client is not initialized")
	}
	
	if key == "" {
		return errors.New("key cannot be empty")
	}
	
	// Convert value to string
	var stringValue string
	switch v := value.(type) {
	case string:
		stringValue = v
	case []byte:
		stringValue = string(v)
	default:
		b, err := json.Marshal(value)
		if err != nil {
			log.Printf("Unable to marshal value: %v", err)
			return errors.New("cannot marshal value: " + err.Error())
		}
		stringValue = string(b)
	}

	_, err := r.Client.Append(key, stringValue).Result()
	return err
}

// Delete removes a key from Redis
func (r *RedisClient) Delete(key string) error {
	if r.Client == nil {
		return errors.New("redis client is not initialized")
	}
	
	if key == "" {
		return errors.New("key cannot be empty")
	}
	
	_, err := r.Client.Del(key).Result()
	return err
}

// GetNumberOfRecords return number of records
func (r *RedisClient) GetNumberOfRecords() int {
	return len(r.Client.Do("KEYS", "*").Args())
}

// GetCapacity method return redis database size
func (r *RedisClient) GetCapacity() (interface{}, error) {
	IntCmd := r.Client.DBSize()

	return IntCmd.Result()
}

// Close method will close redis connection
func (r *RedisClient) Close() error {
	return r.Client.Close()
}
