package tests

import (
	"testing"
	"time"

	"github.com/golang-common-packages/storage"
	"github.com/stretchr/testify/assert"
)

func TestCustomKeyValueOperations(t *testing.T) {
	config := &storage.CustomKeyValue{
		MemorySize:       1024 * 1024 * 10, // 10MB
		CleaningEnable:   true,
		CleaningInterval: 1 * time.Second,
	}

	factory := storage.New(nil, storage.NOSQLKEYVALUE)
	customClient := factory(storage.CUSTOM, &storage.Config{
		CustomKeyValue: *config,
	})

	assert.NotNil(t, customClient, "Custom Key-Value client should not be nil")

	client, ok := customClient.(storage.INoSQLKeyValue)
	assert.True(t, ok, "Should be able to cast to INoSQLKeyValue")

	err := client.Set("test-key", "test-value", 1*time.Hour)
	assert.NoError(t, err, "Set should not return an error")

	result, err := client.Get("test-key")
	assert.NoError(t, err, "Get should not return an error")
	assert.Equal(t, "test-value", result, "Get should return the correct value")

	err = client.Update("test-key", "updated-value", 1*time.Hour)
	assert.NoError(t, err, "Update should not return an error")

	result, err = client.Get("test-key")
	assert.NoError(t, err, "Get should not return an error after update")
	assert.Equal(t, "updated-value", result, "Get should return the updated value")

	err = client.Delete("test-key")
	assert.NoError(t, err, "Delete should not return an error")

	_, err = client.Get("test-key")
	assert.Error(t, err, "Get should return an error after delete")

	err = client.Set("expiring-key", "expiring-value", 100*time.Millisecond)
	assert.NoError(t, err, "Set with short expiration should not return an error")

	time.Sleep(200 * time.Millisecond)

	result, err = client.Get("expiring-key")
	assert.Nil(t, result, "Get should return nil for expired key")

	count := client.GetNumberOfRecords()
	assert.GreaterOrEqual(t, count, 0, "GetNumberOfRecords should return a non-negative number")

	capacity, err := client.GetCapacity()
	assert.NoError(t, err, "GetCapacity should not return an error")
	assert.NotNil(t, capacity, "GetCapacity should return a non-nil value")

	err = client.Close()
	assert.NoError(t, err, "Close should not return an error")
}
