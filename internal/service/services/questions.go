package service

import (
	"context"
	"database/sql"
	"errors"
	"quiz-service/init/logger"
	"strings"
	"sync"

	"quiz-service/internal/entities"
	"quiz-service/internal/repository"
	"quiz-service/pkg/constants"
)

type Questions struct {
	questionRepo repository.QuestionsRepository
	variantRepo  repository.VariantRepository
	testingRepo  repository.TestingRepository

	log logger.Logging

	mu sync.Mutex
}

func NewQuestions(
	questionRepo repository.QuestionsRepository,
	variantRepo repository.VariantRepository,
	testingRepo repository.TestingRepository,
	log logger.Logging) *Questions {
	return &Questions{
		questionRepo: questionRepo,
		variantRepo:  variantRepo,
		testingRepo:  testingRepo,
		log:          log,
		mu:           sync.Mutex{},
	}
}

func (q *Questions) QuestionAdd(ctx context.Context, variantId int, question *entities.Question) error {
	count, err := q.questionRepo.QuestionCount(ctx, variantId)
	if err != nil {
		return err
	}

	if count >= 5 {
		return constants.ErrorQuestionLimitExceeded
	}

	if err := q.questionRepo.QuestionAdd(ctx, variantId, question); err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return constants.ErrorQuestionAlreadyExists
		}
		q.log.ErrorF("QuestionAdd failed: %v", err)
		return err
	}
	return nil
}

func (q *Questions) QuestionRemove(ctx context.Context, variantId int, question *entities.QuestionRemove) error {
	num, err := q.questionRepo.QuestionRemove(ctx, variantId, question.Question)
	if err != nil {
		q.log.ErrorF("QuestionRemove failed: %v", err)
		return err
	}
	if num == 0 {
		return constants.ErrorQuestionNotFound
	}

	return nil
}

func (q *Questions) QuestionGet(ctx context.Context, variantId, questionId int) (*entities.Question, error) {
	questions, err := q.questionRepo.QuestionGet(ctx, variantId, questionId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constants.ErrorQuestionNotFound
		}
		q.log.ErrorF("QuestionGet failed: %v", err)
		return nil, err
	}
	return questions, nil
}

func (q *Questions) QuestionAccept(ctx context.Context, variantId, userId int, answer string) error {
	test, err := q.testingRepo.TestGet(ctx, userId, variantId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return constants.ErrorTestNotFound
		}
		q.log.ErrorF("QuestionAccept-TestGet failed: %v", err)
		return err
	}

	if err := q.questionRepo.QuestionAccept(ctx, answer, userId, test.ID, variantId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return constants.ErrorQuestionNotFound
		}
		q.log.ErrorF("QuestionAccept failed: %v", err)
		return err
	}

	return nil
}
