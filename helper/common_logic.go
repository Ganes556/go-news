package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/news/pkg"
)

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}


func JSONStringify(d any) string {
	dd, _ := json.Marshal(d)
	return string(dd)
}

func GenerateUniqueFileName(fileName string) string {
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))
	// Use a combination of base name, timestamp, and a random string to ensure uniqueness
	randomString := RandomString(6) // Change the length of the random string as needed
	return fmt.Sprintf("%s_%d_%s", fileName, time.Now().UnixNano(), randomString)
}


func SaveInLocal(file *multipart.FileHeader) (string, error) {
	
	if file.Size > 5*1024*1024 { // 5 MB
		return "", errors.New("invalid size")
	}

	fi, err := file.Open()
	if err != nil {
		return "",err
	}
	defer fi.Close() // Close the file

	fReader, err := pkg.Decode(fi)

	if err != nil {
		if err == image.ErrFormat {
			return "",errors.New("invalid format")
		}
		return "",err
	}

	fileName := GenerateUniqueFileName(file.Filename) + ".webp"
	filePath := filepath.Join("./public/img/", fileName)
	fo, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	defer fo.Close()

	if _, err := io.Copy(fo, fReader); err != nil {
		return "", err
	}

	return fileName, nil
}

func DeleteMultiFileInLocal(fileName string) error {
	return os.Remove("./public/img/" + fileName)
}