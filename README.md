# QuizWizard

QuizWizard is an interactive quiz application built with Cobra CLI and the LabStack Echo web framework.

# Getting Started

To get started, follow these steps:

- Run `cd api && go run main.go` to start the API.
- In a separate terminal, run `cd cli && go run main.go` to access the CLI.
- While the API is running, use the commands below within the `cli` directory to interact with QuizWizard.

Retrieve a list of the available categories:
```bash
go run main.go categories
```

Start a quiz with the default category (random):
```bash
go run main.go start
```

Start a quiz with a specified category:
```bash
go run main.go start --category computing
```

# Value Added Extras

- Users can select a quiz category using the `--category` flag.
- The quiz category `random` is selected by default.
- Questions are shuffled to make each execution feel unique.
- An interactive interface is used during the quiz to enhance the user experience.

# Next Steps

- Consider thread safety.
- Increase test coverage.
- Implement a database.
- Containerise and deploy.
- Implement a difficulty setting.
- Shuffle the answer options for each question.
