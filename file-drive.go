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
		panic(err)
	}
	configAsString := hasher.SHA1(string(configAsJSON))

	currentDriveSession := driveClientSessionMapping[configAsString]
	if currentDriveSession == nil {
		currentDriveSession = &DriveServices{nil, nil}

		if config.ByHTTPClient {
			b, err := ioutil.ReadFile(config.Credential)
			if err != nil {
				log.Fatalf("Unable to read client secret file: %v", err)
			}

			// If modifying these scopes, delete your previously saved token.json.
			oauth2Config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope)
			if err != nil {
				log.Fatalf("Unable to parse client secret file to config: %v", err)
			}
			client := getClient(oauth2Config, config.Token)

			srv, err := drive.New(client)
			if err != nil {
				log.Fatalf("Unable to retrieve Drive client: %v", err)
			}

			currentDriveSession.driveService = srv
			currentDriveSession.config = config
			driveClientSessionMapping[configAsString] = currentDriveSession

		} else {
			os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", config.Credential)
			srv, err := drive.NewService(ctx)
			if err != nil {
				log.Fatalf("Unable to retrieve Drive client: %v", err)
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
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
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
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// List all files based on pageSize
func (dr *DriveServices) List(pageSize int64, pageToken ...string) (interface{}, error) {
	var fields googleapi.Field = "nextPageToken, files(id, name, fileExtension, mimeType)"

	if len(pageToken) == 0 {
		return dr.driveService.Files.List().PageSize(pageSize).Fields(fields).Do()
	}

	return dr.driveService.Files.List().PageToken(pageToken[0]).PageSize(pageSize).Fields(fields).Do()
}

// Upload a file to drive
func (dr *DriveServices) Upload(name string, fileContent io.Reader, parents ...string) (interface{}, error) {
	f := &drive.File{
		Name:    name, //should specify a file extension in the name, like Name: "cat.jpg"
		Parents: parents,
	}

	return dr.driveService.Files.Create(f).Media(fileContent).Do()
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

// Download a file based on fileID
func (dr *DriveServices) Download(fileModel interface{}) (interface{}, error) {
	return dr.driveService.Files.Get(fileModel.(*drive.File).Id).Download()
}

// Move a file to new location based on fileID
func (dr *DriveServices) Move(oldParentID, newParentID string, fileModel interface{}) (interface{}, error) {
	return dr.driveService.Files.Update(fileModel.(*drive.File).Id, fileModel.(*drive.File)).RemoveParents(oldParentID).AddParents(newParentID).Do()
}

// Delete a file/folder base on IDs
func (dr *DriveServices) Delete(fileIDs []string) error {
	var mu sync.Mutex
	var errs *multierror.Error
	dwp := workerpool.New(dr.config.PoolSize)

	for _, fileID := range fileIDs {
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
