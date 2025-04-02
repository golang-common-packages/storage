package tests

import (
	"context"
	"testing"
	"time"

	"github.com/golang-common-packages/storage"
	"github.com/stretchr/testify/assert"
)

// TestNewStorage kiểm tra hàm New của package storage
func TestNewStorage(t *testing.T) {
	// Test với các loại storage khác nhau
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
			// Tạo context với timeout
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Gọi hàm New
			factory := storage.New(ctx, tc.storageType)
			
			// Kiểm tra kết quả
			if tc.expectNil {
				assert.Nil(t, factory(0, nil), "Expected nil factory for invalid storage type")
			} else {
				assert.NotNil(t, factory, "Expected non-nil factory for valid storage type")
			}
		})
	}
}

// TestSetContext kiểm tra hàm SetContext và GetContext
func TestContextFunctions(t *testing.T) {
	// Tạo context mới
	ctx := context.WithValue(context.Background(), "test-key", "test-value")
	
	// Thiết lập context
	storage.SetContext(ctx)
	
	// Lấy context và kiểm tra
	retrievedCtx := storage.GetContext()
	assert.Equal(t, "test-value", retrievedCtx.Value("test-key"), "Context value should match")
}
