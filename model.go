package storage

import (
	"time"

	"github.com/allegro/bigcache/v2"
	"google.golang.org/api/drive/v3"
)

// Begin Database Connection Models //

// Config model for database config
type Config struct {
	LIKE        LIKE            `json:"like,omitempty"`
	MongoDB     MongoDB         `json:"mongodb,omitempty"`
	Redis       Redis           `json:"redis,omitempty"`
	CustomCache CustomCache     `json:"customCache,omitempty"`
	BigCache    bigcache.Config `json:"bigCache,omitempty"`
	GoogleDrive GoogleDrive     `json:"googleDrive,omitempty"`
	CustomFile  CustomFile      `json:"customFile,omitempty"`
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

// GoogleDrive config model
type GoogleDrive struct {
	PoolSize     int    `json:"poolSize"`
	ByHTTPClient bool   `json:"byHTTPClient,omitempty"`
	Token        string `json:"token,omitempty"`
	Credential   string `json:"credential,omitempty"`
}

// CustomFile config model
type CustomFile struct {
	PoolSize             int    `json:"poolSize"`
	RootServiceDirectory string `json:"rootDirectory"`
}

// End Database Connection Models //
// -------------------------------------------------------------------------
// Begin Caching Models //

// customCacheItem private model for custom cache record
type customCacheItem struct {
	data    interface{}
	expires int64
}

// End Caching Models //
// -------------------------------------------------------------------------
// Begin File Models //

// GoogleFileListModel for unmarshal object has interface type
type GoogleFileListModel struct {
	drive.FileList
}

// GoogleFileModel for unmarshal object has interface type
type GoogleFileModel struct {
	drive.File
}

// End Caching Models //
// -------------------------------------------------------------------------
