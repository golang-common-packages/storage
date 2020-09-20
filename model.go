package storage

import (
	"time"

	"github.com/allegro/bigcache/v2"
	"google.golang.org/api/drive/v3"
)

// Begin Database Connection Models //

// Config model for database config
type Config struct {
	LIKE        LIKE            `json:"like"`
	MongoDB     MongoDB         `json:"mongodb"`
	Redis       Redis           `json:"redis,omitempty"`
	CustomCache CustomCache     `json:"customCache,omitempty"`
	BigCache    bigcache.Config `json:"bigCache,omitempty"`
	GoogleDrive GoogleDrive     `json:"googleDrive,omitempty"`
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
	ByHTTPClient bool   `json:"byHTTPClient,omitempty"`
	Token        string `json:"token,omitempty"`
	Credential   string `json:"credential,omitempty"`
}

// End Database Connection Models //
// -------------------------------------------------------------------------
// Begin NoSQL Models //

// MatchLookup mongo model
type MatchLookup struct {
	Match  []Match  `json:"match"`
	Lookup []Lookup `json:"lookup"`
}

// Match mongo model
type Match struct {
	Field    string              `json:"field"`
	Operator ComparisonOperators `json:"operator"`
	Value    string              `json:"value"`
}

// Lookup mongo model
type Lookup struct {
	From         string `json:"From"`
	LocalField   string `json:"localField"`
	ForeignField string `json:"foreignField"`
	As           string `json:"as"`
}

// Set mongo model
type Set struct {
	Operator UpdateOperators `json:"operator"`
	Data     interface{}     `json:"data"`
}

///// MongoDB operator model /////

// ComparisonOperators mongodb comparition operation type
type ComparisonOperators string

/*
This is for mongodb comparition operation constant
*/
const (
	Equal                ComparisonOperators = "$eq"
	EqualAny             ComparisonOperators = "$in"
	NotEqual             ComparisonOperators = "$ne"
	NotEqualAnyLL        ComparisonOperators = "$nin"
	GreaterThan          ComparisonOperators = "$gt"
	GreaterThanOrEqualTo ComparisonOperators = "$gte"
	LessThan             ComparisonOperators = "$lt"
	LessThanOrEqualTo    ComparisonOperators = "$lte"
)

// UpdateOperators mongodb update operation type
type UpdateOperators string

/*
This is for mongodb update operation constant
*/
const (
	Replaces UpdateOperators = "$set"
)

// End NoSQL Models //
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
