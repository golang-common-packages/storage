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
	// Init context with default value
	ctx = context.Background()
)

// New database by abstract factory pattern
func New(context context.Context, databaseType int) func(databaseCompany int, config *Config) interface{} {
	SetContext(context)

	switch databaseType {
	case SQLRELATIONAL:
		return newSQLRelational
	case NOSQLDOCUMENT:
		return newNoSQLDocument
	case NOSQLKEYVALUE:
		return newNoSQLKeyValue
	case FILE:
		return newFile
	default:
		return nil
	}
}
