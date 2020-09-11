package database

import (
	"github.com/golang-common-packages/database/sql"
	"github.com/golang-common-packages/database/nosql"
	"github.com/golang-common-packages/database/model"
)

const (
	SQL = iota
	NOSQL
)

// New function for database abstract factory
func New(databaseType int) func(databaseType int, config *model.Config) interface{} {
	switch databaseType {
	case SQL:
		return sql.New
	case NOSQL:
		return nosql.New
	default:
		return nil
	}
}