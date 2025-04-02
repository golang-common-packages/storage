package storage

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"sync"

	"github.com/gammazero/workerpool"
	"github.com/golang-common-packages/hash"
	multierror "github.com/hashicorp/go-multierror"
)

// CustomFileClient manage all custom file action
type CustomFileClient struct {
	config *CustomFile
}

var (
	// customFileClientSessionMapping singleton pattern
	customFileClientSessionMapping = make(map[string]*CustomFileClient)
)

// newCustomFile init new instance
func newCustomFile(config *CustomFile) IFILE {
	hasher := &hash.Client{}
	configAsJSON, err := json.Marshal(config)
	if err != nil {
		log.Fatalln("Unable to marshal local custom file configuration: ", err)
	}
	configAsString := hasher.SHA1(string(configAsJSON))

	currentCustomFileClientSession := customFileClientSessionMapping[configAsString]
	if currentCustomFileClientSession == nil {
		currentCustomFileClientSession = &CustomFileClient{config: config}
		customFileClientSessionMapping[configAsString] = currentCustomFileClientSession
		log.Println("File custom is ready")
	}

	return currentCustomFileClientSession
}

// List all of file and folder
func (cf *CustomFileClient) List(pageSize int64, pageToken ...string) (interface{}, error) {
	f, err := os.Open(cf.config.RootServiceDirectory)
	if err != nil {
		log.Println("Unable to open root servie directory: ", err)
		return nil, err
	}
	files, err := f.Readdir(-1)
	defer f.Close()
	if err != nil {
		log.Println("Unable to read directory: ", err)
		return nil, err
	}

	return files, nil
}

// GetMetaData from file based on fileID (that is file name)
func (cf *CustomFileClient) GetMetaData(fileID string) (interface{}, error) {
	return os.Stat(fileID)
}

// CreateFolder on root service directory
func (cf *CustomFileClient) CreateFolder(name string, parents ...string) (interface{}, error) {
	return nil, os.MkdirAll(cf.config.RootServiceDirectory+name, os.ModePerm)
}

// Upload file to root service directory
func (cf *CustomFileClient) Upload(name string, fileContent io.Reader, parents ...string) (interface{}, error) {
	// Open file using READ & WRITE permission, and check if file exists
	var _, err = os.Stat(cf.config.RootServiceDirectory + name)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(cf.config.RootServiceDirectory + name)
		if err != nil {
			log.Println("Unable to create file: ", err)
			return nil, err
		}
		defer file.Close()

		// Write content to file.
		data, err := streamToByte(fileContent)
		if err != nil {
			log.Println("Unable to read file content: ", err)
			return nil, err
		}
		
		_, err = file.Write(data)
		if err != nil {
			log.Println("Unable to write to file: ", err)
			return nil, err
		}
	}

	return nil, nil
}

// Download feature doesn't implemented for this service
func (cf *CustomFileClient) Download(fileID string) (interface{}, error) {
	return nil, nil
}

// Move file to new location based on fileID, oldParentID, newParentID
func (cf *CustomFileClient) Move(fileID, oldParentID, newParentID string) (interface{}, error) {
	err := os.Rename(oldParentID+fileID, newParentID+fileID)
	if err != nil {
		log.Println("Unable to move file: ", err)
		return nil, err
	}
	return nil, nil
}

// Delete file/folder based on IDs (that is the list of file name)
func (cf *CustomFileClient) Delete(fileIDs []string) error {
	var mu sync.Mutex
	var errs *multierror.Error
	dwp := workerpool.New(cf.config.PoolSize)

	for _, fileID := range fileIDs {
		fileID := fileID
		dwp.Submit(func() {
			if err := os.Remove(fileID); err != nil {
				mu.Lock()
				errs = multierror.Append(errs, err)
				mu.Unlock()
			}
		})
	}

	dwp.StopWait()

	// Return an error if any failed
	if err := errs.ErrorOrNil(); err != nil {
		return err
	}

	return nil
}
