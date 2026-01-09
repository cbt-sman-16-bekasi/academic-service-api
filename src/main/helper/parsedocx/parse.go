package parsedocx

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

func PilihanGanda(fileBytes []byte, filename string) ([]ResultParse, error) {
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

	log.Info().Msgf("Call to parse docx with fileName %s and url http://5.181.217.35:8085/parse-docx", filename)
	resp, err := http.Post("http://5.181.217.35:8085/parse-docx", writer.FormDataContentType(), body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result []ResultParse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	log.Info().Msgf("Parse docx with result size: %v", len(result))

	return result, nil
}

func Essay(fileBytes []byte, filename string) ([]ResultParse, error) {
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

	log.Info().Msgf("Call to parse docx with fileName %s  and url http://5.181.217.35:8085/parse-docx/essay", filename)
	resp, err := http.Post("http://5.181.217.35:8085/parse-docx/essay", writer.FormDataContentType(), body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result []ResultParse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	log.Info().Msgf("Parse docx with result size: %v", len(result))

	return result, nil
}

func StripHTML(input string) string {
	re := regexp.MustCompile(`(?s)<.*?>`)
	output := re.ReplaceAllString(input, "")
	return strings.TrimSpace(output)
}
