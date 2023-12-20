package utils

import "errors"

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
