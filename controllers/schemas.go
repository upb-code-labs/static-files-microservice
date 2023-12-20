package controllers

type DownloadArchiveRequest struct {
	ArchiveUUID string `json:"archive_uuid" validate:"required,uuid4"`
	ArchiveType string `json:"archive_type" validate:"required,oneof=test submission"`
}
