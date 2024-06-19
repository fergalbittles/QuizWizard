package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"quizwizard/cli/config"
	"quizwizard/cli/models"

	"github.com/spf13/cobra"
)

// categoriesCmd represents the categories command
var categoriesCmd = &cobra.Command{
	Use:   "categories",
	Short: "Retrieve a list of available quiz categories",
	Long: `
+++ QuizWizard Categories +++

Reach out to the QuizWizard API to retrieve a 
list of the latest quiz categories.
`,
	Run: func(cmd *cobra.Command, args []string) {
		runCategoriesCommand()
	},
}

func init() {
	rootCmd.AddCommand(categoriesCmd)
}

// runCategoriesCommand will handle all of the steps required to fetch and display the categories
func runCategoriesCommand() {
	fmt.Println("\n+++ QuizWizard Categories +++")

	client := &http.Client{}
	categoryResponse, err := fetchCategories(client)
	if err != nil {
		fmt.Println("\nFailed to fetch categories: " + err.Error())
		return
	}

	err = displayCategories(categoryResponse)
	if err != nil {
		fmt.Println("\nFailed to display categories: " + err.Error())
		return
	}
}

// fetchCategories retrieves the latest list of categories from the API
func fetchCategories(client *http.Client) (*models.CategoriesResponse, error) {
	url := config.ApiUrl + "/categories"
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making categories request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading categories response: %v", err)
	}

	var categoryResponse models.CategoriesResponse
	err = json.Unmarshal([]byte(body), &categoryResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling categories response: %v", err)
	}

	if !categoryResponse.Success {
		return nil, fmt.Errorf("error within categories response: %s", categoryResponse.Message)
	}

	return &categoryResponse, nil
}

// displayCategories outputs the list of categories
func displayCategories(categoryResponse *models.CategoriesResponse) error {
	if categoryResponse == nil {
		return errors.New("categories response is nil")
	}

	if !categoryResponse.Success {
		return fmt.Errorf("error within categories response: %s", categoryResponse.Message)
	}

	if len(categoryResponse.Categories) == 0 {
		fmt.Println("\nNo categories are available at the moment")
		return nil
	}

	fmt.Println()
	for i, category := range categoryResponse.Categories {
		msg := fmt.Sprintf("%d. %s", i+1, category)
		fmt.Println(msg)
	}

	return nil
}
