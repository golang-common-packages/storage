// Package storage provides a unified interface for various storage types
// using the abstract factory pattern.
package storage

import (
	"context"
	"errors"
)

// StorageType defines the type of storage
type StorageType int

const (
	// SQLRELATIONAL is SQL relational type
	SQLRELATIONAL StorageType = iota
	// NOSQLDOCUMENT is NoSQL document type
	NOSQLDOCUMENT
	// NOSQLKEYVALUE is NoSQL key-value type
	NOSQLKEYVALUE
	// FILE is file management
	FILE
)

var (
	// ErrInvalidStorageType is returned when an invalid storage type is provided
	ErrInvalidStorageType = errors.New("invalid storage type")
	// ErrInvalidConfig is returned when an invalid configuration is provided
	ErrInvalidConfig = errors.New("invalid configuration")
	
	// ctx is the default context
	ctx = context.Background()
)

// New creates a new storage instance using the abstract factory pattern
// It returns a function that can be called with a specific database company and config
// to get the concrete implementation.
func New(context context.Context, storageType StorageType) func(databaseCompany int, config *Config) interface{} {
	if context != nil {
		SetContext(context)
	}

	switch storageType {
	case SQLRELATIONAL:
		return newSQLRelational
	case NOSQLDOCUMENT:
		return newNoSQLDocument
	case NOSQLKEYVALUE:
		return newNoSQLKeyValue
	case FILE:
		return newFile
	default:
		return func(_ int, _ *Config) interface{} {
			return nil
		}
	}
}
