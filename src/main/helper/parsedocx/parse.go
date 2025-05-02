package parsedocx

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
)

func ParseDocxPilihanGanda(fileBytes []byte, filename string) ([]PilihanGanda, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	if _, err = part.Write(fileBytes); err != nil {
		return nil, err
	}
	writer.Close()

	resp, err := http.Post("http://172.17.0.1:8085/parse-docx", writer.FormDataContentType(), body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result []PilihanGanda
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
