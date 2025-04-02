package tests

import (
	"testing"
	"time"

	"github.com/golang-common-packages/storage"
	"github.com/stretchr/testify/assert"
)

// TestCustomKeyValueOperations kiểm tra các thao tác cơ bản với Custom Key-Value
func TestCustomKeyValueOperations(t *testing.T) {
	// Tạo cấu hình Custom Key-Value
	config := &storage.CustomKeyValue{
		MemorySize:       1024 * 1024 * 10, // 10MB
		CleaningEnable:   true,
		CleaningInterval: 1 * time.Second,
	}

	// Khởi tạo Custom Key-Value client
	factory := storage.New(nil, storage.NOSQLKEYVALUE)
	customClient := factory(storage.CUSTOM, &storage.Config{
		CustomKeyValue: *config,
	})

	// Kiểm tra client không nil
	assert.NotNil(t, customClient, "Custom Key-Value client should not be nil")

	// Ép kiểu về interface INoSQLKeyValue
	client, ok := customClient.(storage.INoSQLKeyValue)
	assert.True(t, ok, "Should be able to cast to INoSQLKeyValue")

	// Test Set
	err := client.Set("test-key", "test-value", 1*time.Hour)
	assert.NoError(t, err, "Set should not return an error")

	// Test Get
	result, err := client.Get("test-key")
	assert.NoError(t, err, "Get should not return an error")
	assert.Equal(t, "test-value", result, "Get should return the correct value")

	// Test Update
	err = client.Update("test-key", "updated-value", 1*time.Hour)
	assert.NoError(t, err, "Update should not return an error")

	// Test Get sau khi update
	result, err = client.Get("test-key")
	assert.NoError(t, err, "Get should not return an error after update")
	assert.Equal(t, "updated-value", result, "Get should return the updated value")

	// Test Delete
	err = client.Delete("test-key")
	assert.NoError(t, err, "Delete should not return an error")

	// Test Get sau khi delete
	_, err = client.Get("test-key")
	assert.Error(t, err, "Get should return an error after delete")

	// Test Set với expiration ngắn
	err = client.Set("expiring-key", "expiring-value", 100*time.Millisecond)
	assert.NoError(t, err, "Set with short expiration should not return an error")

	// Đợi cho key hết hạn
	time.Sleep(200 * time.Millisecond)

	// Test Get sau khi key hết hạn
	result, err = client.Get("expiring-key")
	assert.Nil(t, result, "Get should return nil for expired key")

	// Test GetNumberOfRecords
	count := client.GetNumberOfRecords()
	assert.GreaterOrEqual(t, count, 0, "GetNumberOfRecords should return a non-negative number")

	// Test GetCapacity
	capacity, err := client.GetCapacity()
	assert.NoError(t, err, "GetCapacity should not return an error")
	assert.NotNil(t, capacity, "GetCapacity should return a non-nil value")

	// Đóng client
	err = client.Close()
	assert.NoError(t, err, "Close should not return an error")
}
