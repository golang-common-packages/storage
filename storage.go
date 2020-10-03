package storage

import "context"

const (
	// SQL type
	SQL = iota
	// NOSQLDOCUMENT is NoSQL document type
	NOSQLDOCUMENT
	// CACHING type
	CACHING
	// FILE type
	FILE
)

var (
	ctx = context.Background()
)

// New database by abstract factory pattern
func New(databaseType int) func(databaseCompany int, config *Config) interface{} {
	switch databaseType {
	case SQL:
		return NewSQL
	case NOSQLDOCUMENT:
		return NewNoSQLDocument
	case CACHING:
		return NewCaching
	case FILE:
		return NewFile
	default:
		return nil
	}
}
