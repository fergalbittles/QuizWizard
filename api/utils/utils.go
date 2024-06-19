package utils

import (
	"errors"
	"fmt"
	"quizwizard/api/globals"
	"quizwizard/api/models"
	"strings"
)

// RandomizeQuestions selects up to five random questions from all categories
func RandomiseQuestions(questions map[string]models.Questions) models.Questions {
	// Aggregate all questions from all categories
	allQuestions := models.Questions{}
	for _, qs := range questions {
		allQuestions = append(allQuestions, qs...)
	}

	// Shuffle the aggregated questions
	shuffledQuestions := allQuestions.ShuffledCopy()

	// Select up to five questions
	if len(shuffledQuestions) > 5 {
		return shuffledQuestions[:5]
	}
	return shuffledQuestions
}

// CalculateScore returns the score of a quiz submission as a string and also a percentage
func CalculateScore(responses []models.QuestionResponse) (string, float64, error) {
	if len(responses) == 0 {
		msg := "no answers were submitted"
		return "", 0, errors.New(msg)
	}

	score := 0
	for _, response := range responses {
		if response.Question == nil {
			msg := "one or more answers were invalid"
			return "", 0, errors.New(msg)
		}
		if response.Question.CorrectAnswerIndex == response.Answer {
			score++
		}
	}

	totalQuestions := len(responses)
	scoreString := fmt.Sprintf("%d/%d", score, totalQuestions)
	scorePercentage := (float64(score) / float64(totalQuestions)) * 100

	return scoreString, scorePercentage, nil
}

// CalculateComparison calculates the percentage of users a score is better than
func CalculateComparison(category string, newScore float64) (float64, error) {
	category = strings.Trim(category, " ")
	category = strings.ToLower(category)

	if len(category) == 0 {
		msg := "a category must be provided"
		return 0.0, errors.New(msg)
	}

	scores, exists := globals.CategoryScores[category]
	if !exists {
		msg := "category '" + category + "' does not exist"
		return 0.0, errors.New(msg)
	}

	if newScore < 0.0 || newScore > 100.0 {
		msg := "score must be a value between 0 and 100"
		return 0.0, errors.New(msg)
	}

	if len(scores) == 0 {
		// First submission for the specified category
		return 0, nil
	}

	betterThanCount := 0
	for _, score := range scores {
		if newScore > score {
			betterThanCount++
		}
	}

	comparisonPercentage := (float64(betterThanCount) / float64(len(scores))) * 100
	return comparisonPercentage, nil
}

// AppendCategoryScore stores a new score for a specific category
func AppendCategoryScore(category string, newScore float64) error {
	category = strings.Trim(category, " ")
	category = strings.ToLower(category)

	if len(category) == 0 {
		msg := "a category must be provided"
		return errors.New(msg)
	}

	if _, ok := globals.CategoryScores[category]; !ok {
		msg := "category '" + category + "' does not exist"
		return errors.New(msg)
	}

	if newScore < 0.0 || newScore > 100.0 {
		msg := "score must be a value between 0 and 100"
		return errors.New(msg)
	}

	globals.CategoryScores[category] = append(globals.CategoryScores[category], newScore)
	return nil
}
