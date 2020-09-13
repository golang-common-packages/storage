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

// NewRedis init new instance
func NewRedis(config *Redis) ICaching {
	hasher := &hash.Client{}
	configAsJSON, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	configAsString := hasher.SHA1(string(configAsJSON))

	currentRedisClientSession := redisClientSessionMapping[configAsString]
	if currentRedisClientSession == nil {
		currentRedisClientSession = &RedisClient{nil}
		client, err := currentRedisClientSession.connect(config)
		if err != nil {
			panic(err)
		} else {
			currentRedisClientSession.Client = client
			redisClientSessionMapping[configAsString] = currentRedisClientSession
			log.Println("Connected to Redis Server")
		}
	}

	return currentRedisClientSession
}

// connect private method establish redis connection
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
			log.Println("Connect to redis fail: ", err)
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
				log.Printf("Can not get accesstoken from redis in echo middleware: %s", err.Error())
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			} else if val == "" {
				return c.NoContent(http.StatusUnauthorized)
			}

			return next(c)
		}
	}
}

// Get return value based on the key provided
func (r *RedisClient) Get(key string) (interface{}, error) {
	return r.Client.Get(key).Result()
}

// Set new record set key and value
func (r *RedisClient) Set(key string, value interface{}, expire time.Duration) error {
	return r.Client.Set(key, value, expire).Err()
}

// Update new value over the key provided
func (r *RedisClient) Update(key string, value interface{}, expire time.Duration) error {
	_, err := r.Client.Get(key).Result()
	if err != nil {
		return err
	}

	return r.Client.Set(key, value, expire).Err()
}

// Append new value over the key provided
func (r *RedisClient) Append(key string, value interface{}) error {
	b, err := json.Marshal(value)
	if err != nil {
		return errors.New("can not marshal value to []byte")
	}

	var v string
	json.Unmarshal(b, &v)

	return r.Client.Append(key, v).Err()
}

// Delete method delete value based on the key provided
func (r *RedisClient) Delete(key string) error {
	return r.Client.Del(key).Err()
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
