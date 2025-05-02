package bucket

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/yon-module/yon-framework/logger"
	"io"
	"mime"
	"os"
	"strings"
)

type MinioConfig struct {
	endpoint    string
	bucket      string
	accessKey   string
	secretKey   string
	useSSL      bool
	minioClient *minio.Client
}

func NewMinio() *MinioConfig {
	return &MinioConfig{
		endpoint:  os.Getenv("MINIO_ENDPOINT"),
		bucket:    os.Getenv("MINIO_BUCKET"),
		accessKey: os.Getenv("MINIO_ACCESS_KEY"),
		secretKey: os.Getenv("MINIO_SECRET_KEY"),
		useSSL:    os.Getenv("MINIO_SSL") == "true",
	}
}

func (conf *MinioConfig) Bucket() string {
	return conf.bucket
}

func (conf *MinioConfig) AccessKey() string {
	return conf.accessKey
}

func (conf *MinioConfig) SecretKey() string {
	return conf.secretKey
}

func (conf *MinioConfig) UseSSL() bool {
	return conf.useSSL
}

func (conf *MinioConfig) Endpoint() string {
	return conf.endpoint
}

func (conf *MinioConfig) createConnection() (*minio.Client, error) {
	minioClient, err := minio.New(conf.endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.accessKey, conf.secretKey, ""),
		Secure: conf.useSSL,
	})
	if err != nil {
		return nil, err
	}

	logger.Log.Info().Msg("Connected to Minio")
	conf.minioClient = minioClient
	return minioClient, nil
}

func (conf *MinioConfig) UploadViaBase64(base64String string, folder string) (minio.UploadInfo, string) {
	connection, err := conf.createConnection()
	fmt.Println(connection)
	fmt.Println(err)
	fileData, contentType, ext, err := conf.parseBase64File(base64String)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	err = conf.ensureBucket(ctx, conf.bucket)
	if err != nil {
		panic(err)
	}
	fileName := fmt.Sprintf(folder+"/file_%s%s", uuid.New().String(), ext)
	uploadInfo, err := conf.minioClient.PutObject(ctx, conf.bucket, fileName, bytes.NewReader(fileData), int64(len(fileData)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		panic(err)
	}

	publicURL := conf.generatePublicURL(fileName)

	return uploadInfo, publicURL
}

// ✨ Helper: Parse base64
func (conf *MinioConfig) parseBase64File(base64Str string) ([]byte, string, string, error) {
	var contentType, extension string

	// Kalau base64 mengandung "data:image/png;base64,..."
	if strings.Contains(base64Str, ",") {
		parts := strings.SplitN(base64Str, ",", 2)
		if len(parts) != 2 {
			return nil, "", "", errors.New("invalid base64 format")
		}

		meta := parts[0] // data:image/png;base64
		data := parts[1]

		// Extract content-type dari metadata
		if strings.HasPrefix(meta, "data:") && strings.Contains(meta, ";base64") {
			contentType = strings.TrimPrefix(meta[:strings.Index(meta, ";")], "data:")
		}

		// Cari file extension dari content-type
		exts, err := mime.ExtensionsByType(contentType)
		if err != nil || len(exts) == 0 {
			return nil, "", "", errors.New("cannot determine file extension")
		}
		extension = exts[0]

		// Decode base64
		decoded, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			return nil, "", "", err
		}

		return decoded, contentType, extension, nil
	}

	// Kalau tidak ada prefix data:..., asumsikan plain base64
	decoded, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, "", "", err
	}

	// Default (kalau kamu perlu), bisa setting contentType ke octet-stream
	contentType = "application/octet-stream"
	extension = ""

	return decoded, contentType, extension, nil
}

// ✨ Helper: Pastikan bucket ada
func (conf *MinioConfig) ensureBucket(ctx context.Context, bucketName string) error {
	exists, err := conf.minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}
	if !exists {
		err = conf.minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
		fmt.Println("Bucket created:", bucketName)
	}
	return nil
}

// ✨ Helper: Generate Public URL
func (conf *MinioConfig) generatePublicURL(objectName string) string {
	scheme := "http"
	if conf.useSSL {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", scheme, conf.endpoint, conf.bucket, objectName)
}

func (conf *MinioConfig) RetrieveObject(bucketName, objectName string) (io.ReadCloser, error) {
	return conf.minioClient.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
}
