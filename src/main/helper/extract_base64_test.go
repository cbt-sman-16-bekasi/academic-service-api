package helper

import (
	"fmt"
	"testing"
)

func TestExtractBase64Images_Success(t *testing.T) {
	html := `<p>ELEPHANT?</p><img src="data:image/jpeg;base64,/9j/4AAQSkZJRgABAQEAYABgAAD/2wBDAAMCAgMCAgMDAwMEAwMEBQgFBQ"/>`

	images := ExtractBase64Images(html)

	for _, img := range images {
		fmt.Println("Found Base64 Image:", img)
	}
}
