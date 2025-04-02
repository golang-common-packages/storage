package tests

import (
	"context"
	"testing"

	"github.com/golang-common-packages/storage"
	"github.com/stretchr/testify/assert"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// TestSQLLikeOperations kiểm tra các thao tác cơ bản với SQL-Like
func TestSQLLikeOperations(t *testing.T) {
	// Bỏ qua test này nếu đang chạy trong CI/CD
	if testing.Short() {
		t.Skip("Skipping SQL-Like tests in short mode")
	}

	// Tạo cấu hình SQL-Like với SQLite
	config := &storage.LIKE{
		DriverName:     "sqlite3",
		DataSourceName: ":memory:", // In-memory SQLite database
	}

	// Khởi tạo SQL-Like client
	factory := storage.New(context.Background(), storage.SQLRELATIONAL)
	sqlClient := factory(storage.SQLLike, &storage.Config{
		LIKE: *config,
	})

	// Kiểm tra client không nil
	assert.NotNil(t, sqlClient, "SQL-Like client should not be nil")

	// Ép kiểu về interface ISQLRelational
	client, ok := sqlClient.(storage.ISQLRelational)
	assert.True(t, ok, "Should be able to cast to ISQLRelational")

	// Tạo bảng test
	createTableQuery := `
	CREATE TABLE test_table (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		value INTEGER
	)
	`
	_, err := client.Execute(createTableQuery, nil)
	assert.NoError(t, err, "Create table should not return an error")

	// Chèn dữ liệu
	insertQuery := `
	INSERT INTO test_table (name, value) VALUES 
	('Test 1', 10),
	('Test 2', 20),
	('Test 3', 30)
	`
	_, err = client.Execute(insertQuery, nil)
	assert.NoError(t, err, "Insert should not return an error")

	// Truy vấn dữ liệu
	selectQuery := `SELECT * FROM test_table`
	
	// Struct để nhận kết quả
	type TestRow struct {
		ID    int
		Name  string
		Value int
	}
	
	result, err := client.Execute(selectQuery, &TestRow{})
	assert.NoError(t, err, "Select should not return an error")
	assert.NotNil(t, result, "Select should return non-nil result")

	// Kiểm tra kết quả
	rows, ok := result.([]interface{})
	assert.True(t, ok, "Result should be a slice of interfaces")
	assert.Equal(t, 3, len(rows), "Should have 3 rows")

	// Cập nhật dữ liệu
	updateQuery := `UPDATE test_table SET value = 15 WHERE name = 'Test 1'`
	_, err = client.Execute(updateQuery, nil)
	assert.NoError(t, err, "Update should not return an error")

	// Kiểm tra cập nhật
	selectUpdatedQuery := `SELECT value FROM test_table WHERE name = 'Test 1'`
	updatedResult, err := client.Execute(selectUpdatedQuery, &struct{ Value int }{})
	assert.NoError(t, err, "Select after update should not return an error")
	
	updatedRows, ok := updatedResult.([]interface{})
	assert.True(t, ok, "Updated result should be a slice of interfaces")
	if len(updatedRows) > 0 {
		row, ok := updatedRows[0].(*struct{ Value int })
		assert.True(t, ok, "Row should be of correct type")
		assert.Equal(t, 15, row.Value, "Value should be updated to 15")
	}

	// Xóa dữ liệu
	deleteQuery := `DELETE FROM test_table WHERE name = 'Test 3'`
	_, err = client.Execute(deleteQuery, nil)
	assert.NoError(t, err, "Delete should not return an error")

	// Kiểm tra xóa
	countQuery := `SELECT COUNT(*) as count FROM test_table`
	countResult, err := client.Execute(countQuery, &struct{ Count int }{})
	assert.NoError(t, err, "Count should not return an error")
	
	countRows, ok := countResult.([]interface{})
	assert.True(t, ok, "Count result should be a slice of interfaces")
	if len(countRows) > 0 {
		row, ok := countRows[0].(*struct{ Count int })
		assert.True(t, ok, "Count row should be of correct type")
		assert.Equal(t, 2, row.Count, "Should have 2 rows after deletion")
	}
}
