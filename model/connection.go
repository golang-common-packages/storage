package model

import "time"

///// Connection Config Model /////

// Config model for database config
type Config struct {
	MongoDB MongoDB `json:"mongodb"`
	LIKE    LIKE    `json:"like"`
}

// LIKE model for SQL-LIKE config
type LIKE struct {
	DriverName            string        `json:"driverName"`
	DataSourceName        string        `json:"dataSourceName"`
	MaxConnectionLifetime time.Duration `json:"maxConnectionLifetime"`
	MaxConnectionIdle     int           `json:"maxConnectionIdle"`
	MaxConnectionOpen     int           `json:"maxConnectionOpen"`
}

// MongoDB model for MongoDB config
type MongoDB struct {
	User     string   `json:"user"`
	Password string   `json:"password"`
	Hosts    []string `json:"hosts"`
	DB       string   `json:"db"`
	Options  []string `json:"options"`
}
