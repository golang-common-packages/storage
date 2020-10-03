package storage

import "context"

const (
	// SQLRELATIONAL is SQL relational type
	SQLRELATIONAL = iota
	// NOSQLDOCUMENT is NoSQL document type
	NOSQLDOCUMENT
	// NOSQLKEYVALUE is NoSQL key-value type
	NOSQLKEYVALUE
	// FILE is file management
	FILE
)

var (
	ctx = context.Background()
)

// New database by abstract factory pattern
func New(databaseType int) func(databaseCompany int, config *Config) interface{} {
	switch databaseType {
	case SQLRELATIONAL:
		return NewSQLRelational
	case NOSQLDOCUMENT:
		return NewNoSQLDocument
	case NOSQLKEYVALUE:
		return NewNoSQLKeyValue
	case FILE:
		return NewFile
	default:
		return nil
	}
}
