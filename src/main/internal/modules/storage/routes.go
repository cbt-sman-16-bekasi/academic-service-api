package storage

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/storage/handler"
	"github.com/gin-gonic/gin"
)

// RegisterRoutesWithDI registers routes with injected handler
func RegisterRoutesWithDI(router *gin.RouterGroup, storageHandler *handler.StorageHandler) {
	academic := router.Group("/academic")

	// File download from Minio bucket
	// Pattern: /academic/{bucketName}/{folder}/{objectName}/download
	academic.GET("/:bucketName/:folder/:objectName/download", storageHandler.DownloadFile)
}
