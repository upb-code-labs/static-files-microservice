package models

import (
	"fmt"

	"github.com/upb-code-labs/static-files-microservice/config"
	"github.com/upb-code-labs/static-files-microservice/utils"
)

func GetTemplate(uuid string) (fileBytes []byte, err error) {
	templatesPath := config.GetEnvironment().TemplatesPath
	file := fmt.Sprintf("%s.zip", uuid)

	return utils.ReadFile(templatesPath, file)
}
