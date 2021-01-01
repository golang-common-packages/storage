# Storage
[![Go Report Card](https://goreportcard.com/badge/github.com/golang-common-packages/storage)](https://goreportcard.com/report/github.com/golang-common-packages/storage)

```go
import "github.com/golang-common-packages/storage"
```

Working with [MongoDB](https://github.com/golang-common-packages/template/blob/master/main.go):

```go
dbConn = storage.New(storage.NOSQLDOCUMENT)(storage.MONGODB, &storage.Config{MongoDB: storage.MongoDB{
		User:     "USERNAME",
		Password: "PASSWORD",
		Hosts:    "STRING_URI_SLICE",
		Options:  "STRING_OPTION_SLICE",
		DB:       "DATABASE_NAME",
	}}).(storage.INoSQLDocument)
```

Working with [Google Drive](https://github.com/golang-common-packages/storage/blob/master/examples/file/main.go):

```go
driveConn := storage.New(storage.FILE)(storage.DRIVE, &storage.Config{GoogleDrive: storage.GoogleDrive{
		PoolSize:     4,
		ByHTTPClient: false,
		Credential:   "credentials.json",
		Token:        "token.json",
	}}).(storage.IFILE)
```

## Note
[How to use this package?](https://github.com/golang-common-packages/template)