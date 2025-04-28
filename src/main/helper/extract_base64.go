package helper

import "regexp"

func ExtractBase64Images(html string) []string {
	re := regexp.MustCompile(`(?i)<img[^>]+src="(data:image\/[^;]+;base64,[^"]+)"`)
	matches := re.FindAllStringSubmatch(html, -1)

	var base64Images []string
	for _, match := range matches {
		if len(match) > 1 {
			base64Images = append(base64Images, match[1])
		}
	}
	return base64Images
}
