package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/upb-code-labs/static-files-microservice/config"
	"github.com/upb-code-labs/static-files-microservice/models"
	"github.com/upb-code-labs/static-files-microservice/utils"
)

func SaveArchiveController(c *gin.Context) {
	// Get the data from the multipart/form-data request
	file, err := c.FormFile("file")
	typeField := c.PostForm("file_type")

	// Check if the fields are valid
	file_is_not_valid := err != nil || file == nil
	if file_is_not_valid || !config.GetValidator().IsArchiveTypeValid(typeField) {
		c.JSON(400, gin.H{
			"message": "Please, make sure you are sending a valid file and a valid file type",
		})
		return
	}

	file_content_type := file.Header.Get("Content-Type")
	if file_content_type != "application/zip" {
		c.JSON(400, gin.H{
			"message": fmt.Sprintf("The file content type must be application/zip, but it is %s", file_content_type),
		})
		return
	}

	// Get the bytes from the file
	file_bytes, err := file.Open()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error while reading the file",
		})
		return
	}
	defer file_bytes.Close()

	// Get the destination folder according to the file type
	destinationFolder, err := utils.GetArchivePathFromFileType(typeField)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error while identifying destination folder from file type",
		})
		return
	}

	// Save the file
	uuid, err := models.SaveArchive(destinationFolder, file_bytes)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error while saving the file",
		})
		return
	}

	c.JSON(200, gin.H{
		"uuid": uuid,
	})
}

func OverwriteArchiveController(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello world",
	})
}

func DeleteArchiveController(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello world",
	})
}

func GetArchiveController(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello world",
	})
}
