package database

import (
	"github.com/golang-common-packages/database/caching"
	"github.com/golang-common-packages/database/model"
	"github.com/golang-common-packages/database/nosql"
	"github.com/golang-common-packages/database/sql"
)

const (
	SQL = iota
	NOSQL
	CACHING
)

// New function for database abstract factory
func New(databaseType int) func(databaseCompany int, config *model.Config) interface{} {
	switch databaseType {
	case SQL:
		return sql.New
	case NOSQL:
		return nosql.New
	case CACHING:
		return caching.New
	default:
		return nil
	}
}
