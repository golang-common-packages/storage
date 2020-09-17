package storage

import "context"

// IFILE factory pattern interface
type IFILE interface {
	List(fileModel *FileModel) (interface{}, error)
	Upload(fileModel *FileModel, parentID string) (interface{}, error)
	Download(fileModel *FileModel) (interface{}, error)
	Delete(fileModel *FileModel) error
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
		return NewDrive()
	}

	return nil
}
