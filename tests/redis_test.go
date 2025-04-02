package tests

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/golang-common-packages/storage"
	"github.com/stretchr/testify/assert"
)

// setupMiniRedis tạo một Redis server giả lập cho việc testing
func setupMiniRedis(t *testing.T) (*miniredis.Miniredis, *storage.Redis) {
	// Tạo Redis server giả lập
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to create miniredis: %v", err)
	}

	// Tạo cấu hình Redis
	config := &storage.Redis{
		Host:       s.Addr(),
		Password:   "",
		DB:         0,
		MaxRetries: 3,
	}

	return s, config
}

// TestRedisOperations kiểm tra các thao tác cơ bản với Redis
func TestRedisOperations(t *testing.T) {
	// Bỏ qua test này nếu đang chạy trong CI/CD
	if testing.Short() {
		t.Skip("Skipping Redis tests in short mode")
	}

	// Thiết lập Redis server giả lập
	s, config := setupMiniRedis(t)
	defer s.Close()

	// Khởi tạo Redis client
	factory := storage.New(nil, storage.NOSQLKEYVALUE)
	redisClient := factory(storage.REDIS, &storage.Config{
		Redis: *config,
	})

	// Kiểm tra client không nil
	assert.NotNil(t, redisClient, "Redis client should not be nil")

	// Ép kiểu về interface INoSQLKeyValue
	client, ok := redisClient.(storage.INoSQLKeyValue)
	assert.True(t, ok, "Should be able to cast to INoSQLKeyValue")

	// Test Set
	err := client.Set("test-key", "test-value", 1*time.Hour)
	assert.NoError(t, err, "Set should not return an error")

	// Kiểm tra giá trị trong miniredis
	val, err := s.Get("test-key")
	assert.NoError(t, err, "Should be able to get value from miniredis")
	assert.Equal(t, "test-value", val, "Value should match")

	// Test Get
	result, err := client.Get("test-key")
	assert.NoError(t, err, "Get should not return an error")
	assert.Equal(t, "test-value", result, "Get should return the correct value")

	// Test Update
	err = client.Update("test-key", "updated-value", 1*time.Hour)
	assert.NoError(t, err, "Update should not return an error")

	// Kiểm tra giá trị đã cập nhật
	val, err = s.Get("test-key")
	assert.NoError(t, err, "Should be able to get updated value")
	assert.Equal(t, "updated-value", val, "Updated value should match")

	// Test Delete
	err = client.Delete("test-key")
	assert.NoError(t, err, "Delete should not return an error")

	// Kiểm tra key đã bị xóa
	exists := s.Exists("test-key")
	assert.False(t, exists, "Key should be deleted")

	// Test Get non-existent key
	_, err = client.Get("non-existent-key")
	assert.Error(t, err, "Get should return an error for non-existent key")
}
