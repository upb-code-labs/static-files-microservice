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

// errorMessages is a map that contains the error messages used in more than one controller
var errorMessages = map[string]string{
	"NOT_FOUND":             "File not found",
	"NO_DESTINATION_FOLDER": "No destination folder found for the given file type",
	"NOT_VALID_REQUEST":     "Please, make sure you are sending a valid request",
	"ERROR_READING_FILE":    "Error while reading the file",
	"ERROR_DETECTING_TYPE":  "Error while detecting the file type",
	"WRONG_CONTENT_TYPE":    "The file content type must be %s, but it is %s",
}

var zipContentType = "application/zip"

func SaveArchiveController(c *gin.Context) {
	// Get the data from the multipart/form-data request
	file, err := c.FormFile("file")
	typeField := c.PostForm("archive_type")

	// Check if the fields are valid
	fileIsNotValid := err != nil || file == nil
	if fileIsNotValid || !config.GetCustomValidator().IsArchiveTypeValid(typeField) {
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
			"message": errorMessages["ERROR_READING_FILE"],
		})
		return
	}
	defer file_bytes.Close()

	// Check if the file is a zip file
	mtype, err := mimetype.DetectReader(file_bytes)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": errorMessages["ERROR_DETECTING_TYPE"],
		})
		return
	}

	if mtype.String() != zipContentType {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf(
				errorMessages["WRONG_CONTENT_TYPE"],
				zipContentType,
				mtype.String(),
			),
		})
		return
	}

	// Get the destination folder according to the file type
	destinationFolder, err := utils.GetArchivePathFromFileType(typeField)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": errorMessages["NO_DESTINATION_FOLDER"],
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
	fileIsNotValid := err != nil || file == nil
	fileTypeNotValid := !config.GetCustomValidator().IsArchiveTypeValid(typeField)
	fileIdNotValid := config.GetGoValidator().Var(fileUUID, "required,uuid4") != nil

	if fileIsNotValid || fileTypeNotValid || fileIdNotValid {
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
			"message": errorMessages["ERROR_READING_FILE"],
		})
		return
	}
	defer file_bytes.Close()

	// Check if the file is a zip file
	mtype, err := mimetype.DetectReader(file_bytes)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": errorMessages["ERROR_DETECTING_TYPE"],
		})
		return
	}

	if mtype.String() != zipContentType {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf(
				errorMessages["WRONG_CONTENT_TYPE"],
				zipContentType,
				mtype.String(),
			),
		})
		return
	}

	// Get the destination folder according to the file type
	destinationFolder, err := utils.GetArchivePathFromFileType(typeField)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": errorMessages["NO_DESTINATION_FOLDER"],
		})
		return
	}

	// Check if the file exists
	if !models.DoesFileExists(destinationFolder, fileUUID) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": errorMessages["NOT_FOUND"],
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
			"message": errorMessages["NOT_VALID_REQUEST"],
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
			"message": errorMessages["NO_DESTINATION_FOLDER"],
		})
		return
	}

	// Check if the file exists
	if !models.DoesFileExists(destinationFolder, request.ArchiveUUID) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": errorMessages["NOT_FOUND"],
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
			"message": errorMessages["NOT_VALID_REQUEST"],
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
			"message": errorMessages["NO_DESTINATION_FOLDER"],
		})
		return
	}

	// Check if the file exists
	if !models.DoesFileExists(destinationFolder, request.ArchiveUUID) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": errorMessages["NOT_FOUND"],
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

	c.Data(http.StatusOK, zipContentType, fileBytes)
}
