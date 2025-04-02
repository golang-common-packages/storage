package tests

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/golang-common-packages/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDriveService là một mock của Google Drive service
type MockDriveService struct {
	mock.Mock
}

// TestFileLocalCustomOperations kiểm tra các thao tác cơ bản với File Local Custom
func TestFileLocalCustomOperations(t *testing.T) {
	// Tạo thư mục tạm thời cho test
	tempDir, err := ioutil.TempDir("", "file-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Tạo cấu hình File Local Custom
	config := &storage.CustomFile{
		PoolSize:             4,
		RootServiceDirectory: tempDir + "/",
	}

	// Khởi tạo File Local Custom client
	factory := storage.New(context.Background(), storage.FILE)
	fileClient := factory(storage.CUSTOMFILE, &storage.Config{
		CustomFile: *config,
	})

	// Kiểm tra client không nil
	assert.NotNil(t, fileClient, "File Local Custom client should not be nil")

	// Ép kiểu về interface IFILE
	client, ok := fileClient.(storage.IFILE)
	assert.True(t, ok, "Should be able to cast to IFILE")

	// Test CreateFolder
	folderName := "test-folder"
	_, err = client.CreateFolder(folderName)
	assert.NoError(t, err, "CreateFolder should not return an error")

	// Kiểm tra thư mục đã được tạo
	_, err = os.Stat(tempDir + "/" + folderName)
	assert.NoError(t, err, "Folder should exist")

	// Test Upload
	fileName := "test-file.txt"
	fileContent := []byte("This is a test file content")
	_, err = client.Upload(fileName, bytes.NewReader(fileContent))
	assert.NoError(t, err, "Upload should not return an error")

	// Kiểm tra file đã được tạo
	_, err = os.Stat(tempDir + "/" + fileName)
	assert.NoError(t, err, "File should exist")

	// Test List
	files, err := client.List(10)
	assert.NoError(t, err, "List should not return an error")
	assert.NotNil(t, files, "List should return non-nil result")

	// Test GetMetaData
	metadata, err := client.GetMetaData(tempDir + "/" + fileName)
	assert.NoError(t, err, "GetMetaData should not return an error")
	assert.NotNil(t, metadata, "GetMetaData should return non-nil result")

	// Test Move
	newFolderPath := tempDir + "/" + folderName + "/"
	_, err = client.Move(fileName, tempDir+"/", newFolderPath)
	assert.NoError(t, err, "Move should not return an error")

	// Kiểm tra file đã được di chuyển
	_, err = os.Stat(newFolderPath + fileName)
	assert.NoError(t, err, "File should exist in new location")
	_, err = os.Stat(tempDir + "/" + fileName)
	assert.Error(t, err, "File should not exist in old location")

	// Test Delete
	err = client.Delete([]string{newFolderPath + fileName})
	assert.NoError(t, err, "Delete should not return an error")

	// Kiểm tra file đã bị xóa
	_, err = os.Stat(newFolderPath + fileName)
	assert.Error(t, err, "File should not exist after deletion")
}
