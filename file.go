package storage

import "io"

// IFILE factory pattern interface
type IFILE interface {
	List(pageSize int64, pageToken ...string) (interface{}, error)
	Upload(name string, fileContent io.Reader, parents ...string) (interface{}, error)
	CreateFolder(name string, parents ...string) (interface{}, error)
	Download(fileModel interface{}) (interface{}, error)
	Move(oldParentID, newParentID string, fileModel interface{}) (interface{}, error)
	Delete(fileIDs []string) error
}

const (
	// DRIVE cloud services
	DRIVE = iota
)

// NewFile Factory Pattern
func NewFile(
	databaseCompany int,
	config *Config) interface{} {

	switch databaseCompany {
	case DRIVE:
		return NewDrive(&config.GoogleDrive)
	}

	return nil
}
