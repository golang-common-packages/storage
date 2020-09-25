package storage

import (
	"encoding/json"
	"io"

	"github.com/golang-common-packages/hash"
)

// CustomFileClient manage all custom file action
type CustomFileClient struct {
	config *CustomFile
}

var (
	// customFileClientSessionMapping singleton pattern
	customFileClientSessionMapping = make(map[string]*CustomFileClient)
)

// NewCustomFile init new instance
func NewCustomFile(config *CustomFile) IFILE {
	hasher := &hash.Client{}
	configAsJSON, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	configAsString := hasher.SHA1(string(configAsJSON))

	currentCustomFileClientSession := customFileClientSessionMapping[configAsString]
	if currentCustomFileClientSession == nil {
		currentCustomFileClientSession = &CustomFileClient{config: config}
		customFileClientSessionMapping[configAsString] = currentCustomFileClientSession
	}

	return currentCustomFileClientSession
}

// List all files based on pageSize
func (cf *CustomFileClient) List(pageSize int64, pageToken ...string) (interface{}, error) {
	return nil, nil
}

// GetMetaData from file based on fileID
func (cf *CustomFileClient) GetMetaData(fileID string) (interface{}, error) {
	return nil, nil
}

// CreateFolder on drive
func (cf *CustomFileClient) CreateFolder(name string, parents ...string) (interface{}, error) {
	return nil, nil
}

// Upload file to drive
func (cf *CustomFileClient) Upload(name string, fileContent io.Reader, parents ...string) (interface{}, error) {
	return nil, nil
}

// Download file based on fileID
func (cf *CustomFileClient) Download(fileID string) (interface{}, error) {
	return nil, nil
}

// Move file to new location based on fileID, oldParentID, newParentID
func (cf *CustomFileClient) Move(fileID, oldParentID, newParentID string) (interface{}, error) {
	return nil, nil
}

// Delete file/folder based on IDs
func (cf *CustomFileClient) Delete(fileIDs []string) error {
	return nil
}
