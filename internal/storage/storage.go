package storage

import "mime/multipart"

// interface to keep files
type Storage interface {
	Delete(storagePath string) error
	Upload(filePath string) (string, error)                                                  // return url address to file storage
	Download(storagePath string, downloadPath string) (string, error)                        // return file location
	SaveImageFromRequest(file multipart.File, handler *multipart.FileHeader) (string, error) // return file location
	ConvertImageFromStorage(storagePath string) (string, error)                              // image -> base64
}
