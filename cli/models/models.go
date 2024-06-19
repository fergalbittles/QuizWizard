package models

// CategoriesResponse represents the response from the get categories API endpoint
type CategoriesResponse struct {
	Success    bool     `json:"success"`
	Message    string   `json:"message"`
	Categories []string `json:"data"`
}

// Question represents a single quiz question
type Question struct {
	ID                 int      `json:"id"`
	Category           string   `json:"category"`
	Question           string   `json:"question"`
	Answers            []string `json:"answers"`
	CorrectAnswerIndex int      `json:"correctAnswerIndex"`
}

// QuestionsResponse represents the response from the get questions API endpoint
type QuestionsResponse struct {
	Success   bool       `json:"success"`
	Message   string     `json:"message"`
	Questions []Question `json:"data"`
}

// QuestionAnswer represents an answer to a quiz question
type QuestionAnswer struct {
	Question *Question `json:"question"`
	Answer   int       `json:"answer"`
}

// QuizSubmission represents a list of question answers
type QuizSubmission struct {
	Category          string           `json:"category"`
	QuestionResponses []QuestionAnswer `json:"questionResponses"`
}

// Results represents the results of a quiz submission
type Results struct {
	Comparison      string  `json:"comparison"`
	ScorePercentage float64 `json:"scorePercentage"`
	ScoreString     string  `json:"scoreString"`
}

// QuizSubmissionResponse represents the response from the submit answers API endpoint
type QuizSubmissionResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Results Results `json:"data"`
}
