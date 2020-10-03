package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gammazero/workerpool"
	multierror "github.com/hashicorp/go-multierror"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"

	"github.com/golang-common-packages/hash"
)

// DriveServices manage all drive action
type DriveServices struct {
	driveService *drive.Service
	config       *GoogleDrive
}

var (
	// driveClientSessionMapping singleton pattern
	driveClientSessionMapping = make(map[string]*DriveServices)
)

// NewDrive init new instance
func NewDrive(config *GoogleDrive) IFILE {
	hasher := &hash.Client{}
	configAsJSON, err := json.Marshal(config)
	if err != nil {
		log.Fatalln("Unable to marshal Drive configuration: ", err)
	}
	configAsString := hasher.SHA1(string(configAsJSON))

	currentDriveSession := driveClientSessionMapping[configAsString]
	if currentDriveSession == nil {
		currentDriveSession = &DriveServices{nil, nil}

		if config.ByHTTPClient {
			b, err := ioutil.ReadFile(config.Credential)
			if err != nil {
				log.Fatalln("Unable to read client secret file: ", err)
			}

			// If modifying these scopes, delete your previously saved token.json.
			oauth2Config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope)
			if err != nil {
				log.Fatalln("Unable to parse client secret file to config: ", err)
			}
			client := getClient(oauth2Config, config.Token)

			srv, err := drive.New(client)
			if err != nil {
				log.Fatalln("Unable to retrieve Drive client: ", err)
			}

			currentDriveSession.driveService = srv
			currentDriveSession.config = config
			driveClientSessionMapping[configAsString] = currentDriveSession

		} else {
			os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", config.Credential)
			srv, err := drive.NewService(ctx)
			if err != nil {
				log.Fatalln("Unable to retrieve Drive client: ", err)
			}

			currentDriveSession.driveService = srv
			currentDriveSession.config = config
			driveClientSessionMapping[configAsString] = currentDriveSession
		}

		log.Println("Connected to Google Drive")
	}

	return currentDriveSession
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config, tokFile string) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalln("Unable to read authorization code: ", err)
	}

	token, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalln("Unable to retrieve token from web: ", err)
	}
	return token
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalln("Unable to cache oauth token: ", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// List all files based on pageSize
func (dr *DriveServices) List(pageSize int64, pageToken ...string) (interface{}, error) {
	var fields googleapi.Field = "nextPageToken, files(id, name, fileExtension, mimeType, parents)"

	if len(pageToken) == 0 {
		return dr.driveService.Files.List().PageSize(pageSize).Fields(fields).Do()
	}

	return dr.driveService.Files.List().PageToken(pageToken[0]).PageSize(pageSize).Fields(fields).Do()
}

// GetMetaData from file based on fileID
func (dr *DriveServices) GetMetaData(fileID string) (interface{}, error) {
	return dr.driveService.Files.Get(fileID).Do()
}

// CreateFolder on drive
func (dr *DriveServices) CreateFolder(name string, parents ...string) (interface{}, error) {
	f := &drive.File{
		Name:     name, //should specify a file extension in the name, like Name: "cat.jpg"
		MimeType: "application/vnd.google-apps.folder",
		Parents:  parents,
	}

	return dr.driveService.Files.Create(f).Do()
}

// Upload file to drive
func (dr *DriveServices) Upload(name string, fileContent io.Reader, parents ...string) (interface{}, error) {
	f := &drive.File{
		Name:    name, //should specify a file extension in the name, like Name: "cat.jpg"
		Parents: parents,
	}

	return dr.driveService.Files.Create(f).Media(fileContent).Do()
}

// Download file based on fileID
func (dr *DriveServices) Download(fileID string) (interface{}, error) {
	return dr.driveService.Files.Get(fileID).Download()
}

// Move file to new location based on fileID, oldParentID, newParentID
func (dr *DriveServices) Move(fileID, oldParentID, newParentID string) (interface{}, error) {
	if _, err := dr.driveService.Files.Update(fileID, nil).RemoveParents(oldParentID).Do(); err != nil {
		log.Println("Unable to move file: ", err)
		return nil, err
	}

	return dr.driveService.Files.Update(fileID, nil).AddParents(newParentID).Do()
}

// Delete file/folder based on IDs
func (dr *DriveServices) Delete(fileIDs []string) error {
	var mu sync.Mutex
	var errs *multierror.Error
	dwp := workerpool.New(dr.config.PoolSize)

	for _, fileID := range fileIDs {
		fileID := fileID
		dwp.Submit(func() {
			if err := dr.driveService.Files.Delete(fileID).Do(); err != nil {
				mu.Lock()
				errs = multierror.Append(errs, err)
				mu.Unlock()
			}
		})
	}

	dwp.StopWait()

	// Return an error if any failed
	if err := errs.ErrorOrNil(); err != nil {
		return err
	}

	return nil
}
