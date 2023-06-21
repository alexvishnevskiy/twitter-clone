package local

import (
	"log"
	"os"
	"testing"
)

func TestLocalStorage(t *testing.T) {
	storage := New("./storage")

	if _, err := os.Create("./test.txt"); err != nil {
		log.Fatal(err)
	}

	path, err := storage.Upload("./test.txt")
	if err != nil {
		t.Errorf("Failed to upload file")
	}
	if _, err = storage.Download(path, "./test_dir/"); err != nil {
		t.Errorf("Failed to download file")
	}

	// clean up everything
	if err = os.Remove("./test.txt"); err != nil {
		t.Error(err)
	}
	if err = os.RemoveAll("./storage"); err != nil {
		t.Error(err)
	}
	if err = os.RemoveAll("./test_dir"); err != nil {
		t.Error(err)
	}
}
