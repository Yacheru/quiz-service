package service

import "quiz-service/internal/repository"

type Testing struct {
	repo repository.TestingRepository
}

func NewTesting(repo repository.TestingRepository) *Testing {
	return &Testing{repo: repo}
}
