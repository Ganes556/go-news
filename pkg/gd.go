package pkg

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/sync/errgroup"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type GoogleDrive interface {
	
}

// googleDriveClient represents a client for interacting with Google Drive.
type googleDriveClient struct {
    service *drive.Service
}

// NewGoogleDriveClient creates a new GoogleDriveClient.
func NewGoogleDriveClient() (GoogleDrive, error) {
    ctx := context.Background()
	credJSON, err := base64.StdEncoding.DecodeString(os.Getenv("GD_SECRET"))
	if err != nil {
		return nil, err
	}
    // Configure OAuth2 with the credentials from the file.
    config, err := google.ConfigFromJSON(credJSON, drive.DriveFileScope)
    if err != nil {
        return nil, fmt.Errorf("unable to parse client secret file to config: %v", err)
    }

    // Create a Drive service.
    client := getClient(ctx, config)
    srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
    if err != nil {
        return nil, fmt.Errorf("unable to retrieve Drive client: %v", err)
    }

    return &googleDriveClient{service: srv}, nil
}

// getClient retrieves a token, saves the token, and returns the generated client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
    tokFile := "token.json"
    tok, err := tokenFromFile(tokFile)
    if err != nil {
        tok = getTokenFromWeb(config)
        saveToken(tokFile, tok)
    }
    return config.Client(ctx, tok)
}

// getTokenFromWeb requests a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
    authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
    fmt.Printf("Go to the following link in your browser then type the "+
        "authorization code: \n%v\n", authURL)

    var authCode string
    if _, err := fmt.Scan(&authCode); err != nil {
        log.Fatalf("Unable to read authorization code: %v", err)
    }

    tok, err := config.Exchange(context.Background(), authCode)
    if err != nil {
        log.Fatalf("Unable to retrieve token from web: %v", err)
    }
    return tok
}

// tokenFromFile retrieves a token from a file.
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

// saveToken saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
    fmt.Printf("Saving credential file to: %s\n", path)
    f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
    if err != nil {
        log.Fatalf("Unable to save token: %v", err)
    }
    defer f.Close()
    json.NewEncoder(f).Encode(token)
}

// UploadFile uploads a file to Google Drive.
func (gd *googleDriveClient) UploadFile(fileHs []*multipart.FileHeader) error {
    g := new(errgroup.Group)
	g.SetLimit(10)
	
    for _, fileH := range fileHs {
        fileH := fileH // Important: create a new variable to avoid concurrency issues
        g.Go(func() error {
            file, err := fileH.Open()
            if err != nil {
                return fmt.Errorf("unable to open file: %v", err)
            }
            defer file.Close()

            f := &drive.File{
                Name: fileH.Filename,
            }
            _, err = gd.service.Files.Create(f).Media(file).Do()
            if err != nil {
                return fmt.Errorf("unable to upload file: %v", err)
            }

            fmt.Printf("File '%s' uploaded successfully.\n", fileH.Filename)
            return nil
        })
    }

    if err := g.Wait(); err != nil {
        return err
    }

    return nil
}

// func main() {
//     // Replace 'credentials.json' with your credentials file.
//     client, err := NewGoogleDriveClient("credentials.json")
//     if err != nil {
//         log.Fatalf("Failed to create Google Drive client: %v", err)
//     }

//     // Replace 'example.txt' with the name of the file you want to upload,
//     // and 'path/to/file.txt' with the path to the file you want to upload.
//     err = client.UploadFile("example.txt", "path/to/file.txt")
//     if err != nil {
//         log.Fatalf("Failed to upload file: %v", err)
//     }
// }
