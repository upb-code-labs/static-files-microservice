package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/upb-code-labs/static-files-microservice/config"
	"github.com/upb-code-labs/static-files-microservice/models"
	"github.com/upb-code-labs/static-files-microservice/utils"
)

func SaveArchiveController(c *gin.Context) {
	// Get the data from the multipart/form-data request
	file, err := c.FormFile("file")
	typeField := c.PostForm("archive_type")

	// Check if the fields are valid
	file_is_not_valid := err != nil || file == nil
	if file_is_not_valid || !config.GetCustomValidator().IsArchiveTypeValid(typeField) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Please, make sure you are sending a valid file and a valid file type",
		})
		return
	}

	// Get the bytes from the file
	file_bytes, err := file.Open()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while reading the file",
		})
		return
	}
	defer file_bytes.Close()

	// Check if the file is a zip file
	mtype, err := mimetype.DetectReader(file_bytes)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while detecting the file type",
		})
		return
	}

	if mtype.String() != "application/zip" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("The file content type must be application/zip, but it is %s", mtype.String()),
		})
		return
	}

	// Get the destination folder according to the file type
	destinationFolder, err := utils.GetArchivePathFromFileType(typeField)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while identifying destination folder from file type",
		})
		return
	}

	// Generate a uuid
	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while generating a uuid",
		})
		return
	}

	// Save the file
	err = models.SaveArchive(destinationFolder, uuid.String(), file_bytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while saving the file",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"uuid": uuid,
	})
}

func OverwriteArchiveController(c *gin.Context) {
	// Get the data from the multipart/form-data request
	file, err := c.FormFile("file")
	typeField := c.PostForm("archive_type")
	fileUUID := c.PostForm("archive_uuid")

	// Check if the fields are valid
	file_is_not_valid := err != nil || file == nil
	file_type_is_not_valid := !config.GetCustomValidator().IsArchiveTypeValid(typeField)
	file_uuid_is_not_valid := config.GetGoValidator().Var(fileUUID, "required,uuid4") != nil

	if file_is_not_valid || file_type_is_not_valid || file_uuid_is_not_valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Please, make sure you are sending a valid file, a valid file type and a valid file uuid",
		})
		return
	}

	// Get the bytes from the file
	file_bytes, err := file.Open()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while reading the file",
		})
		return
	}
	defer file_bytes.Close()

	// Check if the file is a zip file
	mtype, err := mimetype.DetectReader(file_bytes)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while detecting the file type",
		})
		return
	}

	if mtype.String() != "application/zip" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("The file content type must be application/zip, but it is %s", mtype.String()),
		})
		return
	}

	// Get the destination folder according to the file type
	destinationFolder, err := utils.GetArchivePathFromFileType(typeField)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while identifying destination folder from file type",
		})
		return
	}

	// Check if the file exists
	if !models.DoesFileExists(destinationFolder, fileUUID) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "File not found",
		})
		return
	}

	// Overwrite the file
	err = models.OverwriteArchive(destinationFolder, fileUUID, file_bytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while overwriting the file",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func DeleteArchiveController(c *gin.Context) {
	var request DeleteArchiveRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Please, make sure you are sending a valid request",
		})
		return
	}

	// Check if the fields are valid
	if err := config.GetGoValidator().Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Fields validation failed",
			"errors":  err.Error(),
		})
		return
	}

	// Get the destination folder according to the file type
	destinationFolder, err := utils.GetArchivePathFromFileType(request.ArchiveType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while identifying destination folder from file type",
		})
		return
	}

	// Check if the file exists
	if !models.DoesFileExists(destinationFolder, request.ArchiveUUID) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "File not found",
		})
		return
	}

	// Delete the file
	err = models.DeleteArchive(destinationFolder, request.ArchiveUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while deleting the file",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func GetArchiveController(c *gin.Context) {
	var request DownloadArchiveRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Please, make sure you are sending a valid request",
		})
		return
	}

	// Check if the fields are valid
	if err := config.GetGoValidator().Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Fields validation failed",
			"errors":  err.Error(),
		})
		return
	}

	// Get the destination folder according to the file type
	destinationFolder, err := utils.GetArchivePathFromFileType(request.ArchiveType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while identifying destination folder from file type",
		})
		return
	}

	// Check if the file exists
	if !models.DoesFileExists(destinationFolder, request.ArchiveUUID) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "File not found",
		})
		return
	}

	// Get the file
	fileBytes, err := models.GetArchive(destinationFolder, request.ArchiveUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while getting the file",
		})
		return
	}

	c.Data(http.StatusOK, "application/zip", fileBytes)
}
