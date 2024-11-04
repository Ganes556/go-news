package pkg

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	storage "cloud.google.com/go/storage"
	"math/rand"
	"golang.org/x/sync/errgroup"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var (
	InvalidImageFormat = "invalid image format"
	InvalidImageSize   = "image size should be less than 5mb"
)

type Gcloud interface {
	Upload2Storage(ctx context.Context, folderName string, fileHs []*multipart.FileHeader) ([]string, error)
	Upload2StorageWithouCompress(ctx context.Context, folderName, fileName string, file io.Reader) (string, error)
	DeleteInStorage(ctx context.Context, objNames []string) error
	Update2Storage(ctx context.Context, folderName string, newFileH []*multipart.FileHeader, oldObjNames []string) ([]string, error)
	GetQ(q *storage.Query) []string
}

type gcloud struct {
	client *storage.Client
	bucket *storage.BucketHandle
}

func NewGcloud(duration *time.Duration, bucketName string) Gcloud {
	ctxx := context.Background()
	var ctx context.Context
	var cancel context.CancelFunc
	if duration != nil {
		ctx, cancel = context.WithTimeout(ctxx, *duration)
	} else {
		ctx, cancel = context.WithTimeout(ctxx, time.Second*10)
	}

	defer cancel()

	gPrivateKey := os.Getenv("GC_SECRET")
	if gPrivateKey == "" {
		fmt.Println("empty env GCP")
		os.Exit(1)
	}

	privateKey, err := base64.StdEncoding.DecodeString(gPrivateKey)
	if err != nil {
		panic(err)
	}

	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(privateKey))
	if err != nil {
		panic(err)
	}
	b := client.Bucket(bucketName)
	return &gcloud{client, b}
}

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func GenerateUniqueObjectName(baseName, fileName string) string {
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))
	// Use a combination of base name, timestamp, and a random string to ensure uniqueness
	randomString := RandomString(6) // Change the length of the random string as needed
	return fmt.Sprintf("%s/%s_%d_%s", baseName, fileName, time.Now().UnixNano(), randomString)
}

func (g *gcloud) Upload2Storage(ctx context.Context, folderName string, fileHs []*multipart.FileHeader) ([]string, error) {
	if len(fileHs) == 0 {
		return nil, nil
	}
	gr := new(errgroup.Group)
	gr.SetLimit(10)
	var objNames = make([]string, len(fileHs))
	for i, fileH := range fileHs {
		i, fileH := i, fileH // Capture loop variables
		gr.Go(func() error {
			if fileH.Size > 5*1024*1024 { // 5 MB
				return errors.New(InvalidImageSize)
			}
			file, err := fileH.Open()
			if err != nil {
				return err
			}
			fReader, err := Decode(file)
			if err != nil {
				if err == image.ErrFormat {
					return errors.New(InvalidImageFormat)
				}
				return err
			}
			objName, err := g.Upload2StorageWithouCompress(ctx, folderName, fileH.Filename, fReader)
			if err != nil {
				return err
			}
			objNames[i] = objName
			return nil
		})
	}

	if err := gr.Wait(); err != nil {
		return nil, err
	}

	return objNames, nil
}

func (g *gcloud) Upload2StorageWithouCompress(ctx context.Context, folderName, fileName string, file io.Reader) (string, error) {
	objName := GenerateUniqueObjectName(folderName, fileName) + ".webp"
	obj := g.bucket.Object(objName)
	w := obj.NewWriter(ctx)

	if _, err := io.Copy(w, file); err != nil {
		panic(err)
	}

	if err := w.Close(); err != nil {
		panic(err)
	}

	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", err
	}
	return objName, nil
}

func (g *gcloud) DeleteInStorage(ctx context.Context, objNames []string) error {
	if len(objNames) == 0 {
		return nil
	}
	gr := new(errgroup.Group)
	gr.SetLimit(10)
	for _, objName := range objNames {
		objName := objName
		gr.Go(func() error {
			obj := g.bucket.Object(objName)
			if err := obj.Delete(ctx); err != nil {
				return err
			}
			return nil
		})
	}
	if err := gr.Wait(); err != nil {
		return err
	}
	return nil
}

// check if the length of the new file is not same delete if less and upload a new file if greater from the old one
func (g *gcloud) Update2Storage(ctx context.Context, folderName string, newFileH []*multipart.FileHeader, oldObjNames []string) ([]string, error) {
	if len(newFileH) == 0 {
		return nil, nil
	}
	gr := new(errgroup.Group)
	gr.SetLimit(3)
	gr.Go(func() error {
		err := g.DeleteInStorage(ctx, oldObjNames)
		if err != nil {
			panic(err)
		}
		return nil
	})
	var newObjNames []string
	gr.Go(func() error {
		var err error
		newObjNames, err = g.Upload2Storage(ctx, folderName, newFileH)
		if err != nil {
			panic(err)
		}
		return nil
	})

	if err := gr.Wait(); err != nil {
		return nil, err
	}

	return newObjNames, nil
}

func (g *gcloud) GetQ(q *storage.Query) []string {
	ctxx := context.Background()
	ctx, cancel := context.WithTimeout(ctxx, time.Second*10)
	defer cancel()
	it := g.bucket.Objects(ctx, q)
	var list []string
	for {
		objectAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Error iterating objects: %v", err)
		}

		// Print the name of the object
		list = append(list, objectAttrs.Name)
	}
	return list
}
