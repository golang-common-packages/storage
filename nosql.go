package storage

import "reflect"

// INoSQL factory pattern interface
type INoSQL interface {
	Find(databaseName, collectionName string, filter interface{}, limit int64, dataModel reflect.Type) (interface{}, error)
	Insert(databaseName, collectionName string, documents []interface{}) (interface{}, error)
	Update(databaseName, collectionName string, filter, update interface{}) (interface{}, error)
	Delete(databaseName, collectionName string, filter interface{}) (interface{}, error)
}

const (
	// MONGODB database
	MONGODB = iota
)

// NewNoSQL factory pattern
func NewNoSQL(
	databaseCompany int,
	config *Config) interface{} {

	switch databaseCompany {
	case MONGODB:
		return NewMongoDB(&config.MongoDB)
	}

	return nil
}
