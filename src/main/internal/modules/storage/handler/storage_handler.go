package handler

import (
	"io"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/bucket"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/response"
	"github.com/gin-gonic/gin"
)

// StorageHandler handles file storage operations
type StorageHandler struct{}

// NewStorageHandler creates a new StorageHandler
func NewStorageHandler() *StorageHandler {
	return &StorageHandler{}
}

// DownloadFile godoc
// @Summary Download file from storage
// @Description Download file from Minio bucket
// @Tags Storage
// @Accept json
// @Produce octet-stream
// @Param bucketName path string true "Bucket name"
// @Param folder path string true "Folder name"
// @Param objectName path string true "Object/file name"
// @Success 200 {file} file "File downloaded"
// @Failure 500 {object} response.BaseResponse
// @Router /academic/{bucketName}/{folder}/{objectName}/download [get]
func (h *StorageHandler) DownloadFile(c *gin.Context) {
	bucketName := c.Param("bucketName")
	objectName := c.Param("objectName")
	folder := c.Param("folder")

	minioBucket := bucket.NewMinio()
	object, err := minioBucket.RetrieveObject(bucketName, folder+"/"+objectName)
	if err != nil {
		panic(err)
	}
	defer object.Close()

	// Set headers for file download
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	c.Header("Content-Disposition", "attachment; filename="+objectName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Access-Control-Expose-Headers", "Content-Disposition")
	c.Header("Expires", "0")

	// Stream file to response
	if _, err := io.Copy(c.Writer, object); err != nil {
		response.InternalError(c, "Failed download file", nil)
		return
	}
}
