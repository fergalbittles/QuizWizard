package utils

import (
	"quizwizard/api/globals"
	"quizwizard/api/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRandomiseQuestions tests the RandomiseQuestions utility function
func TestRandomiseQuestions(t *testing.T) {
	questions := map[string]models.Questions{
		"science": {
			{ID: 1, Category: "science", Question: "What is the chemical symbol for water?", Answers: []string{"H2O", "O2", "H2O2", "HO"}, CorrectAnswerIndex: 0},
			{ID: 2, Category: "science", Question: "What planet is known as the Red Planet?", Answers: []string{"Earth", "Mars", "Jupiter", "Venus"}, CorrectAnswerIndex: 1},
		},
		"math": {
			{ID: 3, Category: "math", Question: "What is 2 + 2?", Answers: []string{"3", "4", "5", "6"}, CorrectAnswerIndex: 1},
			{ID: 4, Category: "math", Question: "What is the square root of 16?", Answers: []string{"3", "4", "5", "6"}, CorrectAnswerIndex: 1},
		},
		"history": {
			{ID: 5, Category: "history", Question: "Who was the first president of the United States?", Answers: []string{"George Washington", "Thomas Jefferson", "Abraham Lincoln", "John Adams"}, CorrectAnswerIndex: 0},
			{ID: 6, Category: "history", Question: "In which year did the Titanic sink?", Answers: []string{"1912", "1913", "1914", "1915"}, CorrectAnswerIndex: 0},
		},
		"geography": {
			{ID: 7, Category: "geography", Question: "What is the capital of France?", Answers: []string{"Berlin", "Madrid", "Paris", "Lisbon"}, CorrectAnswerIndex: 2},
		},
	}

	tests := []struct {
		name            string
		questions       map[string]models.Questions
		expectedCount   int
		checkCategories bool
	}{
		{
			name:          "more_than_five_total_questions",
			questions:     questions,
			expectedCount: 5,
		},
		{
			name: "less_than_five_total_questions",
			questions: map[string]models.Questions{
				"science": {questions["science"][0]},
				"math":    {questions["math"][0]},
			},
			expectedCount: 2,
		},
		{
			name:            "ensure_categories_are_mixed",
			questions:       questions,
			expectedCount:   5,
			checkCategories: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RandomiseQuestions(tt.questions)
			assert.Equal(t, tt.expectedCount, len(result), "Unexpected number of questions returned")

			if tt.checkCategories {
				categories := map[string]struct{}{}
				for _, question := range result {
					categories[question.Category] = struct{}{}
				}
				assert.Greater(t, len(categories), 1, "Expected questions from multiple categories")
			}
		})
	}
}

// TestCalculateScore tests the CalculateScore utility function
func TestCalculateScore(t *testing.T) {
	questions := []models.Question{
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
	}

	tests := []struct {
		name            string
		responses       []models.QuestionResponse
		expectedString  string
		expectedPercent float64
		expectedError   string
	}{
		{
			name: "success_all_correct_answers",
			responses: []models.QuestionResponse{
				{Question: &questions[0], Answer: 0},
				{Question: &questions[1], Answer: 1},
				{Question: &questions[2], Answer: 0},
				{Question: &questions[3], Answer: 2},
			},
			expectedString:  "4/4",
			expectedPercent: 100.0,
			expectedError:   "",
		},
		{
			name: "success_all_wrong_answers",
			responses: []models.QuestionResponse{
				{Question: &questions[0], Answer: 1},
				{Question: &questions[1], Answer: 0},
				{Question: &questions[2], Answer: 1},
				{Question: &questions[3], Answer: 0},
			},
			expectedString:  "0/4",
			expectedPercent: 0.0,
			expectedError:   "",
		},
		{
			name:            "failure_empty_responses_slice",
			responses:       []models.QuestionResponse{},
			expectedString:  "",
			expectedPercent: 0,
			expectedError:   "no answers were submitted",
		},
		{
			name: "failure_nil_question",
			responses: []models.QuestionResponse{
				{Question: nil, Answer: 0},
				{Question: &questions[1], Answer: 1},
				{Question: &questions[2], Answer: 0},
				{Question: &questions[3], Answer: 2},
			},
			expectedString:  "",
			expectedPercent: 0,
			expectedError:   "one or more answers were invalid",
		},
		{
			name: "positive_answer_out_of_bounds",
			responses: []models.QuestionResponse{
				{Question: &questions[0], Answer: 5},
				{Question: &questions[1], Answer: 1},
				{Question: &questions[2], Answer: 0},
				{Question: &questions[3], Answer: 2},
			},
			expectedString:  "3/4",
			expectedPercent: 75.0,
			expectedError:   "",
		},
		{
			name: "negative_answer_out_of_bounds",
			responses: []models.QuestionResponse{
				{Question: &questions[0], Answer: -3},
				{Question: &questions[1], Answer: 1},
				{Question: &questions[2], Answer: 0},
				{Question: &questions[3], Answer: 2},
			},
			expectedString:  "3/4",
			expectedPercent: 75.0,
			expectedError:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultString, resultPercent, err := CalculateScore(tt.responses)
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedString, resultString)
				assert.Equal(t, tt.expectedPercent, resultPercent)
			}
		})
	}
}

// TestCalculateComparison tests the CalculateComparison utility function
func TestCalculateComparison(t *testing.T) {
	// Save the original CategoryScores map to restore it later
	originalCategoryScores := globals.CategoryScores
	defer func() { globals.CategoryScores = originalCategoryScores }()

	// Mock CategoryScores for testing
	globals.CategoryScores = map[string][]float64{
		"science": {50.0, 60.0, 70.0, 80.0, 90.0},
		"math":    {20.0, 30.0, 40.0, 50.0, 60.0},
		"music":   {},
	}

	tests := []struct {
		name          string
		category      string
		newScore      float64
		expected      float64
		expectedError string
	}{
		{"success_science_category", "science", 75.0, 60.0, ""},
		{"success_math_category_with_whitespace_and_capitals", "  MATH   ", 35.0, 40.0, ""},
		{"success_first_submission", "music", 60.0, 0.0, ""},
		{"failure_invalid_category", "history", 85.0, 0.0, "category 'history' does not exist"},
		{"failure_whitespace_category_string", "    ", 85.0, 0.0, "a category must be provided"},
		{"failure_empty_category_string", "", 85.0, 0.0, "a category must be provided"},
		{"failure_positive_score_out_of_bounds", "science", 105.0, 0.0, "score must be a value between 0 and 100"},
		{"failure_negative_score_out_of_bounds", "science", -5.0, 0.0, "score must be a value between 0 and 100"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CalculateComparison(tt.category, tt.newScore)
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// TestAppendCategoryScore tests the AppendCategoryScore utility function
func TestAppendCategoryScore(t *testing.T) {
	// Save the original CategoryScores map to restore it later
	originalCategoryScores := globals.CategoryScores
	defer func() { globals.CategoryScores = originalCategoryScores }()

	// Mock CategoryScores for testing
	globals.CategoryScores = map[string][]float64{
		"science": {50.0, 60.0, 70.0, 80.0, 90.0},
		"math":    {20.0, 30.0, 40.0, 50.0, 60.0},
		"music":   {},
	}

	tests := []struct {
		name           string
		category       string
		newScore       float64
		expectedError  string
		expectedScores []float64
	}{
		{"success_science_category", "science", 75.0, "", []float64{50.0, 60.0, 70.0, 80.0, 90.0, 75.0}},
		{"success_music_category_with_empty_slice", "music", 35.0, "", []float64{35.0}},
		{"failure_invalid_category", "history", 85.0, "category 'history' does not exist", nil},
		{"failure_empty_category_string", "", 85.0, "a category must be provided", nil},
		{"failure_whitespace_category_string", "     ", 85.0, "a category must be provided", nil},
		{"failure_positive_score_out_of_bounds", "science", 105.0, "score must be a value between 0 and 100", nil},
		{"failure_negative_score_out_of_bounds", "science", -5.0, "score must be a value between 0 and 100", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AppendCategoryScore(tt.category, tt.newScore)
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedScores, globals.CategoryScores[tt.category])
			}
		})
	}
}
