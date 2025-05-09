package parsedocx

import (
	"bytes"
	"encoding/json"
	"github.com/yon-module/yon-framework/logger"
	"mime/multipart"
	"net/http"
)

func ParseDocxPilihanGanda(fileBytes []byte, filename string) ([]ResultParse, error) {
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

	logger.Log.Info().Msgf("Call to parse docx with fileName %s and url http://172.17.0.1:8085/parse-docx", filename)
	resp, err := http.Post("http://172.17.0.1:8085/parse-docx", writer.FormDataContentType(), body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result []ResultParse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	logger.Log.Info().Msgf("Parse docx with result size: %v", len(result))

	return result, nil
}

func ParseDocxEssay(fileBytes []byte, filename string) ([]ResultParse, error) {
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

	logger.Log.Info().Msgf("Call to parse docx with fileName %s  and url http://172.17.0.1:8085/parse-docx/essay", filename)
	resp, err := http.Post("http://172.17.0.1:8085/parse-docx/essay", writer.FormDataContentType(), body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result []ResultParse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	logger.Log.Info().Msgf("Parse docx with result size: %v", len(result))

	return result, nil
}
