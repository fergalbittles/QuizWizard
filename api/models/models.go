package models

import (
	"math/rand"
)

// Question represents a quiz question
type Question struct {
	ID                 int      `json:"id"`
	Category           string   `json:"category"`
	Question           string   `json:"question"`
	Answers            []string `json:"answers"`
	CorrectAnswerIndex int      `json:"correctAnswerIndex"`
}

// Questions represents a group of questions
type Questions []Question

// QuestionResponse represents a response to a quiz question
type QuestionResponse struct {
	Question *Question `json:"question"`
	Answer   int       `json:"answer"`
}

// QuizResponse represents a list of question responses
type QuizResponse struct {
	Category          string             `json:"category"`
	QuestionResponses []QuestionResponse `json:"questionResponses"`
}

// ShuffledCopy returns a shuffled copy of the Questions slice. The original slice remains unchanged.
func (q Questions) ShuffledCopy() Questions {
	cpy := make(Questions, len(q))
	copy(cpy, q)

	for i := range cpy {
		j := rand.Intn(i + 1)
		cpy[i], cpy[j] = cpy[j], cpy[i]
	}

	return cpy
}
