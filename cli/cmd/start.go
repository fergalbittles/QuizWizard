package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"quizwizard/cli/config"
	"quizwizard/cli/models"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var category string

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the quiz",
	Long: `
+++ QuizWizard Start +++

Start the quiz and optionally specify a category.

If no category is specified, "Random" will be
selected by default.

Run the 'categories' command to retrieve a list
of the latest categories.
`,
	Run: func(cmd *cobra.Command, args []string) {
		startQuiz()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringVarP(&category, "category", "c", "random", "Specify the category for the quiz")
}

// startQuiz will handle all of the steps required to take the quiz and display the results
func startQuiz() {
	fmt.Println("\n+++ QuizWizard Starting +++")

	client := &http.Client{}

	questionsResponse, err := fetchQuestions(client)
	if err != nil {
		invalidCategoryError := strings.Contains(err.Error(), "is not a valid category")
		noQuestionsAvailableError := strings.Contains(err.Error(), "no questions available")

		if invalidCategoryError {
			msg := "\nFailure: " + category + " is not a valid category."
			msg += "\n\nUse the 'categories' command for a list of available categories."
			fmt.Println(msg)
			return
		} else if noQuestionsAvailableError {
			msg := "\nCurrently there are no questions available for the " + category + " category."
			msg += "\n\nPlease choose a different category or try again later."
			fmt.Println(msg)
			return
		} else {
			fmt.Println("\nFailed to fetch questions: " + err.Error())
			return
		}
	}

	quizSubmission, err := runQuiz(questionsResponse)
	if err != nil {
		noQuestionsAvailableError := strings.Contains(err.Error(), "no questions available")
		if noQuestionsAvailableError {
			fmt.Println(err.Error())
			return
		} else {
			fmt.Println("\nFailed to display questions: " + err.Error())
			return
		}
	}

	results, err := submitQuiz(quizSubmission, client)
	if err != nil {
		fmt.Println("\nFailed to submit answers: " + err.Error())
		return
	}

	err = displayResults(results)
	if err != nil {
		fmt.Println("\nFailed to display results: " + err.Error())
		return
	}
}

// fetchQuestions retrieves the questions from the API for a specified category
func fetchQuestions(client *http.Client) (*models.QuestionsResponse, error) {
	category = strings.Trim(category, " ")
	category = strings.ToLower(category)

	url := config.ApiUrl + "/questions?category=" + category
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making fetch questions request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading fetch questions response: %v", err)
	}

	var questionsResponse models.QuestionsResponse
	err = json.Unmarshal([]byte(body), &questionsResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling fetch questions response: %v", err)
	}

	if !questionsResponse.Success {
		return nil, fmt.Errorf("error within fetch questions response: %s", questionsResponse.Message)
	}

	return &questionsResponse, nil
}

// runQuiz allows the user to take the quiz using an interactive interface
func runQuiz(questionsResponse *models.QuestionsResponse) (*models.QuizSubmission, error) {
	if questionsResponse == nil {
		return nil, errors.New("fetch questions response is nil")
	}

	if !questionsResponse.Success {
		return nil, fmt.Errorf("error within fetch questions response: %s", questionsResponse.Message)
	}

	if len(questionsResponse.Questions) == 0 {
		msg := "\nCurrently there are no questions available for the " + category + " category."
		msg += "\n\nPlease choose a different category or try again later."
		return nil, errors.New(msg)
	}

	submission := models.QuizSubmission{
		Category:          category,
		QuestionResponses: make([]models.QuestionAnswer, 0, len(questionsResponse.Questions)),
	}

	fmt.Println("\nYou have selected the " + category + " category.")
	fmt.Printf("Please answer all %d questions.\n", len(questionsResponse.Questions))

	for i, question := range questionsResponse.Questions {
		fmt.Printf("\n+++ Question %d: %s +++\n\n", i+1, question.Question)

		for i, answer := range question.Answers {
			fmt.Printf("%d. %s\n", i+1, answer)
		}

		userAnswer, err := promptUser()
		if err != nil {
			userAnswer = -1
		}
		userAnswer--

		if userAnswer == question.CorrectAnswerIndex {
			fmt.Println("\nCorrect! " + question.Answers[question.CorrectAnswerIndex] + " is the right answer.")
		} else if userAnswer >= 0 && userAnswer < len(question.Answers) {
			fmt.Println("\nIncorrect! " + question.Answers[userAnswer] + " is the wrong answer.")
		} else {
			fmt.Println("\nIncorrect! Your selection was invalid.")
			userAnswer = -1
		}

		// Store the question and answer
		qa := models.QuestionAnswer{
			Question: &question,
			Answer:   userAnswer,
		}
		submission.QuestionResponses = append(submission.QuestionResponses, qa)
	}

	return &submission, nil
}

// submitQuiz sends the selected answers for each question to the API
func submitQuiz(quizSubmission *models.QuizSubmission, client *http.Client) (*models.QuizSubmissionResponse, error) {
	if quizSubmission == nil {
		return nil, errors.New("quiz submission is nil")
	}

	jsonData, err := json.Marshal(quizSubmission)
	if err != nil {
		return nil, fmt.Errorf("error marshaling post submission: %v", err)
	}

	url := config.ApiUrl + "/submit"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating post submission request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending post submission request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading post submission response: %v", err)
	}

	var submissionResponse models.QuizSubmissionResponse
	err = json.Unmarshal([]byte(body), &submissionResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling post submission response: %v", err)
	}

	if !submissionResponse.Success {
		return nil, fmt.Errorf("error within post submission response: %s", submissionResponse.Message)
	}

	return &submissionResponse, nil
}

// displayResults outputs the results of the quiz submission
func displayResults(results *models.QuizSubmissionResponse) error {
	if results == nil {
		return errors.New("submission response is nil")
	}

	if !results.Success {
		return fmt.Errorf("error within quiz submission response: %s", results.Message)
	}

	fmt.Println("\n+++ Quiz Results +++")
	fmt.Println("\nRaw score: " + results.Results.ScoreString)
	fmt.Println("Percentage score: " + fmt.Sprintf("%.0f%%", results.Results.ScorePercentage))
	fmt.Println("\n" + results.Results.Comparison)

	return nil
}

// promptUser asks the user to select an answer by entering an option number
func promptUser() (int, error) {
	fmt.Print("\nEnter option number: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())

	intVal, err := strconv.Atoi(input)
	if err != nil {
		return -1, fmt.Errorf("error parsing user input: %v", err)
	}

	return intVal, nil
}
