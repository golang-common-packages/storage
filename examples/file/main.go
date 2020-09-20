package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/golang-common-packages/storage"
)

func main() {
	fileService := storage.New(storage.FILE)(storage.DRIVE, &storage.Config{GoogleDrive: storage.GoogleDrive{
		Credential: "credentials.json",
		Token:      "token.json",
	}}).(storage.IFILE)

	fileListResult, err := fileService.List(10)
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}

	fileListJSON, _ := json.Marshal(fileListResult)
	var fileList storage.GoogleFileListModel

	if err := json.Unmarshal(fileListJSON, &fileList); err != nil {
		log.Fatalf("Unable to unmarshal: %v", err)
	}

	fmt.Println("Files:")
	if len(fileList.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		for _, i := range fileList.Files {
			fmt.Printf("%s --- (%s)\n", i.Name, i.Id)
		}
	}
}
