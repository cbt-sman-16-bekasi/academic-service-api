package request

type UploadBase64Request struct {
	FileData string `json:"file_data" binding:"required"`
}
