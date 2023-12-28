package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/upb-code-labs/static-files-microservice/config"
	"github.com/upb-code-labs/static-files-microservice/models"
	"github.com/upb-code-labs/static-files-microservice/utils"
)

func DownloadTemplateController(c *gin.Context) {
	// Validate the language uuid
	languageUUID := c.Param("language_uuid")
	languageUUIDIsNotValid := config.GetGoValidator().Var(languageUUID, "required,uuid4") != nil
	if languageUUIDIsNotValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Please, make sure you are sending a valid language uuid",
		})
		return
	}

	// Check if the template exists
	path := config.GetEnvironment().TemplatesVolumePath
	file := languageUUID + ".zip"
	templateExists := utils.DoesFileExists(path, file)
	if !templateExists {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Template not found",
		})
		return
	}

	// Get the template
	fileBytes, err := models.GetTemplate(languageUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while getting the template",
		})
		return
	}

	// Send the template
	c.Data(http.StatusOK, "application/zip", fileBytes)
}
