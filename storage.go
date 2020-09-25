package storage

import "context"

/*
	@SQL: SQL database
	@NOSQL: NoSQL database
	@CACHING: Caching database
	@FILE: File service
*/
const (
	SQL = iota
	NOSQL
	CACHING
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
	case NOSQL:
		return NewNoSQL
	case CACHING:
		return NewCaching
	case FILE:
		return NewFile
	default:
		return nil
	}
}
