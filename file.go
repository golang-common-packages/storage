package storage

import "io"

// IFILE factory pattern interface
type IFILE interface {
	List(pageSize int64, pageToken ...string) (interface{}, error)
	GetMetaData(fileID string) (interface{}, error)
	CreateFolder(name string, parents ...string) (interface{}, error)
	Upload(name string, fileContent io.Reader, parents ...string) (interface{}, error)
	Download(fileID string) (interface{}, error)
	Move(fileID, oldParentID, newParentID string) (interface{}, error)
	Delete(fileIDs []string) error
}

const (
	// DRIVE cloud services
	DRIVE = iota
	// CUSTOMFILE file services
	CUSTOMFILE
)

// newFile Factory Pattern
func newFile(
	databaseCompany int,
	config *Config) interface{} {

	switch databaseCompany {
	case DRIVE:
		return newDrive(&config.GoogleDrive)
	case CUSTOMFILE:
		return newCustomFile(&config.CustomFile)
	}

	return nil
}
