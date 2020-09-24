package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/golang-common-packages/storage"
)

func main() {
	fileService := storage.New(storage.FILE)(storage.DRIVE, &storage.Config{GoogleDrive: storage.GoogleDrive{
		PoolSize:     4,
		ByHTTPClient: false,
		Credential:   "credentials.json",
		Token:        "token.json",
	}}).(storage.IFILE)

	//Delete
	fileIDs := []string{""}
	if err := fileService.Delete(fileIDs); err != nil {
		log.Fatalf("Can not delete this file: %v", err)
	}

	// List
	fileListResult, err := fileService.List(100)
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}

	fileListJSON, _ := json.Marshal(fileListResult)
	var fileList storage.GoogleFileListModel
	if err := json.Unmarshal(fileListJSON, &fileList); err != nil {
		log.Fatalf("Unable to unmarshal: %v", err)
	}

	fmt.Println("NextPageToken:")
	fmt.Println(fileList.NextPageToken)

	fmt.Println("Files:")
	if len(fileList.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		for _, i := range fileList.Files {
			fmt.Printf("%s -- (%s) -- (%s) -- (%s)\n", i.Name, i.MimeType, i.Id, i.FileExtension)
		}
	}

	// Upload file
	f, err := os.Open("./test.txt")
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Println("Upload:")
	uploadResult, err := fileService.Upload("test.txt", f)
	if err != nil {
		log.Fatalf("Can not upload: %v", err)
	}

	uploadResultJSON, _ := json.Marshal(uploadResult)
	var fileUpload storage.GoogleFileModel
	if err := json.Unmarshal(uploadResultJSON, &fileUpload); err != nil {
		log.Fatalf("Unable to unmarshal: %v", err)
	}

	fmt.Println(fileUpload.Id)

	//Create folder
	createFolderResult, err := fileService.CreateFolder("golangDemo")
	if err != nil {
		log.Fatalf("Can not create folder: %v", err)
	}

	createFolderResultJSON, _ := json.Marshal(createFolderResult)
	var folderCreated storage.GoogleFileModel
	if err := json.Unmarshal(createFolderResultJSON, &folderCreated); err != nil {
		log.Fatalf("Unable to unmarshal: %v", err)
	}

	fmt.Println(folderCreated.Id)
}
