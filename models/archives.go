package models

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/google/uuid"
	"github.com/upb-code-labs/static-files-microservice/config"
)

func SaveArchive(directory string, file multipart.File) (saved_uuid string, err error) {
	uuid := uuid.New()

	// Create an empty file
	volumePath := config.GetEnvironment().ArchivesVolumePath
	path := fmt.Sprintf("%s/%s/%s.zip", volumePath, directory, uuid.String())

	emptyFile, err := os.Create(path)
	if err != nil {
		return "", errors.New("error while creating the file")
	}

	// Get the bytes from the file
	buffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffer, file); err != nil {
		return "", errors.New("error while reading the file")
	}

	// Copy the file bytes to the empty file
	if _, err := emptyFile.Write(buffer.Bytes()); err != nil {
		return "", errors.New("error while writing the file")
	}

	return uuid.String(), nil
}
