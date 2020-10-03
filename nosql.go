package storage

import "reflect"

// INoSQL factory pattern CRUD interface
type INoSQL interface {
	Create(databaseName, collectionName string, documents []interface{}) (interface{}, error)
	Read(databaseName, collectionName string, filter interface{}, limit int64, dataModel reflect.Type) (interface{}, error)
	Update(databaseName, collectionName string, filter, update interface{}) (interface{}, error)
	Delete(databaseName, collectionName string, filter interface{}) (interface{}, error)
}

const (
	// MONGODB database
	MONGODB = iota
)

// NewNoSQL init instance by factory pattern
func NewNoSQL(databaseCompany int, config *Config) interface{} {

	switch databaseCompany {
	case MONGODB:
		return NewMongoDB(&config.MongoDB)
	}

	return nil
}
