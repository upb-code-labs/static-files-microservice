package models

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/upb-code-labs/static-files-microservice/config"
	"github.com/upb-code-labs/static-files-microservice/utils"
)

func SaveArchive(directory string, uuid string, file multipart.File) (err error) {
	// Create an empty file
	volumePath := config.GetEnvironment().ArchivesVolumePath
	path := fmt.Sprintf("%s/%s/%s.zip", volumePath, directory, uuid)

	emptyFile, err := os.Create(path)
	if err != nil {
		return errors.New("error while creating the file")
	}

	// Get the bytes from the file
	buffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffer, file); err != nil {
		return errors.New("error while reading the file")
	}

	// Copy the file bytes to the empty file
	if _, err := emptyFile.Write(buffer.Bytes()); err != nil {
		return errors.New("error while writing the file")
	}

	return nil
}

func DoesFileExists(directory string, uuid string) bool {
	volumePath := config.GetEnvironment().ArchivesVolumePath
	path := fmt.Sprintf("%s/%s/", volumePath, directory)
	file := fmt.Sprintf("%s.zip", uuid)

	return utils.DoesFileExists(path, file)
}

func GetArchive(directory string, uuid string) (fileBytes []byte, err error) {
	// Get the file path
	volumePath := config.GetEnvironment().ArchivesVolumePath
	path := fmt.Sprintf("%s/%s/", volumePath, directory)
	file := fmt.Sprintf("%s.zip", uuid)

	return utils.ReadFile(path, file)
}

func OverwriteArchive(directory string, uuid string, file multipart.File) (err error) {
	// Delete the file
	err = DeleteArchive(directory, uuid)
	if err != nil {
		return errors.New("error while deleting the file")
	}

	// Save the file
	err = SaveArchive(directory, uuid, file)
	if err != nil {
		return errors.New("error while saving the file")
	}

	return nil
}

func DeleteArchive(directory string, uuid string) (err error) {
	// Get the file path
	volumePath := config.GetEnvironment().ArchivesVolumePath
	path := fmt.Sprintf("%s/%s/%s.zip", volumePath, directory, uuid)

	// Delete the file
	err = os.Remove(path)
	if err != nil {
		return errors.New("error while deleting the file")
	}

	return nil
}
