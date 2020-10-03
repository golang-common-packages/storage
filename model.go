package storage

import (
	"time"

	"github.com/allegro/bigcache/v2"
	"google.golang.org/api/drive/v3"
)

// Begin Database Connection Models //

// Config model for database config
type Config struct {
	LIKE           LIKE            `json:"like,omitempty"`
	MongoDB        MongoDB         `json:"mongodb,omitempty"`
	Redis          Redis           `json:"redis,omitempty"`
	CustomKeyValue CustomKeyValue  `json:"customKeyValue,omitempty"`
	BigCache       bigcache.Config `json:"bigCache,omitempty"`
	GoogleDrive    GoogleDrive     `json:"googleDrive,omitempty"`
	CustomFile     CustomFile      `json:"customFile,omitempty"`
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
	Password   string `json:"password"`
	Host       string `json:"host"`
	DB         int    `json:"db"`
	MaxRetries int    `json:"maxRetries"`
}

// CustomKeyValue config model
type CustomKeyValue struct {
	MemorySize       int64         `json:"memorySize"` // byte
	CleaningEnable   bool          `json:"cleaningEnable"`
	CleaningInterval time.Duration `json:"cleaningInterval"` // nanosecond
}

// GoogleDrive config model
type GoogleDrive struct {
	PoolSize     int    `json:"poolSize"`
	ByHTTPClient bool   `json:"byHTTPClient"`
	Token        string `json:"token"`
	Credential   string `json:"credential"`
}

// CustomFile config model
type CustomFile struct {
	PoolSize             int    `json:"poolSize"`
	RootServiceDirectory string `json:"rootDirectory"`
}

// End Database Connection Models //

// -------------------------------------------------------------------------

// Begin Caching Models //

// customKeyValueItem private model for custom key-value record
type customKeyValueItem struct {
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
