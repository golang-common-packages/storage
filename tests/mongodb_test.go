package tests

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang-common-packages/storage"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

// TestDocument là một struct để test MongoDB
type TestDocument struct {
	ID    string `bson:"_id,omitempty"`
	Name  string `bson:"name"`
	Value int    `bson:"value"`
}

// setupMongoDBTest thiết lập kết nối đến MongoDB cho việc testing
// Lưu ý: Hàm này yêu cầu một MongoDB server đang chạy
// Trong môi trường CI/CD, bạn có thể sử dụng mongodb-memory-server
func setupMongoDBTest(t *testing.T) (storage.INoSQLDocument, func()) {
	// Bỏ qua test này nếu đang chạy trong CI/CD
	if testing.Short() {
		t.Skip("Skipping MongoDB tests in short mode")
	}

	// Tạo cấu hình MongoDB
	config := &storage.MongoDB{
		Hosts:   []string{"localhost:27017"},
		DB:      "test_db",
		Options: []string{},
	}

	// Khởi tạo MongoDB client
	factory := storage.New(context.Background(), storage.NOSQLDOCUMENT)
	mongoClient := factory(storage.MONGODB, &storage.Config{
		MongoDB: *config,
	})

	// Kiểm tra client không nil
	if mongoClient == nil {
		t.Skip("MongoDB client is nil, skipping test (MongoDB server may not be running)")
	}

	// Ép kiểu về interface INoSQLDocument
	client, ok := mongoClient.(storage.INoSQLDocument)
	if !ok {
		t.Fatal("Failed to cast to INoSQLDocument")
	}

	// Hàm cleanup
	cleanup := func() {
		// Xóa database sau khi test
		mongoClient, ok := client.(*storage.MongoClient)
		if ok && mongoClient.Client != nil {
			mongoClient.Client.Database("test_db").Drop(context.Background())
		}
	}

	return client, cleanup
}

// TestMongoDBOperations kiểm tra các thao tác cơ bản với MongoDB
func TestMongoDBOperations(t *testing.T) {
	// Thiết lập MongoDB
	client, cleanup := setupMongoDBTest(t)
	defer cleanup()

	// Kiểm tra client không nil
	assert.NotNil(t, client, "MongoDB client should not be nil")

	// Tạo dữ liệu test
	testDocs := []interface{}{
		TestDocument{Name: "Test 1", Value: 10},
		TestDocument{Name: "Test 2", Value: 20},
		TestDocument{Name: "Test 3", Value: 30},
	}

	// Test Create
	result, err := client.Create("test_db", "test_collection", testDocs)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	assert.NotNil(t, result, "Create should return non-nil result")

	// Test Read
	filter := bson.M{}
	readResult, err := client.Read("test_db", "test_collection", filter, 10, reflect.TypeOf(TestDocument{}))
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	assert.NotNil(t, readResult, "Read should return non-nil result")

	// Test Update
	updateFilter := bson.M{"name": "Test 1"}
	update := bson.M{"$set": bson.M{"value": 15}}
	updateResult, err := client.Update("test_db", "test_collection", updateFilter, update)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	assert.NotNil(t, updateResult, "Update should return non-nil result")

	// Kiểm tra update đã thành công
	readAfterUpdate, err := client.Read("test_db", "test_collection", updateFilter, 1, reflect.TypeOf(TestDocument{}))
	if err != nil {
		t.Fatalf("Read after update failed: %v", err)
	}
	
	// Kiểm tra giá trị đã được cập nhật
	if docs, ok := readAfterUpdate.(*[]TestDocument); ok && len(*docs) > 0 {
		assert.Equal(t, 15, (*docs)[0].Value, "Value should be updated to 15")
	}

	// Test Delete
	deleteFilter := bson.M{"name": "Test 3"}
	deleteResult, err := client.Delete("test_db", "test_collection", deleteFilter)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	assert.NotNil(t, deleteResult, "Delete should return non-nil result")

	// Kiểm tra delete đã thành công
	readAfterDelete, err := client.Read("test_db", "test_collection", deleteFilter, 1, reflect.TypeOf(TestDocument{}))
	if err != nil {
		t.Fatalf("Read after delete failed: %v", err)
	}
	
	// Kiểm tra document đã bị xóa
	if docs, ok := readAfterDelete.(*[]TestDocument); ok {
		assert.Equal(t, 0, len(*docs), "Document should be deleted")
	}
}
