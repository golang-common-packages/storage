package model

import (
	"time"

	"github.com/allegro/bigcache/v2"
)

///// Connection Config Model /////

// Config model for database config
type Config struct {
	MongoDB MongoDB `json:"mongodb"`
	LIKE    LIKE    `json:"like"`
	Caching Caching `json:"caching"`
}

// LIKE model for SQL-LIKE connection config
type LIKE struct {
	DriverName            string        `json:"driverName"`
	DataSourceName        string        `json:"dataSourceName"`
	MaxConnectionLifetime time.Duration `json:"maxConnectionLifetime"`
	MaxConnectionIdle     int           `json:"maxConnectionIdle"`
	MaxConnectionOpen     int           `json:"maxConnectionOpen"`
}

// MongoDB model for MongoDB connection config
type MongoDB struct {
	User     string   `json:"user"`
	Password string   `json:"password"`
	Hosts    []string `json:"hosts"`
	DB       string   `json:"db"`
	Options  []string `json:"options"`
}

// Caching model for database caching connection config
type Caching struct {
	CustomCache CustomCache     `json:"customCache,omitempty"`
	Redis       Redis           `json:"redis,omitempty"`
	BigCache    bigcache.Config `json:"bigCache,omitempty"`
}

// Redis model for redis config
type Redis struct {
	Password   string `json:"password,omitempty"`
	Host       string `json:"host,omitempty"`
	DB         int    `json:"db,omitempty"`
	MaxRetries int    `json:"maxRetries,omitempty"`
}

// CustomCache config model
type CustomCache struct {
	CacheSize        int64         `json:"cacheSize,omitempty"` // byte
	CleaningEnable   bool          `json:"cleaningEnable,omitempty"`
	CleaningInterval time.Duration `json:"cleaningInterval,omitempty"` // nanosecond
}
