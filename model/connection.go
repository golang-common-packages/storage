package model

import (
	"time"

	"github.com/allegro/bigcache/v2"
)

// Config model for database config
type Config struct {
	LIKE        LIKE            `json:"like"`
	MongoDB     MongoDB         `json:"mongodb"`
	Redis       Redis           `json:"redis,omitempty"`
	CustomCache CustomCache     `json:"customCache,omitempty"`
	BigCache    bigcache.Config `json:"bigCache,omitempty"`
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
