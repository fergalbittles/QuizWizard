package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"unicode"

	"quizwizard/api/globals"
	"quizwizard/api/models"
	"quizwizard/api/utils"

	"github.com/labstack/echo"
)

// response represents the payload which is returned by each API endpoint
type response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// GetCategories retrieves and returns a list of the latest quiz categories
func GetCategories(c echo.Context) error {
	questions := globals.Questions

	if len(questions) == 0 {
		msg := "An unexpected error occurred. Please try again later."
		return prepareResponse(c, false, msg, http.StatusInternalServerError, nil)
	}

	categories := make([]string, len(questions))
	i := 0
	for key := range questions {
		runes := []rune(key)
		runes[0] = unicode.ToUpper(runes[0])
		categories[i] = string(runes)
		i++
	}
	sort.Strings(categories)
	categories = append(categories, "Random")

	return prepareResponse(c, true, "Categories retrieved successfully.", http.StatusOK, categories)
}

// GetQuestions retrieves and returns a list of questions for a specified category
func GetQuestions(c echo.Context) error {
	category := c.QueryParam("category")
	category = strings.Trim(category, " ")
	category = strings.ToLower(category)

	if len(globals.Questions) == 0 {
		msg := "An unexpected error occurred. Please try again later."
		return prepareResponse(c, false, msg, http.StatusInternalServerError, nil)
	}

	if len(category) == 0 {
		category = "random" // Select 'random' as the default category
	}

	if _, ok := globals.Questions[category]; !ok && category != "random" {
		msg := category + " is not a valid category."
		return prepareResponse(c, false, msg, http.StatusNotFound, nil)
	}

	var responseQuestions models.Questions
	if category == "random" {
		// Select random questions from all categories
		responseQuestions = utils.RandomiseQuestions(globals.Questions)
	} else {
		// Shuffle the questions from the selected category
		responseQuestions = globals.Questions[category].ShuffledCopy()
	}

	if len(responseQuestions) == 0 {
		msg := "Currently there are no questions available for the " + category + " category. Please choose a different category or try again later."
		return prepareResponse(c, false, msg, http.StatusNotFound, nil)
	}

	msg := "Questions successfully retrieved from the " + category + " category."
	return prepareResponse(c, true, msg, http.StatusOK, responseQuestions)
}

// SubmitAnswers stores a score for a quiz submission and returns the results
func SubmitAnswers(c echo.Context) error {
	if len(globals.CategoryScores) == 0 {
		msg := "An unexpected error occurred. Please try again later."
		return prepareResponse(c, false, msg, http.StatusInternalServerError, nil)
	}

	var quizSubmission models.QuizResponse
	err := c.Bind(&quizSubmission)
	if err != nil {
		msg := "Invalid request format."
		return prepareResponse(c, false, msg, http.StatusBadRequest, nil)
	}

	if len(quizSubmission.QuestionResponses) == 0 {
		return prepareResponse(c, false, "No answers were submitted.", http.StatusBadRequest, nil)
	}

	category := quizSubmission.Category
	category = strings.Trim(category, " ")
	category = strings.ToLower(category)

	if len(category) == 0 {
		return prepareResponse(c, false, "A category must be provided.", http.StatusBadRequest, nil)
	}

	if _, ok := globals.CategoryScores[category]; !ok {
		msg := category + " is not a valid category."
		return prepareResponse(c, false, msg, http.StatusNotFound, nil)
	}

	// Calculate the score
	scoreString, scorePercentage, err := utils.CalculateScore(quizSubmission.QuestionResponses)
	if err != nil {
		msg := "Failed to process submission: " + err.Error()
		return prepareResponse(c, false, msg, http.StatusBadRequest, nil)
	}

	// Calculate the comparison percentage
	comparisonScore, err := utils.CalculateComparison(category, scorePercentage)
	if err != nil {
		msg := "Failed to process submission: " + err.Error()
		return prepareResponse(c, false, msg, http.StatusBadRequest, nil)
	}

	// Update the global scores map
	err = utils.AppendCategoryScore(category, scorePercentage)
	if err != nil {
		msg := "Failed to process submission: " + err.Error()
		return prepareResponse(c, false, msg, http.StatusBadRequest, nil)
	}

	comparisonString := ""
	if len(globals.CategoryScores[category]) <= 1 {
		comparisonString = fmt.Sprintf("You are the first quizzer for the %s category.", category)
	} else {
		comparisonString = fmt.Sprintf("Your score for the %s category was better than %.0f%% of all quizzers.", category, comparisonScore)
	}

	res := map[string]interface{}{
		"scoreString":     scoreString,
		"scorePercentage": scorePercentage,
		"comparison":      comparisonString,
	}
	return prepareResponse(c, true, "Submission processed successfully.", http.StatusOK, res)
}

// prepareResponse prepares the response payload which is returned from each API endpoint
func prepareResponse(c echo.Context, success bool, msg string, statusCode int, data interface{}) error {
	err := &response{
		Success: success,
		Message: msg,
		Data:    data,
	}

	return c.JSON(statusCode, err)
}
