package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/golang-common-packages/storage"
)

func main() {
	// Init file services
	fileService := storage.New(storage.FILE)(storage.DRIVE, &storage.Config{GoogleDrive: storage.GoogleDrive{
		PoolSize:     4,
		ByHTTPClient: false,
		Credential:   "credentials.json",
		Token:        "token.json",
	}}).(storage.IFILE)

	// List
	fmt.Println("List:")
	fileListResult, err := fileService.List(100)
	if err != nil {
		log.Fatalln("Unable to retrieve files: %v", err)
	}

	fileListJSON, _ := json.Marshal(fileListResult)
	var fileList storage.GoogleFileListModel
	if err := json.Unmarshal(fileListJSON, &fileList); err != nil {
		log.Fatalln("Unable to unmarshal: %v", err)
	}

	fmt.Println("NextPageToken:")
	fmt.Println(fileList.NextPageToken)

	if len(fileList.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		for _, i := range fileList.Files {
			fmt.Printf("%s\n -- MineType (%s)\n -- ID (%s)\n -- Parents (%s)\n -- FileExtension (%s)\n", i.Name, i.MimeType, i.Id, i.Parents, i.FileExtension)
		}
	}

	// Delete
	fmt.Println("Delete:")
	fileIDs := []string{"1RkH0-SJT0RJWtxNUUZU1iwv-LW7FUTh7", "1YjynEQJGvq_JnAAxC9pvHmGF6nysF8Yi", "1RkH0-SJT0RJWtxNUUZU1iwv-LW7FUTh8"}
	if err := fileService.Delete(fileIDs); err != nil {
		log.Fatalln("Can not delete: %v", err)
	}

	// Create folder
	fmt.Println("Create:")
	createFolderResult, err := fileService.CreateFolder("golangDemo2")
	if err != nil {
		log.Fatalln("Can not create folder: %v", err)
	}

	createFolderResultJSON, _ := json.Marshal(createFolderResult)
	var folderCreated storage.GoogleFileModel
	if err := json.Unmarshal(createFolderResultJSON, &folderCreated); err != nil {
		log.Fatalln("Unable to unmarshal: %v", err)
	}

	fmt.Println(folderCreated.Id)

	// Move
	fmt.Println("Move:")
	moveResult, err := fileService.Move("1iWusgOJ7yi0Jmx6O4D8L1HK_lw9SXMk4", "1a0MqYDGQ6RQDE_gxz-72-imRgCDyywUn", "1aSwrAe04wy_Bp6eyPRlCLu2_R_m7aF9h")
	if err != nil {
		log.Fatalln("Can not move file: %v", err)
	}

	fmt.Println(moveResult)

	// Upload
	fmt.Println("Upload:")
	f, err := os.Open("./test.txt")
	if err != nil {
		log.Fatalln("%v", err)
	}

	uploadResult, err := fileService.Upload("test.txt", f)
	if err != nil {
		log.Fatalln("Can not upload: %v", err)
	}

	uploadResultJSON, _ := json.Marshal(uploadResult)
	var fileUpload storage.GoogleFileModel
	if err := json.Unmarshal(uploadResultJSON, &fileUpload); err != nil {
		log.Fatalln("Unable to unmarshal: %v", err)
	}

	fmt.Println(fileUpload.Id)
}
