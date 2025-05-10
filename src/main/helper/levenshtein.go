package helper

import (
	"github.com/jdkato/prose/v2"
	"math"
	"regexp"
	"strings"
)

type CosineSimilarity struct {
	ReferenceAnswer string
	Answer          string
	ScoreQuestion   int
}

func NewCosineSimilarity(reference, answer string, score int) *CosineSimilarity {
	return &CosineSimilarity{
		ReferenceAnswer: reference,
		Answer:          answer,
		ScoreQuestion:   score,
	}
}

func (l *CosineSimilarity) EvaluateScoreEssay() int {
	similarity := cosineSimilarity(l.Answer, l.ReferenceAnswer)
	return int(similarity / float64(l.ScoreQuestion))
}

func preprocessText(text string) []string {
	text = strings.ToLower(text)
	re := regexp.MustCompile("[^a-zA-Z0-9 ]+")
	text = re.ReplaceAllString(text, "")
	doc, _ := prose.NewDocument(text)
	var tokens []string
	for _, tok := range doc.Tokens() {
		tokens = append(tokens, tok.Text)
	}
	return tokens
}

func cosineSimilarity(text1, text2 string) float64 {
	tokens1 := preprocessText(text1)
	tokens2 := preprocessText(text2)

	freq1 := make(map[string]int)
	freq2 := make(map[string]int)

	for _, token := range tokens1 {
		freq1[token]++
	}
	for _, token := range tokens2 {
		freq2[token]++
	}

	dotProduct := 0.0
	norm1 := 0.0
	norm2 := 0.0

	for token, freq := range freq1 {
		dotProduct += float64(freq * freq2[token])
		norm1 += float64(freq * freq)
	}

	for _, freq := range freq2 {
		norm2 += float64(freq * freq)
	}

	if norm1 == 0 || norm2 == 0 {
		return 0.0
	}

	return convertScore(dotProduct / (math.Sqrt(norm1) * math.Sqrt(norm2)))
}

func convertScore(score float64) float64 {
	if score > 0.85 {
		return 100.0
	} else if score > 0.70 {
		return 90.0
	} else if score > 0.50 {
		return 75.0
	} else if score > 0.30 {
		return 60.0
	}
	return 40.0
}
