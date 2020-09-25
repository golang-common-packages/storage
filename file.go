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

// NewFile Factory Pattern
func NewFile(
	databaseCompany int,
	config *Config) interface{} {

	switch databaseCompany {
	case DRIVE:
		return NewDrive(&config.GoogleDrive)
	case CUSTOMFILE:
		return NewCustomFile(&config.CustomFile)
	}

	return nil
}
