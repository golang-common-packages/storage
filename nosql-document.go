package storage

import "reflect"

// INoSQLDocument factory pattern CRUD interface
type INoSQLDocument interface {
	Create(databaseName, collectionName string, documents []interface{}) (interface{}, error)
	Read(databaseName, collectionName string, filter interface{}, limit int64, dataModel reflect.Type) (interface{}, error)
	Update(databaseName, collectionName string, filter, update interface{}) (interface{}, error)
	Delete(databaseName, collectionName string, filter interface{}) (interface{}, error)
}

const (
	// MONGODB database
	MONGODB = iota
)

// newNoSQLDocument init instance by factory pattern
func newNoSQLDocument(databaseCompany int, config *Config) interface{} {

	switch databaseCompany {
	case MONGODB:
		return newMongoDB(&config.MongoDB)
	}

	return nil
}
