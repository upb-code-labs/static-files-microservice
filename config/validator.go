package config

import "github.com/go-playground/validator/v10"

// Custom Validator methods
type CustomValidator struct{}

func (v *CustomValidator) IsArchiveTypeValid(archiveType string) bool {
	validFileTypes := []string{"test", "submission"}
	for _, fileType := range validFileTypes {
		if fileType == archiveType {
			return true
		}
	}

	return false
}

// Custom validator instance
var cv *CustomValidator

func GetCustomValidator() *CustomValidator {
	if cv == nil {
		cv = &CustomValidator{}
	}

	return cv
}

// Go validator
var gv *validator.Validate

// Go validator instance
func GetGoValidator() *validator.Validate {
	if gv == nil {
		gv = validator.New()
	}

	return gv
}
