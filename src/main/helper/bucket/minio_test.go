package bucket

import (
	"encoding/base64"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestMinioConfig_UploadViaBase64(t *testing.T) {
	filePath := "./sample.png" // Ganti dengan path file kamu
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	contentType := http.DetectContentType(fileBytes)
	// 2. Encode ke Base64
	base64String := base64.StdEncoding.EncodeToString(fileBytes)
	fullBase64 := fmt.Sprintf("data:%s;base64,%s", contentType, base64String)

	_ = godotenv.Load()
	minioCof := NewMinio()
	log.Println(minioCof.Bucket())
	log.Println(minioCof.Endpoint())
	info, url := minioCof.UploadViaBase64(fullBase64, "try")
	log.Println(info)
	log.Println(url)
}
