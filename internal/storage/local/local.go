package local

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type LocalStorage struct {
	Path string
}

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func New(path string) *LocalStorage {
	fileInfo, err := os.Stat(path)

	if err == nil && !fileInfo.IsDir() {
		log.Fatal("Should be directory")
	} else if os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	}
	return &LocalStorage{Path: path}
}

func (storage *LocalStorage) Upload(filePath string) (string, error) {
	srcFile, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer srcFile.Close()

	destPath := filepath.Join(storage.Path, filepath.Base(filePath))
	destFile, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return "", err
	}

	return destPath, nil
}

func (storage *LocalStorage) Download(storagePath string, downloadPath string) (string, error) {
	srcFile, err := os.Open(storagePath)
	if err != nil {
		return "", err
	}
	defer srcFile.Close()

	// Create necessary directories if they do not exist
	if err := os.MkdirAll(filepath.Dir(downloadPath), 0755); err != nil {
		return "", err
	}

	// If the provided path is a directory, append the file name to the path
	info, err := os.Stat(downloadPath)
	if info.IsDir() {
		downloadPath = filepath.Join(downloadPath, filepath.Base(storagePath))
	} else if err != nil {
		log.Fatal(err)
	}

	dstFile, err := os.Create(downloadPath)
	if err != nil {
		return "", err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return "", err
	}

	return downloadPath, nil
}

// return file location and error
func (storage *LocalStorage) SaveImageFromRequest(file multipart.File, handler *multipart.FileHeader) (string, error) {
	// generate random string to store file
	randString := randomString(5)
	// Create a new file in the local filesystem
	fpath := filepath.Join(storage.Path, randString, handler.Filename)
	dst, err := os.Create(fpath)
	if err != nil {
		return "", fmt.Errorf("failed to write file to local system: %s", err)
	}
	defer dst.Close()

	// Copy the uploaded file to the filesystem at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy the uploaded file to the system: %s", err)
	}
	return fpath, nil
}

// return base64 image from storage
func (storage *LocalStorage) ConvertImageFromStorage(storagePath string) (string, error) {
	imgData, err := ioutil.ReadFile(storagePath)
	if err != nil {
		return "", err
	}

	// Convert the image data to base64.
	imgBase64 := base64.StdEncoding.EncodeToString(imgData)
	return imgBase64, nil
}

func (storage *LocalStorage) Delete(storagePath string) error {
	err := os.Remove(storagePath)
	return err
}
