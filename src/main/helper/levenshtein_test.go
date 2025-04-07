package helper

import (
	"fmt"
	"testing"
)

func TestLevenshtein_EvaluateScoreEssay(t *testing.T) {
	var referenceAnswer = "Pancasila adalah lambang negara"
	answers := []string{
		"Saya tidak tau",
		"Pancasila merupakan burung garuda",
		"Pancasila merupakan lambang negara yang memiliki 5 sila",
		"Pancasila adalah lambang negara indonesia",
	}

	for _, answer := range answers {
		similarity := cosineSimilarity(referenceAnswer, answer)
		fmt.Println("Cosine Similarity:", similarity)

		if similarity > 0.5 {
			fmt.Println("Jawaban benar (akurat 50% atau lebih).")
		} else {
			fmt.Println("Jawaban salah.")
		}
	}
}
