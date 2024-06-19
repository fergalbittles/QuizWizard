package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"quizwizard/api/globals"
	"quizwizard/api/handlers"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	err := loadQuestions("questions.json")
	if err != nil {
		log.Fatalf("Failed to load questions: %v", err)
	}

	initialiseScoresMap()

	err = startServer(":1323")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func loadQuestions(filename string) error {
	log.Println("Preparing to load config...")

	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	if err := json.Unmarshal(data, &globals.Questions); err != nil {
		return fmt.Errorf("failed to unmarshal JSON within file %s: %w", filename, err)
	}

	log.Println("Config loaded successfully")
	return nil
}

func initialiseScoresMap() {
	globals.CategoryScores["random"] = []float64{}
	for category := range globals.Questions {
		globals.CategoryScores[category] = []float64{}
	}
}

func startServer(port string) error {
	log.Println("Preparing to start server...")

	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/categories", handlers.GetCategories)
	e.GET("/questions", handlers.GetQuestions)
	e.POST("/submit", handlers.SubmitAnswers)

	return e.Start(port)
}
