package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// getTestQuestions is a helper function which returns a fresh copy of test questions
func getTestQuestions() Questions {
	return Questions{
		{
			ID:                 1,
			Category:           "science",
			Question:           "What is the chemical symbol for water?",
			Answers:            []string{"H2O", "O2", "H2O2", "HO"},
			CorrectAnswerIndex: 0,
		},
		{
			ID:                 2,
			Category:           "math",
			Question:           "What is 2 + 2?",
			Answers:            []string{"3", "4", "5", "6"},
			CorrectAnswerIndex: 1,
		},
		{
			ID:                 3,
			Category:           "history",
			Question:           "Who was the first president of the United States?",
			Answers:            []string{"George Washington", "Thomas Jefferson", "Abraham Lincoln", "John Adams"},
			CorrectAnswerIndex: 0,
		},
		{
			ID:                 4,
			Category:           "geography",
			Question:           "What is the capital of France?",
			Answers:            []string{"Berlin", "Madrid", "Paris", "Lisbon"},
			CorrectAnswerIndex: 2,
		},
		{
			ID:                 5,
			Category:           "books",
			Question:           "Who wrote 'The Shining'?",
			Answers:            []string{"Stephen King", "Mark Twain", "Jane Austen", "Agatha Christie"},
			CorrectAnswerIndex: 0,
		},
	}
}

// TestShuffledCopyLength checks that the shuffled copy is the same length as the original
func TestShuffledCopyLength(t *testing.T) {
	original := getTestQuestions()

	shuffled := original.ShuffledCopy()

	assert.Equal(t, len(original), len(shuffled), "Shuffled copy should have the same length as the original")
}

// TestShuffledCopySameElements checks that the shuffled copy contains the same elements as the original
func TestShuffledCopySameElements(t *testing.T) {
	original := getTestQuestions()

	shuffled := original.ShuffledCopy()

	assert.ElementsMatch(t, original, shuffled, "Shuffled copy should have the same elements as the original")
}

// TestShuffledCopyDifferentOrder checks that the elements of the shuffled copy are in a different order
func TestShuffledCopyDifferentOrder(t *testing.T) {
	original := getTestQuestions()

	shuffled := original.ShuffledCopy()

	sameOrder := true
	for i := range original {
		if original[i].ID != shuffled[i].ID {
			sameOrder = false
			break
		}
	}

	assert.False(t, sameOrder, "Shuffled copy should be in a different order from the original")
}
