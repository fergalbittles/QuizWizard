package handlers

import (
	"net/http"
	"net/http/httptest"
	"quizwizard/api/globals"
	"quizwizard/api/models"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

// TestGetCategories tests the GetCategories handler function
func TestGetCategories(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name               string
		setup              func()
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "successfully_retrieve_categories",
			setup: func() {
				globals.Questions = map[string]models.Questions{
					"science":   {},
					"math":      {},
					"history":   {},
					"computing": {},
				}
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: `{
                "success": true,
                "message": "Categories retrieved successfully.",
                "data": ["Computing", "History", "Math", "Science", "Random"]
            }`,
		},
		{
			name: "failure_due_to_empty_questions_map",
			setup: func() {
				globals.Questions = map[string]models.Questions{}
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: `{
                "success": false,
                "message": "An unexpected error occurred. Please try again later."
            }`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			req := httptest.NewRequest(http.MethodGet, "/categories", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if assert.NoError(t, GetCategories(c)) {
				assert.Equal(t, tt.expectedStatusCode, rec.Code)
				assert.JSONEq(t, tt.expectedResponse, rec.Body.String())
			}
		})
	}
}

// TestGetQuestions tests the GetQuestions handler function
func TestGetQuestions(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name               string
		setup              func()
		category           string
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "successfully_retrieve_questions_for_science_category",
			setup: func() {
				globals.Questions = map[string]models.Questions{
					"science": {
						{ID: 1, Category: "science", Question: "What is the chemical symbol for water?", Answers: []string{"H2O", "O2", "H2O2", "HO"}, CorrectAnswerIndex: 0},
					},
				}
			},
			category:           "science",
			expectedStatusCode: http.StatusOK,
			expectedResponse: `{
                "success": true,
                "message": "Questions successfully retrieved from the science category.",
                "data": [
                    {
                        "id": 1,
                        "category": "science",
                        "question": "What is the chemical symbol for water?",
                        "answers": ["H2O", "O2", "H2O2", "HO"],
                        "correctAnswerIndex": 0
                    }
                ]
            }`,
		},
		{
			name: "successfully_retrieve_random_questions_by_default",
			setup: func() {
				globals.Questions = map[string]models.Questions{
					"math": {
						{ID: 3, Category: "math", Question: "What is 2 + 2?", Answers: []string{"3", "4", "5", "6"}, CorrectAnswerIndex: 1},
					},
				}
			},
			category:           "", // Random category will be selected by default
			expectedStatusCode: http.StatusOK,
			expectedResponse: `{
                "success": true,
                "message": "Questions successfully retrieved from the random category.",
                "data": [
                    {
                        "id": 3,
                        "category": "math",
                        "question": "What is 2 + 2?",
                        "answers": ["3", "4", "5", "6"],
                        "correctAnswerIndex": 1
                    }
                ]
            }`,
		},
		{
			name: "failure_due_to_invalid_category",
			setup: func() {
				globals.Questions = map[string]models.Questions{
					"science": {
						{ID: 1, Category: "science", Question: "What is the chemical symbol for water?", Answers: []string{"H2O", "O2", "H2O2", "HO"}, CorrectAnswerIndex: 0},
					},
				}
			},
			category:           "history",
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: `{
                "success": false,
                "message": "history is not a valid category."
            }`,
		},
		{
			name: "failure_due_to_no_questions_for_specified_category",
			setup: func() {
				globals.Questions = map[string]models.Questions{
					"history": {},
				}
			},
			category:           "history",
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: `{
                "success": false,
                "message": "Currently there are no questions available for the history category. Please choose a different category or try again later."
            }`,
		},
		{
			name: "failure_due_to_empty_questions_map",
			setup: func() {
				globals.Questions = map[string]models.Questions{}
			},
			category:           "science",
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: `{
                "success": false,
                "message": "An unexpected error occurred. Please try again later."
            }`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			req := httptest.NewRequest(http.MethodGet, "/questions", nil)
			q := req.URL.Query()
			q.Add("category", tt.category)
			req.URL.RawQuery = q.Encode()

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if assert.NoError(t, GetQuestions(c)) {
				assert.Equal(t, tt.expectedStatusCode, rec.Code)
				assert.JSONEq(t, tt.expectedResponse, rec.Body.String())
			}
		})
	}
}

// TestSubmitAnswers tests the SubmitAnswers handler function
func TestSubmitAnswers(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name               string
		setup              func()
		requestBody        string
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "successfully_processed_submission_without_comparison",
			setup: func() {
				globals.Questions = map[string]models.Questions{
					"science": {
						{ID: 1, Category: "science", Question: "What is the chemical symbol for water?", Answers: []string{"H2O", "O2", "H2O2", "HO"}, CorrectAnswerIndex: 0},
						{ID: 2, Category: "science", Question: "What planet is known as the Red Planet?", Answers: []string{"Earth", "Mars", "Jupiter", "Venus"}, CorrectAnswerIndex: 1},
					},
				}
				globals.CategoryScores = map[string][]float64{
					"science": {},
				}
			},
			requestBody: `{
                "category": "science",
                "questionResponses": [
                    {
                        "question": {
                            "id": 1,
                            "category": "science",
                            "question": "What is the chemical symbol for water?",
                            "answers": ["H2O", "O2", "H2O2", "HO"],
                            "correctAnswerIndex": 0
                        },
                        "answer": 0
                    },
                    {
                        "question": {
                            "id": 2,
                            "category": "science",
                            "question": "What planet is known as the Red Planet?",
                            "answers": ["Earth", "Mars", "Jupiter", "Venus"],
                            "correctAnswerIndex": 1
                        },
                        "answer": 1
                    }
                ]
            }`,
			expectedStatusCode: http.StatusOK,
			expectedResponse: `{
                "success": true,
                "message": "Submission processed successfully.",
                "data": {
                    "scoreString": "2/2",
                    "scorePercentage": 100,
                    "comparison": "You are the first quizzer for the science category."
                }
            }`,
		},
		{
			name: "successfully_processed_submission_with_comparison",
			setup: func() {
				globals.Questions = map[string]models.Questions{
					"science": {
						{ID: 1, Category: "science", Question: "What is the chemical symbol for water?", Answers: []string{"H2O", "O2", "H2O2", "HO"}, CorrectAnswerIndex: 0},
						{ID: 2, Category: "science", Question: "What planet is known as the Red Planet?", Answers: []string{"Earth", "Mars", "Jupiter", "Venus"}, CorrectAnswerIndex: 1},
					},
				}
				globals.CategoryScores = map[string][]float64{
					"science": {20, 30},
				}
			},
			requestBody: `{
                "category": "science",
                "questionResponses": [
                    {
                        "question": {
                            "id": 1,
                            "category": "science",
                            "question": "What is the chemical symbol for water?",
                            "answers": ["H2O", "O2", "H2O2", "HO"],
                            "correctAnswerIndex": 0
                        },
                        "answer": 0
                    },
                    {
                        "question": {
                            "id": 2,
                            "category": "science",
                            "question": "What planet is known as the Red Planet?",
                            "answers": ["Earth", "Mars", "Jupiter", "Venus"],
                            "correctAnswerIndex": 1
                        },
                        "answer": 1
                    }
                ]
            }`,
			expectedStatusCode: http.StatusOK,
			expectedResponse: `{
                "success": true,
                "message": "Submission processed successfully.",
                "data": {
                    "scoreString": "2/2",
                    "scorePercentage": 100,
                    "comparison": "Your score for the science category was better than 100% of all quizzers."
                }
            }`,
		},
		{
			name: "failure_due_to_invalid_request_format",
			setup: func() {
				globals.CategoryScores = map[string][]float64{
					"science": {50.0, 60.0},
				}
			},
			requestBody:        `invalid json`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: `{
                "success": false,
                "message": "Invalid request format."
            }`,
		},
		{
			name: "failure_due_to_no_answers_submitted",
			setup: func() {
				globals.CategoryScores = map[string][]float64{
					"science": {50.0, 60.0},
				}
			},
			requestBody: `{
                "category": "science",
                "questionResponses": []
            }`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: `{
                "success": false,
                "message": "No answers were submitted."
            }`,
		},
		{
			name: "failure_due_to_empty_category_string",
			setup: func() {
				globals.CategoryScores = map[string][]float64{
					"science": {50.0, 60.0},
				}
			},
			requestBody: `{
                "category": "",
                "questionResponses": [
                    {
                        "question": {
                            "id": 1,
                            "category": "science",
                            "question": "What is the chemical symbol for water?",
                            "answers": ["H2O", "O2", "H2O2", "HO"],
                            "correctAnswerIndex": 0
                        },
                        "answer": 0
                    }
                ]
            }`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: `{
                "success": false,
                "message": "A category must be provided."
            }`,
		},
		{
			name: "failure_due_to_invalid_category",
			setup: func() {
				globals.CategoryScores = map[string][]float64{
					"science": {50.0, 60.0},
				}
			},
			requestBody: `{
                "category": "history",
                "questionResponses": [
                    {
                        "question": {
                            "id": 1,
                            "category": "science",
                            "question": "What is the chemical symbol for water?",
                            "answers": ["H2O", "O2", "H2O2", "HO"],
                            "correctAnswerIndex": 0
                        },
                        "answer": 0
                    }
                ]
            }`,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: `{
                "success": false,
                "message": "history is not a valid category."
            }`,
		},
		{
			name: "failure_due_to_uninitialized_category_scores",
			setup: func() {
				globals.CategoryScores = map[string][]float64{}
			},
			requestBody: `{
                "category": "science",
                "questionResponses": [
                    {
                        "question": {
                            "id": 1,
                            "category": "science",
                            "question": "What is the chemical symbol for water?",
                            "answers": ["H2O", "O2", "H2O2", "HO"],
                            "correctAnswerIndex": 0
                        },
                        "answer": 0
                    }
                ]
            }`,
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: `{
                "success": false,
                "message": "An unexpected error occurred. Please try again later."
            }`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			req := httptest.NewRequest(http.MethodPost, "/submit", strings.NewReader(tt.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if assert.NoError(t, SubmitAnswers(c)) {
				assert.Equal(t, tt.expectedStatusCode, rec.Code)
				assert.JSONEq(t, tt.expectedResponse, rec.Body.String())
			}
		})
	}
}
