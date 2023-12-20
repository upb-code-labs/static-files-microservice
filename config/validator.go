package config

// Validator methods
type Validator struct{}

func (v *Validator) IsArchiveTypeValid(archiveType string) bool {
	validFileTypes := []string{"test", "submission"}
	for _, fileType := range validFileTypes {
		if fileType == archiveType {
			return true
		}
	}

	return false
}

// Validator instance
var validator *Validator

func GetValidator() *Validator {
	if validator == nil {
		validator = &Validator{}
	}

	return validator
}
