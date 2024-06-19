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

// TestFetchCategories tests the fetchCategories function
func TestFetchCategories(t *testing.T) {
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
			w.Write([]byte(`{"success": true, "message": "Categories retrieved successfully", "data": ["Science", "Math", "History"]}`))
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
		expectedResult *models.CategoriesResponse
	}{
		{
			name:           "successful_response",
			errorScenario:  "",
			expectedError:  "",
			expectedResult: &models.CategoriesResponse{Success: true, Message: "Categories retrieved successfully", Categories: []string{"Science", "Math", "History"}},
		},
		{
			name:           "failure_due_to_read_error",
			errorScenario:  "read_error",
			expectedError:  "error reading categories response",
			expectedResult: nil,
		},
		{
			name:           "failure_due_to_unmarshal_error",
			errorScenario:  "unmarshal_error",
			expectedError:  "error unmarshaling categories response",
			expectedResult: nil,
		},
		{
			name:           "failure_due_to_api_error",
			errorScenario:  "api_error",
			expectedError:  "error within categories response: API error",
			expectedResult: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client.Transport.(*http.Transport).Proxy = func(req *http.Request) (*url.URL, error) {
				req.Header.Set("X-Error-Scenario", tc.errorScenario)
				return url.Parse(mockServer.URL)
			}

			categoryResponse, err := fetchCategories(client)
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

// TestDisplayCategories tests the displayCategories function
func TestDisplayCategories(t *testing.T) {
	tests := []struct {
		name          string
		input         *models.CategoriesResponse
		expectedError string
	}{
		{
			name:          "failure_due_to_nil_category_response",
			input:         nil,
			expectedError: "categories response is nil",
		},
		{
			name: "failure_due_to_unsuccessful_response",
			input: &models.CategoriesResponse{
				Success:    false,
				Message:    "API error",
				Categories: []string{},
			},
			expectedError: "error within categories response: API error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := displayCategories(tc.input)
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
