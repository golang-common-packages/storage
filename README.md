# Storage
[![Go Report Card](https://goreportcard.com/badge/github.com/golang-common-packages/storage)](https://goreportcard.com/report/github.com/golang-common-packages/storage)

```go
import "github.com/golang-common-packages/storage"
```

```go
dbConn = storage.New(storage.NOSQLDOCUMENT)(storage.MONGODB, &storage.Config{MongoDB: storage.MongoDB{
		User:     "USERNAME",
		Password: "PASSWORD",
		Hosts:    "STRING_URI_SLICE",
		Options:  "STRING_OPTION_SLICE",
		DB:       "DATABASE_NAME",
	}}).(storage.INoSQLDocument)
```

## Note
[Check this template for more information and how to use this package](https://github.com/golang-common-packages/template)