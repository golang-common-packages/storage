package tests

import (
	"context"
	"testing"
	"time"

	"github.com/golang-common-packages/storage"
	"github.com/stretchr/testify/assert"
)

func TestNewStorage(t *testing.T) {
	testCases := []struct {
		name        string
		storageType storage.StorageType
		expectNil   bool
	}{
		{
			name:        "SQL Relational",
			storageType: storage.SQLRELATIONAL,
			expectNil:   false,
		},
		{
			name:        "NoSQL Document",
			storageType: storage.NOSQLDOCUMENT,
			expectNil:   false,
		},
		{
			name:        "NoSQL Key-Value",
			storageType: storage.NOSQLKEYVALUE,
			expectNil:   false,
		},
		{
			name:        "File",
			storageType: storage.FILE,
			expectNil:   false,
		},
		{
			name:        "Invalid Storage Type",
			storageType: storage.StorageType(999),
			expectNil:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			factory := storage.New(ctx, tc.storageType)
			
			if tc.expectNil {
				assert.Nil(t, factory(0, nil), "Expected nil factory for invalid storage type")
			} else {
				assert.NotNil(t, factory, "Expected non-nil factory for valid storage type")
			}
		})
	}
}

func TestContextFunctions(t *testing.T) {
	ctx := context.WithValue(context.Background(), "test-key", "test-value")
	
	storage.SetContext(ctx)
	
	retrievedCtx := storage.GetContext()
	assert.Equal(t, "test-value", retrievedCtx.Value("test-key"), "Context value should match")
}
