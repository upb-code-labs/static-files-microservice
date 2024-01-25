package utils

import (
	"errors"
	"log"
	"os"
)

func GetArchivePathFromFileType(fileType string) (string, error) {
	switch fileType {
	case "test":
		return "tests", nil
	case "submission":
		return "submissions", nil
	default:
		return "", errors.New("file type is not valid")
	}
}

func ReadFile(path string, file string) ([]byte, error) {
	fileBytes, err := os.ReadFile(path + "/" + file)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error while reading the file")
	}

	return fileBytes, nil
}

func DoesFileExists(path string, file string) bool {
	if _, err := os.Stat(path + "/" + file); os.IsNotExist(err) {
		return false
	}

	return true
}
