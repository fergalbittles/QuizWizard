package globals

import "quizwizard/api/models"

// Questions stores questions for each category
var Questions = make(map[string]models.Questions)

// CategoryScores stores percentage scores for each category
var CategoryScores = make(map[string][]float64)
