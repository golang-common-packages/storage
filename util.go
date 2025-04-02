package storage

import (
	"context"
	"fmt"
	"hash/fnv"
	"io"
	"io/ioutil"
)

// SetContext sets a new context for the package
func SetContext(context context.Context) {
	if context != nil {
		ctx = context
	}
}

// GetContext returns the current context
func GetContext() context.Context {
	return ctx
}

// generateKey creates a unique hash key from a string
func generateKey(data string) string {
	if data == "" {
		return ""
	}
	
	hash := fnv.New64a()
	_, _ = hash.Write([]byte(data))
	return fmt.Sprint(hash.Sum64())
}

// streamToByte converts an io.Reader to a byte slice
// Note: This function reads the entire stream into memory,
// so it should be used with caution for large streams
func streamToByte(stream io.Reader) ([]byte, error) {
	if stream == nil {
		return nil, nil
	}
	return ioutil.ReadAll(stream)
}

// streamToString converts an io.Reader to a string
// Note: This function reads the entire stream into memory,
// so it should be used with caution for large streams
func streamToString(stream io.Reader) (string, error) {
	if stream == nil {
		return "", nil
	}
	
	bytes, err := ioutil.ReadAll(stream)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
