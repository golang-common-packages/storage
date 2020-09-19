package storage

import "context"

// IFILE factory pattern interface
type IFILE interface {
	List(pageSize int64, pageToken ...string) (interface{}, error)
	Upload(fileModel interface{}) (interface{}, error)
	Download(fileModel *DriveFileModel) (interface{}, error)
	Delete(fileModel *DriveFileModel) error
}

var (
	ctx = context.Background()
)

/*
	@DRIVE: Google Drive
*/
const (
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
