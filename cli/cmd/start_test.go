package cmd

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"quizwizard/cli/config"
	"quizwizard/cli/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFetchQuestions tests the fetchQuestions function
func TestFetchQuestions(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Header.Get("X-Error-Scenario") {
		case "read_error":
			w.Header().Set("Content-Length", "1")
			w.WriteHeader(http.StatusOK)
		case "unmarshal_error":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("asdasda"))
		case "api_error":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"success": false, "message": "API error"}`))
		default:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"success": true, "message": "Questions retrieved successfully", "data": []}`))
		}
	}))
	defer mockServer.Close()

	originalApiUrl := config.ApiUrl
	config.ApiUrl = mockServer.URL
	defer func() { config.ApiUrl = originalApiUrl }()

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return url.Parse(mockServer.URL)
			},
		},
	}

	tests := []struct {
		name           string
		errorScenario  string
		expectedError  string
		expectedResult *models.QuestionsResponse
	}{
		{
			name:           "successful_response",
			errorScenario:  "",
			expectedError:  "",
			expectedResult: &models.QuestionsResponse{Success: true, Message: "Questions retrieved successfully", Questions: []models.Question{}},
		},
		{
			name:           "failure_due_to_read_error",
			errorScenario:  "read_error",
			expectedError:  "error reading fetch questions response",
			expectedResult: nil,
		},
		{
			name:           "failure_due_to_unmarshal_error",
			errorScenario:  "unmarshal_error",
			expectedError:  "error unmarshaling fetch questions response",
			expectedResult: nil,
		},
		{
			name:           "failure_due_to_api_error",
			errorScenario:  "api_error",
			expectedError:  "error within fetch questions response: API error",
			expectedResult: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client.Transport.(*http.Transport).Proxy = func(req *http.Request) (*url.URL, error) {
				req.Header.Set("X-Error-Scenario", tc.errorScenario)
				return url.Parse(mockServer.URL)
			}

			categoryResponse, err := fetchQuestions(client)
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, categoryResponse)
			}
		})
	}
}

// TestRunQuiz tests the runQuiz function
func TestRunQuiz(t *testing.T) {
	tests := []struct {
		name          string
		input         *models.QuestionsResponse
		expectedError string
	}{
		{
			name:          "failure_due_to_nil_questions_response",
			input:         nil,
			expectedError: "fetch questions response is nil",
		},
		{
			name: "failure_due_to_unsuccessful_response",
			input: &models.QuestionsResponse{
				Success: false,
				Message: "API error",
			},
			expectedError: "error within fetch questions response: API error",
		},
		{
			name: "failure_due_to_empty_questions_response",
			input: &models.QuestionsResponse{
				Success:   true,
				Message:   "Success!",
				Questions: []models.Question{},
			},
			expectedError: "\nCurrently there are no questions available for the random category.\n\nPlease choose a different category or try again later.",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := runQuiz(tc.input)
			if tc.expectedError != "" {
				assert.Nil(t, res)
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestSubmitQuiz tests the submitQuiz function
func TestSubmitQuiz(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Header.Get("X-Error-Scenario") {
		case "read_error":
			w.Header().Set("Content-Length", "1")
			w.WriteHeader(http.StatusOK)
		case "unmarshal_error":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("asdasda"))
		case "api_error":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"success": false, "message": "API error"}`))
		default:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"success": true,
				"message": "Submission processed successfully",
				"data": {
					"comparison": "You are better than 90% of users",
					"scorePercentage": 80,
					"scoreString": "4/5"
				}
			}`))
		}
	}))
	defer mockServer.Close()

	originalApiUrl := config.ApiUrl
	config.ApiUrl = mockServer.URL
	defer func() { config.ApiUrl = originalApiUrl }()

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return url.Parse(mockServer.URL)
			},
		},
	}

	tests := []struct {
		name           string
		errorScenario  string
		input          *models.QuizSubmission
		expectedError  string
		expectedResult *models.QuizSubmissionResponse
	}{
		{
			name:          "successful_response",
			errorScenario: "",
			input: &models.QuizSubmission{
				Category: "Science",
				QuestionResponses: []models.QuestionAnswer{
					{Question: &models.Question{ID: 1, CorrectAnswerIndex: 0}, Answer: 0},
				},
			},
			expectedError: "",
			expectedResult: &models.QuizSubmissionResponse{
				Success: true,
				Message: "Submission processed successfully",
				Results: models.Results{
					Comparison:      "You are better than 90% of users",
					ScorePercentage: 80,
					ScoreString:     "4/5",
				},
			},
		},
		{
			name:          "failure_due_to_read_error",
			errorScenario: "read_error",
			input: &models.QuizSubmission{
				Category: "Science",
				QuestionResponses: []models.QuestionAnswer{
					{Question: &models.Question{ID: 1, CorrectAnswerIndex: 0}, Answer: 0},
				},
			},
			expectedError:  "error reading post submission response",
			expectedResult: nil,
		},
		{
			name:          "failure_due_to_unmarshal_error",
			errorScenario: "unmarshal_error",
			input: &models.QuizSubmission{
				Category: "Science",
				QuestionResponses: []models.QuestionAnswer{
					{Question: &models.Question{ID: 1, CorrectAnswerIndex: 0}, Answer: 0},
				},
			},
			expectedError:  "error unmarshaling post submission response",
			expectedResult: nil,
		},
		{
			name:          "failure_due_to_api_error",
			errorScenario: "api_error",
			input: &models.QuizSubmission{
				Category: "Science",
				QuestionResponses: []models.QuestionAnswer{
					{Question: &models.Question{ID: 1, CorrectAnswerIndex: 0}, Answer: 0},
				},
			},
			expectedError:  "error within post submission response: API error",
			expectedResult: nil,
		},
		{
			name:           "nil_submission",
			errorScenario:  "",
			input:          nil,
			expectedError:  "quiz submission is nil",
			expectedResult: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client.Transport.(*http.Transport).Proxy = func(req *http.Request) (*url.URL, error) {
				req.Header.Set("X-Error-Scenario", tc.errorScenario)
				return url.Parse(mockServer.URL)
			}

			submissionResponse, err := submitQuiz(tc.input, client)
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, submissionResponse)
			}
		})
	}
}

// TestDisplayResults tests the displayResults function
func TestDisplayResults(t *testing.T) {
	tests := []struct {
		name          string
		input         *models.QuizSubmissionResponse
		expectedError string
	}{
		{
			name:          "failure_due_to_nil_submission_response",
			input:         nil,
			expectedError: "submission response is nil",
		},
		{
			name: "failure_due_to_unsuccessful_response",
			input: &models.QuizSubmissionResponse{
				Success: false,
				Message: "API error",
			},
			expectedError: "error within quiz submission response: API error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := displayResults(tc.input)
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
