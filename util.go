package storage

import (
	"bytes"
	"hash/fnv"
	"io"
)

func generateKey(data string) string {
	hash := fnv.New64a()
	hash.Write([]byte(data))

	return string(hash.Sum64())
}

func streamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}

func streamToString(stream io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.String()
}
