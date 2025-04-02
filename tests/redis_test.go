package tests

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/golang-common-packages/storage"
	"github.com/stretchr/testify/assert"
)

func setupMiniRedis(t *testing.T) (*miniredis.Miniredis, *storage.Redis) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to create miniredis: %v", err)
	}

	config := &storage.Redis{
		Host:       s.Addr(),
		Password:   "",
		DB:         0,
		MaxRetries: 3,
	}

	return s, config
}

func TestRedisOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Redis tests in short mode")
	}

	s, config := setupMiniRedis(t)
	defer s.Close()

	factory := storage.New(nil, storage.NOSQLKEYVALUE)
	redisClient := factory(storage.REDIS, &storage.Config{
		Redis: *config,
	})

	assert.NotNil(t, redisClient, "Redis client should not be nil")

	client, ok := redisClient.(storage.INoSQLKeyValue)
	assert.True(t, ok, "Should be able to cast to INoSQLKeyValue")

	err := client.Set("test-key", "test-value", 1*time.Hour)
	assert.NoError(t, err, "Set should not return an error")

	val, err := s.Get("test-key")
	assert.NoError(t, err, "Should be able to get value from miniredis")
	assert.Equal(t, "test-value", val, "Value should match")

	result, err := client.Get("test-key")
	assert.NoError(t, err, "Get should not return an error")
	assert.Equal(t, "test-value", result, "Get should return the correct value")

	err = client.Update("test-key", "updated-value", 1*time.Hour)
	assert.NoError(t, err, "Update should not return an error")

	val, err = s.Get("test-key")
	assert.NoError(t, err, "Should be able to get updated value")
	assert.Equal(t, "updated-value", val, "Updated value should match")

	err = client.Delete("test-key")
	assert.NoError(t, err, "Delete should not return an error")

	exists := s.Exists("test-key")
	assert.False(t, exists, "Key should be deleted")

	_, err = client.Get("non-existent-key")
	assert.Nil(t, err, "Get should return a nil for non-existent key")
}
