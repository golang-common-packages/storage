# Storage

[![Go Report Card](https://goreportcard.com/badge/github.com/golang-common-packages/storage)](https://goreportcard.com/report/github.com/golang-common-packages/storage)
[![GoDoc](https://godoc.org/github.com/golang-common-packages/storage?status.svg)](https://godoc.org/github.com/golang-common-packages/storage)
[![License](https://img.shields.io/github/license/golang-common-packages/storage)](LICENSE)

The Storage library provides a unified interface for various storage types, including SQL, NoSQL, and File. The library uses the Abstract Factory design pattern to create different storage objects.

## Installation

```bash
go get github.com/golang-common-packages/storage
```

## Features

- **SQL Relational**: Support for SQL databases through Go's `database/sql` package
- **NoSQL Document**: Support for MongoDB
- **NoSQL Key-Value**: Support for Redis, BigCache, and custom implementations
- **File**: Support for Google Drive and custom implementations

## Usage

### Import

```go
import "github.com/golang-common-packages/storage"
```

### Working with MongoDB

```go
// Initialize MongoDB client
mongoClient := storage.New(context.Background(), storage.NOSQLDOCUMENT)(storage.MONGODB, &storage.Config{
    MongoDB: storage.MongoDB{
        User:     "USERNAME",
        Password: "PASSWORD",
        Hosts:    []string{"localhost:27017"},
        Options:  []string{},
        DB:       "DATABASE_NAME",
    },
}).(storage.INoSQLDocument)

// Create document
documents := []interface{}{
    map[string]interface{}{
        "name": "John Doe",
        "age":  30,
    },
}
result, err := mongoClient.Create("database", "collection", documents)

// Read document
filter := bson.M{"name": "John Doe"}
result, err := mongoClient.Read("database", "collection", filter, 10, reflect.TypeOf(YourModel{}))

// Update document
filter := bson.M{"name": "John Doe"}
update := bson.M{"$set": bson.M{"age": 31}}
result, err := mongoClient.Update("database", "collection", filter, update)

// Delete document
filter := bson.M{"name": "John Doe"}
result, err := mongoClient.Delete("database", "collection", filter)
```

### Working with Redis

```go
// Initialize Redis client
redisClient := storage.New(context.Background(), storage.NOSQLKEYVALUE)(storage.REDIS, &storage.Config{
    Redis: storage.Redis{
        Host:       "localhost:6379",
        Password:   "PASSWORD",
        DB:         0,
        MaxRetries: 3,
    },
}).(storage.INoSQLKeyValue)

// Store value
err := redisClient.Set("key", "value", 1*time.Hour)

// Get value
value, err := redisClient.Get("key")

// Update value
err := redisClient.Update("key", "new-value", 1*time.Hour)

// Delete value
err := redisClient.Delete("key")
```

### Working with Google Drive

```go
// Initialize Google Drive client
driveClient := storage.New(context.Background(), storage.FILE)(storage.DRIVE, &storage.Config{
    GoogleDrive: storage.GoogleDrive{
        PoolSize:     4,
        ByHTTPClient: false,
        Credential:   "credentials.json",
        Token:        "token.json",
    },
}).(storage.IFILE)

// List files
files, err := driveClient.List(100)

// Create folder
folder, err := driveClient.CreateFolder("folder-name")

// Upload file
file, err := os.Open("file.txt")
result, err := driveClient.Upload("file.txt", file)

// Download file
result, err := driveClient.Download("file-id")

// Move file
result, err := driveClient.Move("file-id", "old-parent-id", "new-parent-id")

// Delete file
err := driveClient.Delete([]string{"file-id"})
```

## Development

### Requirements

- Go 1.15+
- MongoDB (for NoSQL Document)
- Redis (for NoSQL Key-Value)
- SQLite (for SQL Relational)

### Running tests

```bash
# Run all tests
make test

# Run short tests (skip tests requiring database connections)
make test-short

# Run tests with coverage
make coverage
```

### Docker

```bash
# Build Docker image
make docker-build

# Run Docker container
make docker-run

# Run with Docker Compose
docker-compose up
```

## Contributing

We welcome contributions from the community. Please create an issue or pull request on GitHub.

## License

This project is distributed under the MIT license. See the [LICENSE](LICENSE) file for more details.
