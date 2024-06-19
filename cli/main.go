package main

import (
	"fmt"
	"log"
	"os"
	"quizwizard/cli/cmd"
	"quizwizard/cli/config"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	config.ApiUrl = os.Getenv("api_url")
	if config.ApiUrl == "" {
		log.Fatal("Failed to load API URL")
	}

	cmd.Execute()
	fmt.Println()
}
