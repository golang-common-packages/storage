package sql

import (
	"github.com/golang-common-packages/database/model"
)

// ISQL interface for SQL factory pattern
type ISQL interface {
	Execute(
		query string,
		dataModel interface{}) (interface{}, error)
}

const (
	SqlLike = iota
)

// New function for SQL factory pattern
func New(
	databaseCompany int,
	config *model.Config) interface{} {

	switch databaseCompany {
	case SqlLike:
		return NewSQLLike(&config.LIKE)
	}

	return nil
}
