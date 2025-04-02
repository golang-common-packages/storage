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

type MockDriveService struct {
	mock.Mock
}

func TestFileLocalCustomOperations(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "file-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	config := &storage.CustomFile{
		PoolSize:             4,
		RootServiceDirectory: tempDir + "/",
	}

	factory := storage.New(context.Background(), storage.FILE)
	fileClient := factory(storage.CUSTOMFILE, &storage.Config{
		CustomFile: *config,
	})

	assert.NotNil(t, fileClient, "File Local Custom client should not be nil")

	client, ok := fileClient.(storage.IFILE)
	assert.True(t, ok, "Should be able to cast to IFILE")

	folderName := "test-folder"
	_, err = client.CreateFolder(folderName)
	assert.NoError(t, err, "CreateFolder should not return an error")

	_, err = os.Stat(tempDir + "/" + folderName)
	assert.NoError(t, err, "Folder should exist")

	fileName := "test-file.txt"
	fileContent := []byte("This is a test file content")
	_, err = client.Upload(fileName, bytes.NewReader(fileContent))
	assert.NoError(t, err, "Upload should not return an error")

	_, err = os.Stat(tempDir + "/" + fileName)
	assert.NoError(t, err, "File should exist")

	files, err := client.List(10)
	assert.NoError(t, err, "List should not return an error")
	assert.NotNil(t, files, "List should return non-nil result")

	metadata, err := client.GetMetaData(tempDir + "/" + fileName)
	assert.NoError(t, err, "GetMetaData should not return an error")
	assert.NotNil(t, metadata, "GetMetaData should return non-nil result")

	newFolderPath := tempDir + "/" + folderName + "/"
	_, err = client.Move(fileName, tempDir+"/", newFolderPath)
	assert.NoError(t, err, "Move should not return an error")

	_, err = os.Stat(newFolderPath + fileName)
	assert.NoError(t, err, "File should exist in new location")
	_, err = os.Stat(tempDir + "/" + fileName)
	assert.Error(t, err, "File should not exist in old location")

	err = client.Delete([]string{newFolderPath + fileName})
	assert.NoError(t, err, "Delete should not return an error")

	_, err = os.Stat(newFolderPath + fileName)
	assert.Error(t, err, "File should not exist after deletion")
}
