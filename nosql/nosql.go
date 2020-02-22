package nosql

import (
	"context"
	"reflect"
)

type INoSQL interface {
	GetALL(
		databaseName,
		collectionName,
		lastID,
		pageSize string,
		dataModel reflect.Type) (results interface{}, err error)

	GetByField(
		databaseName,
		collectionName,
		field,
		value string,
		dataModel reflect.Type) (interface{}, error)

	Create(
		databaseName,
		collectionName string,
		dataModel interface{}) (result interface{}, err error)

	Update(
		databaseName,
		collectionName,
		ID string,
		dataModel interface{}) (result interface{}, err error)

	Delete(
		databaseName,
		collectionName,
		ID string) (result interface{}, err error)

	MatchAndLookup(
		databaseName,
		collectionName string,
		model MatchLookup,
		dataModel reflect.Type) (results interface{}, err error)
}

var ctx = context.Background()

const (
	MONGODB = iota
)

// New function for NoSQL factory pattern
func New(
	databaseType int,
	config *Config) INoSQL {

	switch databaseType {
	case MONGODB:
		return NewMongoDB(&config.MongoDB)
	}

	return nil
}
